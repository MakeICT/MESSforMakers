
--------------------------------------------------------------------------------------------------------------------------------
-- Member
--------------------------------------------------------------------------------------------------------------------------------

CREATE TABLE member_state (
      id SERIAL PRIMARY KEY
    , name TEXT NOT NULL
);
COMMENT ON TABLE member_state IS '';
INSERT INTO member_state (name) VALUES ('active'), ('pass_due'), ('quit'), ('guest');

CREATE TABLE member (
      id SERIAL PRIMARY KEY
    , name TEXT NOT NULL
    , username TEXT NOT NULL -- must be valid email
	, password TEXT NOT NULL
    , dob DATE NOT NULL
    , waiver_signed BOOLEAN NOT NULL DEFAULT FALSE
	, member_state_id INTEGER NOT NULL REFERENCES member_state(id)
    , created_at TIMESTAMP NOT NULL DEFAULT now()
    , updated_at TIMESTAMP
    , UNIQUE (username)
);
COMMENT ON TABLE member IS 'Core table of all members';

CREATE TABLE member_access_token (
      id SERIAL PRIMARY KEY
	, member_id INTEGER NOT NULL REFERENCES member(id)
    , created_at TIMESTAMP NOT NULL DEFAULT now()
    , UNIQUE (member_id)
);
COMMENT ON TABLE member IS 'Tracks token for user to register for first time or restart password';

CREATE TABLE login_status (
      id SERIAL PRIMARY KEY
    , name TEXT NOT NULL
);
INSERT INTO login_status (name) VALUES ('Success'), ('Unknown Username'), ('Wrong Password');

CREATE TABLE login_log (
      id SERIAL PRIMARY KEY
    , username TEXT NOT NULL
    , login_status_id INTEGER NOT NULL REFERENCES login_status(id)
    , created_at TIMESTAMP NOT NULL DEFAULT now()
);
COMMENT ON TABLE login_log IS 'Keep track of member login for trouble shooting and usage data';

CREATE TABLE member_ice (
      id SERIAL PRIMARY KEY
	, member_id INTEGER NOT NULL REFERENCES member(id)
    , name TEXT NOT NULL
	, phone_number TEXT NOT NULL
    -- , relation_id INTEGER NOT NULL REFERENCES relation_type -- how t
);
COMMENT ON TABLE member_ice IS 'Member In case of emergency (ICE)';


-- name bagde gets a qr code
-- look them up in the furture and track them ???


--------------------------------------------------------------------------------------------------------------------------------
-- Event/Classes
--------------------------------------------------------------------------------------------------------------------------------

CREATE TABLE event_type (
      id SERIAL PRIMARY KEY
    , name TEXT NOT NULL
);
INSERT INTO event_type (name) VALUES ('Class'), ('Workshop');

