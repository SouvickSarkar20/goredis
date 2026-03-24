/**
 * ResponseLine
 *
 * Renders a single line (or group of lines) in the terminal output.
 * There are three line types:
 *
 *   "cmd"      — the command the user typed, shown with a green prompt
 *   "pending"  — a "..." placeholder shown while waiting for the server
 *   "response" — the server's response, color-coded by what was returned
 *
 * Color coding mirrors real redis-cli exactly:
 *   Simple string (OK, PONG) → green
 *   Error                    → red
 *   Integer                  → cyan (shown as "(integer) N")
 *   Bulk string (a value)    → no color prefix, white text
 *   Nil                      → dim italic "(nil)"
 *   Array                    → each element numbered and indented
 */
export default function ResponseLine({ line }) {
  if (line.type === 'cmd') {
    return (
      <div style={styles.cmdRow}>
        <span style={styles.prompt}>127.0.0.1:6379&gt;</span>
        <span style={styles.cmdText}>{line.text}</span>
      </div>
    )
  }

  if (line.type === 'pending') {
    return <div style={styles.pending}>...</div>
  }

  if (line.type === 'response') {
    return <ResponseContent data={line.data} />
  }

  return null
}

/**
 * ResponseContent
 *
 * Decides how to render based on what the server returned.
 * Mirrors how redis-cli formats each RESP type.
 */
function ResponseContent({ data }) {
  // ── Error ──────────────────────────────────────────────────────────────
  if (data.error) {
    return (
      <div style={{ ...styles.responseLine, color: 'var(--red)', marginBottom: 6 }}>
        (error) {data.error}
      </div>
    )
  }

  const r = data.result

  // ── Nil ────────────────────────────────────────────────────────────────
  if (r === null || r === undefined) {
    return (
      <div style={{ ...styles.responseLine, color: 'var(--comment)', fontStyle: 'italic', marginBottom: 6 }}>
        (nil)
      </div>
    )
  }

  // ── Integer ────────────────────────────────────────────────────────────
  if (typeof r === 'number') {
    return (
      <div style={{ ...styles.responseLine, color: 'var(--cyan)', marginBottom: 6 }}>
        (integer) {r}
      </div>
    )
  }

  // ── Simple string (OK, PONG) ────────────────────────────────────────────
  if (typeof r === 'string' && (r === 'OK' || r === 'PONG')) {
    return (
      <div style={{ ...styles.responseLine, color: 'var(--green)', marginBottom: 6 }}>
        {r}
      </div>
    )
  }

  // ── Bulk string (any other string value) ───────────────────────────────
  if (typeof r === 'string') {
    return (
      <div style={{ ...styles.responseLine, color: 'var(--bright)', marginBottom: 6 }}>
        &quot;{r}&quot;
      </div>
    )
  }

  // ── Array ──────────────────────────────────────────────────────────────
  if (Array.isArray(r)) {
    if (r.length === 0) {
      return (
        <div style={{ ...styles.responseLine, color: 'var(--comment)', fontStyle: 'italic', marginBottom: 6 }}>
          (empty array)
        </div>
      )
    }
    return (
      <div style={{ marginBottom: 6 }}>
        {r.map((item, i) => (
          <div key={i} style={styles.arrayItem}>
            <span style={styles.arrayIndex}>{i + 1})</span>
            <span style={{ color: 'var(--bright)' }}>&quot;{item}&quot;</span>
          </div>
        ))}
      </div>
    )
  }

  // ── Fallback ────────────────────────────────────────────────────────────
  return (
    <div style={{ ...styles.responseLine, color: 'var(--text)', marginBottom: 6 }}>
      {JSON.stringify(r)}
    </div>
  )
}

const styles = {
  cmdRow: {
    display: 'flex',
    alignItems: 'baseline',
    gap: 10,
    marginTop: 10,
    marginBottom: 2,
  },
  prompt: {
    color: 'var(--green)',
    flexShrink: 0,
    userSelect: 'none',
  },
  cmdText: {
    color: 'var(--bright)',
    wordBreak: 'break-all',
  },
  pending: {
    color: 'var(--comment)',
    paddingLeft: 2,
    marginBottom: 6,
    letterSpacing: 2,
  },
  responseLine: {
    paddingLeft: 2,
    lineHeight: 1.6,
  },
  arrayItem: {
    display: 'flex',
    gap: 10,
    paddingLeft: 2,
    lineHeight: 1.7,
  },
  arrayIndex: {
    color: 'var(--comment)',
    minWidth: 24,
    flexShrink: 0,
  },
}
