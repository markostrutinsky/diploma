import React, { useEffect, useState } from 'react';
import { api, type MyEquipmentItem, ROLE_NAMES } from '../api/client'; // 🔥 Додали ROLE_NAMES
import { useAuth } from '../contexts/AuthContext';
import toast from 'react-hot-toast';
import './Profile.css';

const STATUS_LABELS: Record<string, { label: string, color: string }> = {
  ACTIVE: { label: 'Активний', color: 'success' },
  PENDING: { label: 'Очікує', color: 'warning' },
  BLOCKED: { label: 'Заблокований', color: 'danger' },
};

const PERMISSIONS = {
  manageWarehouses: ['ADMIN', 'REGION_DIRECTOR', 'REGION_LOGISTICIAN', 'BRANCH_MANAGER', 'BRANCH_LOGISTICIAN', 'DEPT_MANAGER', 'TEAM_LEAD'],
  approveRequests: ['ADMIN', 'REGION_DIRECTOR', 'BRANCH_MANAGER', 'DEPT_MANAGER', 'REGION_LOGISTICIAN', 'BRANCH_LOGISTICIAN'],
  createRequests: ['ADMIN', 'REGION_DIRECTOR', 'BRANCH_MANAGER', 'DEPT_MANAGER', 'TEAM_LEAD', 'REGION_LOGISTICIAN', 'BRANCH_LOGISTICIAN', 'DEPT_SUPERVISOR'],
  manageInventory: ['ADMIN', 'REGION_STOREKEEPER', 'BRANCH_STOREKEEPER', 'DEPT_SUPERVISOR'],
  manageUnits: ['ADMIN', 'REGION_DIRECTOR', 'BRANCH_MANAGER', 'DEPT_MANAGER', 'REGION_LOGISTICIAN', 'BRANCH_LOGISTICIAN', 'REGION_STOREKEEPER', 'BRANCH_STOREKEEPER'],
  createUsers: ['ADMIN', 'REGION_DIRECTOR', 'BRANCH_MANAGER', 'DEPT_MANAGER', 'TEAM_LEAD', 'REGION_LOGISTICIAN', 'BRANCH_LOGISTICIAN', 'REGION_STOREKEEPER', 'BRANCH_STOREKEEPER', 'DEPT_SUPERVISOR'],
};

