package cms

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserCreateRequest struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

type UserUpdateRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

func (s *Store) Authenticate(username, password, ip, userAgent string) (User, error) {
	username = strings.TrimSpace(username)
	var user User
	var passwordHash string
	var lastLogin sql.NullString
	var createdAt, updatedAt string
	err := s.db.QueryRow(`SELECT id,username,nickname,email,password_hash,role,status,last_login_at,created_at,updated_at FROM users WHERE username = ?`, username).
		Scan(&user.ID, &user.Username, &user.Nickname, &user.Email, &passwordHash, &user.Role, &user.Status, &lastLogin, &createdAt, &updatedAt)
	if err != nil {
		_ = s.WriteLoginLog(nil, username, ip, userAgent, false, "用户不存在")
		return User{}, errors.New("用户名或密码错误")
	}
	if user.Status != UserActive {
		_ = s.WriteLoginLog(&user.ID, username, ip, userAgent, false, "账号已禁用")
		return User{}, errors.New("账号已禁用")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		_ = s.WriteLoginLog(&user.ID, username, ip, userAgent, false, "密码错误")
		return User{}, errors.New("用户名或密码错误")
	}
	now := time.Now().UTC()
	_, _ = s.db.Exec(`UPDATE users SET last_login_at=?, updated_at=? WHERE id=?`, formatTime(now), formatTime(now), user.ID)
	user.CreatedAt = parseTime(createdAt)
	user.UpdatedAt = parseTime(updatedAt)
	user.LastLoginAt = &now
	_ = s.WriteLoginLog(&user.ID, username, ip, userAgent, true, "success")
	return user, nil
}

func (s *Store) UserByID(id int64) (User, error) {
	var user User
	var lastLogin sql.NullString
	var createdAt, updatedAt string
	err := s.db.QueryRow(`SELECT id,username,nickname,email,role,status,last_login_at,created_at,updated_at FROM users WHERE id = ?`, id).
		Scan(&user.ID, &user.Username, &user.Nickname, &user.Email, &user.Role, &user.Status, &lastLogin, &createdAt, &updatedAt)
	if err != nil {
		return User{}, err
	}
	if lastLogin.Valid {
		parsed := parseTime(lastLogin.String)
		user.LastLoginAt = &parsed
	}
	user.CreatedAt = parseTime(createdAt)
	user.UpdatedAt = parseTime(updatedAt)
	return user, nil
}

