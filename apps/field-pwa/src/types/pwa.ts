/**
 * Field PWA types for offline caching and task synchronization.
 */

export interface FieldTask {
  id: string;
  orderNo: string;
  unitType: string;
  assetId: string;
  priority: "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";
  instructions: string;
  targetLat: number;
  targetLon: number;
  toleranceMeters: number;
  status: "PENDING" | "SYNCED_PROOF" | "OFFLINE_DRAFT";
  ptwRequired: boolean;
}

export interface OfflineProofSubmission {
  taskId: string;
  orderNo: string;
  capturedLat: number;
  capturedLon: number;
  accuracyMeters: number;
  beforePhotoHash: string;
  afterPhotoHash: string;
  technicianNotes: string;
  capturedAt: string;
  synced: boolean;
}