export default function Profile() {
  const { user } = useAuth();
  const [equipment, setEquipment] = useState<MyEquipmentItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'equipment' | 'permissions' | 'settings'>('equipment');
  
  const [reportItem, setReportItem] = useState<MyEquipmentItem | null>(null);
  const [reportReason, setReportReason] = useState('BROKEN');

  const [showPasswordModal, setShowPasswordModal] = useState(false);
  const [passwords, setPasswords] = useState({ old: '', new: '', confirm: '' });
  const [isChangingPass, setIsChangingPass] = useState(false);

  const [formData, setFormData] = useState({
    full_name: '',
    phone: '',
    username: '',
    email: '',
  });
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    if (user) {
      setFormData({
        full_name: user.full_name || '',
        phone: user.phone || '',
        username: user.username || '',
        email: user.email || '',
      });
    }
  }, [user]);

  useEffect(() => {
    loadMyEquipment();
  }, []);

  const loadMyEquipment = async () => {
    try {
      setLoading(true);
      const data = await api.inventory.getMyEquipment();
      setEquipment(data || []);
    } catch (error) {
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleProfileSave = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      setIsSaving(true);
      await api.users.updateProfile({
        full_name: formData.full_name,
        phone: formData.phone || undefined,
        username: formData.username || undefined,
        email: formData.email,
      });
      toast.success('Профіль успішно оновлено! Оновіть сторінку, щоб побачити зміни.');
    } catch (error: any) {
      toast.error(error.message || 'Сталася помилка при збереженні');
    } finally {
      setIsSaving(false);
    }
  };

  const handleReportSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    alert(`Запит щодо ${reportItem?.resource_name} (Причина: ${reportReason}) успішно відправлено!`);
    setReportItem(null);
  };

  const handlePasswordSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (passwords.new !== passwords.confirm) {
      toast.error('Нові паролі не співпадають!');
      return;
    }
    if (passwords.new.length < 8 || !/[A-Z]/.test(passwords.new) || !/[0-9]/.test(passwords.new) || !/[^A-Za-z0-9]/.test(passwords.new)) {
      toast.error('Пароль має містити мін. 8 символів, 1 велику літеру, 1 цифру та 1 спецсимвол');
      return;
    }
    try {
      setIsChangingPass(true);
      await api.auth.updatePassword({ 
        old_password: passwords.old, 
        new_password: passwords.new 
      });
      toast.success('Пароль успішно змінено!');
      setShowPasswordModal(false);
      setPasswords({ old: '', new: '', confirm: '' });
    } catch (error: any) {
      toast.error(error.message || 'Помилка зміни пароля');
    } finally {
      setIsChangingPass(false);
    }
  };

  const formatUnitType = (type: string) => {
    switch(type) {
      case 'PCS': return 'шт';
      case 'KIT': return 'компл';
      case 'KG': return 'кг';
      case 'L': return 'л';
      default: return 'шт';
    }
  };

  const getInitials = (name: string) => {
    if (!name) return '👤';
    const parts = name.split(' ');
    if (parts.length >= 2) return (parts[0][0] + parts[1][0]).toUpperCase();
    return name[0].toUpperCase();
  };

  const hasPerm = (list: string[]) => user?.role && list.includes(user.role);
  const userStatus = user?.status || 'ACTIVE';
  const statusInfo = STATUS_LABELS[userStatus] || STATUS_LABELS['ACTIVE'];

  if (loading) {
    return <div className="page-loading"><div className="spinner" /></div>;
  }

  return (
    <div className="profile-page">
      <div className="profile-header">
        <h1>Особистий кабінет</h1>
        <p className="subtitle">Управління майном та налаштування акаунта</p>
      </div>

      <div className="profile-layout">
        {/* ЛІВА КОЛОНКА */}
        <div className="profile-sidebar">
          <div className="profile-card">
            <div className="avatar-circle">
              {getInitials(user?.full_name || '')}
            </div>
            <h2 className="profile-name">{user?.full_name || 'Ім\'я не вказано'}</h2>
            
            <div className="profile-badges">
              <span className="profile-role-badge">
                {user?.role ? ROLE_NAMES[user.role] || user.role : 'Посада не вказана'}
              </span>
              <span className={`status-badge status-${statusInfo.color}`}>
                {statusInfo.label}
              </span>
            </div>

            <div className="profile-details">
              {user?.username && (
                <div className="detail-item">
                  <span className="detail-label">Логін (Username)</span>
                  <span className="detail-value">@{user.username}</span>
                </div>
              )}
              <div className="detail-item">
                <span className="detail-label">Email</span>
                <span className="detail-value">{user?.email}</span>
              </div>
              <div className="detail-item">
                <span className="detail-label">Телефон</span>
                <span className="detail-value">{user?.phone || 'Не вказано'}</span>
              </div>
              <div className="detail-item">
                <span className="detail-label">Орг. одиниця</span>
                <span className="detail-value font-mono">
                  {user?.unit_id ? `ID: ${user.unit_id}` : 'Без закріплення'}
                </span>
              </div>
            </div>
          </div>
        </div>

        {/* ПРАВА КОЛОНКА */}
        <div className="profile-content">
          <div className="profile-tabs">
            <button 
              className={`ptab-btn ${activeTab === 'equipment' ? 'active' : ''}`}
              onClick={() => setActiveTab('equipment')}
            >
              📦 Моє майно
            </button>
            <button 
              className={`ptab-btn ${activeTab === 'permissions' ? 'active' : ''}`}
              onClick={() => setActiveTab('permissions')}
            >
              🔐 Повноваження
            </button>
            <button 
              className={`ptab-btn ${activeTab === 'settings' ? 'active' : ''}`}
              onClick={() => setActiveTab('settings')}
            >
              ⚙️ Налаштування
            </button>
          </div>

          <div className="tab-body">
            
            {activeTab === 'equipment' && (
              <>
                {equipment.length === 0 ? (
                  <div className="empty-state-card">
                    <span className="empty-icon">📦</span>
                    <h3>У вас немає закріпленого майна</h3>
                    <p>Коли логіст видасть вам робоче майно, воно з'явиться тут.</p>
                  </div>
                ) : (
                  <div className="equipment-grid">
                    {equipment.map((item) => (
                      <div key={item.assignment_id} className="equipment-card">
                        <div className="card-top">
                          <span className="status-indicator">● На руках</span>
                          <span className="issue-date">
                            Видано: {new Date(item.issued_at).toLocaleDateString('uk-UA')}
                          </span>
                        </div>
                        <h3 className="item-name">{item.resource_name}</h3>
                        <div className="item-quantity">
                          <span className="qty-number">{item.quantity}</span>
                          <span className="qty-unit">{formatUnitType(item.unit_type)}</span>
                        </div>
                        <button className="btn-report" onClick={() => setReportItem(item)}>
                          ⚠️ Повідомити про проблему
                        </button>
                      </div>
                    ))}
                  </div>
                )}
              </>
            )}

            {activeTab === 'permissions' && (
              <div className="permissions-panel">
                <h3>Права доступу в системі</h3>
                <p className="text-muted mb-4">
  Ваші можливості в системі визначаються вашою посадою (<strong>{user?.role ? ROLE_NAMES[user.role] || user.role : 'Невідомо'}</strong>).
</p>
                
                <div className="permissions-list">
                  <div className="perm-item">
                    <span className="perm-icon">{hasPerm(PERMISSIONS.approveRequests) ? '✅' : '❌'}</span>
                    <div className="perm-text">
                      <strong>Затвердження заявок</strong>
                      <p>Право розглядати та погоджувати заявки на постачання від співробітників.</p>
                    </div>
                  </div>
                  <div className="perm-item">
                    <span className="perm-icon">{hasPerm(PERMISSIONS.manageWarehouses) ? '✅' : '❌'}</span>
                    <div className="perm-text">
                      <strong>Управління інфраструктурою складів</strong>
                      <p>Право створювати та редагувати інформацію про склади.</p>
                    </div>
                  </div>
                  <div className="perm-item">
                    <span className="perm-icon">{hasPerm(PERMISSIONS.manageInventory) ? '✅' : '❌'}</span>
                    <div className="perm-text">
                      <strong>Управління балансом майна</strong>
                      <p>Право додавати нове майно на склади, видавати та списувати його.</p>
                    </div>
                  </div>
                  <div className="perm-item">
                    <span className="perm-icon">{hasPerm(PERMISSIONS.createUsers) ? '✅' : '❌'}</span>
                    <div className="perm-text">
                      <strong>Управління персоналом</strong>
                      <p>Право створювати профілі для нових співробітників.</p>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {activeTab === 'settings' && (
              <div className="settings-panel">
                <h3 style={{ marginBottom: '24px' }}>Редагування профілю</h3>
                
                <form onSubmit={handleProfileSave} className="profile-form">
                  <div className="form-row">
                    <div className="form-group">
                      <label>ПІБ (Повне ім'я)</label>
                      <input 
                        type="text" 
                        name="full_name"
                        className="erp-input" 
                        value={formData.full_name}
                        onChange={handleInputChange}
                        required
                      />
                    </div>
                    
                    <div className="form-group">
                      <label>Телефон</label>
                      <input 
                        type="tel" 
                        name="phone"
                        className="erp-input" 
                        placeholder="+380..."
                        value={formData.phone}
                        onChange={handleInputChange}
                      />
                    </div>
                  </div>

                  <div className="form-row">
                    <div className="form-group">
                      <label>Логін (Username)</label>
                      <input 
                        type="text" 
                        name="username"
                        className="erp-input" 
                        value={formData.username}
                        onChange={handleInputChange}
                      />
                    </div>

                    <div className="form-group">
                      <label>Email</label>
                      <input 
                        type="email" 
                        name="email"
                        className="erp-input" 
                        value={formData.email}
                        onChange={handleInputChange}
                        required
                      />
                    </div>
                  </div>

                  <div className="form-actions" style={{ marginTop: '24px', textAlign: 'right' }}>
                    <button type="submit" className="btn btn-primary" disabled={isSaving}>
                      {isSaving ? 'Збереження...' : 'Зберегти зміни'}
                    </button>
                  </div>
                </form>

                <hr style={{ margin: '32px 0', borderColor: '#e2e8f0' }} />

                <h3 style={{ color: '#dc2626', marginBottom: '8px' }}>Небезпечна зона</h3>
                <p className="text-muted mb-4">Якщо ви зміните пароль, вам доведеться зайти в систему наново.</p>
                <button 
                  type="button" 
                  className="btn btn-secondary" 
                  onClick={() => setShowPasswordModal(true)}
                >
                  Змінити пароль
                </button>
              </div>
            )}
          </div>
        </div>
      </div>

      {reportItem && (
        <div className="modal-overlay" onClick={() => setReportItem(null)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <h3>Списання / Заміна майна</h3>
            <p className="text-muted text-left" style={{ marginBottom: '15px' }}>
              Майно: <strong>{reportItem.resource_name}</strong>
            </p>
            <form onSubmit={handleReportSubmit}>
              <div className="form-group text-left">
                <label>Причина запиту</label>
                <select value={reportReason} onChange={(e) => setReportReason(e.target.value)} className="erp-input">
                  <option value="BROKEN">Зламано / Пошкоджено під час роботи</option>
                  <option value="LOST">Втрачено</option>
                  <option value="WORN_OUT">Зношено (закінчився термін експлуатації)</option>
                </select>
              </div>
              <div className="warning-box" style={{ marginTop: '16px' }}>
                <p>Цей запит буде автоматично відправлено вашому керівнику та відділу логістики.</p>
              </div>
              <div className="modal-actions" style={{ marginTop: '24px' }}>
                <button type="button" className="btn btn-secondary cancel-margin" onClick={() => setReportItem(null)}>
                  Скасувати
                </button>
                <button type="submit" className="btn btn-danger">Відправити запит</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {showPasswordModal && (
        <div className="modal-overlay" onClick={() => setShowPasswordModal(false)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <h3>Зміна пароля</h3>
            <form onSubmit={handlePasswordSubmit}>
              <div className="form-group text-left" style={{ marginBottom: '15px' }}>
                <label>Поточний пароль</label>
                <input 
                  type="password" 
                  className="erp-input" 
                  value={passwords.old} 
                  onChange={e => setPasswords({...passwords, old: e.target.value})} 
                  required 
                />
              </div>
              <div className="form-group text-left" style={{ marginBottom: '15px' }}>
                <label>Новий пароль (мінімум 8 символів)</label>
                <input 
                  type="password" 
                  className="erp-input" 
                  value={passwords.new} 
                  onChange={e => setPasswords({...passwords, new: e.target.value})} 
                  required 
                  minLength={8}
                />
              </div>
              <div className="form-group text-left" style={{ marginBottom: '20px' }}>
                <label>Підтвердіть новий пароль</label>
                <input 
                  type="password" 
                  className="erp-input" 
                  value={passwords.confirm} 
                  onChange={e => setPasswords({...passwords, confirm: e.target.value})} 
                  required 
                  minLength={8}
                />
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary cancel-margin" onClick={() => setShowPasswordModal(false)}>
                  Скасувати
                </button>
                <button type="submit" className="btn btn-primary" disabled={isChangingPass}>
                  {isChangingPass ? 'Збереження...' : 'Змінити пароль'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}