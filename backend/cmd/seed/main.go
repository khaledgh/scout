package main

import (
	"log"
	"time"

	"kashfi/internal/config"
	"kashfi/internal/db"
	"kashfi/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("db: %v", err)
	}

	if err := db.AutoMigrate(database); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	log.Println("seeding database...")

	hash := func(pw string) string {
		b, _ := bcrypt.GenerateFromPassword([]byte(pw), 10)
		return string(b)
	}
	ptrU := func(u uint) *uint { return &u }
	ptrT := func(t time.Time) *time.Time { return &t }

	loc, _ := time.LoadLocation("Asia/Beirut")
	date := func(y, m, d int) time.Time {
		return time.Date(y, time.Month(m), d, 0, 0, 0, 0, loc)
	}
	dt := func(y, mo, d, h, mi int) time.Time {
		return time.Date(y, time.Month(mo), d, h, mi, 0, 0, loc)
	}

	// ── 1. Users ──────────────────────────────────────────────────────────────
	log.Println("  users...")

	users := []models.User{
		{FullName: "خالد الأمين", Phone: "+96170000001", PasswordHash: hash("Admin@1234"), Role: models.RoleSuperAdmin, IsActive: true},
		{FullName: "سامر حداد", Phone: "+96170000002", PasswordHash: hash("Leader@1234"), Role: models.RoleLeader, IsActive: true},
		{FullName: "ليلى نصر", Phone: "+96170000003", PasswordHash: hash("Leader@1234"), Role: models.RoleLeader, IsActive: true},
		{FullName: "كريم سعد", Phone: "+96170000004", PasswordHash: hash("Assist@1234"), Role: models.RoleAssistant, IsActive: true},
		{FullName: "رنا أبو راشد", Phone: "+96170000005", PasswordHash: hash("Assist@1234"), Role: models.RoleAssistant, IsActive: true},
		{FullName: "علي خوري", Phone: "+96170000006", PasswordHash: hash("Member@1234"), Role: models.RoleMember, IsActive: true},
	}
	for i := range users {
		email := users[i].Phone + "@kashfi.local"
		users[i].Email = &email
		if err := database.Where("phone = ?", users[i].Phone).FirstOrCreate(&users[i]).Error; err != nil {
			log.Fatalf("user: %v", err)
		}
	}

	admin := users[0]
	leaderSamer := users[1]
	leaderLayla := users[2]

	// ── 2. Units ──────────────────────────────────────────────────────────────
	log.Println("  units...")

	units := []models.Unit{
		{Name: "طليعة الأرز", Section: models.SectionKashaf, Motto: "دائماً مستعد", ScoreTotal: 480, Level: 3},
		{Name: "طليعة الشمس", Section: models.SectionKashaf, Motto: "قوة وأمانة", ScoreTotal: 310, Level: 2},
		{Name: "فصيل النسور", Section: models.SectionAshbal, Motto: "نحلق عالياً", ScoreTotal: 220, Level: 2},
		{Name: "فصيل الينابيع", Section: models.SectionAshbal, Motto: "نبع الخير", ScoreTotal: 150, Level: 1},
		{Name: "وحدة الفينيق", Section: models.SectionJawala, Motto: "ننهض من الرماد", ScoreTotal: 560, Level: 4},
	}
	for i := range units {
		if err := database.Where("name = ?", units[i].Name).FirstOrCreate(&units[i]).Error; err != nil {
			log.Fatalf("unit: %v", err)
		}
	}

	unitArz := units[0]
	unitShams := units[1]
	unitNosour := units[2]

	unitLeaders := []models.UnitLeader{
		{UnitID: unitArz.ID, UserID: leaderSamer.ID, RoleInUnit: models.UnitLeaderRoleLeader},
		{UnitID: unitShams.ID, UserID: leaderLayla.ID, RoleInUnit: models.UnitLeaderRoleLeader},
		{UnitID: unitNosour.ID, UserID: users[3].ID, RoleInUnit: models.UnitLeaderRoleAssistant},
	}
	for _, ul := range unitLeaders {
		database.Where("unit_id = ? AND user_id = ?", ul.UnitID, ul.UserID).FirstOrCreate(&ul)
	}

	// ── 3. Members ────────────────────────────────────────────────────────────
	log.Println("  members...")

	type memberSeed struct {
		m      models.Member
		unitID uint
	}

	memberData := []memberSeed{
		// طليعة الأرز — كشاف
		{models.Member{FullName: "يوسف حداد", BirthDate: date(2010, 3, 15), Gender: "male", Section: models.SectionKashaf, JoinDate: date(2020, 9, 1), ParentName: "أنطوان حداد", ParentPhone: "+96171111001", Address: "جبيل، لبنان", XPTotal: 340, Level: 3, Status: models.MemberStatusActive}, unitArz.ID},
		{models.Member{FullName: "لارا معلوف", BirthDate: date(2011, 7, 22), Gender: "female", Section: models.SectionKashaf, JoinDate: date(2021, 9, 1), ParentName: "جورج معلوف", ParentPhone: "+96171111002", Address: "بكرزلا، لبنان", XPTotal: 280, Level: 2, Status: models.MemberStatusActive}, unitArz.ID},
		{models.Member{FullName: "ماريو خوري", BirthDate: date(2010, 11, 8), Gender: "male", Section: models.SectionKashaf, JoinDate: date(2020, 9, 1), ParentName: "سمير خوري", ParentPhone: "+96171111003", Address: "البترون، لبنان", XPTotal: 190, Level: 2, Status: models.MemberStatusActive}, unitArz.ID},
		{models.Member{FullName: "سارة زيادة", BirthDate: date(2012, 4, 3), Gender: "female", Section: models.SectionKashaf, JoinDate: date(2022, 9, 1), ParentName: "بيار زيادة", ParentPhone: "+96171111004", Address: "أنطلياس، لبنان", XPTotal: 90, Level: 1, Status: models.MemberStatusActive}, unitArz.ID},
		{models.Member{FullName: "طوني أبو رزق", BirthDate: date(2011, 1, 17), Gender: "male", Section: models.SectionKashaf, JoinDate: date(2021, 9, 1), ParentName: "نجيب أبو رزق", ParentPhone: "+96171111005", Address: "كسروان، لبنان", XPTotal: 120, Level: 1, Status: models.MemberStatusActive}, unitArz.ID},
		// طليعة الشمس — كشاف
		{models.Member{FullName: "نادين سلامة", BirthDate: date(2010, 6, 10), Gender: "female", Section: models.SectionKashaf, JoinDate: date(2020, 9, 1), ParentName: "إيلي سلامة", ParentPhone: "+96172222001", Address: "جونية، لبنان", XPTotal: 250, Level: 2, Status: models.MemberStatusActive}, unitShams.ID},
		{models.Member{FullName: "رامي فرحات", BirthDate: date(2011, 9, 25), Gender: "male", Section: models.SectionKashaf, JoinDate: date(2021, 9, 1), ParentName: "مارون فرحات", ParentPhone: "+96172222002", Address: "المتن، لبنان", XPTotal: 180, Level: 2, Status: models.MemberStatusActive}, unitShams.ID},
		{models.Member{FullName: "جوليا نصر", BirthDate: date(2012, 2, 14), Gender: "female", Section: models.SectionKashaf, JoinDate: date(2022, 9, 1), ParentName: "وليم نصر", ParentPhone: "+96172222003", Address: "ضبيه، لبنان", XPTotal: 70, Level: 1, Status: models.MemberStatusActive}, unitShams.ID},
		{models.Member{FullName: "شربل عيد", BirthDate: date(2011, 12, 5), Gender: "male", Section: models.SectionKashaf, JoinDate: date(2021, 9, 1), ParentName: "خليل عيد", ParentPhone: "+96172222004", Address: "بيروت، لبنان", XPTotal: 210, Level: 2, Status: models.MemberStatusActive}, unitShams.ID},
		{models.Member{FullName: "ميشال رزق", BirthDate: date(2013, 3, 30), Gender: "male", Section: models.SectionKashaf, JoinDate: date(2023, 9, 1), ParentName: "نبيل رزق", ParentPhone: "+96172222005", Address: "زغرتا، لبنان", XPTotal: 30, Level: 1, Status: models.MemberStatusInactive}, unitShams.ID},
		// فصيل النسور — أشبال
		{models.Member{FullName: "كارلا جبور", BirthDate: date(2014, 5, 20), Gender: "female", Section: models.SectionAshbal, JoinDate: date(2021, 9, 1), ParentName: "فادي جبور", ParentPhone: "+96173333001", Address: "جل الديب، لبنان", XPTotal: 150, Level: 2, Status: models.MemberStatusActive}, unitNosour.ID},
		{models.Member{FullName: "بولس ضومط", BirthDate: date(2015, 8, 11), Gender: "male", Section: models.SectionAshbal, JoinDate: date(2022, 9, 1), ParentName: "جوزيف ضومط", ParentPhone: "+96173333002", Address: "حريصا، لبنان", XPTotal: 80, Level: 1, Status: models.MemberStatusActive}, unitNosour.ID},
		{models.Member{FullName: "ريم بارود", BirthDate: date(2014, 11, 7), Gender: "female", Section: models.SectionAshbal, JoinDate: date(2021, 9, 1), ParentName: "حسان بارود", ParentPhone: "+96173333003", Address: "بيروت، لبنان", XPTotal: 110, Level: 1, Status: models.MemberStatusActive}, unitNosour.ID},
	}

	createdMembers := make([]models.Member, 0, len(memberData))
	for _, ms := range memberData {
		m := ms.m
		if err := database.Where("full_name = ? AND parent_phone = ?", m.FullName, m.ParentPhone).FirstOrCreate(&m).Error; err != nil {
			log.Fatalf("member: %v", err)
		}
		um := models.UnitMember{UnitID: ms.unitID, MemberID: m.ID, IsPrimary: true, JoinedAt: m.JoinDate}
		database.Where("unit_id = ? AND member_id = ?", ms.unitID, m.ID).FirstOrCreate(&um)
		createdMembers = append(createdMembers, m)
	}

	// ── 4. Badges catalog ─────────────────────────────────────────────────────
	log.Println("  badges...")

	badges := []models.Badge{
		{Name: "شارة الإسعافات الأولية", Description: "إتقان مهارات الإسعاف الأولي الأساسية", Category: "صحة وسلامة", XPReward: 50, IsActive: true},
		{Name: "شارة الطبيعة والبيئة", Description: "المعرفة بالنباتات والحياة البرية وحماية البيئة", Category: "طبيعة", XPReward: 40, IsActive: true},
		{Name: "شارة الملاحة", Description: "استخدام البوصلة والخرائط والتوجيه الفلكي", Category: "مهارات كشفية", XPReward: 60, IsActive: true},
		{Name: "شارة القيادة", Description: "مهارات تنظيم المجموعات والقيادة الفعّالة", Category: "مهارات قيادية", XPReward: 80, IsActive: true},
		{Name: "شارة الطبخ في الخلاء", Description: "إعداد الطعام بالطرق الكشفية التقليدية", Category: "مهارات كشفية", XPReward: 35, IsActive: true},
		{Name: "شارة المواطنة", Description: "المشاركة الفاعلة في خدمة المجتمع", Category: "خدمة اجتماعية", XPReward: 45, IsActive: true},
		{Name: "شارة الرياضة", Description: "التفوق في النشاطات الرياضية والبدنية", Category: "رياضة", XPReward: 30, IsActive: true},
		{Name: "شارة الإبداع", Description: "الأعمال الفنية والحرفية والإبداعية", Category: "فنون", XPReward: 25, IsActive: true},
	}
	createdBadges := make([]models.Badge, 0, len(badges))
	for _, b := range badges {
		if err := database.Where("name = ?", b.Name).FirstOrCreate(&b).Error; err != nil {
			log.Fatalf("badge: %v", err)
		}
		createdBadges = append(createdBadges, b)
	}

	// ── 5. Skills catalog ─────────────────────────────────────────────────────
	log.Println("  skills...")

	skills := []models.Skill{
		{Name: "عقد الحبال", Category: "مهارات كشفية", Description: "إتقان أنواع مختلفة من العقد الكشفية", MaxLevel: 5},
		{Name: "الإسعافات الأولية", Category: "صحة", Description: "تقديم الإسعافات الأولية في حالات الطوارئ", MaxLevel: 5},
		{Name: "قراءة الخرائط", Category: "ملاحة", Description: "قراءة الخرائط الطوبوغرافية والتوجيه الجغرافي", MaxLevel: 5},
		{Name: "الطبخ في الطبيعة", Category: "مهارات حياة", Description: "إعداد الطعام في البيئة الطبيعية", MaxLevel: 5},
		{Name: "السباحة", Category: "رياضة", Description: "مهارات السباحة ومستوياتها", MaxLevel: 5},
		{Name: "التسلق", Category: "رياضة", Description: "تسلق الجبال والصخور بأمان", MaxLevel: 5},
		{Name: "الخطابة والتعبير", Category: "تواصل", Description: "مهارات التحدث أمام الجمهور", MaxLevel: 5},
	}
	for _, s := range skills {
		database.Where("name = ?", s.Name).FirstOrCreate(&s)
	}

	// ── 6. Activities ─────────────────────────────────────────────────────────
	log.Println("  activities...")

	activities := []models.Activity{
		{Title: "مخيم جبل لبنان الصيفي", Description: "مخيم صيفي في منطقة الشوف يشمل التسلق والمسير والأنشطة البيئية", Type: models.ActivityTypeCamp, Location: "منطقة الشوف، لبنان", StartsAt: dt(2025, 7, 15, 8, 0), EndsAt: dt(2025, 7, 20, 17, 0), ResponsibleUserID: leaderSamer.ID, Status: models.ActivityStatusCompleted},
		{Title: "مسير في أرز لبنان", Description: "مسير تحت ضوء القمر في محيط أرز الرب المقدسة", Type: models.ActivityTypeHike, Location: "أرز لبنان، الشمال", StartsAt: dt(2025, 10, 5, 7, 0), EndsAt: dt(2025, 10, 5, 18, 0), ResponsibleUserID: leaderLayla.ID, Status: models.ActivityStatusCompleted},
		{Title: "اجتماع شهري — أكتوبر", Description: "الاجتماع الشهري لمراجعة النشاطات وتوزيع المهام", Type: models.ActivityTypeMeeting, Location: "مقر الفوج", StartsAt: dt(2025, 10, 18, 10, 0), EndsAt: dt(2025, 10, 18, 12, 0), ResponsibleUserID: admin.ID, Status: models.ActivityStatusCompleted},
		{Title: "دورة إسعافات أولية", Description: "دورة تدريبية متكاملة في الإسعافات الأولية مع الصليب الأحمر اللبناني", Type: models.ActivityTypeTraining, Location: "مركز الصليب الأحمر، جبيل", StartsAt: dt(2025, 11, 8, 9, 0), EndsAt: dt(2025, 11, 8, 16, 0), ResponsibleUserID: leaderSamer.ID, Status: models.ActivityStatusCompleted},
		{Title: "يوم الخدمة — تنظيف الشاطئ", Description: "حملة تنظيف شاطئ بيروت مع منظمات البيئة اللبنانية", Type: models.ActivityTypeService, Location: "شاطئ رملة البيضا، بيروت", StartsAt: dt(2025, 12, 6, 8, 0), EndsAt: dt(2025, 12, 6, 13, 0), ResponsibleUserID: leaderLayla.ID, Status: models.ActivityStatusCompleted},
		{Title: "اجتماع شهري — يناير 2026", Description: "تخطيط الأنشطة للربع الأول من 2026", Type: models.ActivityTypeMeeting, Location: "مقر الفوج", StartsAt: dt(2026, 1, 18, 10, 0), EndsAt: dt(2026, 1, 18, 12, 0), ResponsibleUserID: admin.ID, Status: models.ActivityStatusCompleted},
		{Title: "مسير ربيعي في قاديشا", Description: "مسير في وادي قاديشا المقدس مع دراسة جيولوجية وتاريخية", Type: models.ActivityTypeHike, Location: "وادي قاديشا، شمال لبنان", StartsAt: dt(2026, 3, 21, 7, 30), EndsAt: dt(2026, 3, 21, 17, 0), ResponsibleUserID: leaderSamer.ID, Status: models.ActivityStatusCompleted},
		{Title: "مخيم نهاية السنة", Description: "مخيم ختامي للسنة الكشفية يتضمن توزيع الشارات والتكريم", Type: models.ActivityTypeCamp, Location: "المنتزه الوطني، بيروت", StartsAt: dt(2026, 7, 10, 8, 0), EndsAt: dt(2026, 7, 14, 17, 0), ResponsibleUserID: admin.ID, Status: models.ActivityStatusPlanned},
	}
	createdActivities := make([]models.Activity, 0, len(activities))
	for _, a := range activities {
		if err := database.Where("title = ? AND starts_at = ?", a.Title, a.StartsAt).FirstOrCreate(&a).Error; err != nil {
			log.Fatalf("activity: %v", err)
		}
		createdActivities = append(createdActivities, a)
	}

	// ── 7. Attendance ─────────────────────────────────────────────────────────
	log.Println("  attendance...")

	manual := models.CheckInMethod("manual")
	statuses := []models.AttendanceStatus{
		models.AttendancePresent, models.AttendancePresent, models.AttendancePresent,
		models.AttendancePresent, models.AttendanceLate, models.AttendancePresent,
		models.AttendancePresent, models.AttendanceAbsent, models.AttendancePresent,
		models.AttendancePresent, models.AttendancePresent, models.AttendanceExcused,
		models.AttendancePresent,
	}
	for i, act := range createdActivities[:7] {
		for j, m := range createdMembers {
			st := statuses[(i+j)%len(statuses)]
			att := models.ActivityAttendance{
				ActivityID:    act.ID,
				MemberID:      m.ID,
				Status:        st,
				CheckInMethod: &manual,
				RecordedBy:    admin.ID,
			}
			if st == models.AttendancePresent || st == models.AttendanceLate {
				t := act.StartsAt.Add(10 * time.Minute)
				att.CheckInAt = &t
			}
			database.Where("activity_id = ? AND member_id = ?", act.ID, m.ID).FirstOrCreate(&att)
		}
	}

	// ── 8. Badge awards ───────────────────────────────────────────────────────
	log.Println("  badge awards...")

	badgeAwards := []struct {
		mi, bi int
		at     time.Time
	}{
		{0, 0, date(2025, 8, 1)},   // يوسف — إسعافات
		{0, 2, date(2025, 11, 10)}, // يوسف — ملاحة
		{0, 3, date(2026, 1, 20)},  // يوسف — قيادة
		{1, 0, date(2025, 9, 15)},  // لارا — إسعافات
		{1, 1, date(2025, 10, 20)}, // لارا — طبيعة
		{2, 4, date(2025, 12, 1)},  // ماريو — طبخ
		{5, 5, date(2025, 10, 15)}, // نادين — مواطنة
		{7, 6, date(2026, 2, 1)},   // شربل — رياضة
		{10, 1, date(2025, 11, 5)}, // كارلا — طبيعة
	}
	for _, ba := range badgeAwards {
		mb := models.MemberBadge{
			MemberID:  createdMembers[ba.mi].ID,
			BadgeID:   createdBadges[ba.bi].ID,
			AwardedAt: ba.at,
			AwardedBy: ptrU(admin.ID),
		}
		database.Where("member_id = ? AND badge_id = ?", mb.MemberID, mb.BadgeID).FirstOrCreate(&mb)
	}

	// ── 9. XP events ──────────────────────────────────────────────────────────
	log.Println("  xp events...")

	xpEvents := []struct {
		mi     int
		source models.XPSource
		points int
		refID  uint
		note   string
	}{
		{0, models.XPSourceAttendance, 10, createdActivities[0].ID, "حضور مخيم الشوف"},
		{0, models.XPSourceBadge, 50, createdBadges[0].ID, "شارة الإسعافات الأولية"},
		{0, models.XPSourceAttendance, 10, createdActivities[1].ID, "حضور مسير الأرز"},
		{0, models.XPSourceBadge, 60, createdBadges[2].ID, "شارة الملاحة"},
		{0, models.XPSourceBadge, 80, createdBadges[3].ID, "شارة القيادة"},
		{1, models.XPSourceAttendance, 10, createdActivities[0].ID, "حضور المخيم"},
		{1, models.XPSourceBadge, 50, createdBadges[0].ID, "شارة الإسعافات"},
		{1, models.XPSourceBadge, 40, createdBadges[1].ID, "شارة الطبيعة"},
		{2, models.XPSourceAttendance, 10, createdActivities[1].ID, "مسير الأرز"},
		{2, models.XPSourceBadge, 35, createdBadges[4].ID, "شارة الطبخ"},
		{5, models.XPSourceBadge, 45, createdBadges[5].ID, "شارة المواطنة"},
		{7, models.XPSourceBadge, 30, createdBadges[6].ID, "شارة الرياضة"},
	}
	for _, xe := range xpEvents {
		rid := xe.refID
		ev := models.XPEvent{
			MemberID: createdMembers[xe.mi].ID,
			Source:   xe.source,
			Points:   xe.points,
			RefID:    &rid,
			Note:     xe.note,
		}
		database.Where("member_id = ? AND source = ? AND ref_id = ?", ev.MemberID, ev.Source, ev.RefID).FirstOrCreate(&ev)
	}

	// ── 10. Announcements ─────────────────────────────────────────────────────
	log.Println("  announcements...")

	t1 := dt(2025, 9, 1, 10, 0)
	t2 := dt(2025, 12, 1, 9, 0)
	t3 := dt(2026, 6, 1, 11, 0)

	for _, ann := range []models.Announcement{
		{Title: "انطلاق السنة الكشفية 2025-2026", Body: "يسرنا الإعلان عن انطلاقة السنة الكشفية الجديدة. ندعو جميع الأعضاء للمشاركة في أنشطة هذا العام الحافل.", Audience: "all", Pinned: true, AuthorID: admin.ID, PublishedAt: &t1},
		{Title: "مخيم الشوف — اللوازم المطلوبة", Body: "يرجى تجهيز الحقيبة وفق القائمة المرفقة قبل موعد المخيم بأسبوع على الأقل.", Audience: "all", AuthorID: leaderSamer.ID, PublishedAt: &t2},
		{Title: "تذكير: الاجتماع الشهري يناير", Body: "نذكّر جميع القادة والمساعدين بحضور الاجتماع الشهري في مقر الفوج.", Audience: "leaders", AuthorID: admin.ID, PublishedAt: &t3},
	} {
		database.Where("title = ?", ann.Title).FirstOrCreate(&ann)
	}

	// ── 11. Training lessons ──────────────────────────────────────────────────
	log.Println("  training lessons...")

	lessons := []models.TrainingLesson{
		{Title: "مقدمة في الإسعافات الأولية", Category: "صحة وسلامة", Content: "يتناول هذا الدرس أساسيات الإسعافات الأولية بما يشمل تقييم الموقف، استدعاء المساعدة، والإسعافات الأولية للجروح والحروق.", OrderIndex: 1, IsPublished: true},
		{Title: "عقد الحبال الكشفية", Category: "مهارات كشفية", Content: "شرح مفصّل لأنواع العقد الكشفية: عقدة الدوران، العقدة المسطحة، عقدة المرساة، وكيفية توظيفها في المخيمات.", OrderIndex: 1, IsPublished: true},
		{Title: "قراءة الخريطة والبوصلة", Category: "ملاحة", Content: "تعلّم كيفية قراءة خرائط الكنتور، تحديد الاتجاهات بالبوصلة، ورسم مسارات المشي.", OrderIndex: 2, IsPublished: true},
		{Title: "حماية البيئة والطبيعة", Category: "بيئة", Content: "مبادئ اللا-أثر في الطبيعة، التعرف على أنواع النباتات اللبنانية، وكيفية التصرف السليم في البيئة.", OrderIndex: 1, IsPublished: true},
		{Title: "مهارات القيادة والعمل الجماعي", Category: "مهارات قيادية", Content: "أنواع القيادة، بناء فريق متماسك، وحل النزاعات داخل المجموعة.", OrderIndex: 2, IsPublished: true},
	}
	createdLessons := make([]models.TrainingLesson, 0, len(lessons))
	for _, l := range lessons {
		if err := database.Where("title = ?", l.Title).FirstOrCreate(&l).Error; err != nil {
			log.Fatalf("lesson: %v", err)
		}
		createdLessons = append(createdLessons, l)
	}

	// ── 12. Quizzes ───────────────────────────────────────────────────────────
	log.Println("  quizzes...")

	quizQuestions := [3][]models.QuizQuestion{
		{
			{Text: "ما هي الخطوة الأولى عند مواجهة حادثة؟", OptionsJSON: `["تقييم الموقف","الاتصال بالإسعاف فوراً","البدء بالإسعاف مباشرة","الانتظار"]`, CorrectIndex: 0, Points: 1},
			{Text: "كم عدد ضغطات الصدر في دورة الإنعاش؟", OptionsJSON: `["10","20","30","40"]`, CorrectIndex: 2, Points: 1},
			{Text: "ماذا تفعل عند تعرض شخص لحرق من الدرجة الأولى؟", OptionsJSON: `["وضع الزيت عليه","تبريده بالماء البارد عشر دقائق","تغطيته بالقماش","تركه وعدم لمسه"]`, CorrectIndex: 1, Points: 1},
		},
		{
			{Text: "ما هي العقدة المستخدمة لربط حبلين معاً؟", OptionsJSON: `["العقدة المسطحة","عقدة الدوران","عقدة الثماني","العقدة البحرية"]`, CorrectIndex: 0, Points: 1},
			{Text: "كم عقدة أساسية يجب على الكشاف إتقانها؟", OptionsJSON: `["3","5","7","10"]`, CorrectIndex: 2, Points: 1},
			{Text: "ما استخدام عقدة المرساة؟", OptionsJSON: `["ربط حبل بعمود","اختصار الحبل","إنقاذ الغارقين","تثبيت الخيمة"]`, CorrectIndex: 0, Points: 1},
		},
		{
			{Text: "ما هو الشمال المغناطيسي؟", OptionsJSON: `["الشمال الجغرافي الحقيقي","اتجاه إبرة البوصلة","اتجاه شروق الشمس","اتجاه القطب الجنوبي"]`, CorrectIndex: 1, Points: 1},
			{Text: "ماذا تعني خطوط الكنتور المتقاربة؟", OptionsJSON: `["أرض مسطحة","منحدر حاد","وادٍ","قمة جبل"]`, CorrectIndex: 1, Points: 1},
			{Text: "كيف تحدد الجنوب بالشمس نهاراً؟", OptionsJSON: `["مغرب الشمس","ظل قصير منتصف النهار","مشرق الشمس","عكس القمر"]`, CorrectIndex: 1, Points: 1},
		},
	}

	for i, lesson := range createdLessons[:3] {
		quiz := models.Quiz{
			LessonID:  lesson.ID,
			Title:     "اختبار: " + lesson.Title,
			PassScore: 70,
			XPReward:  20,
		}
		database.Where("lesson_id = ?", lesson.ID).FirstOrCreate(&quiz)
		for _, q := range quizQuestions[i] {
			q.QuizID = quiz.ID
			database.Where("quiz_id = ? AND text = ?", q.QuizID, q.Text).FirstOrCreate(&q)
		}
		if i == 0 && len(createdMembers) > 0 {
			att := models.QuizAttempt{
				QuizID:      quiz.ID,
				MemberID:    createdMembers[0].ID,
				Score:       100,
				Passed:      true,
				AttemptedAt: dt(2025, 11, 15, 14, 30),
			}
			database.Where("quiz_id = ? AND member_id = ?", quiz.ID, createdMembers[0].ID).FirstOrCreate(&att)
		}
	}

	// ── 13. Evaluations ───────────────────────────────────────────────────────
	log.Println("  evaluations...")

	type evalSeed struct {
		mi                             int
		period                         string
		disc, par, lead, skill, overall int
		notes                          string
	}
	for _, ev := range []evalSeed{
		{0, "2025-Q1", 9, 10, 9, 8, 9, "يوسف طالب متميز، يُظهر روح قيادية ممتازة"},
		{0, "2025-Q2", 10, 9, 9, 9, 9, "تحسّن ملحوظ في المهارات العملية"},
		{1, "2025-Q1", 8, 9, 7, 8, 8, "لارا نشطة ومتحمسة، بحاجة لتطوير مهارات الملاحة"},
		{2, "2025-Q1", 7, 8, 6, 7, 7, "ماريو بحاجة لمزيد من الحضور والمشاركة"},
		{5, "2025-Q1", 8, 8, 8, 9, 8, "نادين متفوقة في المهارات"},
	} {
		e := models.Evaluation{
			MemberID: createdMembers[ev.mi].ID, EvaluatorID: leaderSamer.ID,
			Period: ev.period, Discipline: ev.disc, Participation: ev.par,
			Leadership: ev.lead, Skill: ev.skill, Overall: ev.overall, Notes: ev.notes,
		}
		database.Where("member_id = ? AND period = ?", e.MemberID, e.Period).FirstOrCreate(&e)
	}

	// ── 14. Channels & messages ───────────────────────────────────────────────
	log.Println("  channels & messages...")

	channels := []models.Channel{
		{Name: "عام — الفوج", Type: models.ChannelGroup, UnitID: nil},
		{Name: "طليعة الأرز", Type: models.ChannelUnit, UnitID: &unitArz.ID},
		{Name: "طليعة الشمس", Type: models.ChannelUnit, UnitID: &unitShams.ID},
		{Name: "القادة", Type: models.ChannelGroup, UnitID: nil},
	}
	createdChannels := make([]models.Channel, 0, len(channels))
	for _, ch := range channels {
		database.Where("name = ?", ch.Name).FirstOrCreate(&ch)
		createdChannels = append(createdChannels, ch)
	}

	now := time.Now()
	for i, msg := range []models.Message{
		{ChannelID: createdChannels[0].ID, SenderID: admin.ID, Body: "أهلاً بكم في قناة الفوج العامة! 🌲"},
		{ChannelID: createdChannels[0].ID, SenderID: leaderSamer.ID, Body: "مرحباً بالجميع. لا تنسوا موعد الاجتماع يوم السبت."},
		{ChannelID: createdChannels[1].ID, SenderID: leaderSamer.ID, Body: "طليعة الأرز — جاهزون لمخيم الشوف؟ 💪"},
		{ChannelID: createdChannels[3].ID, SenderID: admin.ID, Body: "اجتماع طارئ للقادة الثلاثاء الساعة 7 مساءً."},
	} {
		msg.CreatedAt = now.Add(-time.Duration(4-i) * 3 * time.Hour)
		database.Where("channel_id = ? AND sender_id = ? AND body = ?", msg.ChannelID, msg.SenderID, msg.Body).FirstOrCreate(&msg)
	}

	// ── 15. Notifications ─────────────────────────────────────────────────────
	log.Println("  notifications...")

	readAt := ptrT(dt(2025, 12, 2, 9, 0))
	for _, uid := range []uint{admin.ID, leaderSamer.ID, leaderLayla.ID} {
		for _, n := range []models.Notification{
			{UserID: uid, Title: "مخيم جبل لبنان", Body: "تذكير: المخيم الصيفي يبدأ خلال 24 ساعة!", Type: "reminder", ReadAt: readAt},
			{UserID: uid, Title: "شارة جديدة", Body: "تهانينا! حصلت طليعتك على شارة المواطنة.", Type: "badge"},
		} {
			database.Where("user_id = ? AND title = ?", n.UserID, n.Title).FirstOrCreate(&n)
		}
	}

	log.Println("✓ seed complete!")
	log.Printf("  users: %d | units: %d | members: %d | badges: %d | activities: %d | lessons: %d",
		len(users), len(units), len(createdMembers), len(createdBadges), len(createdActivities), len(createdLessons))

	_ = ptrT
}
