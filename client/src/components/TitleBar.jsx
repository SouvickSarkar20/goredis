/**
 * TitleBar
 *
 * The top bar of the terminal window. Shows:
 *  - macOS-style traffic light dots (purely decorative — makes it feel like a real terminal)
 *  - The window title
 *  - A live connection status badge (green dot when server is up)
 */
export default function TitleBar({ connected }) {
  return (
    <div style={styles.bar}>

      {/* macOS traffic light dots */}
      <div style={styles.dots}>
        <span style={{ ...styles.dot, background: '#ff5f57' }} />
        <span style={{ ...styles.dot, background: '#febc2e' }} />
        <span style={{ ...styles.dot, background: '#28c840' }} />
      </div>

      {/* centered window title */}
      <span style={styles.title}>
        goredis — playground
      </span>

      {/* connection badge — right side */}
      <div style={styles.badge}>
        <span style={{
          ...styles.connDot,
          background: connected === null
            ? '#565f89'                 // checking — grey
            : connected
              ? '#9ece6a'               // connected — green
              : '#f7768e',              // offline — red
          boxShadow: connected
            ? '0 0 6px #9ece6a88'
            : 'none',
        }} />
        <span style={styles.connText}>
          {connected === null
            ? 'connecting...'
            : connected
              ? 'connected · :6379'
              : 'server offline'}
        </span>
      </div>

    </div>
  )
}

const styles = {
  bar: {
    height: 40,
    background: '#13131a',
    borderBottom: '1px solid #292e42',
    display: 'flex',
    alignItems: 'center',
    padding: '0 16px',
    flexShrink: 0,
    userSelect: 'none',
    gap: 12,
  },
  dots: {
    display: 'flex',
    gap: 6,
  },
  dot: {
    width: 12,
    height: 12,
    borderRadius: '50%',
    display: 'inline-block',
  },
  title: {
    flex: 1,
    textAlign: 'center',
    fontSize: 12,
    color: '#565f89',
    letterSpacing: '0.3px',
  },
  badge: {
    display: 'flex',
    alignItems: 'center',
    gap: 6,
    fontSize: 11,
    color: '#565f89',
  },
  connDot: {
    width: 7,
    height: 7,
    borderRadius: '50%',
    display: 'inline-block',
    transition: 'background 0.3s, box-shadow 0.3s',
  },
  connText: {
    fontFamily: 'inherit',
  },
}
