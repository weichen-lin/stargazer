'use client'

import React, { createContext, useContext, useEffect, useId, useRef, useState } from 'react'
import { AnimatePresence, MotionConfig, Variants, motion } from 'framer-motion'
import { ArrowLeftIcon, Loader2 } from 'lucide-react'
import { cn } from '@/lib/utils'

const TRANSITION = {
  type: 'spring',
  bounce: 0.1,
  duration: 0.4,
}

interface FloatingPanelContextType {
  isOpen: boolean
  openFloatingPanel: (rect: DOMRect) => void
  closeFloatingPanel: () => void
  uniqueId: string
  note: string
  setNote: (note: string) => void
  triggerRect: DOMRect | null
  title: string
  setTitle: (title: string) => void
  error: string | null
  setError: (error: string | null) => void
}

const FloatingPanelContext = createContext<FloatingPanelContextType | undefined>(undefined)

export function useFloatingPanel() {
  const context = useContext(FloatingPanelContext)
  if (!context) {
    throw new Error('useFloatingPanel must be used within a FloatingPanelProvider')
  }
  return context
}

function useFloatingPanelLogic(props: { defaultText?: string } = {}) {
  const { defaultText } = props

  const uniqueId = useId()
  const [isOpen, setIsOpen] = useState(false)
  const [note, setNote] = useState(defaultText ?? '')
  const [triggerRect, setTriggerRect] = useState<DOMRect | null>(null)
  const [title, setTitle] = useState('')
  const [error, setError] = useState<string | null>(null)

  const openFloatingPanel = (rect: DOMRect) => {
    setTriggerRect(rect)
    setIsOpen(true)
  }
  const closeFloatingPanel = () => {
    setIsOpen(false)
  }

  return {
    isOpen,
    openFloatingPanel,
    closeFloatingPanel,
    uniqueId,
    note,
    setNote,
    triggerRect,
    title,
    setTitle,
    error,
    setError,
  }
}

interface FloatingPanelRootProps {
  children: React.ReactNode
  className?: string
  defaultText?: string
}

export function FloatingPanelRoot({ children, className, defaultText }: FloatingPanelRootProps) {
  const floatingPanelLogic = useFloatingPanelLogic({ defaultText })

  return (
    <FloatingPanelContext.Provider value={floatingPanelLogic}>
      <MotionConfig transition={TRANSITION}>
        <div className={cn('relative', className)}>{children}</div>
      </MotionConfig>
    </FloatingPanelContext.Provider>
  )
}

interface FloatingPanelTriggerProps {
  children: React.ReactNode
  className?: string
  title: string
}

export function FloatingPanelTrigger({ children, className, title }: FloatingPanelTriggerProps) {
  const { openFloatingPanel, uniqueId, setTitle } = useFloatingPanel()
  const triggerRef = useRef<HTMLButtonElement>(null)

  const handleClick = () => {
    if (triggerRef.current) {
      openFloatingPanel(triggerRef.current.getBoundingClientRect())
      setTitle(title)
    }
  }

  return (
    <motion.button
      ref={triggerRef}
      layoutId={`floating-panel-trigger-${uniqueId}`}
      className={cn(
        'flex h-9 items-center border border-zinc-950/10 bg-white px-3 text-zinc-950 dark:border-zinc-50/10 dark:bg-zinc-700 dark:text-zinc-50',
        className,
      )}
      style={{ borderRadius: 8 }}
      onClick={handleClick}
      whileHover={{ scale: 1.05 }}
      whileTap={{ scale: 0.95 }}
      aria-haspopup='dialog'
      aria-expanded={false}
    >
      <motion.div layoutId={`floating-panel-label-container-${uniqueId}`} className='flex items-center'>
        <motion.span layoutId={`floating-panel-label-${uniqueId}`} className='text-sm font-semibold'>
          {children}
        </motion.span>
      </motion.div>
    </motion.button>
  )
}

interface FloatingPanelContentProps {
  children: React.ReactNode
  className?: string
}

