export type Role = 'super_admin' | 'leader' | 'assistant' | 'member' | 'parent'
export type Section = 'ashbal' | 'kashaf' | 'jawala' | 'mukashe'
export type MemberStatus = 'active' | 'inactive'
export type ActivityType = 'camp' | 'hike' | 'training' | 'meeting' | 'service'
export type ActivityStatus = 'planned' | 'ongoing' | 'completed' | 'cancelled'
export type AttendanceStatus = 'present' | 'absent' | 'excused' | 'late'
export type CheckInMethod = 'qr' | 'gps' | 'manual'
export type XPSource = 'attendance' | 'badge' | 'quiz' | 'leadership' | 'manual'
export type MediaType = 'image' | 'video'

export interface ApiResponse<T> {
  success: boolean
  data: T
  error?: { code: string; message: string }
  meta?: PaginationMeta
}

export interface PaginationMeta {
  page: number
  page_size: number
  total: number
  total_pages: number
}

export interface User {
  id: number
  full_name: string
  email?: string
  phone: string
  role: Role
  avatar_url?: string
  is_active: boolean
  last_login_at?: string
  created_at: string
}

export interface Member {
  id: number
  user_id?: number
  full_name: string
  birth_date: string
  gender: 'male' | 'female'
  section: Section
  rank_stage: string
  join_date: string
  photo_url?: string
  parent_name: string
  parent_phone: string
  secondary_phone: string
  address: string
  xp_total: number
  level: number
  status: MemberStatus
  created_at: string
  badges?: MemberBadge[]
  skills?: MemberSkill[]
  units?: UnitMember[]
}

export interface MemberMedical {
  member_id: number
  blood_type: string
  allergies: string
  chronic_conditions: string
  medications: string
  emergency_notes: string
}

export interface Unit {
  id: number
  name: string
  section: Section
  motto: string
  emblem_url?: string
  level: number
  score_total: number
  is_active: boolean
  leaders?: UnitLeader[]
  members?: UnitMember[]
}

export interface UnitLeader {
  id: number
  unit_id: number
  user_id: number
  role_in_unit: 'leader' | 'assistant'
  user: User
}

export interface UnitMember {
  id: number
  unit_id: number
  member_id: number
  is_primary: boolean
  joined_at: string
  left_at?: string
  member: Member
  unit?: Unit
}

export interface ActivityMedia {
  id: number
  created_at: string
  activity_id: number
  url: string
  media_type: MediaType
  uploaded_by: number
  caption: string
}

export interface Activity {
  id: number
  title: string
  description: string
  type: ActivityType
  location: string
  location_lat?: number
  location_lng?: number
  starts_at: string
  ends_at: string
  responsible_user_id: number
  unit_id?: number
  status: ActivityStatus
  cover_image_url?: string
  responsible_user?: User
  unit?: Unit
  media?: ActivityMedia[]
}

export interface ActivityAttendance {
  id: number
  activity_id: number
  member_id: number
  status: AttendanceStatus
  check_in_at?: string
  check_in_method?: CheckInMethod
  lat?: number
  lng?: number
  recorded_by: number
  member?: Member
}

export interface Badge {
  id: number
  name: string
  description: string
  category: string
  icon_url?: string
  xp_reward: number
  criteria_json?: string
  is_active: boolean
}

export interface MemberBadge {
  id: number
  member_id: number
  badge_id: number
  awarded_at: string
  awarded_by?: number
  progress: number
  badge?: Badge
  member?: Member
}

export interface Skill {
  id: number
  name: string
  category: string
  description: string
  max_level: number
}

export interface MemberSkill {
  id: number
  member_id: number
  skill_id: number
  level: number
  assessed_by?: number
  assessed_at?: string
  skill?: Skill
}

export interface TrainingLessonMedia {
  id: number
  created_at: string
  lesson_id: number
  url: string
  media_type: MediaType
  uploaded_by: number
  caption: string
}

export interface TrainingLesson {
  id: number
  title: string
  category: string
  content: string
  cover_url?: string
  media_json?: string
  order_index: number
  is_published: boolean
  quizzes?: Quiz[]
  media?: TrainingLessonMedia[]
}

export interface Quiz {
  id: number
  lesson_id: number
  title: string
  pass_score: number
  xp_reward: number
  questions?: QuizQuestion[]
}

export interface QuizQuestion {
  id: number
  quiz_id: number
  text: string
  options_json: string
  correct_index: number
  points: number
}

export interface XPEvent {
  id: number
  member_id: number
  unit_id?: number
  source: XPSource
  points: number
  ref_id?: number
  note: string
  created_at: string
}

export interface Announcement {
  id: number
  title: string
  body: string
  audience: 'all' | 'unit' | 'leaders'
  unit_id?: number
  author_id: number
  pinned: boolean
  published_at?: string
  author?: User
}

export interface Channel {
  id: number
  type: 'unit' | 'group' | 'direct'
  unit_id?: number
  name: string
}

export interface Message {
  id: number
  channel_id: number
  sender_id: number
  body: string
  attachment_url?: string
  created_at: string
  sender?: User
}

export interface Notification {
  id: number
  user_id: number
  title: string
  body: string
  type: string
  data_json: string
  read_at?: string
  created_at: string
}

export interface Equipment {
  id: number
  name: string
  category: string
  quantity_total: number
  quantity_available: number
  condition: string
  notes: string
}

export interface Evaluation {
  id: number
  member_id: number
  evaluator_id: number
  period: string
  discipline: number
  participation: number
  leadership: number
  skill: number
  overall: number
  notes: string
  created_at: string
  evaluator?: User
}

export interface SectionCount {
  section: Section
  count: number
}

export interface DashboardKPIs {
  member_count: number
  active_members: number
  attendance_rate: number
  top_unit?: Unit
  upcoming_activities: Activity[]
  at_risk_members: Member[]
  recent_activities: Activity[]
  members_by_section: SectionCount[]
  top_members: Member[]
  recent_badges: MemberBadge[]
}
