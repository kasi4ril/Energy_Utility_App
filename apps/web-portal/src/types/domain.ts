/**
 * Shared domain entity interfaces for the NESCO Enterprise 4-Zone UI Shell.
 * Reflects models from S01 through S09 bounded contexts.
 */

export type DepartmentSlug =
  | "executive-command"
  | "customer-care"
  | "services-feeder"
  | "services-transformer"
  | "services-metering"
  | "generation-master"
  | "energy-management"
  | "stores"
  | "internal-audit"
  | "billing-commercial";

export interface MetricCard {
  id: string;
  title: string;
  value: string | number;
  unit?: string;
  changePct?: number;
  trend?: "UP" | "DOWN" | "NEUTRAL";
  status?: "NORMAL" | "WARNING" | "CRITICAL";
}

export interface ActionQueueItem {
  id: string;
  referenceNo: string;
  title: string;
  category: string;
  priority: "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";
  status: string;
  createdAt: string;
  sourceDepartment: string;
  slaBreach?: boolean;
}

export interface WorkNote {
  id: string;
  entityType: string;
  entityId: string;
  noteText: string;
  authorStaffId: string;
  authorRole: string;
  noteSeq: number;
  prevHash: string;
  currentHash: string;
  createdAt: string;
}

export interface GridAssistantMessage {
  id: string;
  sender: "USER" | "ASSISTANT" | "SYSTEM";
  senderName?: string;
  content: string;
  timestamp: string;
  metadata?: Record<string, unknown>;
}