export function FloatingPanelContent({ children, className }: FloatingPanelContentProps) {
  const { isOpen, closeFloatingPanel, uniqueId, triggerRect, title } = useFloatingPanel()
  const contentRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (contentRef.current && !contentRef.current.contains(event.target as Node)) {
        closeFloatingPanel()
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [closeFloatingPanel])

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') closeFloatingPanel()
    }
    document.addEventListener('keydown', handleKeyDown)
    return () => document.removeEventListener('keydown', handleKeyDown)
  }, [closeFloatingPanel])

  const variants: Variants = {
    hidden: { opacity: 0, scale: 0.9, y: 10 },
    visible: { opacity: 1, scale: 1, y: 0 },
  }

  return (
    <AnimatePresence>
      {isOpen && (
        <>
          <motion.div
            initial={{ backdropFilter: 'blur(0px)' }}
            animate={{ backdropFilter: 'blur(4px)' }}
            exit={{ backdropFilter: 'blur(0px)' }}
            className='fixed inset-0 z-40'
          />
          <motion.div
            ref={contentRef}
            layoutId={`floating-panel-${uniqueId}`}
            className={cn(
              'fixed z-50 overflow-hidden border border-zinc-950/10 bg-white shadow-lg outline-none dark:border-zinc-50/10 dark:bg-zinc-800',
              className,
            )}
            style={{
              borderRadius: 12,
              left: triggerRect ? triggerRect.left : '50%',
              top: triggerRect ? triggerRect.bottom + 8 : '50%',
              transformOrigin: 'top left',
            }}
            initial='hidden'
            animate='visible'
            exit='hidden'
            variants={variants}
            role='dialog'
            aria-modal='true'
            aria-labelledby={`floating-panel-title-${uniqueId}`}
          >
            <FloatingPanelTitle>{title}</FloatingPanelTitle>
            {children}
          </motion.div>
        </>
      )}
    </AnimatePresence>
  )
}

interface FloatingPanelTitleProps {
  children: React.ReactNode
}

function FloatingPanelTitle({ children }: FloatingPanelTitleProps) {
  const { uniqueId } = useFloatingPanel()

  return (
    <motion.div layoutId={`floating-panel-label-container-${uniqueId}`} className='px-4 py-2 bg-white dark:bg-zinc-800'>
      <motion.div
        layoutId={`floating-panel-label-${uniqueId}`}
        className='font-semibold text-zinc-900 dark:text-zinc-100 text-lg'
        id={`floating-panel-title-${uniqueId}`}
      >
        {children}
      </motion.div>
    </motion.div>
  )
}

interface FloatingPanelFormProps {
  children: React.ReactNode
  onSubmit?: (note: string) => void
  className?: string
}

export function FloatingPanelForm({ children, onSubmit, className }: FloatingPanelFormProps) {
  const { note } = useFloatingPanel()

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit?.(note)
  }

  return (
    <form className={cn('flex h-full flex-col', className)} onSubmit={handleSubmit}>
      {children}
    </form>
  )
}

interface FloatingPanelLabelProps {
  children: React.ReactNode
  htmlFor: string
  className?: string
}

export function FloatingPanelLabel({ children, htmlFor, className }: FloatingPanelLabelProps) {
  const { note } = useFloatingPanel()

  return (
    <motion.label
      htmlFor={htmlFor}
      style={{ opacity: note ? 0 : 1 }}
      className={cn('block mb-2 text-sm font-medium text-zinc-900 dark:text-zinc-100', className)}
    >
      {children}
    </motion.label>
  )
}

interface FloatingPanelTextareaProps {
  className?: string
  id?: string
  disabled?: boolean
  maxLength?: number
}

