-- UUID generator
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ========== FACULTY (accounts) ==========
CREATE TABLE IF NOT EXISTS faculty (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email         text NOT NULL UNIQUE,
    pass_hash     text NOT NULL,
    full_name     text NOT NULL,
    institution   text,
    status        text NOT NULL DEFAULT 'PENDING_VERIFICATION'
    CHECK (status IN ('PENDING_VERIFICATION','VERIFIED')),
    verified_at   timestamptz,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now()
    );

-- verification codes for sign-up
CREATE TABLE IF NOT EXISTS faculty_verifications (
    email       text PRIMARY KEY REFERENCES faculty(email) ON DELETE CASCADE,
    code        text NOT NULL,
    expires_at  timestamptz NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now()
    );

-- ========== COURSES ==========
CREATE TABLE IF NOT EXISTS courses (
                                       id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    code               text NOT NULL,          -- e.g. "CS101"
    term               text NOT NULL,          -- e.g. "Fall 2025"
    title              text NOT NULL,
    section            text,                   -- nullable; some schools use it
    start_date         date,
    end_date           date,
    creator_faculty_id uuid NOT NULL REFERENCES faculty(id) ON DELETE RESTRICT,
    created_at         timestamptz NOT NULL DEFAULT now(),
    updated_at         timestamptz NOT NULL DEFAULT now()
    );

-- ========== GRADERS (TAs/graders per course) ==========
DO $$ BEGIN
CREATE TYPE grader_role AS ENUM ('GRADER','TA');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

CREATE TABLE IF NOT EXISTS course_graders (
                                              course_id  uuid NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    faculty_id uuid NOT NULL REFERENCES faculty(id) ON DELETE CASCADE,
    role       grader_role NOT NULL DEFAULT 'GRADER',
    added_at   timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (course_id, faculty_id)
    );

-- ========== STUDENTS + ENROLLMENTS (roster) ==========
CREATE TABLE IF NOT EXISTS students (
                                        id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email      text NOT NULL UNIQUE,
    full_name  text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
    );

CREATE TABLE IF NOT EXISTS enrollments (
    course_id  uuid NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    student_id uuid NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    added_at   timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (course_id, student_id)
    );

-- ========== ROSTER UPLOAD JOBS (for CSV uploads) ==========
CREATE TABLE IF NOT EXISTS roster_import_jobs (
                                                  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id  uuid NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    filename   text NOT NULL,
    status     text NOT NULL DEFAULT 'QUEUED'
    CHECK (status IN ('QUEUED','PROCESSING','DONE','FAILED')),
    has_header boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    error      text
    );

-- Indexes
CREATE INDEX IF NOT EXISTS courses_creator_idx     ON courses(creator_faculty_id);
CREATE INDEX IF NOT EXISTS graders_course_idx      ON course_graders(course_id, faculty_id);
CREATE INDEX IF NOT EXISTS enrollments_course_idx  ON enrollments(course_id);
