import { useState } from 'react';
import { Modal } from '../Common/Modal';
import { Button } from '../Common/Button';
import { StatusBadge } from '../Board/StatusBadge';
import { ExecutionLogs } from '../Execution/ExecutionLogs';
import { useDeleteTask, useUpdateTaskStatus } from '../../hooks/useTasks';
import type { Task } from '../../types';

interface TaskDetailModalProps {
  task: Task | null;
  isOpen: boolean;
  onClose: () => void;
}

export const TaskDetailModal = ({ task, isOpen, onClose }: TaskDetailModalProps) => {
  const [isLogsExpanded, setIsLogsExpanded] = useState(false);
  const deleteTask = useDeleteTask();
  const updateTaskStatus = useUpdateTaskStatus();

  if (!task) return null;

  const handleDelete = async () => {
    if (confirm('Are you sure you want to delete this task?')) {
      try {
        await deleteTask.mutateAsync(task.id);
        onClose();
      } catch (error) {
        console.error('Failed to delete task:', error);
      }
    }
  };

  const handleRetry = async () => {
    try {
      await updateTaskStatus.mutateAsync({
        id: task.id,
        data: { status: 'queued' },
      });
    } catch (error) {
      console.error('Failed to retry task:', error);
    }
  };

  const getPriorityInfo = (priority: number) => {
    if (priority === 2) return { text: 'High', color: '#DE350B' };
    if (priority === 1) return { text: 'Medium', color: '#FF8B00' };
    return { text: 'Low', color: '#8993A4' };
  };

  const priorityInfo = getPriorityInfo(task.priority);

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Task Details">
      <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
        {/* Task Header */}
        <div>
          <div style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'flex-start',
            marginBottom: '12px',
            gap: '12px'
          }}>
            <h3 style={{
              fontSize: '20px',
              fontWeight: 600,
              color: '#172B4D',
              margin: 0,
              flex: 1,
              lineHeight: 1.3
            }}>
              {task.title}
            </h3>
            <StatusBadge status={task.status} />
          </div>

          {task.description && (
            <p style={{
              fontSize: '14px',
              color: '#5E6C84',
              margin: '0 0 16px 0',
              lineHeight: 1.6
            }}>
              {task.description}
            </p>
          )}

          {/* Task Tags */}
          {task.tags.length > 0 && (
            <div style={{
              display: 'flex',
              flexWrap: 'wrap',
              gap: '8px'
            }}>
              {task.tags.map((tag) => (
                <span
                  key={tag}
                  style={{
                    display: 'inline-flex',
                    alignItems: 'center',
                    padding: '4px 12px',
                    background: '#DEEBFF',
                    color: '#0052CC',
                    borderRadius: '3px',
                    fontSize: '12px',
                    fontWeight: 500
                  }}
                >
                  {tag}
                </span>
              ))}
            </div>
          )}
        </div>

        {/* Task Metadata Grid */}
        <div style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(2, 1fr)',
          gap: '16px',
          padding: '16px',
          background: '#F4F5F7',
          borderRadius: '6px'
        }}>
          <div>
            <div style={{
              fontSize: '11px',
              fontWeight: 600,
              color: '#8993A4',
              textTransform: 'uppercase',
              letterSpacing: '0.5px',
              marginBottom: '4px'
            }}>
              Priority
            </div>
            <div style={{
              fontSize: '14px',
              fontWeight: 600,
              color: priorityInfo.color
            }}>
              {priorityInfo.text}
            </div>
          </div>

          {task.assigned_agent && (
            <div>
              <div style={{
                fontSize: '11px',
                fontWeight: 600,
                color: '#8993A4',
                textTransform: 'uppercase',
                letterSpacing: '0.5px',
                marginBottom: '4px'
              }}>
                Agent
              </div>
              <div style={{
                fontSize: '14px',
                fontWeight: 600,
                color: '#172B4D'
              }}>
                {task.assigned_agent}
              </div>
            </div>
          )}

          <div>
            <div style={{
              fontSize: '11px',
              fontWeight: 600,
              color: '#8993A4',
              textTransform: 'uppercase',
              letterSpacing: '0.5px',
              marginBottom: '4px'
            }}>
              Created
            </div>
            <div style={{
              fontSize: '14px',
              color: '#5E6C84'
            }}>
              {new Date(task.created_at).toLocaleString()}
            </div>
          </div>

          {task.completed_at && (
            <div>
              <div style={{
                fontSize: '11px',
                fontWeight: 600,
                color: '#8993A4',
                textTransform: 'uppercase',
                letterSpacing: '0.5px',
                marginBottom: '4px'
              }}>
                Completed
              </div>
              <div style={{
                fontSize: '14px',
                color: '#5E6C84'
              }}>
                {new Date(task.completed_at).toLocaleString()}
              </div>
            </div>
          )}
        </div>

        {/* Result - Only show for completed tasks */}
        {task.status === 'completed' && task.result && (
          <div style={{
            padding: '16px',
            background: '#E3FCEF',
            borderLeft: '4px solid #00875A',
            borderRadius: '6px'
          }}>
            <div style={{
              fontSize: '12px',
              fontWeight: 600,
              color: '#006644',
              textTransform: 'uppercase',
              letterSpacing: '0.5px',
              marginBottom: '8px'
            }}>
              ✓ Result
            </div>
            <p style={{
              fontSize: '14px',
              color: '#006644',
              margin: 0,
              lineHeight: 1.6,
              wordBreak: 'break-word'
            }}>
              {task.result}
            </p>
          </div>
        )}

        {/* Error - Only show for failed tasks */}
        {task.status === 'failed' && task.error_msg && (
          <div style={{
            padding: '16px',
            background: '#FFEBE6',
            borderLeft: '4px solid #DE350B',
            borderRadius: '6px'
          }}>
            <div style={{
              fontSize: '12px',
              fontWeight: 600,
              color: '#BF2600',
              textTransform: 'uppercase',
              letterSpacing: '0.5px',
              marginBottom: '8px'
            }}>
              ✕ Error
            </div>
            <p style={{
              fontSize: '14px',
              color: '#BF2600',
              margin: 0,
              lineHeight: 1.6,
              wordBreak: 'break-word'
            }}>
              {task.error_msg}
            </p>
          </div>
        )}

        {/* Metadata */}
        {task.metadata && Object.keys(task.metadata).length > 0 && (
          <div style={{
            padding: '16px',
            background: '#F4F5F7',
            borderRadius: '6px'
          }}>
            <div style={{
              fontSize: '12px',
              fontWeight: 600,
              color: '#5E6C84',
              textTransform: 'uppercase',
              letterSpacing: '0.5px',
              marginBottom: '12px'
            }}>
              Task Metadata
            </div>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
              {Object.entries(task.metadata).map(([key, value]) => (
                <div key={key} style={{
                  display: 'grid',
                  gridTemplateColumns: '140px 1fr',
                  gap: '12px',
                  fontSize: '13px'
                }}>
                  <div style={{
                    color: '#8993A4',
                    fontWeight: 500
                  }}>
                    {key.replace(/_/g, ' ')}:
                  </div>
                  <div style={{
                    color: '#172B4D',
                    wordBreak: 'break-word'
                  }}>
                    {Array.isArray(value) ? value.join(', ') : String(value)}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Execution Logs */}
        <div>
          <div
            onClick={() => setIsLogsExpanded(!isLogsExpanded)}
            style={{
              fontSize: '12px',
              fontWeight: 600,
              color: '#5E6C84',
              textTransform: 'uppercase',
              letterSpacing: '0.5px',
              marginBottom: '12px',
              cursor: 'pointer',
              display: 'flex',
              alignItems: 'center',
              gap: '8px',
              userSelect: 'none',
              transition: 'color 0.15s'
            }}
            onMouseEnter={(e) => e.currentTarget.style.color = '#172B4D'}
            onMouseLeave={(e) => e.currentTarget.style.color = '#5E6C84'}
          >
            <span style={{
              display: 'inline-block',
              transition: 'transform 0.2s',
              transform: isLogsExpanded ? 'rotate(90deg)' : 'rotate(0deg)'
            }}>
              ▶
            </span>
            Execution Logs
          </div>
          {isLogsExpanded && <ExecutionLogs taskId={task.id} />}
        </div>

        {/* Actions */}
        <div style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          paddingTop: '16px',
          borderTop: '2px solid #DFE1E6'
        }}>
          <div>
            {task.status === 'failed' && (
              <Button
                variant="primary"
                onClick={handleRetry}
                disabled={updateTaskStatus.isPending}
              >
                Retry Task
              </Button>
            )}
          </div>
          <div style={{ display: 'flex', gap: '8px' }}>
            <Button
              variant="danger"
              onClick={handleDelete}
              disabled={deleteTask.isPending}
            >
              Delete
            </Button>
            <Button variant="secondary" onClick={onClose}>
              Close
            </Button>
          </div>
        </div>
      </div>
    </Modal>
  );
};