export function FloatingPanelTextarea({ className, id, disabled, maxLength }: FloatingPanelTextareaProps) {
  const { note, setNote, error, setError } = useFloatingPanel()

  return (
    <div className=''>
      <textarea
        id={id}
        className={cn('h-full w-full resize-none rounded-md bg-transparent px-0 py-3 text-md outline-none', className)}
        autoFocus
        value={note}
        onChange={e => {
          if (maxLength && e.target.value.length > maxLength) return
          if (error) setError(null)
          setNote(e.target.value)
        }}
        maxLength={maxLength}
        disabled={disabled}
      />
      <span className={cn('text-sm text-slate-300', note.length === maxLength && 'text-destructive/50')}>
        {note.length}/{maxLength}
      </span>
    </div>
  )
}

interface FloatingPanelHeaderProps {
  children: React.ReactNode
  className?: string
}

export function FloatingPanelHeader({ children, className }: FloatingPanelHeaderProps) {
  return (
    <motion.div
      className={cn('px-4 py-2 font-semibold text-zinc-900 dark:text-zinc-100', className)}
      initial={{ opacity: 0, y: -10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.1 }}
    >
      {children}
    </motion.div>
  )
}

interface FloatingPanelBodyProps {
  children: React.ReactNode
  className?: string
}

export function FloatingPanelBody({ children, className }: FloatingPanelBodyProps) {
  return (
    <motion.div
      className={cn('px-4', className)}
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.2 }}
    >
      {children}
    </motion.div>
  )
}

interface FloatingPanelFooterProps {
  children: React.ReactNode
  className?: string
}

export function FloatingPanelFooter({ children, className }: FloatingPanelFooterProps) {
  return (
    <motion.div
      className={cn('flex justify-between px-4 py-3', className)}
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.3 }}
    >
      {children}
    </motion.div>
  )
}

interface FloatingPanelCloseButtonProps {
  className?: string
}

export function FloatingPanelCloseButton({ className }: FloatingPanelCloseButtonProps) {
  const { closeFloatingPanel } = useFloatingPanel()

  return (
    <motion.button
      type='button'
      className={cn('flex items-center', className)}
      onClick={closeFloatingPanel}
      aria-label='Close floating panel'
      whileHover={{ scale: 1.1 }}
      whileTap={{ scale: 0.9 }}
    >
      <ArrowLeftIcon size={16} className='text-zinc-900 dark:text-zinc-100' />
    </motion.button>
  )
}

interface FloatingPanelSubmitButtonProps {
  className?: string
  isLoading: boolean
  text: string
  onClick: () => void
}

export function FloatingPanelSubmitButton({ className, isLoading, text, onClick }: FloatingPanelSubmitButtonProps) {
  const { note, error } = useFloatingPanel()

  const disabledButton = !note || note === '' || isLoading || error

  return (
    <motion.button
      className={cn(
        'bg-transparent px-2 text-sm text-zinc-900 transition-colors hover:bg-zinc-100 hover:text-zinc-800 focus-visible:ring-2 active:scale-[0.98] dark:border-zinc-50/10 dark:text-zinc-50 dark:hover:bg-zinc-800',
        'items-center justify-center border border-zinc-950/60 ',
        'relative ml-1 flex h-8 shrink-0 scale-100 select-none appearance-none rounded-lg',
        className,
        (isLoading || disabledButton) && 'opacity-50 pointer-events-none',
      )}
      type='submit'
      aria-label={text}
      whileHover={{ scale: 1.05 }}
      whileTap={{ scale: 0.95 }}
      onClick={onClick}
    >
      {isLoading ? <Loader2 className='m-2 h-4 w-4 animate-spin' /> : text}
    </motion.button>
  )
}

interface FloatingPanelButtonProps {
  children: React.ReactNode
  onClick?: () => void
  className?: string
}

export function FloatingPanelButton({ children, onClick, className }: FloatingPanelButtonProps) {
  return (
    <motion.button
      className={cn(
        'flex w-full items-center gap-2 rounded-md px-4 py-2 text-left text-sm hover:bg-zinc-100 dark:hover:bg-zinc-700',
        className,
      )}
      onClick={onClick}
      whileHover={{ backgroundColor: 'rgba(0, 0, 0, 0.05)' }}
      whileTap={{ scale: 0.98 }}
    >
      {children}
    </motion.button>
  )
}
