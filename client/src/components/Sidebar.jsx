/**
 * Sidebar
 *
 * A reference panel showing every supported command grouped by data type.
 * Clicking a command fills the input with a template — the user just
 * replaces the placeholder values.
 *
 * Props:
 *   onFill(template) — called with a command template string when clicked
 */
export default function Sidebar({ onFill }) {
  return (
    <aside style={styles.sidebar}>
      <CommandGroup label="Server" tag="srv" color="var(--red)" commands={[
        { name: 'PING',   template: 'PING' },
      ]} onFill={onFill} />

      <CommandGroup label="Strings" tag="str" color="var(--green)" commands={[
        { name: 'SET',    template: 'SET key value' },
        { name: 'GET',    template: 'GET key' },
        { name: 'DEL',    template: 'DEL key' },
        { name: 'EXISTS', template: 'EXISTS key' },
        { name: 'SET EX', template: 'SET key value EX 60' },
        { name: 'TTL',    template: 'TTL key' },
      ]} onFill={onFill} />

      <CommandGroup label="Lists" tag="list" color="var(--cyan)" commands={[
        { name: 'LPUSH',  template: 'LPUSH mylist value' },
        { name: 'LPOP',   template: 'LPOP mylist' },
      ]} onFill={onFill} />

      <CommandGroup label="Hashes" tag="hash" color="var(--yellow)" commands={[
        { name: 'HSET',   template: 'HSET user name Alice' },
        { name: 'HGET',   template: 'HGET user name' },
        { name: 'HDEL',   template: 'HDEL user name' },
      ]} onFill={onFill} />

      <CommandGroup label="Sets" tag="set" color="var(--magenta)" commands={[
        { name: 'SADD',      template: 'SADD myset member' },
        { name: 'SMEMBERS',  template: 'SMEMBERS myset' },
        { name: 'SISMEMBER', template: 'SISMEMBER myset member' },
        { name: 'SREM',      template: 'SREM myset member' },
      ]} onFill={onFill} />
    </aside>
  )
}

/**
 * CommandGroup
 *
 * One collapsible section in the sidebar — e.g. "Strings", "Lists".
 * Each item is a button that calls onFill with its template.
 */
function CommandGroup({ label, tag, color, commands, onFill }) {
  return (
    <div style={styles.group}>
      <div style={styles.groupLabel}>{label}</div>
      {commands.map(cmd => (
        <button
          key={cmd.name}
          style={styles.btn}
          onClick={() => onFill(cmd.template)}
          onMouseEnter={e => {
            e.currentTarget.style.background = 'rgba(122,162,247,0.07)'
            e.currentTarget.style.borderLeftColor = color
            e.currentTarget.style.color = '#c0caf5'
          }}
          onMouseLeave={e => {
            e.currentTarget.style.background = 'none'
            e.currentTarget.style.borderLeftColor = 'transparent'
            e.currentTarget.style.color = 'var(--text)'
          }}
        >
          <span style={styles.cmdName}>{cmd.name}</span>
          <span style={{ ...styles.tag, color, background: hexToAlpha(color, 0.12) }}>
            {tag}
          </span>
        </button>
      ))}
    </div>
  )
}

// converts a CSS var string like "var(--green)" to a usable rgba background
// since we can't use CSS vars inside JS style objects for rgba()
// we hard-code the alpha versions here as a small lookup
function hexToAlpha(colorVar, _alpha) {
  const map = {
    'var(--red)':     'rgba(247,118,142,0.12)',
    'var(--green)':   'rgba(158,206,106,0.12)',
    'var(--cyan)':    'rgba(125,207,255,0.12)',
    'var(--yellow)':  'rgba(224,175,104,0.12)',
    'var(--magenta)': 'rgba(187,154,247,0.12)',
  }
  return map[colorVar] || 'rgba(255,255,255,0.07)'
}

const styles = {
  sidebar: {
    width: 196,
    background: '#1f2335',
    borderRight: '1px solid #292e42',
    overflowY: 'auto',
    flexShrink: 0,
    paddingBottom: 16,
  },
  group: {
    marginBottom: 4,
  },
  groupLabel: {
    fontSize: 10,
    letterSpacing: '1.5px',
    textTransform: 'uppercase',
    color: '#565f89',
    padding: '12px 14px 5px',
  },
  btn: {
    width: '100%',
    background: 'none',
    border: 'none',
    borderLeft: '2px solid transparent',
    padding: '5px 12px 5px 12px',
    fontSize: 12,
    color: 'var(--text)',
    display: 'flex',
    alignItems: 'center',
    gap: 6,
    transition: 'all 0.08s',
    textAlign: 'left',
  },
  cmdName: {
    flex: 1,
    fontFamily: 'inherit',
  },
  tag: {
    fontSize: 10,
    padding: '1px 5px',
    borderRadius: 3,
    fontWeight: 700,
    letterSpacing: '0.3px',
  },
}
