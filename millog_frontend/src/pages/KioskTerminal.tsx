import React, { useState, useRef, useEffect } from 'react';
import { api } from '../api/client';
import { useAuth } from '../contexts/AuthContext';
import type { Warehouse, User } from '../api/client';
import toast from 'react-hot-toast';
import { Html5QrcodeScanner } from 'html5-qrcode';
import './KioskTerminal.css';

interface CartItem {
  resource_id: string;
  name: string;
  barcode: string;
  quantity: number;
  maxQty: number;
}

// Ролі, яким показувати ТІЛЬКИ склади свого unit (не всього tenant'у)
const SCOPED_ROLES = [
  'BRANCH_MANAGER', 'BRANCH_LOGISTICIAN', 'BRANCH_STOREKEEPER',
  'DEPT_MANAGER', 'DEPT_SUPERVISOR', 'TEAM_LEAD',
  'REGION_STOREKEEPER', 'REGION_LOGISTICIAN', 'REGION_DIRECTOR',
  'EMPLOYEE', 'CONTRACTOR',
];

export default function KioskTerminal() {
  const { user } = useAuth();

  // ── Крок 0: вибір складу ──────────────────────────────────────
  const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
  const [selectedWarehouse, setSelectedWarehouse] = useState<Warehouse | null>(null);
  const [warehousesLoading, setWarehousesLoading] = useState(true);

  // ── Крок 1: вибір отримувача ──────────────────────────────────
  const [recipient, setRecipient] = useState<User | null>(null);
  const [visibleUsers, setVisibleUsers] = useState<User[]>([]);
  const [usersLoading, setUsersLoading] = useState(false);
  const [recipientSearch, setRecipientSearch] = useState('');
  const recipientInputRef = useRef<HTMLInputElement>(null);

  // ── Крок 2: сканування ───────────────────────────────────────
  const [barcodeInput, setBarcodeInput] = useState('');
  const [cart, setCart] = useState<CartItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [isIssuing, setIsIssuing] = useState(false);
  const [isCameraOpen, setIsCameraOpen] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  // Завантаження складів при відкритті сторінки
  useEffect(() => {
    const loadWarehouses = async () => {
      setWarehousesLoading(true);
      try {
        const all = await api.warehouses.list();
        const isScoped = user && SCOPED_ROLES.includes(user.role);
        const filtered = isScoped && user?.unit_id
          ? all.filter(w => w.unit_id === user.unit_id)
          : all;
        setWarehouses(filtered);
        if (filtered.length === 1) {
          setSelectedWarehouse(filtered[0]);
        }
      } catch {
        toast.error('Помилка завантаження складів');
      } finally {
        setWarehousesLoading(false);
      }
    };
    loadWarehouses();
  }, [user]);

  // Завантаження видимих користувачів при переході до кроку вибору отримувача
  useEffect(() => {
    if (!selectedWarehouse || recipient) return;
    const loadUsers = async () => {
      setUsersLoading(true);
      try {
        const allUsers = await api.users.getVisible();

        // Фільтруємо по підрозділу складу + завжди включаємо поточного користувача
        const filtered = !selectedWarehouse.unit_id
          ? allUsers
          : allUsers.filter(u => u.unit_id === selectedWarehouse.unit_id || u.id === user?.id);

        setVisibleUsers(filtered);
      } catch {
        toast.error('Помилка завантаження списку особового складу');
      } finally {
        setUsersLoading(false);
      }
    };
    loadUsers();
  }, [selectedWarehouse, recipient, user]);

  // Тримаємо фокус для сканера
  useEffect(() => {
    if (selectedWarehouse && recipient && !isIssuing) inputRef.current?.focus();
  }, [cart, selectedWarehouse, recipient, isIssuing]);

  // Фокус на полі пошуку отримувача
  useEffect(() => {
    if (selectedWarehouse && !recipient) {
      setTimeout(() => recipientInputRef.current?.focus(), 100);
    }
  }, [selectedWarehouse, recipient]);

  useEffect(() => {
    if (isCameraOpen) {
      const scanner = new Html5QrcodeScanner(
        "kiosk-qr-reader",
        { fps: 15, qrbox: { width: 250, height: 250 } },
        false
      );
      scanner.render(
        (decodedText) => {
          let finalCode = decodedText;
          if (decodedText.startsWith('Omnilog-resource:')) {
            finalCode = decodedText.split(':')[1];
          }
          scanner.clear();
          setIsCameraOpen(false);
          processScannedCode(finalCode);
        },
        () => {}
      );
      return () => { scanner.clear().catch(() => {}); };
    }
  }, [isCameraOpen]);

  const filteredUsers = visibleUsers.filter(u => {
    const q = recipientSearch.toLowerCase();
    return (
      u.full_name.toLowerCase().includes(q) ||
      (u.username || '').toLowerCase().includes(q) ||
      u.email.toLowerCase().includes(q)
    );
  });

  const processScannedCode = async (inputCode: string) => {
    if (!selectedWarehouse) return;
    const input = inputCode.trim();
    if (!input) return;

    setLoading(true);
    try {
      const unitId = user?.unit_id ?? undefined;
      const resources = await api.inventory.listResources(unitId);
      const warehouseResources = resources.filter(
        r => r.warehouse_id === selectedWarehouse.id
      );

      const item = warehouseResources.find(r =>
        r.id === input ||
        (r.barcode && r.barcode === input) ||
        r.name.toLowerCase().includes(input.toLowerCase())
      );

      if (!item) {
        const anyItem = resources.find(r =>
          r.id === input || (r.barcode && r.barcode === input) ||
          r.name.toLowerCase().includes(input.toLowerCase())
        );
        if (anyItem) {
          toast.error(`"${anyItem.name}" є на складі "${anyItem.warehouse_name}", але не на "${selectedWarehouse.name}"`);
        } else {
          toast.error(`"${input}" не знайдено у системі`);
        }
        return;
      }

      if (item.quantity <= 0) {
        toast.error(`"${item.name}" відсутній на складі "${selectedWarehouse.name}"`);
        return;
      }

      setCart(prev => {
        const existing = prev.find(p => p.resource_id === item.id);
        if (existing) {
          if (existing.quantity >= item.quantity) {
            toast.error(`Максимум доступно: ${item.quantity} шт`);
            return prev;
          }
          return prev.map(p =>
            p.resource_id === item.id
              ? { ...p, quantity: p.quantity + 1, maxQty: item.quantity }
              : p
          );
        }
        return [...prev, {
          resource_id: item.id,
          name: item.name,
          barcode: item.barcode || input,
          quantity: 1,
          maxQty: item.quantity,
        }];
      });

      toast.success(`${item.name} додано до кошика`);
    } catch {
      toast.error('Помилка пошуку ресурсу');
    } finally {
      setLoading(false);
      inputRef.current?.focus();
    }
  };

  const handleFormSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await processScannedCode(barcodeInput);
    setBarcodeInput('');
  };

  const changeQty = (id: string, delta: number) => {
    setCart(prev =>
      prev.map(item => {
        if (item.resource_id !== id) return item;
        const newQty = item.quantity + delta;
        if (newQty < 1) return item;
        if (newQty > item.maxQty) {
          toast.error(`Максимум доступно: ${item.maxQty} шт`);
          return item;
        }
        return { ...item, quantity: newQty };
      })
    );
    inputRef.current?.focus();
  };

  const removeFromCart = (id: string) => {
    setCart(prev => prev.filter(item => item.resource_id !== id));
    inputRef.current?.focus();
  };

  const clearCart = () => {
    setCart([]);
    inputRef.current?.focus();
  };

  const handleIssueAll = async () => {
    if (cart.length === 0 || !selectedWarehouse || !recipient) return;
    setIsIssuing(true);
    const loadingToast = toast.loading('Перевіряємо залишки...');
    try {
      const unitId = user?.unit_id ?? undefined;
      const resources = await api.inventory.listResources(unitId);
      const warehouseResources = resources.filter(r => r.warehouse_id === selectedWarehouse.id);
      const errors: string[] = [];

      for (const cartItem of cart) {
        const freshItem = warehouseResources.find(r => r.id === cartItem.resource_id);
        if (!freshItem) {
          errors.push(`"${cartItem.name}" не знайдено на складі "${selectedWarehouse.name}"`);
        } else if (freshItem.quantity < cartItem.quantity) {
          errors.push(`"${cartItem.name}": запит ${cartItem.quantity} шт, доступно ${freshItem.quantity} шт`);
        }
      }

      if (errors.length > 0) {
        toast.error(errors.join('\n'), { id: loadingToast, duration: 6000 });
        setCart(prev =>
          prev.map(cartItem => {
            const fresh = warehouseResources.find(r => r.id === cartItem.resource_id);
            if (!fresh) return cartItem;
            return { ...cartItem, maxQty: fresh.quantity, quantity: Math.min(cartItem.quantity, fresh.quantity) };
          })
        );
        return;
      }

      toast.loading('Проводимо видачу...', { id: loadingToast });
      const failed: string[] = [];
      for (const item of cart) {
        try {
          await api.inventory.issueResource({
            resource_id: item.resource_id,
            user_id: recipient.id,
            quantity: item.quantity,
            warehouse_id: selectedWarehouse.id,
          });
        } catch {
          failed.push(`${item.name} ×${item.quantity}`);
        }
      }

      if (failed.length > 0) {
        toast.error(`Не вдалося видати:\n${failed.join('\n')}`, { id: loadingToast, duration: 6000 });
      } else {
        toast.success(
          `Видано ${cart.reduce((s, i) => s + i.quantity, 0)} шт → ${recipient.full_name}`,
          { id: loadingToast }
        );
        // Після успішної видачі скидаємо сесію (залишаємо склад, скидаємо отримувача і кошик)
        setCart([]);
        setRecipient(null);
        setRecipientSearch('');
      }
    } catch {
      toast.error('Помилка під час видачі', { id: loadingToast });
    } finally {
      setIsIssuing(false);
    }
  };

  const totalQty = cart.reduce((acc, curr) => acc + curr.quantity, 0);

  // ── Екран 0: вибір складу ─────────────────────────────────────
  if (!selectedWarehouse) {
    return (
      <div className="kiosk-mode">
        <div className="kiosk-header">
          <h1>📦 Термінал Видачі Майна</h1>
          <p>Оберіть склад, з якого будете видавати майно</p>
        </div>
        <div className="warehouse-select-screen">
          {warehousesLoading ? (
            <div className="warehouse-loading">⏳ Завантаження складів...</div>
          ) : warehouses.length === 0 ? (
            <div className="warehouse-empty">
              <p>Немає доступних складів для вашого підрозділу.</p>
              <p>Зверніться до адміністратора.</p>
            </div>
          ) : (
            <div className="warehouse-grid">
              {warehouses.map(w => (
                <button
                  key={w.id}
                  className="warehouse-card"
                  onClick={() => setSelectedWarehouse(w)}
                >
                  <div className="warehouse-icon">
                    {w.location_type === 'MOBILE' ? '🚛' : '🏭'}
                  </div>
                  <div className="warehouse-name">{w.name}</div>
                  <div className="warehouse-type">
                    {w.location_type === 'MOBILE' ? 'Мобільний' : 'Стаціонарний'}
                  </div>
                </button>
              ))}
            </div>
          )}
        </div>
      </div>
    );
  }

  // ── Екран 1: вибір отримувача ─────────────────────────────────
  if (!recipient) {
    return (
      <div className="kiosk-mode">
        <div className="kiosk-header">
          <h1>📦 Термінал Видачі Майна</h1>
          <div className="active-warehouse-badge">
            <span className="warehouse-badge-icon">🏭</span>
            <span className="warehouse-badge-name">{selectedWarehouse.name}</span>
            <button
              className="warehouse-change-btn"
              onClick={() => { setSelectedWarehouse(null); setCart([]); setRecipient(null); setRecipientSearch(''); }}
              title="Змінити склад"
            >
              ↩ Змінити
            </button>
          </div>
        </div>
        <div className="recipient-select-screen">
          <h2>👤 Оберіть отримувача майна</h2>
          <p className="recipient-hint">Майно буде закріплено за обраною особою у системі обліку</p>
          <div className="recipient-search-box">
            <input
              ref={recipientInputRef}
              type="text"
              value={recipientSearch}
              onChange={e => setRecipientSearch(e.target.value)}
              placeholder="Пошук за ПІБ, логіном або email..."
              autoComplete="off"
            />
          </div>
          {usersLoading ? (
            <div className="warehouse-loading">⏳ Завантаження особового складу...</div>
          ) : (
            <div className="recipient-list">
              {filteredUsers.length === 0 ? (
                <div className="warehouse-empty">Особовий склад не знайдено</div>
              ) : (
                filteredUsers.map(u => (
                  <button
                    key={u.id}
                    className="recipient-card"
                    onClick={() => { setRecipient(u); setRecipientSearch(''); }}
                  >
                    <div className="recipient-avatar">👤</div>
                    <div className="recipient-info">
                      <strong>{u.full_name}</strong>
                      <span>{u.username ? `@${u.username}` : u.email}</span>
                      <span className="recipient-role">{u.role}</span>
                    </div>
                  </button>
                ))
              )}
            </div>
          )}
        </div>
      </div>
    );
  }

  // ── Екран 2: сканування та видача ─────────────────────────────
  return (
    <div className="kiosk-mode">
      <div className="kiosk-header">
        <h1>📦 Термінал Видачі Майна</h1>
        <div className="kiosk-session-info">
          <div className="active-warehouse-badge">
            <span className="warehouse-badge-icon">🏭</span>
            <span className="warehouse-badge-name">{selectedWarehouse.name}</span>
            <button
              className="warehouse-change-btn"
              onClick={() => { setSelectedWarehouse(null); setCart([]); setRecipient(null); setRecipientSearch(''); }}
              title="Змінити склад"
            >
              ↩ Змінити
            </button>
          </div>
          <div className="active-recipient-badge">
            <span className="recipient-badge-icon">👤</span>
            <span className="recipient-badge-name">{recipient.full_name}</span>
            <button
              className="warehouse-change-btn"
              onClick={() => { setRecipient(null); setCart([]); setRecipientSearch(''); }}
              title="Змінити отримувача"
            >
              ↩ Змінити
            </button>
          </div>
        </div>
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
              disabled={loading || isCameraOpen || isIssuing}
              autoFocus
              autoComplete="off"
            />
            <button type="submit" disabled={loading || isIssuing || !barcodeInput}>
              {loading ? '...' : 'Знайти [Enter]'}
            </button>
          </form>

          <button
            type="button"
            className="btn btn-secondary"
            style={{ width: '100%', marginTop: '16px', display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '8px' }}
            onClick={() => setIsCameraOpen(!isCameraOpen)}
            disabled={isIssuing}
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
            <p>4. Для зміни кількості — кнопки <strong>−</strong> / <strong>+</strong> у кошику.</p>
          </div>
        </div>

        <div className="cart-section">
          <div className="cart-header">
            <h2>До видачі ({totalQty} шт)</h2>
            {cart.length > 0 && (
              <button className="clear-cart-btn" onClick={clearCart} disabled={isIssuing} title="Очистити кошик">
                🗑 Очистити
              </button>
            )}
          </div>

          <div className="cart-list">
            {cart.length === 0 ? (
              <div className="empty-cart">Кошик порожній</div>
            ) : (
              cart.map(item => (
                <div key={item.resource_id} className="cart-item">
                  <div className="item-info">
                    <strong>{item.name}</strong>
                    <span className="item-stock">залишок: {item.maxQty} шт</span>
                  </div>
                  <div className="item-actions">
                    <button
                      className="qty-btn"
                      onClick={() => changeQty(item.resource_id, -1)}
                      disabled={item.quantity <= 1 || isIssuing}
                      aria-label="Зменшити"
                    >−</button>
                    <span className="qty-badge">{item.quantity}</span>
                    <button
                      className="qty-btn"
                      onClick={() => changeQty(item.resource_id, 1)}
                      disabled={item.quantity >= item.maxQty || isIssuing}
                      aria-label="Збільшити"
                    >+</button>
                    <button
                      className="remove-btn"
                      onClick={() => removeFromCart(item.resource_id)}
                      disabled={isIssuing}
                      aria-label="Видалити"
                    >✕</button>
                  </div>
                </div>
              ))
            )}
          </div>

          <div className="cart-footer">
            <button
              className="issue-btn"
              disabled={cart.length === 0 || isIssuing}
              onClick={handleIssueAll}
            >
              {isIssuing ? 'Виконується...' : `Видати → ${recipient.full_name}`}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
