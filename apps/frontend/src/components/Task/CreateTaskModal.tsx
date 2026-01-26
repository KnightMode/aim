import { useState } from 'react';
import { Modal } from '../Common/Modal';
import { Button } from '../Common/Button';
import { useCreateTask } from '../../hooks/useTasks';
import { TASK_TAGS, PRIORITY_OPTIONS } from '../../utils/constants';
import type { CreateTaskRequest } from '../../types';

interface CreateTaskModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const CreateTaskModal = ({ isOpen, onClose }: CreateTaskModalProps) => {
  const createTask = useCreateTask();
  const [formData, setFormData] = useState<CreateTaskRequest>({
    title: '',
    description: '',
    tags: [],
    priority: 1,
    metadata: {},
  });

  const selectedTag = formData.tags[0] || '';

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      await createTask.mutateAsync(formData);
      onClose();
      // Reset form
      setFormData({
        title: '',
        description: '',
        tags: [],
        priority: 1,
        metadata: {},
      });
    } catch (error) {
      console.error('Failed to create task:', error);
    }
  };

  const handleTagChange = (tag: string) => {
    setFormData({ ...formData, tags: [tag], metadata: {} });
  };

  const handleMetadataChange = (key: string, value: string) => {
    setFormData({
      ...formData,
      metadata: { ...formData.metadata, [key]: value },
    });
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Create New Task">
      <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-lg)' }}>
        <div>
          <label>Title *</label>
          <input
            type="text"
            className="input"
            value={formData.title}
            onChange={(e) => setFormData({ ...formData, title: e.target.value })}
            required
          />
        </div>

        <div>
          <label>Description</label>
          <textarea
            className="textarea"
            rows={3}
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
          />
        </div>

        <div>
          <label>Task Type *</label>
          <select
            className="input"
            value={selectedTag}
            onChange={(e) => handleTagChange(e.target.value)}
            required
          >
            <option value="">Select task type</option>
            <option value={TASK_TAGS.CODING}>Coding</option>
            <option value={TASK_TAGS.DOCUMENTATION}>Documentation</option>
          </select>
        </div>

        <div>
          <label>Priority</label>
          <select
            className="input"
            value={formData.priority}
            onChange={(e) => setFormData({ ...formData, priority: parseInt(e.target.value) })}
          >
            {PRIORITY_OPTIONS.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </div>

        {/* Coding Task Metadata */}
        {selectedTag === TASK_TAGS.CODING && (
          <div style={{
            display: 'flex',
            flexDirection: 'column',
            gap: 'var(--space-lg)',
            paddingTop: 'var(--space-lg)',
            borderTop: '1px solid var(--border-default)'
          }}>
            <h3 style={{
              fontSize: 'var(--font-size-base)',
              fontWeight: 'var(--font-weight-semibold)',
              color: 'var(--text-primary)',
              margin: 0
            }}>
              Coding Task Details
            </h3>

            <div>
              <label>Repository URL *</label>
              <input
                type="url"
                className="input"
                placeholder="https://github.com/username/repo"
                value={(formData.metadata.repo_url as string) || ''}
                onChange={(e) => handleMetadataChange('repo_url', e.target.value)}
                required
              />
            </div>

            <div>
              <label>Branch (optional)</label>
              <input
                type="text"
                className="input"
                placeholder="main"
                value={(formData.metadata.branch as string) || ''}
                onChange={(e) => handleMetadataChange('branch', e.target.value)}
              />
            </div>

            <div className="info-box">
              Use the main Description field above to provide task details. The AI agent will automatically determine which files need to be modified.
            </div>
          </div>
        )}

        {/* Documentation Task Metadata */}
        {selectedTag === TASK_TAGS.DOCUMENTATION && (
          <div style={{
            display: 'flex',
            flexDirection: 'column',
            gap: 'var(--space-lg)',
            paddingTop: 'var(--space-lg)',
            borderTop: '1px solid var(--border-default)'
          }}>
            <h3 style={{
              fontSize: 'var(--font-size-base)',
              fontWeight: 'var(--font-weight-semibold)',
              color: 'var(--text-primary)',
              margin: 0
            }}>
              Documentation Task Details
            </h3>

            <div>
              <label>Confluence Space *</label>
              <input
                type="text"
                className="input"
                placeholder="DEV"
                value={(formData.metadata.confluence_space as string) || ''}
                onChange={(e) => handleMetadataChange('confluence_space', e.target.value)}
                required
              />
            </div>

            <div>
              <label>Page Title *</label>
              <input
                type="text"
                className="input"
                placeholder="API Documentation"
                value={(formData.metadata.page_title as string) || ''}
                onChange={(e) => handleMetadataChange('page_title', e.target.value)}
                required
              />
            </div>

            <div>
              <label>Content Source *</label>
              <select
                className="input"
                value={(formData.metadata.content_source as string) || 'text'}
                onChange={(e) => handleMetadataChange('content_source', e.target.value)}
              >
                <option value="text">Text</option>
                <option value="file">File</option>
                <option value="url">URL</option>
              </select>
            </div>

            <div>
              <label>Source Value *</label>
              <textarea
                className="textarea"
                rows={3}
                placeholder="File path, URL, or text content"
                value={(formData.metadata.source_value as string) || ''}
                onChange={(e) => handleMetadataChange('source_value', e.target.value)}
                required
              />
            </div>

            <div className="info-box">
              Use the main Description field above to provide documentation instructions.
            </div>
          </div>
        )}

        <div style={{
          display: 'flex',
          justifyContent: 'flex-end',
          gap: 'var(--space-md)',
          paddingTop: 'var(--space-lg)',
          borderTop: '1px solid var(--border-default)'
        }}>
          <Button type="button" variant="secondary" onClick={onClose}>
            Cancel
          </Button>
          <Button type="submit" variant="primary" disabled={createTask.isPending}>
            {createTask.isPending ? 'Creating...' : 'Create Task'}
          </Button>
        </div>
      </form>
    </Modal>
  );
};
