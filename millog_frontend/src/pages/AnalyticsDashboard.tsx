import React, { useEffect, useState } from 'react';
import { ResponsiveContainer, LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, BarChart, Bar, Cell } from 'recharts';
import './AnalyticsDashboard.css';

const TCO_COLORS = ['#3b82f6', '#f59e0b', '#10b981', '#6366f1', '#ec4899'];

// Словники для волонтерів
const requestLabels: Record<string, string> = { 'OPEN': 'Відкриті', 'IN_PROGRESS': 'В роботі', 'COMPLETED': 'Виконані', 'CANCELLED': 'Скасовані' };
const requestColors: Record<string, string> = { 'OPEN': '#3b82f6', 'IN_PROGRESS': '#f59e0b', 'COMPLETED': '#10b981', 'CANCELLED': '#ef4444' };

const AnalyticsDashboard: React.FC = () => {
  const [data, setData] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState('overview');

  const defaultEnd = new Date().toISOString().split('T')[0];
  const defaultStart = new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0];
  
  const [startDate, setStartDate] = useState(defaultStart);
  const [endDate, setEndDate] = useState(defaultEnd);

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        const token = localStorage.getItem('token');
        const response = await fetch(`/api/analytics/dashboard?start=${startDate}&end=${endDate}`, {
          headers: { 'Authorization': `Bearer ${token}` }
        });
        if (!response.ok) throw new Error('Помилка завантаження даних');
        setData(await response.json());
      } catch (err: any) {
        console.error(err);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, [startDate, endDate]);

  if (loading && !data) return <div className="loading-state"><div className="spinner"></div> Синхронізація даних...</div>;
  if (!data) return <div className="error-state">Немає даних</div>;

  // Рахуємо загальну кількість заявок для відсотків
  const totalVolRequests = data.volunteer_funnel?.reduce((acc: number, curr: any) => acc + curr.count, 0) || 1;

  return (
    <div className="analytics-erp-container">
      {/* ХЕДЕР */}
      <div className="erp-header">
        <div>
          <h2>Командна панель (Analytics)</h2>
          <p className="subtitle">Інтелектуальне управління логістикою та автопарком</p>
        </div>
        <div className="date-filter-group">
          <div className="date-input-wrapper"><label>Від:</label><input type="date" value={startDate} onChange={(e) => setStartDate(e.target.value)} className="erp-date-input" /></div>
          <div className="date-input-wrapper"><label>До:</label><input type="date" value={endDate} onChange={(e) => setEndDate(e.target.value)} className="erp-date-input" /></div>
        </div>
      </div>

      {/* ТОП МЕТРИКИ */}
      <div className="kpi-row">
        <div className="kpi-card"><div className="kpi-info"><span className="kpi-label">Активна техніка</span><span className="kpi-val text-blue">{data.active_vehicles}</span></div></div>
        <div className="kpi-card"><div className="kpi-info"><span className="kpi-label">Критичні залишки майна</span><span className="kpi-val text-warning">{data.critical_resources}</span></div></div>
        <div className="kpi-card danger-card"><div className="kpi-info"><span className="kpi-label">Аномалії ГСМ за період</span><span className="kpi-val text-danger">{data.fuel_anomalies}</span></div></div>
      </div>

      {/* ВКЛАДКИ */}
      <div className="erp-tabs">
        <button className={`tab-btn ${activeTab === 'overview' ? 'active' : ''}`} onClick={() => setActiveTab('overview')}>📊 Зведення</button>
        <button className={`tab-btn ${activeTab === 'logistics' ? 'active' : ''}`} onClick={() => setActiveTab('logistics')}>🛡️ Майно</button>
        <button className={`tab-btn ${activeTab === 'fleet' ? 'active' : ''}`} onClick={() => setActiveTab('fleet')}>🚙 Автопарк</button>
        <button className={`tab-btn ${activeTab === 'volunteers' ? 'active' : ''}`} onClick={() => setActiveTab('volunteers')}>🤝 Волонтери</button>
      </div>

      {/* ЗМІСТ ВКЛАДОК */}
      <div className="tab-content">
        {/* ВКЛАДКА 1: ЗВЕДЕННЯ */}
        {activeTab === 'overview' && (
          <div className="grid-layout">
            <div className="erp-widget col-span-2">
              <div className="widget-header"><h3>Забезпеченість підрозділів</h3><span className="info-badge">Готовність</span></div>
              <div className="scroll-container">
                {data.unit_readiness?.length === 0 ? <p className="empty">Немає даних.</p> : (
                  data.unit_readiness?.map((u: any, idx: number) => {
                    const score = u.readiness_score;
                    const colorClass = score >= 80 ? 'bg-green-500' : score >= 50 ? 'bg-yellow-500' : 'bg-red-500';
                    return (
                      <div key={idx} className="readiness-bar-item">
                        <div className="readiness-text"><span className="unit-title">{u.unit_name}</span><span className="unit-score">{score}%</span></div>
                        <div className="progress-track"><div className={`progress-fill ${colorClass}`} style={{ width: `${score}%` }}></div></div>
                        <div className="unit-subtext">Норма: {u.ready_resources} з {u.total_resources} позицій</div>
                      </div>
                    )
                  })
                )}
              </div>
            </div>
            <div className="erp-widget col-span-1">
              <div className="widget-header"><h3>Ефективність волонтерів (SLA)</h3></div>
              <div className="sla-card">
                <div className="sla-big-number">{data.volunteer_sla?.average_days.toFixed(1)} <span className="sla-unit">днів</span></div>
                <p className="sla-desc">Середній час закриття однієї заявки цивільним сектором</p>
                <div className="sla-footer">Виконано за період: <strong>{data.volunteer_sla?.completed_count} шт</strong></div>
              </div>
            </div>
          </div>
        )}

        {/* ВКЛАДКА 2: МАЙНО (ЛОГІСТИКА) */}
        {activeTab === 'logistics' && (
          <div className="grid-layout">
            <div className="erp-widget col-span-full">
              <div className="widget-header"><h3>Алгоритмічний прогноз вичерпання запасів</h3><span className="info-badge">За обраний період</span></div>
              <div className="scroll-container predict-list">
                {(!data.predictive_burn_rate || data.predictive_burn_rate.length === 0) ? <p className="empty">Немає витрат для прогнозу.</p> : (
                  data.predictive_burn_rate.map((item: any, idx: number) => (
                    <div key={idx} className="predict-item">
                      <div className="predict-info"><strong>{item.resource_name}</strong><span>Залишок: {item.current_stock} шт (Сер. витрата: {item.daily_burn_rate.toFixed(1)}/день)</span></div>
                      <div className="predict-status">
                        {item.days_left <= 7 ? <div className="alert-box alert-red">Вистачить на {item.days_left} дн.</div> : 
                         item.days_left <= 20 ? <div className="alert-box alert-yellow">Вистачить на {item.days_left} дн.</div> : 
                         <div className="alert-box alert-green">Запас &gt; {item.days_left} дн.</div>}
                      </div>
                    </div>
                  ))
                )}
              </div>
            </div>
          </div>
        )}

        {/* ВКЛАДКА 3: АВТОПАРК */}
        {activeTab === 'fleet' && (
          <div className="grid-layout">
            <div className="erp-widget col-span-2">
              <div className="widget-header"><h3>Аналіз споживання пального та виявлення фроду</h3><span className="info-badge">По днях</span></div>
              <div className="chart-container">
                {data.fuel_history?.length === 0 ? <p className="empty">Немає заправок.</p> : (
                  <ResponsiveContainer width="100%" height="100%">
                    <LineChart data={data.fuel_history} margin={{ top: 10, right: 10, left: -20, bottom: 0 }}>
                      <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#e2e8f0" />
                      <XAxis dataKey="month" tick={{ fill: '#64748b', fontSize: 12 }} axisLine={false} tickLine={false} />
                      <YAxis yAxisId="left" tick={{ fill: '#94a3b8', fontSize: 12 }} axisLine={false} tickLine={false} />
                      <YAxis yAxisId="right" orientation="right" hide={true} />
                      <Tooltip contentStyle={{ borderRadius: '8px', border: '1px solid #e2e8f0' }} />
                      <Legend iconType="circle" wrapperStyle={{ fontSize: '13px', paddingTop: '10px' }} />
                      <Line yAxisId="left" type="monotone" dataKey="total_liters" name="Витрата (л)" stroke="#3b82f6" strokeWidth={3} dot={false} />
                      <Line yAxisId="right" type="stepAfter" dataKey="anomalies" name="Аномалії" stroke="#ef4444" strokeWidth={2} dot={{ r: 4 }} />
                    </LineChart>
                  </ResponsiveContainer>
                )}
              </div>
            </div>

            <div className="erp-widget col-span-1">
              <div className="widget-header"><h3>Прогноз ТО</h3></div>
              <div className="scroll-container fraud-list">
                {data.maintenance_predict?.length === 0 ? <p className="empty">Немає авто.</p> : (
                  data.maintenance_predict?.map((m: any, idx: number) => (
                    <div key={idx} className="fraud-card">
                      <div className="fraud-head"><h4>{m.vehicle_name}</h4>
                        {m.km_left <= 0 ? <div className="risk-badge bg-red-500">Прострочено!</div> : 
                         m.km_left <= 1000 ? <div className="risk-badge bg-orange-500">Скоро ТО</div> : <div className="risk-badge" style={{background: '#10b981'}}>В нормі</div>}
                      </div>
                      <div className="fraud-stats">
                        <div className="stat-line"><span>Поточний пробіг:</span> <strong>{m.current_odo} км</strong></div>
                        <div className="stat-line"><span>Залишилось до ТО:</span> <strong className={m.km_left <= 0 ? 'text-danger' : ''}>{m.km_left} км</strong></div>
                      </div>
                    </div>
                  ))
                )}
              </div>
            </div>

            <div className="erp-widget col-span-2">
              <div className="widget-header"><h3>Фінансова ефективність (Ремонти)</h3><span className="info-badge">TCO</span></div>
              <div className="chart-container">
                {(!data.fleet_tco || data.fleet_tco.length === 0) ? <p className="empty">Немає витрат на ремонт.</p> : (
                  <ResponsiveContainer width="100%" height="100%">
                    <BarChart data={data.fleet_tco} layout="vertical" margin={{ top: 5, right: 30, left: 20, bottom: 5 }}>
                      <CartesianGrid strokeDasharray="3 3" horizontal={false} />
                      <XAxis type="number" tick={{fill: '#64748b'}} />
                      <YAxis dataKey="vehicle_brand" type="category" tick={{fill: '#0f172a', fontWeight: 600}} width={100} />
                      <Tooltip formatter={(value: number) => [`${value} грн`, 'Витрачено']} cursor={{fill: '#f8fafc'}} />
                      <Bar dataKey="total_cost" radius={[0, 4, 4, 0]} barSize={25}>
                        {data.fleet_tco.map((_: any, index: number) => <Cell key={`cell-${index}`} fill={TCO_COLORS[index % TCO_COLORS.length]} />)}
                      </Bar>
                    </BarChart>
                  </ResponsiveContainer>
                )}
              </div>
            </div>
            
            <div className="erp-widget col-span-1">
              <div className="widget-header"><h3>Рейтинг ризику ГСМ</h3></div>
              <div className="scroll-container fraud-list">
                {data.fleet_risk?.map((f: any, idx: number) => (
                  <div key={idx} className="fraud-card">
                    <div className="fraud-head"><h4>{f.vehicle_name}</h4><div className={`risk-badge ${f.risk_score > 50 ? 'bg-red-500' : 'bg-orange-500'}`}>Ризик: {f.risk_score}%</div></div>
                    <div className="fraud-stats">
                      <div className="stat-line"><span>Усього заправок:</span> <strong>{f.total_refuels}</strong></div>
                      <div className="stat-line"><span>Підозрілі:</span> <strong className="text-danger">{f.anomalies}</strong></div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}

        {/* ВКЛАДКА 4: ВОЛОНТЕРИ (ПОВНІСТЮ РОБОЧА) */}
        {activeTab === 'volunteers' && (
          <div className="grid-layout">
            
            <div className="erp-widget col-span-1">
              <div className="widget-header">
                <h3>Статуси заявок</h3>
                <span className="info-badge">Воронка</span>
              </div>
              <div className="scroll-container">
                {data.volunteer_funnel?.length === 0 ? <p className="empty">Заявок немає.</p> : (
                  data.volunteer_funnel?.map((item: any, idx: number) => {
                    const percentage = Math.round((item.count / totalVolRequests) * 100);
                    const color = requestColors[item.status] || '#64748b';
                    return (
                      <div key={idx} className="readiness-bar-item">
                        <div className="readiness-text">
                          <span className="unit-title">{requestLabels[item.status] || item.status}</span>
                          <span className="unit-score">{item.count} шт ({percentage}%)</span>
                        </div>
                        <div className="progress-track">
                          <div className="progress-fill" style={{ width: `${percentage}%`, backgroundColor: color }}></div>
                        </div>
                      </div>
                    )
                  })
                )}
              </div>
            </div>

            <div className="erp-widget col-span-2">
              <div className="widget-header">
                <h3>Динаміка формування потреб</h3>
                <span className="info-badge">Кількість нових запитів по днях</span>
              </div>
              <div className="chart-container">
                {data.volunteer_timeline?.length === 0 ? <p className="empty">Нових потреб не виникало.</p> : (
                  <ResponsiveContainer width="100%" height="100%">
                    <BarChart data={data.volunteer_timeline} margin={{ top: 10, right: 10, left: -20, bottom: 0 }}>
                      <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#e2e8f0" />
                      <XAxis dataKey="date" tick={{ fill: '#64748b', fontSize: 12 }} axisLine={false} tickLine={false} />
                      <YAxis tick={{ fill: '#94a3b8', fontSize: 12 }} axisLine={false} tickLine={false} />
                      <Tooltip contentStyle={{ borderRadius: '8px', border: '1px solid #e2e8f0' }} cursor={{fill: '#f8fafc'}} />
                      <Bar dataKey="count" name="Створено заявок" fill="#3b82f6" radius={[4, 4, 0, 0]} barSize={30} />
                    </BarChart>
                  </ResponsiveContainer>
                )}
              </div>
            </div>

          </div>
        )}

      </div>
    </div>
  );
};

export default AnalyticsDashboard;