CREATE TABLE event (
      id SERIAL PRIMARY KEY
    , name TEXT NOT NULL
    , day DATE NOT NULL
    , time TIME NOT NULL
    , duration INTERVAL NOT NULL
    , location TEXT NOT NULL -- TODO fixme
    , created_by INTEGER NOT NULL REFERENCES member(id)
    , created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE event_host_rel (
      id SERIAL PRIMARY KEY
	, member_id INTEGER NOT NULL REFERENCES member(id)
	, event_id INTEGER NOT NULL REFERENCES event(id)
    , role TEXT NOT NULL --- TODO fixme why they are doing at the event like Instructor, Check-in, Assistand, ???
    , created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE event_attendee_rel (
      id SERIAL PRIMARY KEY
	, member_id INTEGER NOT NULL REFERENCES member(id)
	, event_id INTEGER NOT NULL REFERENCES event(id)
    , created_at TIMESTAMP NOT NULL DEFAULT now()
);

--------------------------------------------------------------------------------------------------------------------------------
-- Tables need for control of member in the site
--------------------------------------------------------------------------------------------------------------------------------

-- read-list, read, write
-- maps to the http routes
-- GET list, GET object, PUT:POST:DELETE
CREATE TABLE rbac_permission_access (
      id SERIAL PRIMARY KEY
    , name TEXT NOT NULL
    , UNIQUE (name)
);
COMMENT ON TABLE rbac_permission_access IS 'The actions that are possible on a resource: list, read, write';

INSERT INTO rbac_permission_access (name) VALUES ('list'), ('read'), ('write');

CREATE TABLE rbac_permission (
      id SERIAL PRIMARY KEY
	, rbac_permission_access_id INTEGER NOT NULL REFERENCES rbac_permission_access(id)
    , name TEXT NOT NULL
    , created_at TIMESTAMP NOT NULL DEFAULT now()
    , UNIQUE (name)
);

CREATE TABLE rbac_role (
      id SERIAL PRIMARY KEY
    , name TEXT NOT NULL
    , created_at TIMESTAMP NOT NULL DEFAULT now()
    , UNIQUE (name)
);

CREATE TABLE rbac_role_permission_rel (
      id SERIAL PRIMARY KEY
    , rbac_role INTEGER NOT NULL REFERENCES rbac_role(id)
    , rbac_permission_id INTEGER NOT NULL REFERENCES rbac_permission(id)
    , UNIQUE (rbac_role, rbac_permission_id)
);

CREATE TABLE rbac_member_role_rel (
      id SERIAL PRIMARY KEY
    , member_id INTEGER NOT NULL REFERENCES member(id)
    , rbac_role_id INTEGER NOT NULL REFERENCES rbac_role(id)
    , UNIQUE (member_id, rbac_role_id)
);
COMMENT ON TABLE rbac_member_role_rel IS 'Defines the role a member is a part of';


--------------------------------------------------------------------------------------------------------------------------------
-- Area
--------------------------------------------------------------------------------------------------------------------------------

-- Is the maker space an area ???
CREATE TABLE area (
      id SERIAL PRIMARY KEY
    , name TEXT NOT NULL
    , UNIQUE (name)
);

CREATE TABLE equipment (
      id SERIAL PRIMARY KEY
    , area_id INTEGER NOT NULL REFERENCES area(id)
    , name TEXT NOT NULL
    , tag TEXT NOT NULL -- like a qr sticker on the equipment
    , brought_at TIMESTAMP
    , created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- keeps it running
CREATE TABLE equipment_maintainer_member_rel (
      id SERIAL PRIMARY KEY
    , member_id INTEGER NOT NULL REFERENCES member(id)
    , equipment_id INTEGER NOT NULL REFERENCES equipment(id)
);

-- can use the it
CREATE TABLE equipment_auth_member_rel (
      id SERIAL PRIMARY KEY
    , member_id INTEGER NOT NULL REFERENCES member(id)
    , equipment_id INTEGER NOT NULL REFERENCES equipment(id)
    , created_at TIMESTAMP NOT NULL DEFAULT now()
    , UNIQUE (member_id, equipment_id)
);

CREATE TABLE equipment_reservation (
      id SERIAL PRIMARY KEY
);

--------------------------------------------------------------------------------------------------------------------------------
-- Ticket system
--------------------------------------------------------------------------------------------------------------------------------

CREATE TABLE ticket_type (
      id SERIAL PRIMARY KEY
    , name TEXT NOT NULL
);
INSERT INTO ticket_type (name) VALUES ('Membership'), ('Area'), ('Website');

CREATE TABLE ticket_category (
      id SERIAL PRIMARY KEY
    , ticket_type_id INTEGER NOT NULL REFERENCES ticket_type(id)
    , name TEXT NOT NULL
);

-- somehow need to route the ticket to groups
-- this should be done with the ares lead

CREATE TABLE ticket_state (
      id SERIAL PRIMARY KEY
    , name TEXT NOT NULL
);
INSERT INTO ticket_state (name) VALUES ('Open'), ('Close');

CREATE TABLE ticket (
      id SERIAL PRIMARY KEY
    , ticket_category_id INTEGER NOT NULL REFERENCES ticket_category(id)
    , ticket_state_id INTEGER NOT NULL REFERENCES ticket_state(id)
	, title text NOT NULL
	, member_id INTEGER NOT NULL REFERENCES member(id)
    , created_at TIMESTAMP
    , updated_at TIMESTAMP
);

CREATE TABLE ticket_version (
      id SERIAL PRIMARY KEY
    , ticket_id INTEGER NOT NULL REFERENCES ticket(id)
	, member_id INTEGER NOT NULL REFERENCES member(id)
    , comments TEXT NOT NULL
    , created_at TIMESTAMP
    , updated_at TIMESTAMP
);
--------------------------------------------------------------------------------------------------------------------------------
-- Door Access
--------------------------------------------------------------------------------------------------------------------------------

-- a member can have more then on
CREATE TABLE access_key (
      id SERIAL PRIMARY KEY
	, member_id INTEGER NOT NULL REFERENCES member(id)
    --, access_key_state_id
    , sercert TEXT NOT NULL
);

CREATE TABLE access_key_area_rel (
      id SERIAL PRIMARY KEY
	, access_key_id INTEGER NOT NULL REFERENCES access_key(id)
	, area_id INTEGER NOT NULL REFERENCES area(id)
);

CREATE TABLE access_log (
      id SERIAL PRIMARY KEY
	, access_key_id INTEGER NOT NULL REFERENCES access_key(id)
    --, access_status_id
    , created_at TIMESTAMP
);
