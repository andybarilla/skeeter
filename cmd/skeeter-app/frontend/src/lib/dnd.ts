import { MoveTask } from '../../wailsjs/go/main/App';
import { refreshBoard } from './stores/board';
import { notify, notifyError } from './stores/notifications';

const MIME = 'application/x-skeeter-task';

export function handleDragStart(e: DragEvent, taskId: string) {
  if (!e.dataTransfer) return;
  e.dataTransfer.setData(MIME, taskId);
  e.dataTransfer.effectAllowed = 'move';
  const el = e.target as HTMLElement;
  el.style.opacity = '0.4';
}

export function handleDragEnd(e: DragEvent) {
  const el = e.target as HTMLElement;
  el.style.opacity = '1';
}

export function handleDragOver(e: DragEvent) {
  e.preventDefault();
  if (e.dataTransfer) {
    e.dataTransfer.dropEffect = 'move';
  }
}

export function handleDragEnter(e: DragEvent) {
  const el = (e.currentTarget as HTMLElement);
  el.classList.add('drop-target');
}

export function handleDragLeave(e: DragEvent) {
  const el = (e.currentTarget as HTMLElement);
  // Only remove if leaving the column itself, not a child
  const related = e.relatedTarget as HTMLElement | null;
  if (related && el.contains(related)) return;
  el.classList.remove('drop-target');
}

export async function handleDrop(e: DragEvent, targetStatus: string) {
  e.preventDefault();
  const el = (e.currentTarget as HTMLElement);
  el.classList.remove('drop-target');

  if (!e.dataTransfer) return;
  const taskId = e.dataTransfer.getData(MIME);
  if (!taskId) return;

  try {
    await MoveTask(taskId, targetStatus);
    await refreshBoard();
    notify('success', `Moved ${taskId} to ${targetStatus}`);
  } catch (err) {
    notifyError(err);
  }
}
