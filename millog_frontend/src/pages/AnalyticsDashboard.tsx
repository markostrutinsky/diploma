import React, { useEffect, useState } from 'react';
import { ResponsiveContainer, LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, BarChart, Bar, Cell, PieChart, Pie } from 'recharts';
import { toast } from 'react-hot-toast';
import jsPDF from 'jspdf';
import html2canvas from 'html2canvas';
import './AnalyticsDashboard.css';
import { api, Unit } from '../api/client';
import { usePermissions } from '../hooks/usePermissions';
import { PaywallBadge } from '../components/FeatureGate';

const TCO_COLORS = ['#3b82f6', '#f59e0b', '#10b981', '#6366f1', '#ec4899'];
const requestLabels: Record<string, string> = { 'OPEN': 'Відкриті', 'IN_PROGRESS': 'В роботі', 'COMPLETED': 'Виконані', 'CANCELLED': 'Скасовані', 'ESCALATED': 'Ескальовані (SLA)' };
const requestColors: Record<string, string> = { 'OPEN': '#3b82f6', 'IN_PROGRESS': '#f59e0b', 'COMPLETED': '#10b981', 'CANCELLED': '#ef4444', 'ESCALATED': '#9f1239' };

const AnalyticsDashboard: React.FC = () => {
  const [data, setData] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState('overview');

  // Стейт для орг. одиниць
  const [units, setUnits] = useState<Unit[]>([]);
  const [selectedUnit, setSelectedUnit] = useState<string>('');

  // Модалка Smart Поповнення
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [orderConfig, setOrderConfig] = useState<Record<string, string>>({});

  const defaultEnd = new Date().toISOString().split('T')[0];
  const defaultStart = new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0];
  const [startDate, setStartDate] = useState(defaultStart);
  const [endDate, setEndDate] = useState(defaultEnd);

  const [isSlaChecking, setIsSlaChecking] = useState(false);

  const perms = usePermissions();
  const hasSmartReplenish = perms.hasFeature('smart_replenish');
  const canTriggerSla = perms.can('sla_trigger');

  // Завантажуємо список орг. одиниць при старті
  useEffect(() => {
    api.units.list()
      .then(res => setUnits(Array.isArray(res) ? res : []))
      .catch(err => console.error("Помилка завантаження орг. одиниць", err));
  }, []);

  const fetchData = async () => {
    setLoading(true);
    try {
      const result = await api.analytics.getDashboard(startDate, endDate, selectedUnit);
      setData(result || {});
      
      if (result?.deficit_resources) {
        const initialConfig: Record<string, string> = {};
        result.deficit_resources.forEach((r: any) => {
          initialConfig[r.id] = 'WAREHOUSE';
        });
        setOrderConfig(initialConfig);
      }
    } catch (err: any) {
      toast.error('Помилка синхронізації з сервером');
      setData({});
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [startDate, endDate, selectedUnit]);

  // УНІВЕРСАЛЬНА ГЕНЕРАЦІЯ СУЧАСНОГО PDF (EXECUTIVE SUMMARY)
  const handleExportPDF = async () => {
    const input = document.getElementById('official-pdf-report');
    if (!input) return;
    
    // Жорстко фіксуємо розміри перед рендером
    input.style.display = 'block';
    input.style.position = 'absolute';
    input.style.top = '0';
    input.style.left = '-9999px'; 
    input.style.width = '210mm'; // Строго формат А4
    
    toast.loading("Формування детального звіту...", { id: 'pdf-toast' });
    try {
      // Чекаємо трохи довше, щоб усі шрифти підтягнулися
      await new Promise(resolve => setTimeout(resolve, 800));

      const canvas = await html2canvas(input, { 
        scale: 2, 
        useCORS: true,
        windowWidth: 800 // Фіксоване вікно рендеру
      });
      
      const imgData = canvas.toDataURL('image/jpeg', 1.0);
      const pdf = new jsPDF('p', 'mm', 'a4');
      const pdfWidth = pdf.internal.pageSize.getWidth();
      const pageHeight = pdf.internal.pageSize.getHeight();
      const imgHeight = (canvas.height * pdfWidth) / canvas.width;
      
      let heightLeft = imgHeight;
      let position = 0;

      pdf.addImage(imgData, 'JPEG', 0, position, pdfWidth, imgHeight);
      heightLeft -= pageHeight;

      while (heightLeft > 0) {
        position = position - pageHeight;
        pdf.addPage();
        pdf.addImage(imgData, 'JPEG', 0, position, pdfWidth, imgHeight);
        heightLeft -= pageHeight;
      }
      
      const unitName = selectedUnit ? units.find(u => u.id.toString() === selectedUnit)?.name : 'All';
      pdf.save(`OmniLog_Analytics_${unitName}_${endDate}.pdf`);
      
      toast.success("Звіт успішно збережено", { id: 'pdf-toast' });
    } catch (e) {
      toast.error("Помилка генерації документа", { id: 'pdf-toast' });
    } finally {
      input.style.display = 'none';
    }
  };

  const handleExportInventory = async () => {
    const toastId = toast.loading('Формування Excel звіту (Залишки)...');
    try {
      const unitIdToExport = selectedUnit ? parseInt(selectedUnit, 10) : undefined;
      const { blob, filename } = await api.analytics.exportInventory(unitIdToExport); 
      
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', filename);
      document.body.appendChild(link);
      link.click();
      link.parentNode?.removeChild(link);
      window.URL.revokeObjectURL(url);
      
      toast.success('Excel завантажено!', { id: toastId });
    } catch (error: any) {
      toast.error(error.message || 'Не вдалося завантажити звіт', { id: toastId });
    }
  };

  const handleExportFuel = async () => {
    const toastId = toast.loading('Формування Excel звіту (Пальне)...');
    try {
      const { blob, filename } = await api.analytics.exportFuel(startDate, endDate); 
      
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', filename);
      document.body.appendChild(link);
      link.click();
      link.parentNode?.removeChild(link);
      window.URL.revokeObjectURL(url);
      
      toast.success('Excel завантажено!', { id: toastId });
    } catch (error: any) {
      toast.error(error.message || 'Не вдалося завантажити звіт', { id: toastId });
    }
  };
  const handleTriggerSLA = async () => {
  setIsSlaChecking(true);
  const toastId = toast.loading("Перевірка термінів виконання...");
  try {
    const res = await api.admin.triggerSLACheck();
    const newlyEscalated = res.escalated_count ?? 0;
    const existing = res.existing_escalated ?? 0;
    const pendingSoon = res.pending_near_sla ?? 0;

    // Формуємо зрозуміле повідомлення з контекстом
    if (newlyEscalated > 0) {
      toast.success(
        `Ескальовано ${newlyEscalated} нових заявок. Всього у статусі ESCALATED: ${existing}.`,
        { id: toastId, duration: 6000 }
      );
    } else if (existing > 0) {
      toast(
        `Нових порушень SLA не знайдено.\nАле у системі вже ${existing} ескальованих заявок — опрацюйте їх.` +
          (pendingSoon > 0 ? `\nЩе ${pendingSoon} заявок наближаються до дедлайну.` : ''),
        { id: toastId, icon: 'ℹ️', duration: 7000 }
      );
    } else if (pendingSoon > 0) {
      toast.success(
        `Порушень SLA немає. ${pendingSoon} заявок наближаються до дедлайну.`,
        { id: toastId, duration: 5000 }
      );
    } else {
      toast.success("Всі заявки в межах норми SLA", { id: toastId });
    }
    fetchData(); // Оновлюємо цифри на дашборді
  } catch (err) {
    toast.error("Не вдалося запустити монітор", { id: toastId });
  } finally {
    setIsSlaChecking(false);
  }
};

  const submitSmartReplenish = async () => {
    const payloadItems = (data?.deficit_resources || [])
      .filter((r: any) => orderConfig[r.id] && orderConfig[r.id] !== 'NONE')
      .map((r: any) => ({
        resource_id: r.id,
        name: r.name,
        quantity: r.needed,
        target: orderConfig[r.id]
      }));

    if (payloadItems.length === 0) {
      toast.error("Не обрано жодного ресурсу для замовлення");
      return;
    }

    toast.loading("Формування заявок...", { id: 'replenish-toast' });
    try {
      const resData: any = await api.analytics.smartReplenish(payloadItems);
      toast.success(`Сформовано ${resData.count} запитів!`, { id: 'replenish-toast' });
      setIsModalOpen(false);
      fetchData(); 
    } catch (err) {
      toast.error("Помилка створення заявок", { id: 'replenish-toast' });
    }
  };

  if (loading && !data) return <div className="loading-state"><div className="spinner"></div> Завантаження аналітики...</div>;
  if (!data) return <div className="error-state">Немає даних</div>;

  // БЕЗПЕЧНИЙ РОЗРАХУНОК ВОРОНКИ
  const contractorFunnel = Array.isArray(data?.CONTRACTOR_funnel) ? data.CONTRACTOR_funnel : [];
  const totalcontractorReqs = contractorFunnel.reduce((acc: number, curr: any) => acc + (curr.count || 0), 0) || 1;

  return (
    <div className="analytics-erp-container">


      {/* --- МОДАЛЬНЕ ВІКНО SMART ПОПОВНЕННЯ --- */}
      {isModalOpen && (
        <div className="modal-overlay">
          <div className="modal-content">
            <div className="modal-header">
              <h3>⚡ Smart Поповнення (Критичний дефіцит)</h3>
              <button className="close-btn" onClick={() => setIsModalOpen(false)}>×</button>
            </div>
            <div className="modal-body">
              <p className="modal-desc">Оберіть джерело постачання для майна, запаси якого впали нижче критичної норми.</p>
              {(!data?.deficit_resources || data.deficit_resources.length === 0) ? (
                <div className="empty">Наразі дефіциту майна не виявлено. Всі запаси в нормі.</div>
              ) : (
                <table className="deficit-table">
                  <thead>
                    <tr>
                      <th>Майно</th>
                      <th>Залишок</th>
                      <th>Потрібно дозамовити</th>
                      <th>Джерело запиту</th>
                    </tr>
                  </thead>
                  <tbody>
                    {data.deficit_resources.map((res: any) => (
                      <tr key={res.id}>
                        <td style={{fontWeight: 600}}>{res.name}</td>
                        <td className="text-danger">{res.current} з {res.min} (мін)</td>
                        <td style={{fontWeight: 600}}>+{res.needed} шт.</td>
                        <td>
                          <select 
                            className="source-select"
                            value={orderConfig[res.id]} 
                            onChange={(e) => setOrderConfig(prev => ({ ...prev, [res.id]: e.target.value }))}
                          >
                            <option value="WAREHOUSE">🏢 Зі складу (Внутрішній)</option>
                            <option value="CONTRACTOR">🤝 Зовнішній запит</option>
                            <option value="NONE">❌ Не замовляти зараз</option>
                          </select>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}
            </div>
            <div className="modal-footer">
              <button className="btn-secondary" onClick={() => setIsModalOpen(false)}>Скасувати</button>
              <button className="btn-primary" onClick={submitSmartReplenish} disabled={!data?.deficit_resources || data.deficit_resources.length === 0}>
                Підтвердити формування
              </button>
            </div>
          </div>
        </div>
      )}

      {/* ХЕДЕР ДАШБОРДУ З ФІЛЬТРАМИ ТА КНОПКАМИ */}
      <div className="erp-header">
        <div>
          <h2>Аналітична панель (Дашборд)</h2>
          <p className="subtitle">Інтелектуальне управління ресурсами та активами</p>
        </div>
        <div className="action-bar">
          
          <div className="unit-filter">
            <select 
              value={selectedUnit} 
              onChange={(e) => setSelectedUnit(e.target.value)}
              className="erp-date-input"
              style={{ minWidth: '200px', cursor: 'pointer' }}
            >
              <option value="">Вся організація (Зведення)</option>
              {units.map(u => (
                <option key={u.id} value={u.id}>{u.name}</option>
              ))}
            </select>
          </div>

          <div className="date-filter-group">
            <input type="date" value={startDate} onChange={(e) => setStartDate(e.target.value)} className="erp-date-input" />
            <span style={{color: '#94a3b8', margin: '0 8px'}}>-</span>
            <input type="date" value={endDate} onChange={(e) => setEndDate(e.target.value)} className="erp-date-input" />
          </div>
          <div className="button-group">
            <button onClick={handleExportInventory} className="btn-secondary" title="Залишки на складах">📊 Excel (Склади)</button>
            <button onClick={handleExportFuel} className="btn-secondary" title="Історія пального">⛽ Excel (Пальне)</button>
            <button onClick={handleExportPDF} className="btn-secondary">📄 Звіт (А4)</button>
            {canTriggerSla && (
              <button 
                onClick={handleTriggerSLA} 
                className="btn-danger-action" 
                disabled={isSlaChecking}
              >
                {isSlaChecking ? '⌛ Перевірка...' : '⚡ Перевірити SLA'}
              </button>
            )}
            {hasSmartReplenish ? (
              <button onClick={() => setIsModalOpen(true)} className="btn-primary">⚡ Smart Поповнення</button>
            ) : (
              <button
                className="btn-primary"
                onClick={() =>
                  toast('Smart Поповнення доступне на тарифі PRO. Зверніться до білінгу.', {
                    icon: '🔒',
                    duration: 5000,
                  })
                }
                style={{ opacity: 0.6, cursor: 'not-allowed' }}
                title="Доступно на тарифі PRO"
              >
                🔒 Smart Поповнення <PaywallBadge feature="smart_replenish" compact />
              </button>
            )}
          </div>
        </div>
      </div>

      {/* ТОП МЕТРИКИ */}
      <div className="kpi-row">
        <div className="kpi-card"><div className="kpi-info"><span className="kpi-label">Експлуатована техніка</span><span className="kpi-val text-blue">{data?.active_vehicles || 0}</span></div></div>
        <div className="kpi-card"><div className="kpi-info"><span className="kpi-label">Дефіцитні позиції ресурсів</span><span className="kpi-val text-warning">{data?.critical_resources || 0}</span></div></div>
        <div className="kpi-card danger-card"><div className="kpi-info"><span className="kpi-label">Інциденти перевитрат (Аномалії)</span><span className="kpi-val text-danger">{data?.fuel_anomalies || 0}</span></div></div>
      </div>

      {/* ВКЛАДКИ */}
      <div className="erp-tabs">
        <button className={`tab-btn ${activeTab === 'overview' ? 'active' : ''}`} onClick={() => setActiveTab('overview')}>📊 Зведення</button>
        <button className={`tab-btn ${activeTab === 'logistics' ? 'active' : ''}`} onClick={() => setActiveTab('logistics')}>🛡️ Ресурси</button>
        <button className={`tab-btn ${activeTab === 'fleet' ? 'active' : ''}`} onClick={() => setActiveTab('fleet')}>🚙 Активи (Автопарк)</button>
        <button className={`tab-btn ${activeTab === 'CONTRACTORs' ? 'active' : ''}`} onClick={() => setActiveTab('CONTRACTORs')}>🤝 Зовнішні запити</button>
      </div>

      {/* ЗМІСТ ВКЛАДОК */}
      <div className="tab-content">
        {/* ВКЛАДКА 1: ЗВЕДЕННЯ */}
        {activeTab === 'overview' && (
          <div className="grid-layout">
            {canTriggerSla && (
              <div className="erp-widget col-span-full sla-control-panel">
                <div className="widget-header">
                  <h3>🛠️ Панель оперативного контролю</h3>
                </div>
                
                <div className="sla-control-body">
                  <div className="sla-control-text">
                    <p className="sla-control-title">Моніторинг «завислих» заявок</p>
                    <p className="sla-control-desc">
                      Система автоматично перевіряє заявки кожну годину. Натисніть кнопку, щоб запустити перевірку та ескалацію негайно.
                    </p>
                  </div>
                  
                  <button 
                    onClick={handleTriggerSLA} 
                    className="btn-danger-action" 
                    disabled={isSlaChecking}
                  >
                    Запустити монітор
                  </button>
                </div>
              </div>
            )}

            <div className="erp-widget col-span-2">
              <div className="widget-header"><h3>Рівень забезпечення орг. одиниць</h3><span className="info-badge">Готовність</span></div>
              <div className="scroll-container">
                {(!data?.unit_readiness || data.unit_readiness.length === 0) ? <p className="empty">Немає даних.</p> : (
                  data.unit_readiness.map((u: any, idx: number) => {
                    const score = u.readiness_score;
                    const colorClass = score >= 80 ? 'bg-green-500' : score >= 50 ? 'bg-yellow-500' : 'bg-red-500';
                    return (
                      <div key={idx} className="readiness-bar-item">
                        <div className="readiness-text"><span className="unit-title">{u.unit_name}</span><span className="unit-score">{score}%</span></div>
                        <div className="progress-track"><div className={`progress-fill ${colorClass}`} style={{ width: `${score}%` }}></div></div>
                        <div className="unit-subtext">Норматив: {u.ready_resources} з {u.total_resources} позицій</div>
                      </div>
                    )
                  })
                )}
              </div>
            </div>
            
            <div className="erp-widget col-span-1">
              <div className="widget-header"><h3>Критичний дефіцит</h3></div>
              <div className="scroll-container">
                {(!data?.deficit_resources || data.deficit_resources.length === 0) ? (
                   <p className="empty">Всі ресурси в нормі</p>
                ) : (
                  data.deficit_resources.map((res: any, idx: number) => (
                    <div key={idx} className="fraud-card" style={{borderLeft: '4px solid #ef4444'}}>
                      <div className="fraud-head">
                        <h4>{res.name}</h4>
                        <div className="risk-badge bg-red-500">Залишок: {res.current}</div>
                      </div>
                      <div style={{fontSize: '0.85rem', color: '#64748b', marginTop: '4px'}}>
                        Мінімум: {res.min} | Треба дозамовити: <strong style={{color: '#ef4444'}}>+{res.needed}</strong>
                      </div>
                    </div>
                  ))
                )}
              </div>
            </div>

            {/* ОНОВЛЕНИЙ ВІДЖЕТ SLA (З додатковими метриками) */}
            <div className="erp-widget col-span-full">
              <div className="widget-header">
                <h3>Ефективність зовнішніх запитів (SLA)</h3>
                <span className="info-badge">Швидкість постачання</span>
              </div>
              
              <div className="sla-banner">
                <div className="sla-banner-left">
                  <span className="sla-banner-value">
                    {data?.CONTRACTOR_sla?.average_days?.toFixed(1) || "0.0"}
                  </span>
                  <span className="sla-banner-unit">днів</span>
                </div>
                
                <div className="sla-banner-divider"></div>
                
                <div className="sla-banner-right">
                  <p className="sla-banner-desc">
                    Середній час закриття одного запиту зовнішніми постачальниками (від моменту створення до статусу "Виконані").
                  </p>
                  <div className="sla-banner-stats">
                    <span>Успішно виконано за обраний період:</span>
                    <strong>{data?.CONTRACTOR_sla?.completed_count || 0} шт</strong>
                  </div>
                </div>
              </div>

              {/* НОВИЙ БЛОК: Динамічні метрики SLA */}
              <div className="sla-extra-metrics">
                <div className="sla-metric-box">
                  <span className="sla-metric-title">On-Time Delivery (OTD)</span>
                  <span className={`sla-metric-val ${
                    (data?.CONTRACTOR_sla?.otd_percentage || 0) >= 90 ? 'text-success' : 
                    (data?.CONTRACTOR_sla?.otd_percentage || 0) >= 75 ? 'text-warning' : 'text-danger'
                  }`}>
                    {data?.CONTRACTOR_sla?.otd_percentage || 0}%
                  </span>
                  <span className="sla-metric-sub">Рівень вчасного виконання</span>
                </div>
                
                <div className="sla-metric-box">
                  <span className="sla-metric-title">Прострочені завдання</span>
                  <span className={`sla-metric-val ${
                    (data?.CONTRACTOR_sla?.overdue_count || 0) > 0 ? 'text-danger' : 'text-success'
                  }`}>
                    {data?.CONTRACTOR_sla?.overdue_count || 0}
                  </span>
                  <span className="sla-metric-sub">Критичні порушення дедлайнів</span>
                </div>
                
                <div className="sla-metric-box">
                  <span className="sla-metric-title">Найшвидше виконання</span>
                  <span className="sla-metric-val text-blue">
                    {data?.CONTRACTOR_sla?.fastest_days?.toFixed(1) || "0.0"} дн
                  </span>
                  <span className="sla-metric-sub">Рекорд за поточний період</span>
                </div>
              </div>
            </div>
            
            <div className="erp-widget col-span-full" style={{ height: 'auto', minHeight: '180px' }}>
              <div className="widget-header">
                <h3>🔄 Життєвий цикл та рух активів</h3>
                <span className="info-badge">За весь час</span>
              </div>
              
              <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '20px', marginTop: '10px' }}>
                <div style={{ background: 'var(--bg-card)', padding: '20px', borderRadius: '8px', border: '1px solid var(--border)', display: 'flex', flexDirection: 'column' }}>
                  <span style={{ fontSize: '0.85rem', color: 'var(--text-muted)', fontWeight: 600, textTransform: 'uppercase' }}>Списано (Використано)</span>
                  <div style={{ fontSize: '2rem', fontWeight: 700, color: 'var(--text)', marginTop: '10px' }}>
                    {data?.written_off_resources || 0} <span style={{fontSize: '1rem', color: 'var(--text-muted)', fontWeight: 500}}>позицій</span>
                  </div>
                </div>

                <div style={{ background: 'var(--bg-card)', padding: '20px', borderRadius: '8px', border: '1px solid var(--border)', display: 'flex', flexDirection: 'column' }}>
                  <span style={{ fontSize: '0.85rem', color: 'var(--text-muted)', fontWeight: 600, textTransform: 'uppercase' }}>Успішні переміщення</span>
                  <div style={{ fontSize: '2rem', fontWeight: 700, color: 'var(--success)', marginTop: '10px' }}>
                    {data?.completed_requests || 0} <span style={{fontSize: '1rem', color: 'var(--text-muted)', fontWeight: 500}}>заявок</span>
                  </div>
                </div>

                <div style={{ background: 'rgba(251, 191, 36, 0.1)', padding: '20px', borderRadius: '8px', border: '1px solid rgba(251, 191, 36, 0.3)', display: 'flex', flexDirection: 'column' }}>
                  <span style={{ fontSize: '0.85rem', color: 'var(--warning)', fontWeight: 600, textTransform: 'uppercase' }}>Автопарк: В ремонті</span>
                  <div style={{ fontSize: '2rem', fontWeight: 700, color: 'var(--warning)', marginTop: '10px' }}>
                    {data?.in_repair_vehicles || 0} <span style={{fontSize: '1rem', color: 'var(--warning)', opacity: 0.7, fontWeight: 500}}>ТЗ</span>
                  </div>
                </div>

                <div style={{ background: 'rgba(239, 68, 68, 0.1)', padding: '20px', borderRadius: '8px', border: '1px solid rgba(239, 68, 68, 0.3)', display: 'flex', flexDirection: 'column' }}>
                  <span style={{ fontSize: '0.85rem', color: 'var(--error)', fontWeight: 600, textTransform: 'uppercase' }}>Виведено з експлуатації</span>
                  <div style={{ fontSize: '2rem', fontWeight: 700, color: 'var(--error)', marginTop: '10px' }}>
                    {data?.inactive_vehicles || 0} <span style={{fontSize: '1rem', color: 'var(--error)', opacity: 0.7, fontWeight: 500}}>ТЗ</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* ВКЛАДКА 2: МАЙНО (ЛОГІСТИКА) */}
        {activeTab === 'logistics' && (
          <div className="grid-layout">
            
            <div className="erp-widget col-span-2">
              <div className="widget-header"><h3>Завантаженість складів (Топ-5)</h3></div>
              <div className="chart-container">
                {(!data?.warehouse_load || data.warehouse_load.length === 0) ? <p className="empty">Немає майна на складах.</p> : (
                  <ResponsiveContainer width="100%" height="100%">
                    <BarChart data={data.warehouse_load} margin={{ top: 10, right: 30, left: 0, bottom: 0 }}>
                      <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#e2e8f0" />
                      <XAxis dataKey="warehouse_name" tick={{ fill: '#64748b', fontSize: 12 }} axisLine={false} tickLine={false} />
                      <YAxis tick={{ fill: '#94a3b8', fontSize: 12 }} axisLine={false} tickLine={false} />
                      <Tooltip contentStyle={{ borderRadius: '8px', border: '1px solid #e2e8f0' }} cursor={{ fill: '#f8fafc' }} />
                      <Bar dataKey="total_items" name="Одиниць майна" fill="#6366f1" radius={[4, 4, 0, 0]} barSize={40} />
                    </BarChart>
                  </ResponsiveContainer>
                )}
              </div>
            </div>

            <div className="erp-widget col-span-1">
              <div className="widget-header"><h3>Топ-5 затребуваних ресурсів</h3></div>
              <div className="chart-container">
                {(!data?.top_resources || data.top_resources.length === 0) ? <p className="empty">Замовлення відсутні.</p> : (
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie 
                        data={data.top_resources} 
                        dataKey="total_ordered" 
                        nameKey="resource_name" 
                        cx="50%" 
                        cy="50%" 
                        innerRadius={60} 
                        outerRadius={90} 
                        paddingAngle={5}
                      >
                        {data.top_resources.map((_: any, index: number) => (
                          <Cell key={`cell-${index}`} fill={TCO_COLORS[index % TCO_COLORS.length]} />
                        ))}
                      </Pie>
                      <Tooltip contentStyle={{ borderRadius: '8px', border: '1px solid #e2e8f0' }} />
                      <Legend verticalAlign="bottom" height={36} iconType="circle" wrapperStyle={{ fontSize: '12px' }}/>
                    </PieChart>
                  </ResponsiveContainer>
                )}
              </div>
            </div>

            <div className="erp-widget col-span-full" style={{ height: '350px' }}>
              <div className="widget-header"><h3>Прогноз вичерпання ресурсів</h3><span className="info-badge">За обраний період</span></div>
              <div className="scroll-container predict-list">
                {(!data?.predictive_burn_rate || data.predictive_burn_rate.length === 0) ? <p className="empty">Немає витрат для прогнозу.</p> : (
                  data.predictive_burn_rate.map((item: any, idx: number) => (
                    <div key={idx} className="predict-item">
                      <div className="predict-info"><strong>{item.resource_name}</strong><span>Залишок: {item.current_stock} шт (Сер. витрата: {(item.daily_burn_rate ?? 0).toFixed(1)}/день)</span></div>
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

            {/* ФІНАНСОВІ МЕТРИКИ */}
            <div className="erp-widget col-span-1">
              <div className="widget-header"><h3>💰 Загальна вартість майна</h3></div>
              <div style={{ padding: '24px', textAlign: 'center' }}>
                <div style={{ fontSize: '2rem', fontWeight: 700, color: 'var(--primary)', marginBottom: '8px' }}>
                  {(data?.inventory_total_value || 0).toLocaleString('uk-UA', { minimumFractionDigits: 2, maximumFractionDigits: 2 })} ₴
                </div>
                <div style={{ color: 'var(--text-secondary)', fontSize: '0.9rem' }}>Активних залишків на складах</div>
                {(data?.write_off_total_value || 0) > 0 && (
                  <div style={{ marginTop: '16px', padding: '12px', background: 'var(--surface)', borderRadius: '8px' }}>
                    <div style={{ fontSize: '0.85rem', color: 'var(--text-secondary)' }}>Списано за період</div>
                    <div style={{ fontSize: '1.3rem', fontWeight: 600, color: 'var(--error)' }}>
                      {(data.write_off_total_value || 0).toLocaleString('uk-UA', { minimumFractionDigits: 2, maximumFractionDigits: 2 })} ₴
                    </div>
                  </div>
                )}
              </div>
            </div>

            <div className="erp-widget col-span-2">
              <div className="widget-header"><h3>Вартість майна по складах (Топ-5)</h3></div>
              <div className="chart-container">
                {(!data?.warehouse_value_stats || data.warehouse_value_stats.length === 0) ? <p className="empty">Немає оцінених залишків.</p> : (
                  <ResponsiveContainer width="100%" height="100%">
                    <BarChart data={data.warehouse_value_stats} margin={{ top: 10, right: 30, left: 20, bottom: 0 }}>
                      <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#e2e8f0" />
                      <XAxis dataKey="warehouse_name" tick={{ fill: '#64748b', fontSize: 12 }} axisLine={false} tickLine={false} />
                      <YAxis tick={{ fill: '#94a3b8', fontSize: 12 }} axisLine={false} tickLine={false} tickFormatter={(v) => `${(v/1000).toFixed(0)}к`} />
                      <Tooltip contentStyle={{ borderRadius: '8px', border: '1px solid #e2e8f0' }} formatter={(v: number) => [`${v.toLocaleString('uk-UA')} ₴`, 'Вартість']} cursor={{ fill: '#f8fafc' }} />
                      <Bar dataKey="total_value" name="Вартість (₴)" fill="#10b981" radius={[4, 4, 0, 0]} barSize={40} />
                    </BarChart>
                  </ResponsiveContainer>
                )}
              </div>
            </div>

            <div className="erp-widget col-span-full">
              <div className="widget-header"><h3>Топ-5 найдорожчих позицій</h3></div>
              {(!data?.top_costly_resources || data.top_costly_resources.length === 0) ? <p className="empty" style={{ padding: '16px' }}>Немає оцінених ресурсів.</p> : (
                <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '14px' }}>
                  <thead>
                    <tr style={{ borderBottom: '2px solid var(--border)' }}>
                      <th style={{ padding: '10px 16px', textAlign: 'left', color: 'var(--text-secondary)', fontWeight: 600 }}>Найменування</th>
                      <th style={{ padding: '10px 16px', textAlign: 'right', color: 'var(--text-secondary)', fontWeight: 600 }}>К-сть</th>
                      <th style={{ padding: '10px 16px', textAlign: 'right', color: 'var(--text-secondary)', fontWeight: 600 }}>Ціна за од., ₴</th>
                      <th style={{ padding: '10px 16px', textAlign: 'right', color: 'var(--text-secondary)', fontWeight: 600 }}>Загальна вартість, ₴</th>
                    </tr>
                  </thead>
                  <tbody>
                    {data.top_costly_resources.map((r: any, idx: number) => (
                      <tr key={idx} style={{ borderBottom: '1px solid var(--border)' }}>
                        <td style={{ padding: '10px 16px', fontWeight: 500 }}>{r.resource_name}</td>
                        <td style={{ padding: '10px 16px', textAlign: 'right' }}>{r.quantity}</td>
                        <td style={{ padding: '10px 16px', textAlign: 'right' }}>{r.unit_price.toLocaleString('uk-UA', { minimumFractionDigits: 2 })}</td>
                        <td style={{ padding: '10px 16px', textAlign: 'right', fontWeight: 700, color: 'var(--primary)' }}>{r.total_value.toLocaleString('uk-UA', { minimumFractionDigits: 2 })}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}
            </div>

          </div>
        )}
        {activeTab === 'fleet' && (
          <div className="grid-layout">
            <div className="erp-widget col-span-2">
              <div className="widget-header"><h3>Аналіз споживання палива та виявлення аномалій</h3></div>
              <div className="chart-container">
                {(!data?.fuel_history || data.fuel_history.length === 0) ? <p className="empty">Немає транзакцій.</p> : (
                  <ResponsiveContainer width="100%" height="100%">
                    <LineChart data={data.fuel_history} margin={{ top: 10, right: 10, left: -20, bottom: 0 }}>
                      <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#e2e8f0" />
                      <XAxis dataKey="month" tick={{ fill: '#64748b', fontSize: 12 }} axisLine={false} tickLine={false} />
                      <YAxis yAxisId="left" tick={{ fill: '#94a3b8', fontSize: 12 }} axisLine={false} tickLine={false} />
                      <YAxis yAxisId="right" orientation="right" hide={true} />
                      <Tooltip contentStyle={{ borderRadius: '8px', border: '1px solid #e2e8f0' }} />
                      <Legend iconType="circle" wrapperStyle={{ fontSize: '13px', paddingTop: '10px' }} />
                      <Line yAxisId="left" type="monotone" dataKey="total_liters" name="Витрата (л)" stroke="#3b82f6" strokeWidth={3} dot={false} />
                      <Line yAxisId="right" type="stepAfter" dataKey="anomalies" name="Інциденти" stroke="#ef4444" strokeWidth={2} dot={{ r: 4 }} />
                    </LineChart>
                  </ResponsiveContainer>
                )}
              </div>
            </div>
            <div className="erp-widget col-span-1">
              <div className="widget-header"><h3>Прогноз Сервісного Обслуговування</h3></div>
              <div className="scroll-container fraud-list">
                {(!data?.maintenance_predict || data.maintenance_predict.length === 0) ? <p className="empty">Немає активів.</p> : (
                  data.maintenance_predict.map((m: any, idx: number) => (
                    <div key={idx} className="fraud-card">
                      <div className="fraud-head"><h4>{m.vehicle_name}</h4>
                        {m.km_left <= 0 ? <div className="risk-badge bg-red-500">Прострочено!</div> : 
                         m.km_left <= 1000 ? <div className="risk-badge bg-orange-500">Скоро ТО</div> : <div className="risk-badge" style={{background: '#10b981'}}>В нормі</div>}
                      </div>
                      <div className="fraud-stats">
                        <div className="stat-line"><span>Поточний пробіг:</span> <strong>{m.current_odo} км</strong></div>
                        <div className="stat-line"><span>До сервісу:</span> <strong className={m.km_left <= 0 ? 'text-danger' : ''}>{m.km_left} км</strong></div>
                      </div>
                    </div>
                  ))
                )}
              </div>
            </div>
            <div className="erp-widget col-span-2">
              <div className="widget-header"><h3>Фінансова ефективність (TCO)</h3></div>
              <div className="chart-container">
                {(!data?.fleet_tco || data.fleet_tco.length === 0) ? <p className="empty">Немає витрат на ремонт.</p> : (
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
              <div className="widget-header"><h3>Рейтинг ризику перевитрат</h3></div>
              <div className="scroll-container fraud-list">
                {(!data?.fleet_risk || data.fleet_risk.length === 0) ? <p className="empty">Аномалій не зафіксовано</p> : (
                  data.fleet_risk.map((f: any, idx: number) => (
                    <div key={idx} className="fraud-card">
                      <div className="fraud-head"><h4>{f.vehicle_name}</h4><div className={`risk-badge ${f.risk_score > 50 ? 'bg-red-500' : 'bg-orange-500'}`}>Ризик: {f.risk_score}%</div></div>
                      <div className="fraud-stats">
                        <div className="stat-line"><span>Транзакцій:</span> <strong>{f.total_refuels}</strong></div>
                        <div className="stat-line"><span>Підозрілі:</span> <strong className="text-danger">{f.anomalies}</strong></div>
                      </div>
                    </div>
                  ))
                )}
              </div>
            </div>
          </div>
        )}

        {/* ВКЛАДКА 4: ПІДРЯДНИКИ (ЗОВНІШНІ ЗАПИТИ) */}
        {activeTab === 'CONTRACTORs' && (
          <div className="grid-layout">
            <div className="erp-widget col-span-1">
              <div className="widget-header"><h3>Статуси запитів</h3></div>
              <div className="scroll-container">
                {contractorFunnel.length === 0 ? <p className="empty">Запитів немає.</p> : (
                  contractorFunnel.map((item: any, idx: number) => {
                    const percentage = Math.round((item.count / totalcontractorReqs) * 100);
                    const color = requestColors[item.status] || '#64748b';
                    return (
                      <div key={idx} className="readiness-bar-item">
                        <div className="readiness-text"><span className="unit-title">{requestLabels[item.status] || item.status}</span><span className="unit-score">{item.count} шт ({percentage}%)</span></div>
                        <div className="progress-track"><div className="progress-fill" style={{ width: `${percentage}%`, backgroundColor: color }}></div></div>
                      </div>
                    )
                  })
                )}
              </div>
            </div>
            <div className="erp-widget col-span-2">
              <div className="widget-header"><h3>Динаміка формування потреб</h3></div>
              <div className="chart-container">
                {(!data?.CONTRACTOR_timeline || data.CONTRACTOR_timeline.length === 0) ? <p className="empty">Нових потреб не виникало.</p> : (
                  <ResponsiveContainer width="100%" height="100%">
                    <BarChart data={data.CONTRACTOR_timeline} margin={{ top: 10, right: 10, left: -20, bottom: 0 }}>
                      <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#e2e8f0" />
                      <XAxis dataKey="date" tick={{ fill: '#64748b', fontSize: 12 }} axisLine={false} tickLine={false} />
                      <YAxis tick={{ fill: '#94a3b8', fontSize: 12 }} axisLine={false} tickLine={false} />
                      <Tooltip contentStyle={{ borderRadius: '8px', border: '1px solid #e2e8f0' }} cursor={{fill: '#f8fafc'}} />
                      <Bar dataKey="count" name="Створено запитів" fill="#3b82f6" radius={[4, 4, 0, 0]} barSize={30} />
                    </BarChart>
                  </ResponsiveContainer>
                )}
              </div>
            </div>
          </div>
        )}
      </div>

      {/* =========================================================================
          ОНОВЛЕНИЙ PDF ШАБЛОН (СУВОРІ ДЕТАЛЬНІ ТАБЛИЦІ + ЖИТТЄВИЙ ЦИКЛ)
      ========================================================================= */}
      <div id="official-pdf-report" style={{ display: 'none', background: '#fff', color: '#000', width: '210mm', padding: '20mm', boxSizing: 'border-box', fontFamily: 'sans-serif' }}>
        
        {/* ШАПКА ЗВІТУ */}
        <div style={{ borderBottom: '2px solid #000', paddingBottom: '10px', marginBottom: '20px' }}>
          <table style={{ width: '100%', border: 'none' }}>
            <tbody>
              <tr>
                <td style={{ verticalAlign: 'top', width: '60%' }}>
                  <strong style={{ fontSize: '14pt' }}>OMNILOG - ОПЕРАЦІЙНА ЛОГІСТИКА</strong><br/>
                  <span style={{ fontSize: '10pt', color: '#555' }}>Деталізований аналітичний звіт стану активів</span>
                </td>
                <td style={{ verticalAlign: 'top', textAlign: 'right', fontSize: '10pt' }}>
                  <strong>Документ №:</strong> {Math.floor(Math.random() * 10000)}<br/>
                  <strong>Дата:</strong> {new Date().toLocaleDateString('uk-UA')}<br/>
                  <strong>Період:</strong> {startDate} — {endDate}
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <h1 style={{ textAlign: 'center', fontSize: '16pt', marginBottom: '20px' }}>
          ЗВЕДЕННЯ: {selectedUnit ? units.find(u => u.id.toString() === selectedUnit)?.name?.toUpperCase() : "ВСЯ ОРГАНІЗАЦІЯ"}
        </h1>

        {/* СЕКЦІЯ 1: ГОЛОВНІ ПОКАЗНИКИ */}
        <h2 style={{ fontSize: '12pt', backgroundColor: '#f0f0f0', padding: '8px', borderLeft: '4px solid #3b82f6', marginBottom: '15px' }}>1. ЗАГАЛЬНІ ПОКАЗНИКИ ДІЯЛЬНОСТІ</h2>
        <table style={{ width: '100%', borderCollapse: 'collapse', marginBottom: '30px', fontSize: '11pt' }}>
          <tbody>
            <tr>
              <td style={{ border: '1px solid #ccc', padding: '10px', width: '33%' }}><strong>Експлуатовані активи:</strong><br/>{data?.active_vehicles || 0} од.</td>
              <td style={{ border: '1px solid #ccc', padding: '10px', width: '33%' }}><strong>Критичний дефіцит:</strong><br/>{data?.critical_resources || 0} позицій</td>
              <td style={{ border: '1px solid #ccc', padding: '10px', width: '33%' }}><strong>Аномалії транзакцій:</strong><br/>{data?.fuel_anomalies || 0} інцидентів</td>
            </tr>
          </tbody>
        </table>

        {/* СЕКЦІЯ 2: ЖИТТЄВИЙ ЦИКЛ ТА РУХ */}
        <h2 style={{ fontSize: '12pt', backgroundColor: '#f0f0f0', padding: '8px', borderLeft: '4px solid #10b981', marginBottom: '15px' }}>2. РУХ АКТИВІВ ТА СПИСАННЯ (LIFECYCLE)</h2>
        <table style={{ width: '100%', borderCollapse: 'collapse', marginBottom: '30px', fontSize: '11pt' }}>
          <tbody>
            <tr>
              <td style={{ border: '1px solid #ccc', padding: '10px', width: '25%' }}><strong>Списано майна:</strong><br/>{data?.written_off_resources || 0} позицій</td>
              <td style={{ border: '1px solid #ccc', padding: '10px', width: '25%' }}><strong>Виконано переміщень:</strong><br/>{data?.completed_requests || 0} заявок</td>
              <td style={{ border: '1px solid #ccc', padding: '10px', width: '25%' }}><strong>Техніка в ремонті:</strong><br/>{data?.in_repair_vehicles || 0} од.</td>
              <td style={{ border: '1px solid #ccc', padding: '10px', width: '25%' }}><strong>Списані авто:</strong><br/>{data?.inactive_vehicles || 0} од.</td>
            </tr>
          </tbody>
        </table>

        {/* СЕКЦІЯ 3: ПОЛНИЙ СПИСОК ДЕФІЦИТУ */}
        <h2 style={{ fontSize: '12pt', backgroundColor: '#f0f0f0', padding: '8px', borderLeft: '4px solid #ef4444', marginBottom: '15px' }}>3. ДЕТАЛІЗАЦІЯ ДЕФІЦИТНИХ РЕСУРСІВ (ПОТРЕБА)</h2>
        <table style={{ width: '100%', borderCollapse: 'collapse', marginBottom: '30px', fontSize: '10pt' }}>
          <thead>
            <tr style={{ backgroundColor: '#e2e8f0' }}>
              <th style={{ border: '1px solid #94a3b8', padding: '8px', textAlign: 'left' }}>Найменування майна</th>
              <th style={{ border: '1px solid #94a3b8', padding: '8px', textAlign: 'center' }}>Факт. залишок</th>
              <th style={{ border: '1px solid #94a3b8', padding: '8px', textAlign: 'center' }}>Мінімальна норма</th>
              <th style={{ border: '1px solid #94a3b8', padding: '8px', textAlign: 'center' }}>Витрата (шт/день)</th>
              <th style={{ border: '1px solid #94a3b8', padding: '8px', textAlign: 'center' }}>Прогноз вичерпання</th>
            </tr>
          </thead>
          <tbody>
            {(!data?.predictive_burn_rate || data.predictive_burn_rate.length === 0) ? (
              <tr><td colSpan={5} style={{ border: '1px solid #ccc', padding: '8px', textAlign: 'center' }}>Показники в межах норми</td></tr>
            ) : (
              data.predictive_burn_rate.map((item: any, idx: number) => (
                <tr key={idx}>
                  <td style={{ border: '1px solid #ccc', padding: '8px' }}>{item.resource_name}</td>
                  <td style={{ border: '1px solid #ccc', padding: '8px', textAlign: 'center', fontWeight: 'bold' }}>{item.current_stock}</td>
                  <td style={{ border: '1px solid #ccc', padding: '8px', textAlign: 'center' }}>{item.min_quantity ?? '---'}</td>
                  <td style={{ border: '1px solid #ccc', padding: '8px', textAlign: 'center' }}>{(item.daily_burn_rate ?? 0).toFixed(1)}</td>
                  <td style={{ border: '1px solid #ccc', padding: '8px', textAlign: 'center', color: item.days_left <= 7 ? 'red' : 'black' }}>{item.days_left} дн.</td>
                </tr>
              ))
            )}
          </tbody>
        </table>

        {/* СЕКЦІЯ 4: РИЗИКИ АВТОПАРКУ */}
        <h2 style={{ fontSize: '12pt', backgroundColor: '#f0f0f0', padding: '8px', borderLeft: '4px solid #f59e0b', marginBottom: '15px' }}>4. ВИЯВЛЕНІ АНОМАЛІЇ ТА РИЗИКИ АВТОПАРКУ</h2>
        <table style={{ width: '100%', borderCollapse: 'collapse', marginBottom: '40px', fontSize: '10pt' }}>
          <thead>
            <tr style={{ backgroundColor: '#e2e8f0' }}>
              <th style={{ border: '1px solid #94a3b8', padding: '8px', textAlign: 'left' }}>Транспортний засіб (Ідентифікатор)</th>
              <th style={{ border: '1px solid #94a3b8', padding: '8px', textAlign: 'center' }}>Всього транзакцій (за період)</th>
              <th style={{ border: '1px solid #94a3b8', padding: '8px', textAlign: 'center' }}>Аномальних (Перевитрата)</th>
              <th style={{ border: '1px solid #94a3b8', padding: '8px', textAlign: 'center' }}>Рейтинг ризику</th>
            </tr>
          </thead>
          <tbody>
            {(!data?.fleet_risk || data.fleet_risk.length === 0) ? (
              <tr><td colSpan={4} style={{ border: '1px solid #ccc', padding: '8px', textAlign: 'center' }}>Аномалій не зафіксовано</td></tr>
            ) : (
              data.fleet_risk.map((risk: any, idx: number) => (
                <tr key={idx}>
                  <td style={{ border: '1px solid #ccc', padding: '8px', fontWeight: 'bold' }}>{risk.vehicle_name}</td>
                  <td style={{ border: '1px solid #ccc', padding: '8px', textAlign: 'center' }}>{risk.total_refuels}</td>
                  <td style={{ border: '1px solid #ccc', padding: '8px', textAlign: 'center', color: 'red', fontWeight: 'bold' }}>{risk.anomalies}</td>
                  <td style={{ border: '1px solid #ccc', padding: '8px', textAlign: 'center' }}>{risk.risk_score}%</td>
                </tr>
              ))
            )}
          </tbody>
        </table>

        {/* ПІДПИСИ */}
        <div style={{ marginTop: '40px', fontSize: '11pt' }}>
          <table style={{ width: '100%', border: 'none' }}>
            <tbody>
              <tr>
                <td style={{ width: '50%' }}>
                  <strong>Особа, що сформувала звіт:</strong><br/><br/>
                  ________________________ / ________________ /
                </td>
                <td style={{ width: '50%', textAlign: 'right' }}>
                  <strong>Керівник організації:</strong><br/><br/>
                  ________________________ / ________________ /
                </td>
              </tr>
            </tbody>
          </table>
        </div>

      </div>
    </div>
  );
};

export default AnalyticsDashboard;