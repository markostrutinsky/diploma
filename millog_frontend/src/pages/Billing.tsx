import { useAuth } from '../contexts/AuthContext';
import toast, { Toaster } from 'react-hot-toast';
import './Billing.css';

export default function Billing() {
  const { user } = useAuth();
  
  const currentTier = user?.effective_subscription_tier || 'BASIC';
  const isPro = currentTier === 'PRO';

  const handleUpgradeRequest = (tier: string) => {
    toast.success(`Запит на оновлення до ${tier} надіслано! Наш менеджер зв'яжеться з вами для оформлення договору.`, {
      duration: 5000,
      icon: '📩',
    });
  };

  return (
    <div className="billing-page">
      <Toaster position="top-right" />
      <div className="page-header">
        <h1>Тарифні плани</h1>
        <p className="subtitle" style={{ color: '#64748b' }}>
          Оберіть рівень можливостей для вашої логістичної мережі. 
          Підписка автоматично розповсюджується на всі підпорядковані філії.
        </p>
      </div>

      <div className="plans-grid">
        {/* ПЛАН BASIC */}
        <div className={`plan-card ${currentTier === 'BASIC' ? 'current' : ''}`}>
          <div className="plan-header">
            <h3>BASIC</h3>
            <div className="price">Безкоштовно</div>
          </div>
          <ul className="features-list">
            <li>✅ До 10 складів</li>
            <li>✅ До 100 товарних позицій</li>
            <li>✅ До 50 користувачів</li>
            <li>✅ До 5 одиниць транспорту</li>
            <li>✅ Базовий облік та звітність</li>
            <li>✅ Ручне формування рейсів</li>
            <li>✅ Журнал аудиту (30 днів)</li>
          </ul>
          <button className="btn btn-secondary" disabled={currentTier === 'BASIC'}>
            {currentTier === 'BASIC' ? 'Поточний тариф' : 'Базовий план'}
          </button>
        </div>

        {/* ПЛАН PRO */}
        <div className={`plan-card pro ${currentTier === 'PRO' ? 'current' : ''}`}>
          <div className="popular-badge">Рекомендовано для бізнесу</div>
          <div className="plan-header">
            <h3>PRO</h3>
            <div className="price">4 999 грн<span>/міс</span></div>
          </div>
          <ul className="features-list">
            <li>📦 До 100 складів, 1000 товарів, 500 юзерів, 50 авто</li>
            <li>✨ <strong>Smart Dispatch</strong> — оптимізація маршрутів</li>
            <li>📊 <strong>Advanced Analytics</strong> — KPI, SLA, TCO, прогнози</li>
            <li>🔮 <strong>Predictive Maintenance</strong> — прогноз ТО транспорту</li>
            <li>🛡️ <strong>Fuel Anti-Fraud</strong> — виявлення крадіжок пального</li>
            <li>🌍 <strong>GPS Tracking</strong> — відстеження флоту в реальному часі</li>
            <li>📥 Excel масовий імпорт/експорт</li>
            <li>👨‍💻 Пріоритетна підтримка</li>
          </ul>
          <button 
            className="btn btn-primary" 
            style={{ background: isPro ? '#10b981' : 'linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)', border: 'none' }}
            onClick={() => handleUpgradeRequest('PRO')}
            disabled={isPro || currentTier === 'ENTERPRISE'}
          >
            {isPro ? '✅ Поточний тариф' : currentTier === 'ENTERPRISE' ? 'Доступно у вашому плані' : 'Запросити оновлення до PRO'}
          </button>
        </div>

        {/* ПЛАН ENTERPRISE */}
        <div className={`plan-card enterprise ${currentTier === 'ENTERPRISE' ? 'current' : ''}`}>
          <div className="popular-badge" style={{ background: 'linear-gradient(135deg, #f59e0b 0%, #dc2626 100%)' }}>Преміум</div>
          <div className="plan-header">
            <h3>ENTERPRISE</h3>
            <div className="price">Індивідуально</div>
          </div>
          <ul className="features-list">
            <li>♾️ <strong>Безлімітні</strong> ресурси (склади, товари, користувачі, транспорт)</li>
            <li>✨ Всі функції PRO</li>
            <li>🆘 <strong>Підтримка 24/7</strong></li>
            <li>📜 <strong>SLA гарантії</strong></li>
            <li>👨‍� Персональний менеджер</li>
            <li>🔧 Кастомні інтеграції</li>
            <li>🎓 Навчання команди</li>
          </ul>
          <button 
            className="btn btn-primary" 
            style={{ background: currentTier === 'ENTERPRISE' ? '#10b981' : 'linear-gradient(135deg, #f59e0b 0%, #dc2626 100%)', border: 'none' }}
            onClick={() => handleUpgradeRequest('ENTERPRISE')}
            disabled={currentTier === 'ENTERPRISE'}
          >
            {currentTier === 'ENTERPRISE' ? '✅ Поточний тариф' : 'Зв\'язатися з відділом продажів'}
          </button>
        </div>
      </div>
    </div>
  );
}