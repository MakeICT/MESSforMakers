--------------------------------------------------------------------------------------------------------------------------------
-- Tables need for control of member in the site
--------------------------------------------------------------------------------------------------------------------------------

-- read-list, read, write
-- maps to the http routes
-- GET list, GET object, PUT:POST:DELETE
-- TODO all RBAC stuff needs to be added to the wireframes
CREATE TABLE rbac_permission_access (
      id SERIAL PRIMARY KEY
    , name TEXT NOT NULL
    , UNIQUE (name)
);
COMMENT ON TABLE rbac_permission_access IS 'The actions that are possible on a resource: list, read, write';

INSERT INTO rbac_permission_access (name) VALUES ('list'), ('read'), ('write');

CREATE TABLE rbac_permission (
      id SERIAL PRIMARY KEY
	, rbac_permission_access_id INTEGER NOT NULL REFERENCES rbac_permission_access(id) ON DELETE RESTRICT
    , name TEXT NOT NULL
    , created_at TIMESTAMP NOT NULL DEFAULT now()
    , UNIQUE (name)
);
COMMENT ON TABLE rbac_permission IS 'Links a type of access to a permission';

CREATE TABLE rbac_role (
      id SERIAL PRIMARY KEY
    , name TEXT NOT NULL
    , created_at TIMESTAMP NOT NULL DEFAULT now()
    , UNIQUE (name)
);
COMMENT ON TABLE rbac_role IS 'Groups permissions to assign to a member';

INSERT INTO rbac_role (name) VALUES ('guest');

CREATE TABLE rbac_role_permission_rel (
      id SERIAL PRIMARY KEY
    , rbac_role_id INTEGER NOT NULL REFERENCES rbac_role(id) ON DELETE CASCADE
    , rbac_permission_id INTEGER NOT NULL REFERENCES rbac_permission(id) ON DELETE RESTRICT
    , UNIQUE (rbac_role_id, rbac_permission_id)
);
COMMENT ON TABLE rbac_role_permission_rel IS 'Links various permissions to a single role';

CREATE TABLE rbac_group(
	id SERIAL PRIMARY KEY
	, name TEXT NOT NULL
	, UNIQUE (name)
);
COMMENT ON TABLE rbac_group IS 'Groups of permissions to make assignment easier';

CREATE TABLE rbac_role_group_rel (
	id SERIAL PRIMARY KEY
	, rbac_role_id INTEGER NOT NULL REFERENCES rbac_role(id) ON DELETE CASCADE
	, rbac_group_id INTEGER NOT NULL REFERENCES rbac_group(id) ON DELETE CASCADE
	, UNIQUE (rbac_role_id, rbac_group_id)
);
COMMENT ON TABLE rbac_role_group_rel IS 'links the group of permissions to a specific role';

CREATE TABLE rbac_group_permission_rel (
	id SERIAL PRIMARY KEY
	, rbac_group_id INTEGER NOT NULL REFERENCES rbac_group(id) ON DELETE CASCADE
	, rbac_permission_id INTEGER NOT NULL REFERENCES rbac_permission(id) ON DELETE CASCADE
);
COMMENT ON TABLE rbac_group_permission_rel IS 'links a specific permission to a group';

--------------------------------------------------------------------------------------------------------------------------------
-- Member
--------------------------------------------------------------------------------------------------------------------------------

CREATE TABLE membership_status (
      id SERIAL PRIMARY KEY
    , name TEXT NOT NULL
    , UNIQUE (name)
);
COMMENT ON TABLE membership_status IS 'Holds values to indicate whether a person is an active member of MakeICT or not';
INSERT INTO membership_status (name) VALUES ('guest'), ('active'), ('past_due'), ('quit');

CREATE TABLE membership_options (
	id SERIAL PRIMARY KEY
	, name TEXT NOT NULL
	, is_recurring BOOLEAN NOT NULL DEFAULT 'f'
	, period INTERVAL NOT NULL DEFAULT '1 month'
	, UNIQUE (name)
);
COMMENT ON TABLE membership_options IS 'for members only, not guests, defines when member will be charged for dues';
INSERT INTO membership_options (name, is_recurring, period) VALUES ('One month', 'f', '0 months'), ('Recurring - monthly', 't', '1 month'), ('Recurring - 6 months', 't', '6 months'), ('Recurring - 12 months', 't', '12 months');
-- TODO implement some interface to allow admins to add or remove options

