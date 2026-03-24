import { useState, useEffect, useRef } from 'react'

/**
 * useRedis — the brain of the app.
 *
 * Responsibilities:
 *  - Tracks connection status (pings server on mount)
 *  - Maintains the command history (for arrow-up/down recall)
 *  - Sends commands to the Go server via POST /api/command
 *  - Maintains the list of output lines shown in the terminal
 *
 * Every component that needs data gets it from here.
 * No component fetches directly — all network calls live in this hook.
 */
export function useRedis() {
  const [lines, setLines]         = useState([])          // terminal output entries
  const [connected, setConnected] = useState(null)        // null=checking, true, false
  const [cmdHistory, setCmdHistory] = useState([])        // previously run commands
  const [historyIdx, setHistoryIdx] = useState(-1)        // current arrow-up position
  const counterRef = useRef(0)                             // unique id per entry

  // ── Ping on mount to check server health ────────────────────────────────
  useEffect(() => {
    ping()
  }, [])

  async function ping() {
    try {
      const res = await fetch('/api/command', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ command: 'PING' }),
      })
      setConnected(res.ok)
    } catch {
      setConnected(false)
    }
  }

  // ── Send a command string to the server ─────────────────────────────────
  async function runCommand(raw) {
    if (!raw.trim()) return

    const id = ++counterRef.current

    // Add the command line to output immediately (before response arrives)
    addLine({ id, type: 'cmd', text: raw })

    // Save to history and reset index
    setCmdHistory(prev => [...prev, raw])
    setHistoryIdx(-1)

    // Add a placeholder "..." while we wait
    addLine({ id: id + 0.5, type: 'pending', parentId: id })

    try {
      const res = await fetch('/api/command', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ command: raw }),
      })
      const data = await res.json()

      // Replace the pending placeholder with the real response
      replacePending(id + 0.5, { id: id + 0.5, type: 'response', data })
    } catch {
      replacePending(id + 0.5, {
        id: id + 0.5,
        type: 'response',
        data: { error: 'network error — is the server running?' },
      })
    }
  }

  function addLine(line) {
    setLines(prev => [...prev, line])
  }

  function replacePending(pendingId, newLine) {
    setLines(prev => prev.map(l => l.id === pendingId ? newLine : l))
  }

  function clearLines() {
    setLines([])
  }

  // ── Arrow key history navigation (called from InputBar) ──────────────────
  function historyUp(currentInput) {
    if (cmdHistory.length === 0) return currentInput
    const newIdx = historyIdx < cmdHistory.length - 1 ? historyIdx + 1 : historyIdx
    setHistoryIdx(newIdx)
    return cmdHistory[cmdHistory.length - 1 - newIdx]
  }

  function historyDown() {
    if (historyIdx <= 0) {
      setHistoryIdx(-1)
      return ''
    }
    const newIdx = historyIdx - 1
    setHistoryIdx(newIdx)
    return cmdHistory[cmdHistory.length - 1 - newIdx]
  }

  return {
    lines,
    connected,
    runCommand,
    clearLines,
    historyUp,
    historyDown,
  }
}
