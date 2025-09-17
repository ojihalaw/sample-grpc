
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    sku VARCHAR(50) UNIQUE NOT NULL,
    variant VARCHAR(20) NOT NULL,
    price INT NOT NULL DEFAULT 0,
    stock INT NOT NULL DEFAULT 0,
    description TEXT,
    star DOUBLE PRECISION DEFAULT 5.0,
    image_url VARCHAR(255) NOT NULL,
    category_id UUID NOT NULL,
    special_type VARCHAR(50),
    is_special BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);