CREATE TABLE member (
      id SERIAL PRIMARY KEY
    , first_name TEXT NOT NULL
    , last_name TEXT NOT NULL
    , username TEXT NOT NULL -- must be valid email
	, password TEXT NOT NULL
    , dob DATE NOT NULL
    , phone TEXT NOT NULL
	, membership_status_id INTEGER NOT NULL REFERENCES membership_status(id)
	, membership_expires DATE 
	, membership_option INTEGER REFERENCES membership_options(id)
	, rbac_role_id INTEGER NOT NULL REFERENCES rbac_role(id)
    , created_at TIMESTAMP NOT NULL DEFAULT now()
    , updated_at TIMESTAMP NOT NULL DEFAULT now()
    , UNIQUE (username)  -- members cannot sign up for multiple accounts with the same email
);
COMMENT ON TABLE member IS 'Core table of all members and guests';

CREATE TABLE member_address (
	id SERIAL PRIMARY KEY
	, member_id INTEGER NOT NULL REFERENCES member(id) ON DELETE CASCADE
	, addr_type TEXT NOT NULL --home or billing
	, addr1 TEXT NOT NULL
	, addr2 TEXT
	, city TEXT NOT NULL
	, state TEXT NOT NULL
	, zip TEXT NOT NULL
	, UNIQUE(member_id, addr_type)
);
COMMENT ON TABLE member_address IS 'Home and billing addresses for members and guests';

CREATE TABLE member_access_token (
      id SERIAL PRIMARY KEY
	, member_id INTEGER NOT NULL REFERENCES member(id) ON DELETE CASCADE
    , created_at TIMESTAMP NOT NULL DEFAULT now()
    , expires_at TIMESTAMP NOT NULL DEFAULT (now() + INTERVAL '24 hours')
    , token TEXT NOT NULL 
    , UNIQUE (member_id)
);
COMMENT ON TABLE member_access_token IS 'Tracks token for user to register for first time or reset password';

CREATE TABLE login_status (
      id SERIAL PRIMARY KEY
    , name TEXT NOT NULL
);
INSERT INTO login_status (name) VALUES ('Success'), ('Unknown Username'), ('Wrong Password');

CREATE TABLE login_log (
      id SERIAL PRIMARY KEY
    , username TEXT NOT NULL
    , login_status_id INTEGER NOT NULL REFERENCES login_status(id) ON DELETE RESTRICT
    , created_at TIMESTAMP NOT NULL DEFAULT now()
);
COMMENT ON TABLE login_log IS 'Keep track of member login attempts for troubleshooting and usage data';

CREATE TABLE session (
	  id SERIAL PRIMARY KEY
	, member_id INTEGER NOT NULL REFERENCES member(id) ON DELETE CASCADE
	, authtoken VARCHAR(64) NOT NULL
	, originated TIMESTAMP NOT NULL DEFAULT now()
	, last_seen TIMESTAMP NOT NULL DEFAULT now()
	, last_ip VARCHAR(46) NOT NULL --46 characters will allow for storing ipv6 addresses or ipv4
	, agent VARCHAR(100) NOT NULL --100 characters is enough to get a good idea of the user agent
	, UNIQUE (authtoken)
);
COMMENT ON TABLE session IS 'Keep track of user sessions, allow them to be deleted, expired, and investigated.';

CREATE TABLE member_ice (
      id SERIAL PRIMARY KEY
	, member_id INTEGER NOT NULL REFERENCES member(id) ON DELETE CASCADE
    , name TEXT NOT NULL
	, phone_number TEXT NOT NULL
    , relationship TEXT NOT NULL  -- just raw text of whatevr the member entered. No need to list out every possible relationship choice.
);
COMMENT ON TABLE member_ice IS 'Member In case of emergency (ICE)';

CREATE TABLE waivers (
	id SERIAL PRIMARY KEY
	, filename TEXT NOT NULL
	, date_signed DATE NOT NULL DEFAULT now()
	, valid BOOLEAN NOT NULL DEFAULT 'f'  
	, member_id INTEGER NOT NULL REFERENCES member(id) ON DELETE CASCADE
	, UNIQUE (filename)
);
COMMENT ON TABLE waivers IS 'Location and date for all member waivers to allow area and equipment access';

