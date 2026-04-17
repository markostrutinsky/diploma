import { useEffect, useState } from 'react';
import { api } from '../api/client';
import toast, { Toaster } from 'react-hot-toast';
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

const AuditLogs = () => {
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [isRefreshing, setIsRefreshing] = useState(false);

  // Стейт для фільтрів
  const [actionFilter, setActionFilter] = useState('ALL');
  const [entityFilter, setEntityFilter] = useState('ALL');
  const [searchQuery, setSearchQuery] = useState('');

  const fetchLogs = async (showToast = false) => {
    if (showToast) setIsRefreshing(true);
    try {
      if (!showToast) setLoading(true);
      const data = await api.admin.getAuditLogs();
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
  }, []);

  // Генерація кольорових бейджів дій
  const getActionBadge = (action: string) => {
    const act = action.toUpperCase();
    if (act.includes('CREATE')) return <span className="audit-badge badge-create">Створення</span>;
    if (act.includes('UPDATE')) return <span className="audit-badge badge-update">Оновлення</span>;
    if (act.includes('DELETE')) return <span className="audit-badge badge-delete">Видалення</span>;
    if (act.includes('WRITE_OFF')) return <span className="audit-badge badge-warning">Списання</span>;
    if (act.includes('ASSIGN')) return <span className="audit-badge badge-primary">Видача персоналу</span>;
    return <span className="audit-badge badge-neutral">{action}</span>;
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

  if (loading) {
    return (
      <div className="page-loading">
        <div className="spinner" />
        <p>Завантаження історії дій...</p>
      </div>
    );
  }

  return (
    <div className="audit-page">
      <Toaster position="top-right" />
      
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
            <option value="CREATE">Створення</option>
            <option value="UPDATE">Оновлення</option>
            <option value="DELETE">Видалення</option>
            <option value="WRITE_OFF">Списання</option>
          </select>
          <select className="erp-input" value={entityFilter} onChange={(e) => setEntityFilter(e.target.value)}>
            <option value="ALL">Всі об'єкти</option>
            <option value="VEHICLE">Автопарк (VEHICLE)</option>
            <option value="WAREHOUSE">Склади (WAREHOUSE)</option>
            <option value="RESOURCE">Майно (RESOURCE)</option>
            <option value="SUPPLY_REQUEST">Заявки (SUPPLY_REQUEST)</option>
            <option value="UNIT">Орг. одиниці (UNIT)</option>
            <option value="USER">Користувачі (USER)</option>
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
          <table className="data-table table-audit">
            <thead>
              <tr>
                <th style={{ width: '15%' }}>Дата та Час</th>
                <th style={{ width: '25%' }}>Ініціатор (Користувач)</th>
                <th style={{ width: '10%' }}>Дія</th>
                <th style={{ width: '20%' }}>Об'єкт системи</th>
                <th style={{ width: '30%' }}>Деталі операції</th>
              </tr>
            </thead>
            <tbody>
              {filteredLogs.map((log) => (
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
                      <span className="user-role-sub">Роль: {log.user_role}</span>
                    </div>
                  </td>
                  <td>{getActionBadge(log.action_type)}</td>
                  <td>
                    <div className="entity-tag">
                      {log.entity_type}
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
        )}
      </div>

      <div style={{ marginTop: '20px', fontSize: '12px', color: '#94a3b8', textAlign: 'center' }}>
        Записи зберігаються автоматично та не підлягають редагуванню згідно з політикою безпеки проекту Millog.
      </div>
    </div>
  );
};

export default AuditLogs;