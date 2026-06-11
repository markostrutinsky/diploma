import { useEffect, useState } from 'react';
import { api, ROLE_NAMES } from '../api/client';
import toast from 'react-hot-toast';
import Pagination from '../components/Pagination';
import { useAuth } from '../contexts/AuthContext';
import './AuditLogs.css';

interface AuditLog {
  id: number;
  user_email: string;
  user_role: string;
  action_type: string;
  entity_type: string;
  entity_id: string;
  details: string;
  created_at: string;
}

// Українські переклади дій у системі
const ACTION_LABELS: Record<string, { label: string; cls: string; icon: string }> = {
  CREATE: { label: 'Створення', cls: 'badge-create', icon: '✨' },
  UPDATE: { label: 'Оновлення', cls: 'badge-update', icon: '✏️' },
  DELETE: { label: 'Видалення', cls: 'badge-delete', icon: '🗑️' },
  EXPORT: { label: 'Експорт', cls: 'badge-primary', icon: '📤' },
  REPORT: { label: 'Рапорт', cls: 'badge-warning', icon: '📝' },
  REFRESH: { label: 'Оновлення токена', cls: 'badge-neutral', icon: '🔁' },
  INVENTORY_AUDIT: { label: 'Переоблік', cls: 'badge-primary', icon: '📋' },
  WRITE_OFF: { label: 'Списання', cls: 'badge-warning', icon: '📉' },
  ASSIGN: { label: 'Видача персоналу', cls: 'badge-primary', icon: '👤' },
  RETURN: { label: 'Повернення на склад', cls: 'badge-primary', icon: '↩️' },
  APPROVE: { label: 'Погодження', cls: 'badge-create', icon: '✅' },
  REJECT: { label: 'Відхилення', cls: 'badge-delete', icon: '❌' },
  CANCEL: { label: 'Скасування', cls: 'badge-neutral', icon: '🚫' },
  DISPATCH: { label: 'Відправка рейсу', cls: 'badge-primary', icon: '🚚' },
  DELIVER: { label: 'Доставка', cls: 'badge-create', icon: '📦' },
  LOGIN: { label: 'Авторизація', cls: 'badge-neutral', icon: '🔑' },
  LOGOUT: { label: 'Вихід', cls: 'badge-neutral', icon: '🚪' },
  SLA_VIOLATION: { label: 'Порушення SLA', cls: 'badge-warning', icon: '⚠️' },
  UNAUTHORIZED_PREMIUM_ACCESS: { label: 'Спроба premium доступу', cls: 'badge-delete', icon: '🔒' },
};

// Українські переклади сутностей
const ENTITY_LABELS: Record<string, string> = {
  VEHICLE: 'Автомобіль',
  WAREHOUSE: 'Склад',
  TENANT: 'Організація',
  RESOURCE: 'Майно / Ресурс',
  SUPPLY_REQUEST: 'Заявка на постачання',
  CONTRACTOR_REQUEST: 'Заявка підряднику',
  CONTRACTOR_MEMBERSHIP: 'Співпраця з підрядником',
  REQUEST: 'Заявка',
  UNIT: 'Орг. одиниця',
  USER: 'Користувач',
  TOKEN: 'Токен доступу',
  SHIPMENT: 'Рейс / Відправка',
  SHIPMENT_REFUEL: 'Дозаправка рейсу',
  FUEL_LOG: 'Журнал пального',
  FUEL_RECORD: 'Запис пального',
  FUEL: 'Пальне',
  MAINTENANCE: 'ТО автомобіля',
  INVENTORY: 'Складські залишки',
  AUDIT: 'Інвентаризація',
  SECURITY: 'Безпека',
  CATEGORY: 'Категорія',
  RESOURCE_ASSIGNMENT: 'Видача майна',
  SMART_REPLENISH: 'Автопоповнення складу',
};

const translateAction = (action: string): { label: string; cls: string; icon: string } => {
  const direct = ACTION_LABELS[action?.toUpperCase()];
  if (direct) return direct;
  const act = (action || '').toUpperCase();
  // Фолбек на підрядки, якщо система десь зберегла дію у нестандартному форматі
  for (const key of Object.keys(ACTION_LABELS)) {
    if (act.includes(key)) return ACTION_LABELS[key];
  }
  return { label: action || '—', cls: 'badge-neutral', icon: '📝' };
};

const translateEntity = (entity: string): string => {
  return ENTITY_LABELS[entity?.toUpperCase()] || entity || '—';
};

