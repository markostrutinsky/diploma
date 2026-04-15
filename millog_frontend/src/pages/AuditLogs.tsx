import React, { useEffect, useState } from 'react';
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

const AuditLogs: React.FC = () => {
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchLogs = async () => {
    try {
      setLoading(true);
      const data = await api.admin.getAuditLogs();
      setLogs(Array.isArray(data) ? data : []);
    } catch (err: any) {
      toast.error(err.message || 'Помилка доступу до журналу аудиту');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchLogs();
  }, []);

  // Функція для генерації кольорових бейджів дій
  const getActionBadge = (action: string) => {
    switch (action) {
      case 'DELETE': 
        return <span className="badge badge-critical">Видалення</span>;
      case 'WRITE_OFF': 
        return <span className="badge badge-warning">Списання</span>;
      case 'UPDATE_ROLE': 
        return <span className="badge" style={{ background: '#c4b5fd', color: '#4c1d95' }}>Права доступу</span>;
      case 'CREATE': 
        return <span className="badge badge-success">Створення</span>;
      case 'ASSIGN':
        return <span className="badge badge-primary">Видача о/с</span>;
      default: 
        return <span className="badge badge-neutral">{action}</span>;
    }
  };

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
          <button className="btn btn-secondary" onClick={fetchLogs}>
            🔄 Оновити дані
          </button>
        </div>
      </div>

      <div className="card card-table">
        {logs.length === 0 ? (
          <div className="empty-state">
            <div style={{ fontSize: '48px', marginBottom: '16px' }}>📝</div>
            <h3>Журнал порожній</h3>
            <p>Жодних системних подій ще не зафіксовано.</p>
          </div>
        ) : (
          <table className="data-table table-audit">
            <thead>
              <tr>
                <th>Дата та Час</th>
                <th>Ініціатор (Користувач)</th>
                <th>Дія</th>
                <th>Об'єкт системи</th>
                <th>Деталі операції</th>
              </tr>
            </thead>
            <tbody>
              {logs.map((log) => (
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
                        #{log.entity_id.split('-')[0].toUpperCase()}
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