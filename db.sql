DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'message_status') THEN
        CREATE TYPE message_status AS ENUM ('pending', 'failed', 'sent');
    END IF;
END$$;

CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipient_phone VARCHAR(20) NOT NULL,
    content VARCHAR(255) NOT NULL,
    status message_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    sent_at TIMESTAMP
);

INSERT INTO messages (recipient_phone, content, status)
VALUES
('+905301234001', 'Message content 1', 'pending'),
('+905301234002', 'Message content 2', 'pending'),
('+905301234003', 'Message content 3', 'pending'),
('+905301234004', 'Message content 4', 'pending'),
('+905301234005', 'Message content 5', 'pending'),
('+905301234006', 'Message content 6', 'pending'),
('+905301234007', 'Message content 7', 'pending'),
('+905301234008', 'Message content 8', 'pending'),
('+905301234009', 'Message content 9', 'pending'),
('+905301234010', 'Message content 10', 'pending'),
('+905301234011', 'Message content 11', 'pending'),
('+905301234012', 'Message content 12', 'pending'),
('+905301234013', 'Message content 13', 'pending'),
('+905301234014', 'Message content 14', 'pending'),
('+905301234015', 'Message content 15', 'pending'),
('+905301234016', 'Message content 16', 'pending'),
('+905301234017', 'Message content 17', 'pending'),
('+905301234018', 'Message content 18', 'pending'),
('+905301234019', 'Message content 19', 'pending'),
('+905301234020', 'Message content 20', 'pending'),
('+905301234021', 'Message content 21', 'pending'),
('+905301234022', 'Message content 22', 'pending'),
('+905301234023', 'Message content 23', 'pending'),
('+905301234024', 'Message content 24', 'pending'),
('+905301234025', 'Message content 25', 'pending'),
('+905301234026', 'Message content 26', 'pending'),
('+905301234027', 'Message content 27', 'pending'),
('+905301234028', 'Message content 28', 'pending'),
('+905301234029', 'Message content 29', 'pending'),
('+905301234030', 'Message content 30', 'pending'),
('+905301234031', 'Message content 31', 'pending'),
('+905301234032', 'Message content 32', 'pending'),
('+905301234033', 'Message content 33', 'pending'),
('+905301234034', 'Message content 34', 'pending'),
('+905301234035', 'Message content 35', 'pending'),
('+905301234036', 'Message content 36', 'pending'),
('+905301234037', 'Message content 37', 'pending'),
('+905301234038', 'Message content 38', 'pending'),
('+905301234039', 'Message content 39', 'pending'),
('+905301234040', 'Message content 40', 'pending'),
('+905301234041', 'Message content 41', 'pending'),
('+905301234042', 'Message content 42', 'pending'),
('+905301234043', 'Message content 43', 'pending'),
('+905301234044', 'Message content 44', 'pending'),
('+905301234045', 'Message content 45', 'pending'),
('+905301234046', 'Message content 46', 'pending'),
('+905301234047', 'Message content 47', 'pending'),
('+905301234048', 'Message content 48', 'pending'),
('+905301234049', 'Message content 49', 'pending'),
('+905301234050', 'Message content 50', 'pending');
