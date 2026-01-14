import { useState } from 'react';
import {
  DndContext,
  DragOverlay,
  PointerSensor,
  useSensor,
  useSensors,
} from '@dnd-kit/core';
import type { DragEndEvent, DragStartEvent } from '@dnd-kit/core';
import { Column } from './Column';
import { TaskCard } from './TaskCard';
import { useTasks, useUpdateTaskStatus } from '../../hooks/useTasks';
import type { Task, TaskStatus } from '../../types';
import { TASK_STATUSES, STATUS_LABELS } from '../../utils/constants';

interface KanbanBoardProps {
  onTaskClick: (task: Task) => void;
}

export const KanbanBoard = ({ onTaskClick }: KanbanBoardProps) => {
  const { data: tasks = [], isLoading, error } = useTasks();
  const updateTaskStatus = useUpdateTaskStatus();
  const [activeTask, setActiveTask] = useState<Task | null>(null);

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8,
      },
    })
  );

  const columns: TaskStatus[] = [
    TASK_STATUSES.TODO,
    TASK_STATUSES.QUEUED,
    TASK_STATUSES.IN_PROGRESS,
    TASK_STATUSES.COMPLETED,
    TASK_STATUSES.FAILED,
  ];

  const getTasksByStatus = (status: TaskStatus) => {
    return tasks.filter((task) => task.status === status);
  };

  const handleDragStart = (event: DragStartEvent) => {
    const task = tasks.find((t) => t.id === event.active.id);
    setActiveTask(task || null);
  };

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    setActiveTask(null);

    if (!over || active.id === over.id) {
      return;
    }

    const task = tasks.find((t) => t.id === active.id);
    const newStatus = over.id as TaskStatus;

    if (task && task.status !== newStatus) {
      updateTaskStatus.mutate({
        id: task.id,
        data: { status: newStatus },
      });
    }
  };

  if (isLoading) {
    return (
      <div style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        minHeight: '400px'
      }}>
        <div style={{ textAlign: 'center' }}>
          <div className="loading-spinner" style={{ margin: '0 auto var(--space-md) auto' }} />
          <div style={{
            color: 'var(--text-secondary)',
            fontSize: 'var(--font-size-base)'
          }}>
            Loading tasks...
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        minHeight: '400px'
      }}>
        <div className="error-box" style={{ maxWidth: '400px' }}>
          <div style={{
            fontWeight: 'var(--font-weight-semibold)',
            marginBottom: 'var(--space-sm)'
          }}>
            Error loading tasks
          </div>
          <div>
            Please check your connection and try refreshing the page.
          </div>
        </div>
      </div>
    );
  }

  return (
    <div style={{ position: 'relative' }}>
      <DndContext
        sensors={sensors}
        onDragStart={handleDragStart}
        onDragEnd={handleDragEnd}
      >
        <div style={{
          display: 'flex',
          gap: '10px',
          overflowX: 'auto',
          paddingBottom: '12px'
        }}>
          {columns.map((status) => (
            <Column
              key={status}
              status={status}
              title={STATUS_LABELS[status]}
              tasks={getTasksByStatus(status)}
              onTaskClick={onTaskClick}
            />
          ))}
        </div>

        <DragOverlay>
          {activeTask ? (
            <div style={{ opacity: 0.9, transform: 'rotate(2deg)' }}>
              <TaskCard task={activeTask} onClick={() => {}} />
            </div>
          ) : null}
        </DragOverlay>
      </DndContext>
    </div>
  );
};