CREATE TABLE addon_types (
	id SERIAL PRIMARY KEY
	, name TEXT NOT NULL
	, monthly_cost MONEY NOT NULL
	, UNIQUE (name)
);
COMMENT ON TABLE addon_types IS 'Types of additional services that can be purchased along with membership dues';
INSERT INTO addon_types (name, monthly_cost) VALUES ('Locker', '$5.00'), ('Studio', '$75.00');
-- TODO add interface for admin to add and remove addons

CREATE TABLE member_addon_rel (
	id SERIAL PRIMARY KEY
	, member_id INTEGER NOT NULL REFERENCES member(id) ON DELETE CASCADE
	, addon_id INTEGER NOT NULL REFERENCES addon_types(id) ON DELETE RESTRICT
	, created_at TIMESTAMP NOT NULL DEFAULT now()
	, UNIQUE (member_id, addon_id)
);
COMMENT ON TABLE member_addon_rel IS 'Links a member to the additional services they have purchased';

CREATE TABLE member_recurring_donation (
	id SERIAL PRIMARY KEY
	, member_id INTEGER NOT NULL REFERENCES member(id) ON DELETE CASCADE
	, amount MONEY NOT NULL DEFAULT '$0.00'
	, created_at TIMESTAMP NOT NULL DEFAULT now()
	, updated_at TIMESTAMP
	, UNIQUE (member_id)
);
COMMENT ON TABLE member_recurring_donation IS 'If a member chooses to support with an additionaly monthly donations, it is stored here';

CREATE TABLE locker (
	id serial PRIMARY KEY
	, locker_id TEXT NOT NULL
	, UNIQUE (locker_id)
);
COMMENT ON TABLE locker IS 'list of lockers that can be rented';
-- TODO add interface so that admin can manage lockers

CREATE TABLE member_locker_rel (
	id SERIAL PRIMARY KEY
	, member_id INTEGER NOT NULL REFERENCES member(id) ON DELETE CASCADE
	, locker_id INTEGER NOT NULL REFERENCES locker(id) ON DELETE RESTRICT
	, UNIQUE (member_id, locker_id)
);
COMMENT ON TABLE member_locker_rel IS 'Records which member has which locker';

--------------------------------------------------------------------------------------------------------------------------------
-- Invoice and Payment
--------------------------------------------------------------------------------------------------------------------------------

-- I think this whole finance section is woefully inadequate, but I don't know how to fix it

CREATE TABLE invoice_status (
	id SERIAL PRIMARY KEY
	, name TEXT NOT NULL
	, UNIQUE (name)
);
COMMENT ON TABLE invoice_status IS 'Contains the status types for an invoice';
INSERT INTO invoice_status(name) VALUES ('paid'), ('unpaid'), ('cancelled');

CREATE TABLE payment_method (
	id SERIAL PRIMARY KEY
	, name TEXT NOT NULL
	, UNIQUE (name)
);
COMMENT ON TABLE payment_method IS 'Values for what kind of payment was made';
INSERT INTO payment_method(name) VALUES ('cash'), ('check'), ('online');

CREATE TABLE payment (
	id SERIAL PRIMARY KEY
	, amount MONEY NOT NULL DEFAULT 0.00
	, member_id INTEGER NOT NULL REFERENCES member(id)
	, payment_method_id INTEGER NOT NULL REFERENCES payment_method(id)
	, created_at TIMESTAMP NOT NULL DEFAULT now()  -- TODO changes diagram
);
COMMENT ON TABLE payment IS 'Holds payment history for members';

CREATE TABLE invoice (
	id SERIAL PRIMARY KEY
	, amount MONEY NOT NULL DEFAULT 0.00
	, description TEXT NOT NULL
	, member_id INTEGER NOT NULL REFERENCES member(id) 
	, payment_id INTEGER REFERENCES payment(id)
	, status_id INTEGER REFERENCES invoice_status(id)
	, created_at TIMESTAMP NOT NULL DEFAULT now()   
);
COMMENT ON TABLE invoice IS 'Records an amount that a member owes for dues, class fees, or other fees';

--------------------------------------------------------------------------------------------------------------------------------
-- Locations, Areas and Equipment
--------------------------------------------------------------------------------------------------------------------------------