func (s *Store) ListUsers(page, pageSize int) ([]User, int, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	var total int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := s.db.Query(`SELECT id,username,nickname,email,role,status,last_login_at,created_at,updated_at FROM users ORDER BY id ASC LIMIT ? OFFSET ?`, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	users := make([]User, 0)
	for rows.Next() {
		var user User
		var lastLogin sql.NullString
		var createdAt, updatedAt string
		if err := rows.Scan(&user.ID, &user.Username, &user.Nickname, &user.Email, &user.Role, &user.Status, &lastLogin, &createdAt, &updatedAt); err != nil {
			return nil, 0, err
		}
		if lastLogin.Valid {
			parsed := parseTime(lastLogin.String)
			user.LastLoginAt = &parsed
		}
		user.CreatedAt = parseTime(createdAt)
		user.UpdatedAt = parseTime(updatedAt)
		users = append(users, user)
	}
	return users, total, nil
}

func (s *Store) CreateUser(request UserCreateRequest) (User, error) {
	request.Username = strings.TrimSpace(request.Username)
	request.Nickname = strings.TrimSpace(request.Nickname)
	request.Email = strings.TrimSpace(request.Email)
	request.Role = normalizeRole(request.Role)
	request.Status = normalizeStatus(request.Status)
	if request.Username == "" || request.Password == "" {
		return User{}, errors.New("用户名和密码不能为空")
	}
	if request.Nickname == "" {
		request.Nickname = request.Username
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	now := formatTime(time.Now().UTC())
	result, err := s.db.Exec(`INSERT INTO users (username,nickname,email,password_hash,role,status,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?)`,
		request.Username, request.Nickname, request.Email, string(hash), request.Role, request.Status, now, now)
	if err != nil {
		return User{}, err
	}
	id, _ := result.LastInsertId()
	return s.UserByID(id)
}

func (s *Store) UpdateUser(id int64, request UserUpdateRequest) (User, error) {
	request.Role = normalizeRole(request.Role)
	request.Status = normalizeStatus(request.Status)
	_, err := s.db.Exec(`UPDATE users SET nickname=?, email=?, role=?, status=?, updated_at=? WHERE id=?`,
		strings.TrimSpace(request.Nickname), strings.TrimSpace(request.Email), request.Role, request.Status, formatTime(time.Now().UTC()), id)
	if err != nil {
		return User{}, err
	}
	return s.UserByID(id)
}

func (s *Store) DeleteUser(id int64) error {
	var count int
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM users WHERE role = ? AND status = ?`, RoleSuperAdmin, UserActive).Scan(&count)
	user, err := s.UserByID(id)
	if err != nil {
		return err
	}
	if user.Role == RoleSuperAdmin && count <= 1 {
		return errors.New("不能删除最后一个启用的超级管理员")
	}
	_, err = s.db.Exec(`DELETE FROM users WHERE id = ?`, id)
	return err
}

func (s *Store) ResetPassword(id int64, password string) error {
	if strings.TrimSpace(password) == "" {
		return errors.New("密码不能为空")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`UPDATE users SET password_hash=?, updated_at=? WHERE id=?`, string(hash), formatTime(time.Now().UTC()), id)
	return err
}

func (s *Store) SetUserStatus(id int64, status string) (User, error) {
	user, err := s.UserByID(id)
	if err != nil {
		return User{}, err
	}
	status = normalizeStatus(status)
	if user.Role == RoleSuperAdmin && status == UserDisabled {
		var count int
		_ = s.db.QueryRow(`SELECT COUNT(*) FROM users WHERE role = ? AND status = ?`, RoleSuperAdmin, UserActive).Scan(&count)
		if count <= 1 {
			return User{}, errors.New("不能禁用最后一个超级管理员")
		}
	}
	_, err = s.db.Exec(`UPDATE users SET status=?, updated_at=? WHERE id=?`, status, formatTime(time.Now().UTC()), id)
	if err != nil {
		return User{}, err
	}
	return s.UserByID(id)
}

func (s *Store) ListLoginLogs(page, pageSize int) ([]LoginLog, int, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	var total int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM login_logs`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := s.db.Query(`SELECT id,user_id,username,ip,user_agent,success,reason,created_at FROM login_logs ORDER BY id DESC LIMIT ? OFFSET ?`, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	logs := make([]LoginLog, 0)
	for rows.Next() {
		var log LoginLog
		var userID sql.NullInt64
		var success int
		var createdAt string
		if err := rows.Scan(&log.ID, &userID, &log.Username, &log.IP, &log.UserAgent, &success, &log.Reason, &createdAt); err != nil {
			return nil, 0, err
		}
		if userID.Valid {
			log.UserID = &userID.Int64
		}
		log.Success = success == 1
		log.CreatedAt = parseTime(createdAt)
		logs = append(logs, log)
	}
	return logs, total, nil
}

func (s *Store) WriteLoginLog(userID *int64, username, ip, userAgent string, success bool, reason string) error {
	var uid any
	if userID != nil {
		uid = *userID
	}
	_, err := s.db.Exec(`INSERT INTO login_logs (user_id,username,ip,user_agent,success,reason,created_at) VALUES (?,?,?,?,?,?,?)`,
		uid, username, ip, userAgent, boolInt(success), reason, formatTime(time.Now().UTC()))
	return err
}

func normalizeRole(role string) string {
	switch role {
	case RoleSuperAdmin, RoleAdmin, RoleEditor:
		return role
	default:
		return RoleEditor
	}
}

func normalizeStatus(status string) string {
	if status == UserDisabled {
		return UserDisabled
	}
	return UserActive
}
