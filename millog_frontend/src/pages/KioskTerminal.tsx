import React, { useState, useRef, useEffect } from 'react';
import { api } from '../api/client';
import toast from 'react-hot-toast';
import { Html5QrcodeScanner } from 'html5-qrcode';
import './KioskTerminal.css'; 

interface CartItem {
  resource_id: string; // Використовуємо UUID
  name: string;
  barcode: string;
  quantity: number;
}

export default function KioskTerminal() {
  const [barcodeInput, setBarcodeInput] = useState('');
  const [cart, setCart] = useState<CartItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [isCameraOpen, setIsCameraOpen] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  // Тримаємо фокус для сканера
  useEffect(() => {
    inputRef.current?.focus();
  }, [cart]);

  useEffect(() => {
    if (isCameraOpen) {
      const scanner = new Html5QrcodeScanner(
        "kiosk-qr-reader",
        { fps: 15, qrbox: { width: 250, height: 250 } },
        false
      );

      scanner.render(
        (decodedText) => {
          // Обробка твоїх фірмових QR-кодів (якщо вони мають префікс)
          let finalCode = decodedText;
          if (decodedText.startsWith('Omnilog-resource:')) {
            finalCode = decodedText.split(':')[1];
          }
          
          scanner.clear();
          setIsCameraOpen(false);
          processScannedCode(finalCode);
        },
        () => { /* ігноруємо помилки під час пошуку коду в кадрі */ }
      );

      return () => {
        scanner.clear().catch(e => console.error("Помилка зупинки сканера", e));
      };
    }
  }, [isCameraOpen]);

  // Спільна функція для обробки будь-якого отриманого коду
  const processScannedCode = async (inputCode: string) => {
    const input = inputCode.trim();
    if (!input) return;

    setLoading(true);
    try {
      const resources = await api.inventory.listResources();
      
      // Пошук за UUID, Barcode або назвою
      const item = resources.find(r => 
        r.id === input || 
        (r.barcode && r.barcode === input) || 
        r.name.toLowerCase().includes(input.toLowerCase())
      );

      if (!item) {
        toast.error(`Код ${input} не знайдено!`);
        return;
      }

      if (item.quantity <= 0) {
        toast.error(`"${item.name}" відсутній на складі`);
        return;
      }

      setCart(prev => {
        const existing = prev.find(p => p.resource_id === item.id);
        if (existing) {
          if (existing.quantity >= item.quantity) {
             toast.error('Недостатньо на складі');
             return prev;
          }
          return prev.map(p => p.resource_id === item.id ? { ...p, quantity: p.quantity + 1 } : p);
        }
        return [...prev, { resource_id: item.id, name: item.name, barcode: input, quantity: 1 }];
      });

      toast.success(`${item.name} в кошику`);
    } catch (error) {
      toast.error('Помилка пошуку');
    } finally {
      setLoading(false);
      inputRef.current?.focus();
    }
  };

  // Викликається при натисканні Enter у полі вводу
  const handleFormSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await processScannedCode(barcodeInput);
    setBarcodeInput('');
  };

  const handleIssueAll = async () => {
    if (cart.length === 0) return;
    const loadingToast = toast.loading('Проводимо видачу...');
    try {
      for (const item of cart) {
         // Викликаємо існуючий метод списання зі складу
         await api.inventory.writeOffResource(item.resource_id, item.quantity);
      }
      toast.success('Успішно видано!', { id: loadingToast });
      setCart([]);
    } catch (error) {
      toast.error('Помилка під час видачі', { id: loadingToast });
    }
    inputRef.current?.focus();
  };

  const removeFromCart = (id: string) => {
    setCart(prev => prev.filter(item => item.resource_id !== id));
    inputRef.current?.focus();
  };

  return (
    <div className="kiosk-mode">

      <div className="kiosk-header">
        <h1>📦 Термінал Швидкої Видачі</h1>
        <p>Готовий до сканування. Використовуйте сканер, QR або введіть назву вручну.</p>
      </div>

      <div className="kiosk-layout">
        <div className="scanner-section">
          <div className="scanner-status">
            <div className="pulse-dot"></div>
            <span>Сканер активний</span>
          </div>
          
          <form onSubmit={handleFormSubmit} className="scan-form">
            <input
              ref={inputRef}
              type="text"
              value={barcodeInput}
              onChange={(e) => setBarcodeInput(e.target.value)}
              placeholder="Штрих-код / QR / Назва..."
              disabled={loading || isCameraOpen}
              autoFocus
              autoComplete="off"
            />
            <button type="submit" disabled={loading || !barcodeInput}>
              {loading ? '...' : 'Знайти [Enter]'}
            </button>
          </form>

          {/* КНОПКА ТА КОНТЕЙНЕР КАМЕРИ */}
          <button 
            type="button"
            className="btn btn-secondary" 
            style={{ width: '100%', marginTop: '16px', display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '8px' }}
            onClick={() => setIsCameraOpen(!isCameraOpen)}
          >
            {isCameraOpen ? '❌ Закрити камеру' : '📷 Сканувати камерою'}
          </button>

          {isCameraOpen && (
            <div style={{ marginTop: '20px', background: 'var(--bg-input)', borderRadius: '8px', overflow: 'hidden', border: '1px solid var(--border)' }}>
              <div id="kiosk-qr-reader"></div>
            </div>
          )}
          
          <div className="instructions">
            <p>1. Використовуйте апаратний сканер.</p>
            <p>2. Або натисніть кнопку камери вище.</p>
            <p>3. Або введіть назву/ID вручну.</p>
          </div>
        </div>

        <div className="cart-section">
          <h2>До видачі ({cart.reduce((acc, curr) => acc + curr.quantity, 0)} шт)</h2>
          
          <div className="cart-list">
            {cart.length === 0 ? (
              <div className="empty-cart">Кошик порожній</div>
            ) : (
              cart.map(item => (
                <div key={item.resource_id} className="cart-item">
                  <div className="item-info">
                    <strong>{item.name}</strong>
                  </div>
                  <div className="item-actions">
                    <span className="qty-badge">x {item.quantity}</span>
                    <button className="remove-btn" onClick={() => removeFromCart(item.resource_id)}>✕</button>
                  </div>
                </div>
              ))
            )}
          </div>

          <div className="cart-footer">
            <button className="issue-btn" disabled={cart.length === 0} onClick={handleIssueAll}>
              Підтвердити Видачу
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}