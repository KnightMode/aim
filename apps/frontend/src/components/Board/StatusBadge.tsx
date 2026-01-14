import { STATUS_LABELS } from '../../utils/constants';
import type { TaskStatus } from '../../types';

interface StatusBadgeProps {
  status: TaskStatus;
}

export const StatusBadge = ({ status }: StatusBadgeProps) => {
  return (
    <span className={`status-badge status-${status}`}>
      {STATUS_LABELS[status]}
    </span>
  );
};
