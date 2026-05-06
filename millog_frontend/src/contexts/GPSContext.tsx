/**
 * GPSContext — глобальний GPS-трекінг для водіїв.
 * Живе на рівні App, не прив'язаний до конкретної сторінки.
 * Автоматично стартує/зупиняється залежно від наявності IN_TRANSIT рейсу.
 */
import { createContext, useContext, useEffect, useRef, useState, ReactNode } from 'react';
import { getInMemoryToken } from '../api/client';
import { useAuth } from './AuthContext';

export type GpsStatus = 'idle' | 'active' | 'error' | 'no_permission' | 'unavailable' | 'no_shipment';

interface GPSContextType {
  gpsStatus: GpsStatus;
  lastCoords: { lat: number; lng: number; speed: number | null } | null;
  hasActiveShipment: boolean;
}

const GPSContext = createContext<GPSContextType>({
  gpsStatus: 'idle',
  lastCoords: null,
  hasActiveShipment: false,
});

export function GPSProvider({ children }: { children: ReactNode }) {
  const { user } = useAuth();
  const [gpsStatus, setGpsStatus] = useState<GpsStatus>('idle');
  const [lastCoords, setLastCoords] = useState<{ lat: number; lng: number; speed: number | null } | null>(null);
  const [hasActiveShipment, setHasActiveShipment] = useState(false);

  const watchIdRef = useRef<number | null>(null);
  const pingIntervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const checkIntervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const lastPositionRef = useRef<GeolocationPosition | null>(null);

  // Перевіряємо активний рейс кожні 30 секунд
  useEffect(() => {
    if (!user) {
      stopGpsTracking();
      setHasActiveShipment(false);
      return;
    }

    // Негайна перевірка при логіні/рестарті
    checkActiveShipment();

    // Повторна перевірка кожні 30с (на випадок якщо рейс стартував/завершився)
    checkIntervalRef.current = setInterval(checkActiveShipment, 30_000);

    return () => {
      if (checkIntervalRef.current) clearInterval(checkIntervalRef.current);
    };
  }, [user]);

  // Стартуємо/зупиняємо GPS залежно від стану рейсу
  useEffect(() => {
    if (hasActiveShipment) {
      startGpsTracking();
    } else {
      stopGpsTracking();
    }
  }, [hasActiveShipment]);

  const checkActiveShipment = async () => {
    const authToken = getInMemoryToken();
    if (!authToken) return;
    try {
      const res = await fetch('/api/gps/driver/active-shipment', {
        headers: { 'Authorization': `Bearer ${authToken}` },
        credentials: 'include',
      });
      if (res.ok) {
        const data = await res.json();
        setHasActiveShipment(!!data.active);
      }
    } catch {
      // мовчки ігноруємо
    }
  };

  const sendPing = async (pos: GeolocationPosition) => {
    const authToken = getInMemoryToken();
    if (!authToken) return;
    try {
      const res = await fetch('/api/gps/driver/ping', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${authToken}` },
        credentials: 'include',
        body: JSON.stringify({
          latitude: pos.coords.latitude,
          longitude: pos.coords.longitude,
          altitude: pos.coords.altitude ?? 0,
          speed: pos.coords.speed ?? 0,
          heading: pos.coords.heading ?? 0,
          accuracy: pos.coords.accuracy,
        }),
      });
      if (res.ok) {
        setGpsStatus('active');
        setLastCoords({
          lat: pos.coords.latitude,
          lng: pos.coords.longitude,
          speed: pos.coords.speed !== null ? Math.round((pos.coords.speed ?? 0) * 3.6) : null,
        });
      } else if (res.status === 403) {
        // Рейс завершився — зупиняємо GPS і перевіряємо стан
        setHasActiveShipment(false);
        setGpsStatus('no_shipment');
      }
    } catch {
      // мовчки ігноруємо мережеву помилку
    }
  };

  const startGpsTracking = () => {
    // Вже запущено — не дублюємо
    if (watchIdRef.current !== null) return;
    if (!navigator.geolocation) {
      setGpsStatus('unavailable');
      return;
    }

    navigator.geolocation.getCurrentPosition(
      (pos) => {
        lastPositionRef.current = pos;
        sendPing(pos);

        watchIdRef.current = navigator.geolocation.watchPosition(
          (p) => { lastPositionRef.current = p; },
          (e) => setGpsStatus(e.code === 1 ? 'no_permission' : 'unavailable'),
          { enableHighAccuracy: true, timeout: 15000, maximumAge: 5000 }
        );

        pingIntervalRef.current = setInterval(() => {
          if (lastPositionRef.current) sendPing(lastPositionRef.current);
        }, 10_000);
      },
      (e) => setGpsStatus(e.code === 1 ? 'no_permission' : 'unavailable'),
      { enableHighAccuracy: true, timeout: 15000 }
    );
  };

  const stopGpsTracking = () => {
    if (watchIdRef.current !== null) {
      navigator.geolocation.clearWatch(watchIdRef.current);
      watchIdRef.current = null;
    }
    if (pingIntervalRef.current !== null) {
      clearInterval(pingIntervalRef.current);
      pingIntervalRef.current = null;
    }
    if (gpsStatus === 'active') setGpsStatus('idle');
  };

  return (
    <GPSContext.Provider value={{ gpsStatus, lastCoords, hasActiveShipment }}>
      {children}
    </GPSContext.Provider>
  );
}

export function useGPS() {
  return useContext(GPSContext);
}
