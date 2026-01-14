import { Button } from '../Common/Button';
import { useWebSocket } from '../../hooks/useWebSocket';

interface HeaderProps {
  onNewTask: () => void;
}

export const Header = ({ onNewTask }: HeaderProps) => {
  const { isConnected } = useWebSocket();

  return (
    <header style={{
      background: '#0052CC',
      borderBottom: '1px solid #0747A6',
      boxShadow: '0 1px 2px rgba(0, 0, 0, 0.08)'
    }}>
      <div style={{
        maxWidth: '1600px',
        margin: '0 auto',
        padding: '0 24px',
        height: '56px',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between'
      }}>
        {/* Left Side - Logo and Title */}
        <div style={{
          display: 'flex',
          alignItems: 'center',
          gap: '12px'
        }}>
          <svg width="28" height="28" viewBox="0 0 28 28" fill="none">
            <rect width="28" height="28" rx="4" fill="white"/>
            <path d="M14 8L19 11V17L14 20L9 17V11L14 8Z" fill="#0052CC"/>
            <circle cx="14" cy="14" r="2.5" fill="white"/>
          </svg>
          <h1 style={{
            fontSize: '16px',
            fontWeight: 600,
            color: 'white',
            lineHeight: 1,
            margin: 0,
            letterSpacing: '-0.01em'
          }}>
            AI Task Manager
          </h1>
        </div>

        {/* Right Side - Actions */}
        <div style={{
          display: 'flex',
          alignItems: 'center',
          gap: '12px'
        }}>
          <div style={{
            display: 'flex',
            alignItems: 'center',
            gap: '6px',
            padding: '6px 12px',
            background: 'rgba(255, 255, 255, 0.12)',
            borderRadius: '4px'
          }}>
            <div style={{
              width: '6px',
              height: '6px',
              background: isConnected ? '#36d399' : '#ef4444',
              borderRadius: '50%',
              boxShadow: isConnected 
                ? '0 0 6px rgba(54, 211, 153, 0.6)' 
                : '0 0 6px rgba(239, 68, 68, 0.6)'
            }} />
            <span style={{
              fontSize: '11px',
              fontWeight: 600,
              color: 'rgba(255, 255, 255, 0.95)',
              textTransform: 'uppercase',
              letterSpacing: '0.5px'
            }}>
              {isConnected ? 'System Active' : 'Disconnected'}
            </span>
          </div>
          <Button
            variant="primary"
            onClick={onNewTask}
            style={{
              background: 'white',
              color: '#0052CC',
              fontWeight: 600,
              fontSize: '14px',
              padding: '8px 16px',
              height: '36px',
              borderRadius: '4px',
              border: 'none',
              boxShadow: 'none'
            }}
          >
            + Create Task
          </Button>
        </div>
      </div>
    </header>
  );
};
