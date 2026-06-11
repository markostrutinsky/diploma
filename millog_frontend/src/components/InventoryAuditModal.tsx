import { useState, useEffect } from 'react';
import { Html5QrcodeScanner, Html5Qrcode } from 'html5-qrcode'; // 🔥 Додали Html5Qrcode
import { api, type InventoryItem } from '../api/client';
import toast from 'react-hot-toast';
import ModalPortal from './ModalPortal';
import './InventoryAuditModal.css';

interface AuditItem extends InventoryItem {
  actual_qty: number;
  is_verified: boolean;
}

interface InventoryAuditModalProps {
  warehouseId: string;
  warehouseName: string;
  onClose: () => void;
}

export default function InventoryAuditModal({ warehouseId, warehouseName, onClose }: InventoryAuditModalProps) {
  const [items, setItems] = useState<AuditItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [isProcessing, setIsProcessing] = useState(false);

  // 1. Завантажуємо Snapshot (книжкові залишки)
  useEffect(() => {
    api.inventory.getByWarehouse(warehouseId)
      .then(res => {
        const snapshot = (res || []).map((item) => {
          const qty = Number(item.available) || 0;
          return { ...item, available: qty, actual_qty: qty, is_verified: false };
        });
        setItems(snapshot);
        setLoading(false);
      })
      .catch(() => {
        toast.error("Помилка завантаження залишків складу");
        onClose();
      });
  }, [warehouseId, onClose]);

  // 2. Ініціалізація сканера ЖИВОЇ КАМЕРИ
  useEffect(() => {
    if (!loading) {
      const scanner = new Html5QrcodeScanner(
        "audit-qr-reader", 
        { fps: 10, qrbox: { width: 250, height: 250 } }, 
        false
      );
      
      scanner.render(
        (text) => {
          if (text.startsWith('Omnilog-resource:')) {
            const id = text.split(':')[1];
            handleMatch(id);
          }
        }, 
        () => {} // ігноруємо спам від помилок фокусування камери
      );

      return () => { 
        scanner.clear().catch(e => console.error("Scanner cleanup error", e)); 
      };
    }
  }, [loading]);

  // 3. Обробка успішного сканування (з камери або з файлу)
  const handleMatch = (resourceId: string) => {
    setItems(prev => {
      const found = prev.find(i => i.id === resourceId);
      
      if (found) {
        toast.success(`Знайдено: ${found.name}`, { icon: '🎯', id: `toast-${resourceId}-${Date.now()}` });
        return prev.map(item =>
          item.id === resourceId
            ? { ...item, is_verified: true }
            : item
        );
      }
      
      toast.error("Це майно не належить даному складу!", { id: 'wrong-item' });
      return prev;
    });
  };

  // 🔥 НОВЕ: 4. Обробка завантаженого ФОТО з QR-кодом
  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    const toastId = toast.loading('Аналіз зображення...');
    try {
      // Використовуємо прихований div для сканування статичного файлу
      const html5QrCode = new Html5Qrcode("file-qr-reader");
      const decodedText = await html5QrCode.scanFile(file, false);

      if (decodedText.startsWith('Omnilog-resource:')) {
        const id = decodedText.split(':')[1];
        handleMatch(id);
        toast.success('QR-код успішно зчитано з фото!', { id: toastId });
      } else {
        toast.error('Невідомий формат QR-коду', { id: toastId });
      }

      html5QrCode.clear(); // Очищаємо екземпляр після успіху
    } catch (err) {
      toast.error('Не вдалося знайти або розпізнати QR-код на цьому фото', { id: toastId });
    } finally {
      // Скидаємо інпут, щоб можна було завантажити те саме фото ще раз, якщо треба
      e.target.value = ''; 
    }
  };

  // 5. Ручна зміна кількості (якщо факт не збігається з базою)
  const handleQuantityChange = (resourceId: string, newQty: string) => {
    const val = parseInt(newQty, 10) || 0;
    setItems(prev => prev.map(item => 
      item.id === resourceId ? { ...item, actual_qty: val, is_verified: true } : item
    ));
  };

  // 6. Фіналізація (Відправка на бекенд)
  const handleFinish = async () => {
    setIsProcessing(true);
    const loadingToast = toast.loading('Збереження результатів інвентаризації...');

    try {
      const discrepancies = items
        .filter(item => item.actual_qty !== item.available)
        .map(item => ({
          resource_id: item.id,
          book_quantity: item.available,
          actual_quantity: item.actual_qty,
          difference: item.actual_qty - item.available
        }));

      await api.inventory.submitAudit(warehouseId, discrepancies);
      
      toast.success("Акт переобліку успішно сформовано, залишки оновлено!", { id: loadingToast });
      onClose();
    } catch (err: any) {
      toast.error(err.message || "Помилка збереження результатів", { id: loadingToast });
    } finally {
      setIsProcessing(false);
    }
  };

  if (loading) {
    return (
      <ModalPortal>
        <div className="audit-modal-wrapper">
          <div className="spinner" />
        </div>
      </ModalPortal>
    );
  }

  return (
    <ModalPortal>
      <div className="audit-modal-wrapper">
        <div className="audit-modal-content">
        
        <div className="audit-header">
          <h3>📋 Переоблік: <span className="warehouse-name-highlight">{warehouseName}</span></h3>
          <button className="btn btn-secondary" onClick={onClose} disabled={isProcessing}>
            Закрити без збереження
          </button>
        </div>

        <div className="audit-body">
          {/* ЛІВА ЧАСТИНА: Камера та Кнопка файлу */}
          <div className="scanner-panel">
            <div className="scanner-container">
              <div id="audit-qr-reader"></div>
            </div>

            {/* 🔥 ПРИХОВАНИЙ КОНТЕЙНЕР ДЛЯ СКАНЕРА ФАЙЛІВ */}
            <div id="file-qr-reader" style={{ display: 'none' }}></div>

            {/* 🔥 НОВА КНОПКА ЗАВАНТАЖЕННЯ */}
            <label className="btn btn-secondary" style={{ 
              width: '100%', 
              display: 'flex', 
              justifyContent: 'center', 
              cursor: 'pointer', 
              margin: 0, 
              height: '42px', 
              alignItems: 'center',
              backgroundColor: 'var(--bg-input)',
              border: '1px dashed #94a3b8'
            }}>
              <input
                type="file"
                accept="image/*"
                style={{ display: 'none' }}
                onChange={handleFileUpload}
              />
              📁 Завантажити фото з QR-кодом
            </label>

            <div className="audit-hint-box">
              <strong>💡 Як це працює:</strong><br/>
              Скануйте QR-коди майна через камеру або завантажте фото з галереї. Якщо система знаходить збіг, рядок підсвічується зеленим.<br/><br/>
              Якщо фактична кількість відрізняється від бази — впишіть правильну цифру в колонку <b>«Факт»</b>.
            </div>
          </div>

          {/* ПРАВА ЧАСТИНА: Таблиця звірки */}
          <div className="table-panel">
            <table className="audit-table">
              <thead>
                <tr>
                  <th>Найменування майна</th>
                  <th className="text-center">База</th>
                  <th className="text-center">Факт</th>
                  <th className="text-center">Різниця</th>
                </tr>
              </thead>
              <tbody>
                {items.length === 0 ? (
                  <tr>
                    <td colSpan={4} className="text-center" style={{ padding: '30px', color: 'var(--text-muted)' }}>
                      На цьому складі ще немає майна.
                    </td>
                  </tr>
                ) : (
                  items.map(item => {
                    const diff = item.actual_qty - item.available;
                    return (
                      <tr key={item.id} className={item.is_verified ? 'row-verified' : 'row-pending'}>
                        <td>
                          <span className="status-icon">{item.is_verified ? '✅' : '⏳'}</span>
                          {item.name}
                        </td>
                        <td className="text-center text-muted">{item.available}</td>
                        <td className="text-center">
                          <input 
                            type="number" 
                            min="0"
                            className="qty-input" 
                            value={item.actual_qty}
                            onChange={(e) => handleQuantityChange(item.id, e.target.value)}
                          />
                        </td>
                        <td className="text-center">
                          <span className={`badge ${diff === 0 ? 'badge-success' : diff < 0 ? 'badge-critical' : 'badge-warning'}`}>
                            {diff > 0 ? `+${diff}` : diff}
                          </span>
                        </td>
                      </tr>
                    );
                  })
                )}
              </tbody>
            </table>
          </div>
        </div>

        <div className="audit-footer">
          <button className="btn btn-primary" onClick={handleFinish} disabled={isProcessing}>
            {isProcessing ? 'Обробка...' : '✅ Затвердити результати (Сформувати Акт)'}
          </button>
        </div>
        </div>
      </div>
    </ModalPortal>
  );
}
