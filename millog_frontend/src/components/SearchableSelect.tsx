import { useState, useRef, useEffect, useCallback } from 'react'
import { createPortal } from 'react-dom'
import './SearchableSelect.css'

export interface SelectOption {
  value: string
  label: string
}

interface SearchableSelectProps {
  options: SelectOption[]
  value: string
  onChange: (value: string) => void
  placeholder?: string
  emptyLabel?: string
  disabled?: boolean
  searchPlaceholder?: string
  className?: string
  /** Рендерити список інлайн у потоці DOM (без portal/fixed).
   *  Використовуй там де дропдаун має скролитись разом з контейнером. */
  inlineDropdown?: boolean
}

export default function SearchableSelect({
  options,
  value,
  onChange,
  placeholder = 'Оберіть...',
  emptyLabel,
  disabled = false,
  searchPlaceholder = 'Пошук...',
  className = '',
  inlineDropdown = false,
}: SearchableSelectProps) {
  const [open, setOpen] = useState(false)
  const [query, setQuery] = useState('')
  const [dropdownStyle, setDropdownStyle] = useState<React.CSSProperties>({})
  const containerRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLInputElement>(null)

  const selectedLabel = value === ''
    ? (emptyLabel ?? placeholder)
    : (options.find(o => o.value === value)?.label ?? placeholder)

  const filtered = query.trim()
    ? options.filter(o => o.label.toLowerCase().includes(query.toLowerCase()))
    : options

  const updateDropdownPosition = useCallback(() => {
    if (inlineDropdown) return
    if (!containerRef.current) return
    const rect = containerRef.current.getBoundingClientRect()
    const spaceBelow = window.innerHeight - rect.bottom
    const dropdownHeight = Math.min(300, 44 + Math.min(filtered.length + (emptyLabel !== undefined ? 1 : 0), 6) * 37)

    if (spaceBelow >= dropdownHeight || spaceBelow >= 150) {
      setDropdownStyle({
        position: 'fixed',
        top: rect.bottom + 4,
        left: rect.left,
        width: rect.width,
        zIndex: 99999,
      })
    } else {
      setDropdownStyle({
        position: 'fixed',
        bottom: window.innerHeight - rect.top + 4,
        left: rect.left,
        width: rect.width,
        zIndex: 99999,
      })
    }
  }, [filtered.length, emptyLabel, inlineDropdown])

  const handleToggle = () => {
    if (disabled) return
    if (!open && !inlineDropdown) updateDropdownPosition()
    setOpen(prev => !prev)
  }

  const handleSelect = useCallback((val: string) => {
    onChange(val)
    setOpen(false)
    setQuery('')
  }, [onChange])

  // Закриваємо при кліку поза компонентом
  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        if (!inlineDropdown) {
          const dropdown = document.getElementById('ss-portal-dropdown')
          if (dropdown && dropdown.contains(e.target as Node)) return
        }
        setOpen(false)
        setQuery('')
      }
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [inlineDropdown])

  // При ресайзі — закриваємо portal-дропдаун
  useEffect(() => {
    if (!open || inlineDropdown) return
    const handleResize = () => { setOpen(false); setQuery('') }
    window.addEventListener('resize', handleResize)
    return () => window.removeEventListener('resize', handleResize)
  }, [open, inlineDropdown])

  // Focus search input when opened
  useEffect(() => {
    if (open) {
      setTimeout(() => inputRef.current?.focus(), 50)
    }
  }, [open])

  const dropdownContent = open ? (
    <div
      id={inlineDropdown ? undefined : 'ss-portal-dropdown'}
      className={`ss-dropdown${inlineDropdown ? ' ss-dropdown-inline' : ''}`}
      style={inlineDropdown ? undefined : dropdownStyle}
    >
      <div className="ss-search-wrap">
        <input
          ref={inputRef}
          type="text"
          className="ss-search"
          placeholder={searchPlaceholder}
          value={query}
          onChange={e => setQuery(e.target.value)}
          onClick={e => e.stopPropagation()}
        />
        {query && (
          <button type="button" className="ss-clear-btn" onClick={() => setQuery('')}>✕</button>
        )}
      </div>
      <ul className="ss-list">
        {emptyLabel !== undefined && (
          <li
            className={`ss-item ss-item-empty ${value === '' ? 'ss-item-selected' : ''}`}
            onClick={() => handleSelect('')}
          >
            {emptyLabel}
          </li>
        )}
        {filtered.length === 0 && (
          <li className="ss-item ss-item-none">Нічого не знайдено</li>
        )}
        {filtered.map(o => (
          <li
            key={o.value}
            className={`ss-item ${o.value === value ? 'ss-item-selected' : ''}`}
            onClick={() => handleSelect(o.value)}
          >
            {o.label}
          </li>
        ))}
      </ul>
    </div>
  ) : null

  return (
    <div
      ref={containerRef}
      className={`ss-container ${open ? 'ss-open' : ''} ${disabled ? 'ss-disabled' : ''} ${className}`}
    >
      <button
        type="button"
        className="ss-trigger erp-input"
        onClick={handleToggle}
        disabled={disabled}
        tabIndex={0}
      >
        <span className="ss-trigger-label">{selectedLabel}</span>
        <span className="ss-arrow">{open ? '▲' : '▼'}</span>
      </button>

      {inlineDropdown
        ? dropdownContent
        : createPortal(dropdownContent, document.body)
      }
    </div>
  )
}