CREATE TABLE location (
	id SERIAL PRIMARY KEY
	, name TEXT NOT NULL
	, address_1 TEXT NOT NULL
	, address_2 TEXT NOT NULL
	, city TEXT NOT NULL
	, state TEXT NOT NULL
	, zip TEXT NOT NULL
	, UNIQUE (name)
);
COMMENT ON TABLE location IS 'General locations, not specific rooms';
-- admin interface for creating/removing locations

CREATE TABLE area (
      id SERIAL PRIMARY KEY
    , name TEXT NOT NULL
    , UNIQUE (name)
);
COMMENT ON TABLE area IS 'Lists all the areas available for reservation or use. Needs to have a location';
-- TODO add area management to the wireframes

CREATE TABLE equipment (
      id SERIAL PRIMARY KEY
    , area_id INTEGER NOT NULL REFERENCES area(id) ON DELETE RESTRICT
    , name TEXT NOT NULL
    , brought_at TIMESTAMP
    , created_at TIMESTAMP NOT NULL DEFAULT now()
);
COMMENT ON TABLE equipment IS 'Lists all the equipment owned by the organization and what area it is in';
-- TODO add equipment management to the wireframes

--------------------------------------------------------------------------------------------------------------------------------
-- Certifications
--------------------------------------------------------------------------------------------------------------------------------

CREATE TABLE certification (
	id SERIAL PRIMARY KEY
	, name TEXT NOT NULL
	, created_at TIMESTAMP NOT NULL DEFAULT now()
	, UNIQUE (name)
);
COMMENT ON TABLE certification IS 'Records the different kinds of certifications that a member can have';
-- Is certification management part of equipment management (e.g. you indicate that equipment needs cert from the equipment edit page), or separate?

CREATE TABLE member_certification (  -- does this need a status-active/status-revoked?
	id SERIAL PRIMARY KEY
	, member_id INTEGER NOT NULL REFERENCES member(id) ON DELETE CASCADE
	, certification_id INTEGER NOT NULL REFERENCES certification(id) ON DELETE RESTRICT
	, created_at TIMESTAMP NOT NULL DEFAULT now()
	, UNIQUE (member_id, certification_id)
);
COMMENT ON TABLE member_certification IS 'Records the certifications that a member has';

CREATE TABLE certification_approval (
	id SERIAL PRIMARY KEY
	, member_certification_id INTEGER NOT NULL REFERENCES member_certification(id) ON DELETE CASCADE
	, approver_id INTEGER NOT NULL REFERENCES member(id) ON DELETE RESTRICT
	, created_at TIMESTAMP NOT NULL DEFAULT now()
	, UNIQUE (member_certification_id, approver_id)
);
COMMENT ON TABLE certification_approval IS 'Records the person that approved the certification. membver certification is not valid without an approval';

CREATE TABLE certification_resource_rel (
	id SERIAL PRIMARY KEY
	, certification_id INTEGER NOT NULL REFERENCES certification(id) ON DELETE CASCADE
	, equipment_id INTEGER NOT NULL REFERENCES equipment(id) ON DELETE RESTRICT
	, area_id INTEGER NOT NULL REFERENCES area(id)
);
COMMENT ON TABLE certification_resource_rel IS 'Links a certification to the resources to which it grants access';

-- event checkin and certification approval need to be figured out in the wireframes

--------------------------------------------------------------------------------------------------------------------------------
-- Event/Classes
--------------------------------------------------------------------------------------------------------------------------------

CREATE TABLE event (
      id SERIAL PRIMARY KEY
    , name TEXT NOT NULL
    , during TSRANGE NOT NULL  -- TODO used range instead of start-stop or start-length to match the reservation tables. but maybe start-length is better?
    , location_id INTEGER NOT NULL REFERENCES location(id) ON DELETE RESTRICT  
    -- TODO need location management added to the wireframes
    , created_by INTEGER NOT NULL REFERENCES member(id) ON DELETE RESTRICT
    , created_at TIMESTAMP NOT NULL DEFAULT now()
    , description TEXT NOT NULL
    -- The following options determine whether the software will look for information in the corresponding tables.
    -- There are too many permutations of event types to represent them using different types of events.
	, requires_materials BOOLEAN NOT NULL DEFAULT 'f'
	, requires_registration BOOLEAN NOT NULL DEFAULT 'f'
	, requires_fees BOOLEAN NOT NULL DEFAULT 'f'
	, requires_prerequisites BOOLEAN NOT NULL DEFAULT 'f'
    , grants_certifications BOOLEAN NOT NULL DEFAULT 'f'
);
COMMENT ON TABLE event IS 'Master Event List';

