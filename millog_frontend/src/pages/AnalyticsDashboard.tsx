import React, { useEffect, useState } from 'react';
import { 
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer,
  PieChart, Pie, Cell, Line, ComposedChart
} from 'recharts';
import './AnalyticsDashboard.css';

const COLORS = ['#10b981', '#f59e0b', '#ef4444', '#64748b', '#3b82f6'];

// Словники перекладу
const conditionLabels: Record<string, string> = {
  'NEW': 'Нове', 'USED': 'Вживане', 'WRITTEN_OFF': 'Списане'
};
const locationLabels: Record<string, string> = {
  'STATIONARY': 'Стаціонарні', 'MOBILE': 'Мобільні', 'UNASSIGNED': 'На руках'
};
const vehicleStatusLabels: Record<string, string> = {
  'ACTIVE': 'На ходу', 'IN_REPAIR': 'В ремонті', 'INACTIVE': 'Резерв'
};
const requestLabels: Record<string, string> = {
  'OPEN': 'Відкриті', 'IN_PROGRESS': 'В роботі', 'COMPLETED': 'Виконані', 'CANCELLED': 'Скасовані'
};

const AnalyticsDashboard: React.FC = () => {
  const [data, setData] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchData = async () => {
      try {
        const token = localStorage.getItem('token');
        const response = await fetch('/api/analytics/dashboard', {
          headers: { 'Authorization': `Bearer ${token}` }
        });
        if (!response.ok) throw new Error('Помилка завантаження даних');
        setData(await response.json());
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  if (loading) return <div className="loading-state">Завантаження аналітики...</div>;
  if (error) return <div className="error-state">{error}</div>;
  if (!data) return null;

  // Форматування даних для Recharts
  const tacticalData = data.tactical_stats?.map((item: any) => ({
    name: locationLabels[item.location_type] || item.location_type,
    Нове: item.new_items,
    Вживане: item.used_items
  })) || [];

  const burnRateData = data.burn_rate?.map((item: any) => ({
    name: conditionLabels[item.condition] || item.condition, value: item.count
  })) || [];

  const fleetHealthData = data.fleet_health?.map((item: any) => ({
    name: vehicleStatusLabels[item.condition] || item.condition, value: item.count
  })) || [];

  const funnelData = data.volunteer_funnel?.map((item: any) => ({
    name: requestLabels[item.status] || item.status, count: item.count
  })) || [];

  return (
    <div className="analytics-pro-container">
      <div className="page-header">
        <h2>Оперативна аналітика (Ситуаційний центр)</h2>
      </div>

      <div className="dashboard-grid">
        {/* ТОП МЕТРИКИ */}
        <div className="metrics-cards-row top-metrics">
          <div className="funnel-card">
            <span className="funnel-value text-success">{data.active_vehicles}</span>
            <span className="funnel-label">Авто на ходу</span>
          </div>
          <div className="funnel-card">
            <span className="funnel-value text-warning">{data.critical_resources}</span>
            <span className="funnel-label">Критичні залишки (поз.)</span>
          </div>
          <div className="funnel-card alert-bg">
            <span className="funnel-value text-danger">{data.fuel_anomalies}</span>
            <span className="funnel-label text-danger font-bold">Аномалії списання ГСМ</span>
          </div>
        </div>

        {/* БЛОК 1: ТАКТИЧНА ЛОГІСТИКА */}
        <section className="dashboard-section">
          <h3>🛡️ Матеріальне забезпечення</h3>
          <div className="charts-row">
            <div className="chart-item">
              <p className="chart-label">Ешелонування майна</p>
              <ResponsiveContainer width="100%" height={250}>
                <BarChart data={tacticalData}>
                  <XAxis dataKey="name" />
                  <YAxis />
                  <Tooltip />
                  <Legend />
                  <Bar dataKey="Нове" stackId="a" fill="#10b981" />
                  <Bar dataKey="Вживане" stackId="a" fill="#f59e0b" />
                </BarChart>
              </ResponsiveContainer>
            </div>
            <div className="chart-item">
              <p className="chart-label">Зношеність майна</p>
              <ResponsiveContainer width="100%" height={250}>
                <PieChart>
                  {/* ТУТ ВИПРАВЛЕНО ТИПИ */}
                  <Pie data={burnRateData} innerRadius={60} outerRadius={80} dataKey="value" label={({ name, percent }: { name: string; percent: number }) => `${name} ${(percent * 100).toFixed(0)}%`}>
                    {burnRateData.map((_: any, index: number) => <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />)}
                  </Pie>
                  <Tooltip />
                </PieChart>
              </ResponsiveContainer>
            </div>
          </div>
          
          <div className="critical-list">
            <h4>⚠️ Червона зона (Топ-5 критичних залишків)</h4>
            {data.critical_items?.length === 0 ? <p>Критичних залишків немає.</p> : null}
            {data.critical_items?.map((item: any) => (
              <div key={item.id} className="critical-item">
                <span>{item.name}</span>
                <span className="badge-critical">{item.quantity} / {item.min_quantity} (мін)</span>
              </div>
            ))}
          </div>
        </section>

        {/* БЛОК 2: АВТОПАРК */}
        <section className="dashboard-section">
          <h3>🚙 Моніторинг автопарку та ГСМ</h3>
          <div className="charts-row">
            <div className="chart-item">
              <p className="chart-label">Статус автомобілів</p>
              <ResponsiveContainer width="100%" height={250}>
                <PieChart>
                  {/* ТУТ ВИПРАВЛЕНО ТИПИ */}
                  <Pie data={fleetHealthData} outerRadius={80} dataKey="value" label={({ name, percent }: { name: string; percent: number }) => `${name} ${(percent * 100).toFixed(0)}%`}>
                    {fleetHealthData.map((_: any, index: number) => <Cell key={`cell-${index}`} fill={COLORS[(index + 2) % COLORS.length]} />)}
                  </Pie>
                  <Tooltip />
                </PieChart>
              </ResponsiveContainer>
            </div>
            
            <div className="chart-item">
              <p className="chart-label">Витрати пального та аномалії</p>
              <ResponsiveContainer width="100%" height={250}>
                <ComposedChart data={data.fuel_history}>
                  <CartesianGrid strokeDasharray="3 3" vertical={false} />
                  <XAxis dataKey="month" />
                  <YAxis yAxisId="left" orientation="left" stroke="#3b82f6" />
                  <YAxis yAxisId="right" orientation="right" stroke="#ef4444" />
                  <Tooltip />
                  <Legend />
                  <Bar yAxisId="left" dataKey="total_liters" fill="#3b82f6" name="Витрата (л)" opacity={0.3} />
                  <Line yAxisId="right" type="monotone" dataKey="anomalies" stroke="#ef4444" name="К-сть аномалій" strokeWidth={3} dot={{ r: 6, fill: '#ef4444' }} />
                </ComposedChart>
              </ResponsiveContainer>
            </div>
          </div>
        </section>
        
        {/* БЛОК 3: ВОЛОНТЕРИ */}
        <section className="dashboard-section">
          <h3>🤝 Статуси заявок волонтерам</h3>
          <div className="metrics-cards-row">
            {funnelData.map((v: any, index: number) => (
              <div key={v.name} className="funnel-card border-card">
                <span className="funnel-value" style={{color: COLORS[index % COLORS.length]}}>{v.count}</span>
                <span className="funnel-label">{v.name}</span>
              </div>
            ))}
          </div>
        </section>
      </div>
    </div>
  );
};

export default AnalyticsDashboard;