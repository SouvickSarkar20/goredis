import { useState, useRef, useEffect } from 'react'

/**
 * InputBar
 *
 * The command input row pinned to the bottom of the terminal.
 * Looks exactly like the redis-cli input line.
 *
 * Features:
 *   - Enter to run
 *   - Arrow up/down for command history
 *   - Auto-focus on mount
 *   - "CLEAR" button to wipe the terminal
 *
 * Props:
 *   onRun(command)    — called with the command string on Enter
 *   onClear()         — called when CLEAR is clicked
 *   onHistoryUp()     — returns the previous command string
 *   onHistoryDown()   — returns the next command string (or empty)
 *   fillValue         — when sidebar fills the input, this prop changes
 */
export default function InputBar({ onRun, onClear, onHistoryUp, onHistoryDown, fillValue }) {
  const [value, setValue] = useState('')
  const inputRef = useRef(null)

  // Auto-focus on mount
  useEffect(() => {
    inputRef.current?.focus()
  }, [])

  // When sidebar sets a fill value, update the input
  useEffect(() => {
    if (fillValue) {
      setValue(fillValue)
      inputRef.current?.focus()
    }
  }, [fillValue])

  function handleKeyDown(e) {
    if (e.key === 'Enter') {
      const cmd = value.trim()
      if (!cmd) return
      onRun(cmd)
      setValue('')
    } else if (e.key === 'ArrowUp') {
      e.preventDefault()   // prevent cursor moving to start of input
      const prev = onHistoryUp()
      if (prev !== undefined) setValue(prev)
    } else if (e.key === 'ArrowDown') {
      e.preventDefault()
      const next = onHistoryDown()
      if (next !== undefined) setValue(next)
    }
  }

  return (
    <div style={styles.bar}>

      {/* The redis-cli style prompt */}
      <span style={styles.prompt} aria-hidden>
        127.0.0.1:6379&gt;
      </span>

      {/* The actual input */}
      <input
        ref={inputRef}
        type="text"
        value={value}
        onChange={e => setValue(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder="type a command..."
        autoComplete="off"
        autoCorrect="off"
        autoCapitalize="off"
        spellCheck={false}
        style={styles.input}
        aria-label="Redis command input"
      />

      {/* Clear button */}
      <button
        style={styles.clearBtn}
        onClick={onClear}
        onMouseEnter={e => e.currentTarget.style.color = 'var(--text)'}
        onMouseLeave={e => e.currentTarget.style.color = 'var(--comment)'}
      >
        clear
      </button>

    </div>
  )
}

const styles = {
  bar: {
    display: 'flex',
    alignItems: 'center',
    gap: 10,
    padding: '10px 20px',
    borderTop: '1px solid #292e42',
    background: '#1a1b26',
    flexShrink: 0,
  },
  prompt: {
    color: 'var(--green)',
    flexShrink: 0,
    userSelect: 'none',
    fontSize: 13,
  },
  input: {
    flex: 1,
    background: 'none',
    border: 'none',
    outline: 'none',
    color: 'var(--bright)',
    fontSize: 13,
    caretColor: 'var(--cyan)',
    lineHeight: 1.5,
  },
  clearBtn: {
    background: 'none',
    border: 'none',
    color: 'var(--comment)',
    fontSize: 11,
    letterSpacing: '0.5px',
    padding: '4px 8px',
    borderRadius: 3,
    transition: 'color 0.1s',
    flexShrink: 0,
  },
}