CREATE TABLE event_materials (
	id SERIAL PRIMARY KEY
	, event_id INTEGER NOT NULL REFERENCES event(id) ON DELETE CASCADE
	, item TEXT NOT NULL 
	, UNIQUE (event_id, item)  -- this could be easy for users to duplicate. Check client side and server side before INSERT
);
COMMENT ON TABLE event_materials IS 'List of materials needed for each event';

CREATE TABLE event_registration (
	id SERIAL PRIMARY KEY
	, event_id INTEGER NOT NULL REFERENCES event(id) ON DELETE CASCADE
	, registration_begins DATE NOT NULL DEFAULT now()
	, registration_ends DATE -- TODO practice with checks to see if default can be the event start time from event table, and checks to enforce that end time is after start time
	, min_attendees SMALLINT
	, max_attendees SMALLINT
	, UNIQUE (event_id)
);
COMMENT ON TABLE event_registration IS 'Registration requirements for a specific event';

CREATE TABLE fee_structures (
	id SERIAL PRIMARY KEY
	, name TEXT NOT NULL
	, facilities_fee MONEY NOT NULL DEFAULT '$0.00'
	, authorization_fee MONEY NOT NULL DEFAULT '$0.00'
	, UNIQUE (name)
);
COMMENT ON TABLE fee_structures IS 'Sets the Authorization and Facilities fees for non-members. Material and class fees are stored in event_fees';
INSERT INTO fee_structures (name, facilities_fee, authorization_fee) VALUES ('Authorization', '$0.00', '$20.00'), ('Workshop', '$5.00', '$0.00'), ('Class', '$5.00', '$0.00');
-- TODO admin interface to manage fee structures

CREATE TABLE refund_policy (
	id SERIAL PRIMARY KEY
	, name TEXT NOT NULL
	, description TEXT NOT NULL
	, UNIQUE (name)
);
COMMENT ON TABLE refund_policy IS 'Text description of how and when a member might get a refund if they do not attend a paid-for event';
-- TODO add refund policy management to the wireframes

CREATE TABLE event_fees (
	id SERIAL PRIMARY KEY
	, event_id INTEGER NOT NULL REFERENCES event(id) ON DELETE CASCADE
	, fee_structure_id INTEGER NOT NULL REFERENCES fee_structures(id) ON DELETE RESTRICT
	, refund_policy_id INTEGER NOT NULL REFERENCES refund_policy(id) ON DELETE RESTRICT
	, class_material_fee MONEY NOT NULL DEFAULT '$0.00'
	, UNIQUE (event_id)
);
COMMENT ON TABLE event_fees IS 'Sets up the fees for a specific event';

CREATE TABLE event_prerequisites_rel (
	id SERIAL PRIMARY KEY
	, certification_id INTEGER NOT NULL REFERENCES certification(id) ON DELETE RESTRICT
	, event_id INTEGER NOT NULL REFERENCES event(id) ON DELETE CASCADE
	, UNIQUE (certification_id, event_id)
);
COMMENT ON TABLE event_prerequisites_rel IS 'List of all the certifications a member must have before registering for this event';

CREATE TABLE event_certifications_rel (
	id SERIAL PRIMARY KEY
	, certification_id INTEGER NOT NULL REFERENCES certification(id) ON DELETE RESTRICT
	, event_id INTEGER NOT NULL REFERENCES event(id) ON DELETE CASCADE
	, UNIQUE (certification_id, event_id)
);
COMMENT ON TABLE event_certifications_rel IS 'List of all the certifications given to a member when they complete this event';

CREATE TABLE checked_in_status (
	id SERIAL PRIMARY KEY
	, status TEXT NOT NULL
	, UNIQUE (status)
);
COMMENT ON TABLE checked_in_status IS 'List of the different stages of check in';
INSERT INTO checked_in_status (status) VALUES ('Not Checked In'), ('Checked In'), ('Cancelled');