const AuditLogs = () => {
  const { user } = useAuth();
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [isRefreshing, setIsRefreshing] = useState(false);

  // Стейт для фільтрів
  const [actionFilter, setActionFilter] = useState('ALL');
  const [entityFilter, setEntityFilter] = useState('ALL');
  const [searchQuery, setSearchQuery] = useState('');
  const [logsPage, setLogsPage] = useState(0);
  const LOGS_PAGE_SIZE = 50;

  const fetchLogs = async (showToast = false) => {
    if (showToast) setIsRefreshing(true);
    try {
      if (!showToast) setLoading(true);
      const data = user?.role === 'SYSTEM_ADMIN'
        ? await api.platform.getAuditLogs()
        : await api.admin.getAuditLogs();
      setLogs(Array.isArray(data) ? data : []);
      if (showToast) toast.success('Журнал оновлено');
    } catch (err: any) {
      toast.error(err.message || 'Помилка доступу до журналу аудиту');
    } finally {
      setLoading(false);
      setIsRefreshing(false);
    }
  };

  useEffect(() => {
    fetchLogs();
  }, [user?.role]);

  useEffect(() => { setLogsPage(0); }, [searchQuery, actionFilter, entityFilter]);

  // Генерація кольорових бейджів дій
  const getActionBadge = (action: string) => {
    const { label, cls, icon } = translateAction(action);
    return (
      <span className={`audit-badge ${cls}`} title={action}>
        <span className="badge-icon">{icon}</span>
        {label}
      </span>
    );
  };

  // Скорочення UUID до красивого хешу
  const getEntityHash = (entityId: string) => {
    if (!entityId) return '';
    return `#${entityId.split('-')[0].toUpperCase()}`;
  };

  // Застосування фільтрів
  const filteredLogs = logs.filter(log => {
    const matchesAction = actionFilter === 'ALL' || log.action_type.toUpperCase().includes(actionFilter);
    const matchesEntity = entityFilter === 'ALL' || log.entity_type.toUpperCase() === entityFilter;
    
    const searchStr = searchQuery.toLowerCase();
    const matchesSearch = searchStr === '' || 
      (log.user_email && log.user_email.toLowerCase().includes(searchStr)) || 
      (log.details && log.details.toLowerCase().includes(searchStr));

    return matchesAction && matchesEntity && matchesSearch;
  });

  // Статистика для заголовку
  const stats = {
    total: logs.length,
    today: logs.filter(l => {
      const d = new Date(l.created_at);
      const now = new Date();
      return d.toDateString() === now.toDateString();
    }).length,
    critical: logs.filter(l => /DELETE|WRITE_OFF|REJECT|SLA_VIOLATION|UNAUTHORIZED/i.test(l.action_type)).length,
  };

  if (loading) {
    return (
      <div className="page-loading">
        <div className="spinner" />
        <p>Завантаження історії дій...</p>
      </div>
    );
  }

  const logsTotalPages = Math.max(1, Math.ceil(filteredLogs.length / LOGS_PAGE_SIZE));
  const safeLogsPage = Math.min(logsPage, logsTotalPages - 1);
  const pagedLogs = filteredLogs.slice(safeLogsPage * LOGS_PAGE_SIZE, (safeLogsPage + 1) * LOGS_PAGE_SIZE);

  return (
    <div className="audit-page">

      
      <div className="page-header">
        <div className="header-title-block">
          <h1>🛡️ Журнал аудиту</h1>
          <p className="audit-subtitle">Протоколювання критичних операцій та дій адміністраторів</p>
        </div>
        <div className="page-actions">
          <button className="btn btn-secondary" onClick={() => fetchLogs(true)} disabled={isRefreshing}>
            {isRefreshing ? '⏳ Оновлення...' : '🔄 Оновити дані'}
          </button>
        </div>
      </div>

      {/* Статистика */}
      <div className="audit-stats-grid">
        <div className="audit-stat-card">
          <div className="audit-stat-icon" style={{ background: 'rgba(59, 130, 246, 0.12)', color: '#2563eb' }}>📚</div>
          <div>
            <div className="audit-stat-val">{stats.total}</div>
            <div className="audit-stat-lbl">Всього записів</div>
          </div>
        </div>
        <div className="audit-stat-card">
          <div className="audit-stat-icon" style={{ background: 'rgba(16, 185, 129, 0.12)', color: '#16a34a' }}>📅</div>
          <div>
            <div className="audit-stat-val">{stats.today}</div>
            <div className="audit-stat-lbl">Дій за сьогодні</div>
          </div>
        </div>
        <div className="audit-stat-card">
          <div className="audit-stat-icon" style={{ background: 'rgba(239, 68, 68, 0.12)', color: '#dc2626' }}>⚠️</div>
          <div>
            <div className="audit-stat-val">{stats.critical}</div>
            <div className="audit-stat-lbl">Критичних дій</div>
          </div>
        </div>
        <div className="audit-stat-card">
          <div className="audit-stat-icon" style={{ background: 'rgba(245, 158, 11, 0.12)', color: '#d97706' }}>🔍</div>
          <div>
            <div className="audit-stat-val">{filteredLogs.length}</div>
            <div className="audit-stat-lbl">У вибірці</div>
          </div>
        </div>
      </div>

      {/* Панель фільтрів */}
      <div className="card audit-filters-card">
        <div className="audit-filters-grid">
          <input 
            type="text" 
            className="erp-input" 
            placeholder="🔍 Пошук по користувачу або деталях..." 
            value={searchQuery} 
            onChange={(e) => setSearchQuery(e.target.value)}
          />
          <select className="erp-input" value={actionFilter} onChange={(e) => setActionFilter(e.target.value)}>
            <option value="ALL">Всі дії</option>
            <option value="CREATE">✨ Створення</option>
            <option value="UPDATE">✏️ Оновлення</option>
            <option value="DELETE">🗑️ Видалення</option>
            <option value="EXPORT">📤 Експорт</option>
            <option value="WRITE_OFF">📉 Списання</option>
            <option value="INVENTORY_AUDIT">📋 Переоблік</option>
            <option value="ASSIGN">👤 Видача персоналу</option>
            <option value="APPROVE">✅ Погодження</option>
            <option value="REJECT">❌ Відхилення</option>
            <option value="SLA_VIOLATION">⚠️ Порушення SLA</option>
          </select>
          <select className="erp-input" value={entityFilter} onChange={(e) => setEntityFilter(e.target.value)}>
            <option value="ALL">Всі об'єкти</option>
            <option value="VEHICLE">🚗 Автомобіль</option>
            <option value="WAREHOUSE">🏬 Склад</option>
            <option value="TENANT">🏛️ Організація</option>
            <option value="RESOURCE">📦 Майно / Ресурс</option>
            <option value="SUPPLY_REQUEST">📝 Заявка на постачання</option>
            <option value="CONTRACTOR_REQUEST">🤝 Заявка підряднику</option>
            <option value="CONTRACTOR_MEMBERSHIP">🤝 Співпраця з підрядником</option>
            <option value="UNIT">🏢 Орг. одиниця</option>
            <option value="USER">👤 Користувач</option>
            <option value="SHIPMENT">🚚 Рейс / Відправка</option>
            <option value="FUEL_RECORD">⛽ Запис пального</option>
            <option value="SHIPMENT_REFUEL">⛽ Дозаправка рейсу</option>
            <option value="INVENTORY">📦 Складські залишки</option>
            <option value="CATEGORY">🏷️ Категорія</option>
            <option value="SECURITY">🛡️ Безпека</option>
          </select>
        </div>
      </div>

      <div className="card card-table">
        {filteredLogs.length === 0 ? (
          <div className="empty-state">
            <div style={{ fontSize: '48px', marginBottom: '16px' }}>📝</div>
            <h3>Записів не знайдено</h3>
            <p>Змініть параметри фільтрації або оновіть сторінку.</p>
          </div>
        ) : (
          <>
          <table className="data-table table-audit">
            <thead>
              <tr>
                <th style={{ width: '15%' }}>Дата та Час</th>
                <th style={{ width: '22%' }}>Ініціатор (Користувач)</th>
                <th style={{ width: '14%' }}>Дія</th>
                <th style={{ width: '19%' }}>Об'єкт системи</th>
                <th style={{ width: '30%' }}>Деталі операції</th>
              </tr>
            </thead>
            <tbody>
              {pagedLogs.map((log) => (
                <tr key={log.id}>
                  <td className="timestamp-cell">
                    {new Date(log.created_at).toLocaleString('uk-UA', {
                      day: '2-digit',
                      month: '2-digit',
                      year: 'numeric',
                      hour: '2-digit',
                      minute: '2-digit',
                      second: '2-digit'
                    })}
                  </td>
                  <td>
                    <div className="user-info-stack">
                      <span className="user-email-main">{log.user_email}</span>
                      <span className="user-role-sub">Роль: {ROLE_NAMES[log.user_role as keyof typeof ROLE_NAMES] || log.user_role}</span>
                    </div>
                  </td>
                  <td>{getActionBadge(log.action_type)}</td>
                  <td>
                    <div className="entity-tag">
                      <span className="entity-ua">{translateEntity(log.entity_type)}</span>
                      <span className="entity-id-short">
                        {getEntityHash(log.entity_id)}
                      </span>
                    </div>
                  </td>
                  <td className="details-text">
                    {log.details}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          <Pagination
            currentPage={safeLogsPage}
            totalPages={logsTotalPages}
            onPageChange={setLogsPage}
            totalItems={filteredLogs.length}
            itemsPerPage={LOGS_PAGE_SIZE}
          />
          </>
        )}
      </div>

      <div className="audit-footer-note">
        🔒 Записи зберігаються автоматично та не підлягають редагуванню згідно з політикою безпеки проекту Omnilog.
      </div>
    </div>
  );
};

export default AuditLogs;
