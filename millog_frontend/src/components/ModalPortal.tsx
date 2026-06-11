import { useEffect, useMemo } from 'react';
import { createPortal } from 'react-dom';
import type { ReactNode } from 'react';

interface ModalPortalProps {
  children: ReactNode;
}

let lockCount = 0;
let savedScrollY = 0;
let savedBodyPaddingRight = '';

const lockDocumentScroll = () => {
  const root = document.documentElement;
  const body = document.body;

  if (lockCount === 0) {
    savedScrollY = window.scrollY;
    savedBodyPaddingRight = body.style.paddingRight;
  }

  const scrollbarWidth = window.innerWidth - root.clientWidth;

  root.classList.add('app-modal-open');
  body.classList.add('app-modal-open');
  body.style.paddingRight = scrollbarWidth > 0 ? `${scrollbarWidth}px` : savedBodyPaddingRight;
  body.style.position = 'fixed';
  body.style.top = `-${savedScrollY}px`;
  body.style.left = '0';
  body.style.right = '0';
  body.style.width = '100%';
};

const unlockDocumentScroll = () => {
  const root = document.documentElement;
  const body = document.body;

  root.classList.remove('app-modal-open');
  body.classList.remove('app-modal-open');
  body.style.paddingRight = savedBodyPaddingRight;
  body.style.position = '';
  body.style.top = '';
  body.style.left = '';
  body.style.right = '';
  body.style.width = '';
  window.scrollTo(0, savedScrollY);
};

export default function ModalPortal({ children }: ModalPortalProps) {
  const container = useMemo(() => {
    if (typeof document === 'undefined') {
      return null;
    }

    const element = document.createElement('div');
    element.setAttribute('data-modal-portal', 'true');
    return element;
  }, []);

  useEffect(() => {
    if (!container) return;

    document.body.appendChild(container);
    lockCount += 1;

    if (lockCount === 1) {
      lockDocumentScroll();
    }

    return () => {
      if (container.parentNode) {
        container.parentNode.removeChild(container);
      }

      lockCount = Math.max(0, lockCount - 1);
      if (lockCount === 0) {
        unlockDocumentScroll();
      }
    };
  }, [container]);

  if (!container) return null;

  return createPortal(children, container);
}
