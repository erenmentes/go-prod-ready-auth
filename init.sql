/* you can add more column if you want, for ex : bio, profile picture */
CREATE TABLE IF NOT EXISTS users (
    ID int PRIMARY KEY,
    Username varchar(50) UNIQUE NOT NULL,
    Email varchar(150) UNIQUE NOT NULL,
    UserRole varchar(25) NOT NULL DEFAULT 'User',
    PhoneNumber varchar(30) UNIQUE,
    Pass text NOT NULL,
    RefreshToken varchar(255) UNIQUE NOT NULL,
    IsEmailVerified boolean NOT NULL DEFAULT false,
    IsTwoFactorVerificationActivated boolean NOT NULL DEFAULT false,
    EmailVerificationCode varchar(255) UNIQUE NOT NULL,
    EmailVerificationCodeExpiryDate timestamptz NOT NULL,
    TwoFactorVerificationCode varchar(255) UNIQUE
)