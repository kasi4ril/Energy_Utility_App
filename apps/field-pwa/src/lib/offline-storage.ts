/**
 * Local offline storage manager for PWA field tasks and proof queues.
 */

import { FieldTask, OfflineProofSubmission } from "../types/pwa";

const TASKS_KEY = "nesco_pwa_cached_tasks";
const PROOFS_KEY = "nesco_pwa_offline_proofs";

export class OfflineStorage {
  public static getCachedTasks(): FieldTask[] {
    if (typeof window === "undefined") return [];
    const raw = localStorage.getItem(TASKS_KEY);
    return raw ? JSON.parse(raw) : [];
  }

  public static setCachedTasks(tasks: FieldTask[]): void {
    if (typeof window === "undefined") return;
    localStorage.setItem(TASKS_KEY, JSON.stringify(tasks));
  }

  public static queueOfflineProof(proof: OfflineProofSubmission): void {
    if (typeof window === "undefined") return;
    const existing = this.getQueuedProofs();
    existing.push(proof);
    localStorage.setItem(PROOFS_KEY, JSON.stringify(existing));
  }

  public static getQueuedProofs(): OfflineProofSubmission[] {
    if (typeof window === "undefined") return [];
    const raw = localStorage.getItem(PROOFS_KEY);
    return raw ? JSON.parse(raw) : [];
  }

  public static clearQueuedProof(taskId: string): void {
    if (typeof window === "undefined") return;
    const remaining = this.getQueuedProofs().filter((p) => p.taskId !== taskId);
    localStorage.setItem(PROOFS_KEY, JSON.stringify(remaining));
  }
}

