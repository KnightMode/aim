import { useState } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Header } from './components/Layout/Header';
import { KanbanBoard } from './components/Board/KanbanBoard';
import { CreateTaskModal } from './components/Task/CreateTaskModal';
import { TaskDetailModal } from './components/Task/TaskDetailModal';
import { useWebSocket } from './hooks/useWebSocket';
import type { Task } from './types';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
});

function AppContent() {
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [selectedTask, setSelectedTask] = useState<Task | null>(null);

  // Connect to WebSocket for real-time updates
  useWebSocket();

  const handleTaskClick = (task: Task) => {
    setSelectedTask(task);
  };

  return (
    <div className="min-h-screen" style={{ background: '#F4F5F7' }}>
      <Header onNewTask={() => setIsCreateModalOpen(true)} />

      <main className="max-w-[1600px] mx-auto px-6 py-4">
        <KanbanBoard onTaskClick={handleTaskClick} />
      </main>

      <CreateTaskModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
      />

      <TaskDetailModal
        task={selectedTask}
        isOpen={!!selectedTask}
        onClose={() => setSelectedTask(null)}
      />
    </div>
  );
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AppContent />
    </QueryClientProvider>
  );
}

export default App;
