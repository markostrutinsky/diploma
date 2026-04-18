import { useAuth } from '../contexts/AuthContext';
import toast, { Toaster } from 'react-hot-toast';
import './Billing.css';

export default function Billing() {
  const { user } = useAuth();
  
  // Перевіряємо, чи є у користувача або його відділу вже PRO-план
  const isPro = user?.effective_subscription_tier === 'PRO' || user?.effective_subscription_tier === 'ENTERPRISE';

  const handleUpgradeRequest = () => {
    toast.success('Запит надіслано! Наш менеджер зв’яжеться з вами для оформлення договору та виставлення рахунку.', {
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
        <div className={`plan-card ${!isPro ? 'current' : ''}`}>
          <div className="plan-header">
            <h3>Standard (Basic)</h3>
            <div className="price">0 грн<span>/міс</span></div>
          </div>
          <ul className="features-list">
            <li>✅ Повний облік майна та складів</li>
            <li>✅ Базова звітність (PDF/Excel)</li>
            <li>✅ Ручне формування рейсів</li>
            <li>✅ Журнал аудиту (7 днів)</li>
          </ul>
          <button className="btn btn-secondary" disabled={!isPro}>
            {!isPro ? 'Ваш поточний тариф' : 'Базовий план'}
          </button>
        </div>

        {/* ПЛАН PRO */}
        <div className={`plan-card pro ${isPro ? 'current' : ''}`}>
          <div className="popular-badge">Рекомендовано для B2B</div>
          <div className="plan-header">
            <h3>Enterprise (PRO)</h3>
            <div className="price">4 999 грн<span>/міс</span></div>
          </div>
          <ul className="features-list">
            <li>✨ <strong>Smart Розподіл</strong> (Алгоритм оптимізації рейсів)</li>
            <li>📊 Розширена аналітика (SLA, Ризики, TCO)</li>
            <li>🛡️ Антифрод-система для контролю пального</li>
            <li>⛓️ <strong>Успадкування підписки</strong> для всіх дочірніх філій</li>
            <li>👨‍💻 Пріоритетна технічна підтримка</li>
          </ul>
          <button 
            className="btn btn-primary" 
            style={{ background: isPro ? '#10b981' : 'linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)', border: 'none' }}
            onClick={handleUpgradeRequest}
            disabled={isPro}
          >
            {isPro ? '✅ Тариф активовано' : 'Запросити оновлення до PRO'}
          </button>
        </div>
      </div>
    </div>
  );
}