import { PartialCard } from "./Card";
import { Tag } from "./Tags";

export interface Task {
  id: number;
  card_pk: number;
  user_id: number;
  scheduled_date: Date | null;
  dueDate: Date | null;
  created_at: Date;
  updated_at: Date;
  completed_at: Date | null;
  title: string;
  is_complete: boolean;
  is_deleted: boolean;
  card: PartialCard | null;
  tags: Tag[];
}

export const emptyTask: Task = {
  id: 0,
  card_pk: 0,
  user_id: 0,
  created_at: new Date(0),
  updated_at: new Date(0),
  dueDate: null,
  scheduled_date: new Date(),
  completed_at: null,
  title: "",
  is_complete: false,
  is_deleted: false,
  card: null,
  tags: [],
};
