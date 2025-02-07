CREATE TABLE admins (
    id UUID PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE users (
    id UUID PRIMARY KEY,
    student_id VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    grade  VARCHAR(100) NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    profile_image VARCHAR(100) NULL,
    status boolean DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE semesters (
    id UUID PRIMARY KEY,
    semester_number VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE courses (
    id UUID PRIMARY KEY,
    course_title VARCHAR(100) NOT NULL,
    semester_id UUID REFERENCES semesters("id"),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE publications (
    id UUID PRIMARY KEY,
    course_id UUID REFERENCES courses("id"),
    title VARCHAR(100) NOT NULL,
    description VARCHAR(100) NULL,
    tags VARCHAR(100) NULL,
    image_id  VARCHAR(100) NOT NULL,
    file_id  VARCHAR(100) NOT NULL,
    contributor_id UUID REFERENCES users("id"),
    status VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE likes (
    id UUID PRIMARY KEY,
    count NUMERIC NOT NULL,
    publication_id UUID REFERENCES publications("id"),
    contributor_id UUID REFERENCES users("id"),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE downloads (
    id UUID PRIMARY KEY,
    count NUMERIC NOT NULL,
    publication_id UUID REFERENCES publications("id"),
    contributor_id UUID REFERENCES users("id"),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);



CREATE TABLE notifications (
    id UUID PRIMARY KEY,
    user_image VARCHAR(100) NOT NULL,
    message VARCHAR(100) NOT NULL,
    file VARCHAR(100) NOT NULL,
    username VARCHAR(100) NOT NULL,
    message_type VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
