import { useState } from 'react'
import { useRedis } from './hooks/useRedis.js'
import TitleBar from './components/TitleBar.jsx'
import Sidebar from './components/Sidebar.jsx'
import Terminal from './components/Terminal.jsx'
import InputBar from './components/InputBar.jsx'

/**
 * App
 *
 * The root component. Its only job is to:
 *   1. Call useRedis() to get all state and actions
 *   2. Pass them down to the right child components
 *   3. Handle the sidebar "fill" — when a sidebar button is clicked,
 *      App passes the template string to InputBar via fillValue state
 *
 * Layout:
 *   ┌─────────────────────────────────┐
 *   │           TitleBar              │
 *   ├──────────┬──────────────────────┤
 *   │          │      Terminal        │
 *   │ Sidebar  │  (scrollable output) │
 *   │          ├──────────────────────┤
 *   │          │      InputBar        │
 *   └──────────┴──────────────────────┘
 */
export default function App() {
  const {
    lines,
    connected,
    runCommand,
    clearLines,
    historyUp,
    historyDown,
  } = useRedis()

  // fillValue is how the Sidebar tells InputBar what to put in the input.
  // We use an object { value, ts } so that clicking the same command twice
  // still triggers the useEffect in InputBar (object reference changes).
  const [fillValue, setFillValue] = useState(null)

  function handleSidebarFill(template) {
    setFillValue({ value: template, ts: Date.now() })
  }

  return (
    <div style={styles.app}>

      <TitleBar connected={connected} />

      <div style={styles.body}>

        <Sidebar onFill={handleSidebarFill} />

        <div style={styles.main}>
          <Terminal lines={lines} />
          <InputBar
            onRun={runCommand}
            onClear={clearLines}
            onHistoryUp={historyUp}
            onHistoryDown={historyDown}
            fillValue={fillValue?.value}
            fillTs={fillValue?.ts}
          />
        </div>

      </div>

    </div>
  )
}

const styles = {
  app: {
    height: '100vh',
    display: 'flex',
    flexDirection: 'column',
    overflow: 'hidden',
    background: 'var(--bg)',
  },
  body: {
    flex: 1,
    display: 'flex',
    overflow: 'hidden',
  },
  main: {
    flex: 1,
    display: 'flex',
    flexDirection: 'column',
    overflow: 'hidden',
  },
}