CREATE TABLE member_event_registration ( 
	id SERIAL PRIMARY KEY
	, member_id INTEGER NOT NULL REFERENCES member(id) ON DELETE CASCADE
	, event_id INTEGER NOT NULL REFERENCES event(id) ON DELETE RESTRICT
	, checked_in_status_id INTEGER NOT NULL REFERENCES checked_in_status(id) ON DELETE RESTRICT
	, payment_status TEXT NOT NULL DEFAULT 'unpaid' -- paid or unpaid  TODO make a table to hold these strings
	, payment_id INTEGER REFERENCES payment(id) ON DELETE RESTRICT
	, created_at TIMESTAMP NOT NULL DEFAULT now()
	, updated_at TIMESTAMP 
	, UNIQUE (member_id, event_id)
);
COMMENT ON TABLE event_registration IS 'Shows a users registration for event, whether they paid, and whether they attended';

CREATE TABLE event_template (
	id SERIAL PRIMARY KEY
	, name TEXT NOT NULL -- name of the template, defaults to name of event 
	, event_id INTEGER NOT NULL REFERENCES event(id) ON DELETE RESTRICT
	, created_at TIMESTAMP NOT NULL DEFAULT now()
	, UNIQUE (event_id)
);
COMMENT ON TABLE event_template IS 'List of events that can be used as a template for creating a new event';
-- TODO add dedicated template management, eventually. Until then, manage from event creation page only

CREATE TABLE event_area_rel (
	id SERIAL PRIMARY KEY
	, event_id INTEGER NOT NULL REFERENCES event(id) ON DELETE CASCADE
	, area_id INTEGER NOT NULL REFERENCES area(id) ON DELETE RESTRICT
);
COMMENT ON TABLE event_area_rel IS 'An area required by the event';

CREATE TABLE event_equipment_rel (
	id SERIAL PRIMARY KEY
	, event_id INTEGER NOT NULL REFERENCES event(id) ON DELETE CASCADE
	, equipment_id INTEGER NOT NULL REFERENCES equipment(id) ON DELETE RESTRICT
);
COMMENT ON TABLE event_equipment_rel IS 'An equipment required by the event';

CREATE TABLE host_role (
	id SERIAL PRIMARY KEY
	, name TEXT NOT NULL
	, UNIQUE (name)
);
COMMENT ON TABLE host_role IS 'The different kinds of roles that a host might have';
INSERT INTO host_role (name) VALUES ('Host'), ('Instructor'), ('Check-In Volunteer');

CREATE TABLE event_host (
	id SERIAL PRIMARY KEY
	, member_id INTEGER NOT NULL REFERENCES member(id) ON DELETE CASCADE
	, event_id INTEGER NOT NULL REFERENCES event(id) ON DELETE CASCADE
	, host_role_id INTEGER NOT NULL REFERENCES host_role(id) ON DELETE RESTRICT
	, created_at TIMESTAMP NOT NULL DEFAULT now()
	, UNIQUE (member_id, event_id)
);
COMMENT ON TABLE event_host IS 'Lists all the people associated with a specific event, and how they are related to it.';
-- TODO future feature: hosts may only be selected from a table of people that have chosen to be hosts, including preferred contact methods

CREATE TABLE area_reservation (
	id SERIAL PRIMARY KEY
	, area_id INTEGER NOT NULL REFERENCES area(id) ON DELETE RESTRICT
	, during TSRANGE NOT NULL  -- uses timestamp range to make overlap exclusion possible
	, event_id INTEGER NOT NULL REFERENCES event(id) ON DELETE CASCADE
	, EXCLUDE USING gist (area_id WITH =, during WITH &&)
);
COMMENT ON TABLE area_reservation IS 'Reserves a room for an event during a specific time period';
-- TODO future feature: add calendars and reservations for users to reserve equipment or areas during specific times

CREATE TABLE equipment_reservation (
	id SERIAL PRIMARY KEY
	, equipment_id INTEGER NOT NULL REFERENCES equipment(id) ON DELETE RESTRICT
	, during TSRANGE NOT NULL    -- uses timestamp range to make overlap exclusion possible
	, event_id INTEGER NOT NULL REFERENCES event(id) ON DELETE CASCADE
	, EXCLUDE USING gist (equipment_id WITH =, during WITH &&)
);
COMMENT ON TABLE equipment_reservation IS 'Reserves a piece of equipment for an event during a specific time period';
