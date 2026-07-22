/* you can add more column if you want, for ex : bio, profile picture */
CREATE TABLE IF NOT EXISTS users (
    ID int PRIMARY KEY,
    Username varchar(50) UNIQUE NOT NULL,
    Email varchar(150) UNIQUE NOT NULL,
    PhoneNumber varchar(30) UNIQUE,
    Pass text NOT NULL,
    EmailVerificationCode varchar(255) UNIQUE NOT NULL,
    EmailVerificationCodeExpiryDate date NOT NULL,
    TwoFactorVerificationCode varchar(255) UNIQUE
)