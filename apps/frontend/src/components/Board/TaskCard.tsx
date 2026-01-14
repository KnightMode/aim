import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { StatusBadge } from './StatusBadge';
import type { Task } from '../../types';

interface TaskCardProps {
  task: Task;
  onClick: () => void;
}

export const TaskCard = ({ task, onClick }: TaskCardProps) => {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: task.id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  const getPriorityLabel = (priority: number) => {
    if (priority === 2) return { text: 'High', class: 'priority-high' };
    if (priority === 1) return { text: 'Medium', class: 'priority-medium' };
    return { text: 'Low', class: 'priority-low' };
  };

  const priorityInfo = getPriorityLabel(task.priority);

  return (
    <div
      ref={setNodeRef}
      style={{
        ...style,
        background: 'white',
        border: '1px solid #DFE1E6',
        borderRadius: '6px',
        padding: '12px',
        cursor: 'pointer',
        transition: 'all 0.2s ease',
        boxShadow: '0 1px 2px rgba(9, 30, 66, 0.08)'
      }}
      {...attributes}
      {...listeners}
      onClick={onClick}
      onMouseEnter={(e) => {
        e.currentTarget.style.background = '#FAFBFC';
        e.currentTarget.style.boxShadow = '0 4px 8px rgba(9, 30, 66, 0.12)';
        e.currentTarget.style.borderColor = '#C1C7D0';
        e.currentTarget.style.transform = 'translateY(-1px)';
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.background = 'white';
        e.currentTarget.style.boxShadow = '0 1px 2px rgba(9, 30, 66, 0.08)';
        e.currentTarget.style.borderColor = '#DFE1E6';
        e.currentTarget.style.transform = 'translateY(0)';
      }}
    >
      {/* Task Title and Status */}
      <div style={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'flex-start',
        marginBottom: '6px',
        gap: '8px'
      }}>
        <h3 style={{
          fontSize: '14px',
          fontWeight: 600,
          color: '#172B4D',
          margin: 0,
          flex: 1,
          lineHeight: 1.4
        }}>
          {task.title}
        </h3>
        <StatusBadge status={task.status} />
      </div>

      {/* Task Description */}
      {task.description && (
        <p style={{
          fontSize: '13px',
          color: '#5E6C84',
          margin: '0 0 8px 0',
          lineHeight: 1.5
        }} className="line-clamp-2">
          {task.description}
        </p>
      )}

      {/* Task Tags */}
      {task.tags.length > 0 && (
        <div style={{
          display: 'flex',
          flexWrap: 'wrap',
          gap: '6px',
          marginBottom: '8px'
        }}>
          {task.tags.map((tag) => (
            <span
              key={tag}
              style={{
                display: 'inline-flex',
                alignItems: 'center',
                padding: '2px 8px',
                background: '#F4F5F7',
                color: '#5E6C84',
                borderRadius: '3px',
                fontSize: '11px',
                fontWeight: 500
              }}
            >
              {tag}
            </span>
          ))}
        </div>
      )}

      {/* Task Metadata */}
      <div style={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        paddingTop: '6px',
        borderTop: '1px solid #DFE1E6',
        fontSize: '11px'
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
          <span style={{ color: '#8993A4' }}>Priority:</span>
          <span className={priorityInfo.class} style={{ fontWeight: 600 }}>
            {priorityInfo.text}
          </span>
        </div>

        {task.assigned_agent && (
          <span style={{
            color: '#5E6C84',
            fontSize: '11px',
            fontWeight: 500
          }}>
            {task.assigned_agent}
          </span>
        )}
      </div>

      {/* Task Result - Only show for completed tasks */}
      {task.status === 'completed' && task.result && (
        <div style={{
          marginTop: '8px',
          padding: '8px 10px',
          background: '#E3FCEF',
          borderLeft: '3px solid #00875A',
          borderRadius: '3px'
        }}>
          <div style={{
            fontSize: '11px',
            fontWeight: 600,
            marginBottom: '4px',
            color: '#006644'
          }}>
            ✓ Result
          </div>
          <div className="line-clamp-2" style={{
            fontSize: '12px',
            color: '#006644',
            lineHeight: 1.4
          }}>
            {task.result}
          </div>
        </div>
      )}

      {/* Task Error - Only show for failed tasks */}
      {task.status === 'failed' && task.error_msg && (
        <div style={{
          marginTop: '8px',
          padding: '8px 10px',
          background: '#FFEBE6',
          borderLeft: '3px solid #DE350B',
          borderRadius: '3px'
        }}>
          <div style={{
            fontSize: '11px',
            fontWeight: 600,
            marginBottom: '4px',
            color: '#BF2600'
          }}>
            ✕ Error
          </div>
          <div className="line-clamp-2" style={{
            fontSize: '12px',
            color: '#BF2600',
            lineHeight: 1.4
          }}>
            {task.error_msg}
          </div>
        </div>
      )}
    </div>
  );
};
