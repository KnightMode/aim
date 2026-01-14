import { useDroppable } from '@dnd-kit/core';
import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable';
import { TaskCard } from './TaskCard';
import type { Task, TaskStatus } from '../../types';

interface ColumnProps {
  status: TaskStatus;
  title: string;
  tasks: Task[];
  onTaskClick: (task: Task) => void;
}

export const Column = ({ status, title, tasks, onTaskClick }: ColumnProps) => {
  const { setNodeRef, isOver } = useDroppable({
    id: status,
  });

  const taskIds = tasks.map((task) => task.id);

  return (
    <div style={{
      flex: '1 1 280px',
      minWidth: '280px',
      maxWidth: '320px',
      display: 'flex',
      flexDirection: 'column',
      background: '#EBECF0',
      borderRadius: '8px',
      overflow: 'hidden'
    }}>
      {/* Column Header */}
      <div style={{
        padding: '12px 10px 8px 10px',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between'
      }}>
        <span style={{
          fontSize: '12px',
          fontWeight: 600,
          color: '#5E6C84',
          textTransform: 'uppercase',
          letterSpacing: '0.04em'
        }}>
          {title}
        </span>
        <span style={{
          display: 'inline-flex',
          alignItems: 'center',
          justifyContent: 'center',
          minWidth: '24px',
          height: '24px',
          padding: '0 6px',
          background: 'rgba(9, 30, 66, 0.08)',
          borderRadius: '12px',
          fontSize: '11px',
          fontWeight: 700,
          color: '#5E6C84'
        }}>
          {tasks.length}
        </span>
      </div>

      {/* Column Body */}
      <div
        ref={setNodeRef}
        style={{
          flex: 1,
          padding: '6px',
          minHeight: '200px',
          background: isOver ? 'rgba(9, 30, 66, 0.08)' : 'transparent',
          transition: 'background 0.2s ease'
        }}
      >
        <div style={{
          display: 'flex',
          flexDirection: 'column',
          gap: '6px'
        }}>
          <SortableContext items={taskIds} strategy={verticalListSortingStrategy}>
            {tasks.map((task) => (
              <TaskCard key={task.id} task={task} onClick={() => onTaskClick(task)} />
            ))}
          </SortableContext>

          {tasks.length === 0 && (
            <div style={{
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              justifyContent: 'center',
              padding: '32px 16px',
              textAlign: 'center',
              minHeight: '100px'
            }}>
              <div style={{
                fontSize: '48px',
                opacity: 0.15,
                marginBottom: '8px',
                fontWeight: 300
              }}>
                ○
              </div>
              <div style={{
                fontSize: '13px',
                color: '#8993A4',
                fontWeight: 500
              }}>
                No tasks
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
