--
-- PostgreSQL database dump
--

\restrict 4IA6MIzqh232qwcJZIukbtnlb0enZHiG1RSuwZRJwfn8tQlvm8EL84pjAZWDmjt

-- Dumped from database version 16.13
-- Dumped by pg_dump version 16.13

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

DROP POLICY IF EXISTS tenant_isolation ON public.warehouses;
DROP POLICY IF EXISTS tenant_isolation ON public.vehicles;
DROP POLICY IF EXISTS tenant_isolation ON public.users;
DROP POLICY IF EXISTS tenant_isolation ON public.units;
DROP POLICY IF EXISTS tenant_isolation ON public.supply_requests;
DROP POLICY IF EXISTS tenant_isolation ON public.shipments;
DROP POLICY IF EXISTS tenant_isolation ON public.shipment_refuels;
DROP POLICY IF EXISTS tenant_isolation ON public.resources;
DROP POLICY IF EXISTS tenant_isolation ON public.notifications;
DROP POLICY IF EXISTS tenant_isolation ON public.gps_locations;
DROP POLICY IF EXISTS tenant_isolation ON public.geofences;
DROP POLICY IF EXISTS tenant_isolation ON public.fuel_records;
DROP POLICY IF EXISTS tenant_isolation ON public.contractor_requests;
DROP POLICY IF EXISTS tenant_isolation ON public.audit_logs;
ALTER TABLE IF EXISTS ONLY public.warehouses DROP CONSTRAINT IF EXISTS warehouses_unit_id_fkey;
ALTER TABLE IF EXISTS ONLY public.warehouses DROP CONSTRAINT IF EXISTS warehouses_tenant_id_fkey;
ALTER TABLE IF EXISTS ONLY public.vehicles DROP CONSTRAINT IF EXISTS vehicles_tenant_id_fkey;
ALTER TABLE IF EXISTS ONLY public.vehicles DROP CONSTRAINT IF EXISTS vehicles_home_warehouse_id_fkey;
ALTER TABLE IF EXISTS ONLY public.vehicles DROP CONSTRAINT IF EXISTS vehicles_driver_id_fkey;
ALTER TABLE IF EXISTS ONLY public.vehicles DROP CONSTRAINT IF EXISTS vehicles_current_warehouse_id_fkey;
ALTER TABLE IF EXISTS ONLY public.vehicle_driver_history DROP CONSTRAINT IF EXISTS vehicle_driver_history_vehicle_id_fkey;
ALTER TABLE IF EXISTS ONLY public.vehicle_driver_history DROP CONSTRAINT IF EXISTS vehicle_driver_history_driver_id_fkey;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_tenant_id_fkey;
ALTER TABLE IF EXISTS ONLY public.units DROP CONSTRAINT IF EXISTS units_tenant_id_fkey;
ALTER TABLE IF EXISTS ONLY public.units DROP CONSTRAINT IF EXISTS units_parent_id_fkey;
ALTER TABLE IF EXISTS ONLY public.supply_requests DROP CONSTRAINT IF EXISTS supply_requests_tenant_id_fkey;
ALTER TABLE IF EXISTS ONLY public.supply_requests DROP CONSTRAINT IF EXISTS supply_requests_target_warehouse_id_fkey;
ALTER TABLE IF EXISTS ONLY public.supply_requests DROP CONSTRAINT IF EXISTS supply_requests_resource_id_fkey;
ALTER TABLE IF EXISTS ONLY public.supply_requests DROP CONSTRAINT IF EXISTS supply_requests_resource_category_id_fkey;
ALTER TABLE IF EXISTS ONLY public.supply_requests DROP CONSTRAINT IF EXISTS supply_requests_created_by_fkey;
ALTER TABLE IF EXISTS ONLY public.supply_requests DROP CONSTRAINT IF EXISTS supply_requests_approved_by_fkey;
ALTER TABLE IF EXISTS ONLY public.shipments DROP CONSTRAINT IF EXISTS shipments_vehicle_id_fkey;
ALTER TABLE IF EXISTS ONLY public.shipments DROP CONSTRAINT IF EXISTS shipments_to_warehouse_id_fkey;
ALTER TABLE IF EXISTS ONLY public.shipments DROP CONSTRAINT IF EXISTS shipments_tenant_id_fkey;
ALTER TABLE IF EXISTS ONLY public.shipments DROP CONSTRAINT IF EXISTS shipments_from_warehouse_id_fkey;
ALTER TABLE IF EXISTS ONLY public.shipment_refuels DROP CONSTRAINT IF EXISTS shipment_refuels_vehicle_id_fkey;
ALTER TABLE IF EXISTS ONLY public.shipment_refuels DROP CONSTRAINT IF EXISTS shipment_refuels_tenant_id_fkey;
ALTER TABLE IF EXISTS ONLY public.shipment_refuels DROP CONSTRAINT IF EXISTS shipment_refuels_shipment_id_fkey;
ALTER TABLE IF EXISTS ONLY public.shipment_refuels DROP CONSTRAINT IF EXISTS shipment_refuels_created_by_fkey;
ALTER TABLE IF EXISTS ONLY public.shipment_items DROP CONSTRAINT IF EXISTS shipment_items_shipment_id_fkey;
ALTER TABLE IF EXISTS ONLY public.shipment_items DROP CONSTRAINT IF EXISTS shipment_items_resource_id_fkey;
ALTER TABLE IF EXISTS ONLY public.shipment_items DROP CONSTRAINT IF EXISTS shipment_items_request_id_fkey;
ALTER TABLE IF EXISTS ONLY public.resources DROP CONSTRAINT IF EXISTS resources_warehouse_id_fkey;
ALTER TABLE IF EXISTS ONLY public.resources DROP CONSTRAINT IF EXISTS resources_unit_id_fkey;
ALTER TABLE IF EXISTS ONLY public.resources DROP CONSTRAINT IF EXISTS resources_tenant_id_fkey;
ALTER TABLE IF EXISTS ONLY public.resources DROP CONSTRAINT IF EXISTS resources_category_id_fkey;
ALTER TABLE IF EXISTS ONLY public.resource_categories DROP CONSTRAINT IF EXISTS resource_categories_tenant_id_fkey;
ALTER TABLE IF EXISTS ONLY public.resource_assignments DROP CONSTRAINT IF EXISTS resource_assignments_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.resource_assignments DROP CONSTRAINT IF EXISTS resource_assignments_resource_id_fkey;
ALTER TABLE IF EXISTS ONLY public.refresh_tokens DROP CONSTRAINT IF EXISTS refresh_tokens_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.notifications DROP CONSTRAINT IF EXISTS notifications_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.notifications DROP CONSTRAINT IF EXISTS notifications_tenant_id_fkey;
ALTER TABLE IF EXISTS ONLY public.maintenance_records DROP CONSTRAINT IF EXISTS maintenance_records_vehicle_id_fkey;
ALTER TABLE IF EXISTS ONLY public.maintenance_records DROP CONSTRAINT IF EXISTS maintenance_records_driver_id_fkey;
ALTER TABLE IF EXISTS ONLY public.invite_tokens DROP CONSTRAINT IF EXISTS invite_tokens_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.inventory_checks DROP CONSTRAINT IF EXISTS inventory_checks_warehouse_id_fkey;
ALTER TABLE IF EXISTS ONLY public.inventory_checks DROP CONSTRAINT IF EXISTS inventory_checks_created_by_fkey;
ALTER TABLE IF EXISTS ONLY public.inventory_check_items DROP CONSTRAINT IF EXISTS inventory_check_items_resource_id_fkey;
ALTER TABLE IF EXISTS ONLY public.inventory_check_items DROP CONSTRAINT IF EXISTS inventory_check_items_check_id_fkey;
ALTER TABLE IF EXISTS ONLY public.gps_locations DROP CONSTRAINT IF EXISTS gps_locations_vehicle_id_fkey;
ALTER TABLE IF EXISTS ONLY public.gps_locations DROP CONSTRAINT IF EXISTS gps_locations_unit_id_fkey;
ALTER TABLE IF EXISTS ONLY public.gps_locations DROP CONSTRAINT IF EXISTS gps_locations_tenant_id_fkey;
ALTER TABLE IF EXISTS ONLY public.geofences DROP CONSTRAINT IF EXISTS geofences_unit_id_fkey;
ALTER TABLE IF EXISTS ONLY public.geofences DROP CONSTRAINT IF EXISTS geofences_tenant_id_fkey;
ALTER TABLE IF EXISTS ONLY public.geofence_alerts DROP CONSTRAINT IF EXISTS geofence_alerts_vehicle_id_fkey;
ALTER TABLE IF EXISTS ONLY public.geofence_alerts DROP CONSTRAINT IF EXISTS geofence_alerts_geofence_id_fkey;
ALTER TABLE IF EXISTS ONLY public.fuel_records DROP CONSTRAINT IF EXISTS fuel_records_vehicle_id_fkey;
ALTER TABLE IF EXISTS ONLY public.fuel_records DROP CONSTRAINT IF EXISTS fuel_records_tenant_id_fkey;
ALTER TABLE IF EXISTS ONLY public.fuel_records DROP CONSTRAINT IF EXISTS fuel_records_created_by_fkey;
ALTER TABLE IF EXISTS ONLY public.contractor_requests DROP CONSTRAINT IF EXISTS contractor_requests_unit_id_fkey;
ALTER TABLE IF EXISTS ONLY public.contractor_requests DROP CONSTRAINT IF EXISTS contractor_requests_tenant_id_fkey;
ALTER TABLE IF EXISTS ONLY public.contractor_requests DROP CONSTRAINT IF EXISTS contractor_requests_target_warehouse_id_fkey;
ALTER TABLE IF EXISTS ONLY public.contractor_requests DROP CONSTRAINT IF EXISTS contractor_requests_taken_by_fkey;
ALTER TABLE IF EXISTS ONLY public.contractor_requests DROP CONSTRAINT IF EXISTS contractor_requests_created_by_fkey;
ALTER TABLE IF EXISTS ONLY public.contractor_memberships DROP CONSTRAINT IF EXISTS contractor_memberships_tenant_id_fkey;
ALTER TABLE IF EXISTS ONLY public.contractor_memberships DROP CONSTRAINT IF EXISTS contractor_memberships_decided_by_fkey;
ALTER TABLE IF EXISTS ONLY public.contractor_memberships DROP CONSTRAINT IF EXISTS contractor_memberships_contractor_id_fkey;
ALTER TABLE IF EXISTS ONLY public.audit_logs DROP CONSTRAINT IF EXISTS audit_logs_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.audit_logs DROP CONSTRAINT IF EXISTS audit_logs_tenant_id_fkey;
DROP INDEX IF EXISTS public.idx_warehouses_tenant;
DROP INDEX IF EXISTS public.idx_vehicles_warehouse;
DROP INDEX IF EXISTS public.idx_vehicles_tenant;
DROP INDEX IF EXISTS public.idx_vehicles_home_warehouse;
DROP INDEX IF EXISTS public.idx_vehicle_driver_history_vehicle_id;
DROP INDEX IF EXISTS public.idx_vehicle_driver_history_driver_id;
DROP INDEX IF EXISTS public.idx_vehicle_driver_history_assigned_at;
DROP INDEX IF EXISTS public.idx_users_tenant;
DROP INDEX IF EXISTS public.idx_users_status;
DROP INDEX IF EXISTS public.idx_users_role;
DROP INDEX IF EXISTS public.idx_units_tenant;
DROP INDEX IF EXISTS public.idx_tenants_slug;
DROP INDEX IF EXISTS public.idx_supply_requests_tenant;
DROP INDEX IF EXISTS public.idx_supply_requests_status;
DROP INDEX IF EXISTS public.idx_supply_requests_created_by;
DROP INDEX IF EXISTS public.idx_shipments_to;
DROP INDEX IF EXISTS public.idx_shipments_tenant;
DROP INDEX IF EXISTS public.idx_shipments_from;
DROP INDEX IF EXISTS public.idx_shipments_direction;
DROP INDEX IF EXISTS public.idx_shipment_refuels_vehicle;
DROP INDEX IF EXISTS public.idx_shipment_refuels_shipment;
DROP INDEX IF EXISTS public.idx_shipment_items_shipment;
DROP INDEX IF EXISTS public.idx_resources_warehouse_id;
DROP INDEX IF EXISTS public.idx_resources_unit;
DROP INDEX IF EXISTS public.idx_resources_tenant;
DROP INDEX IF EXISTS public.idx_resources_category;
DROP INDEX IF EXISTS public.idx_resource_categories_tenant;
DROP INDEX IF EXISTS public.idx_resource_assignments_user;
DROP INDEX IF EXISTS public.idx_resource_assignments_resource;
DROP INDEX IF EXISTS public.idx_refresh_tokens_user;
DROP INDEX IF EXISTS public.idx_refresh_tokens_hash;
DROP INDEX IF EXISTS public.idx_notifications_user_id;
DROP INDEX IF EXISTS public.idx_notifications_tenant_id;
DROP INDEX IF EXISTS public.idx_notifications_is_read;
DROP INDEX IF EXISTS public.idx_notifications_created_at;
DROP INDEX IF EXISTS public.idx_maintenance_records_vehicle_id;
DROP INDEX IF EXISTS public.idx_invite_tokens_hash;
DROP INDEX IF EXISTS public.idx_gps_locations_vehicle_time;
DROP INDEX IF EXISTS public.idx_gps_locations_vehicle;
DROP INDEX IF EXISTS public.idx_gps_locations_unit;
DROP INDEX IF EXISTS public.idx_gps_locations_timestamp;
DROP INDEX IF EXISTS public.idx_gps_locations_tenant;
DROP INDEX IF EXISTS public.idx_geofences_unit;
DROP INDEX IF EXISTS public.idx_geofences_tenant;
DROP INDEX IF EXISTS public.idx_geofence_alerts_vehicle;
DROP INDEX IF EXISTS public.idx_geofence_alerts_geofence;
DROP INDEX IF EXISTS public.idx_geofence_alerts_created;
DROP INDEX IF EXISTS public.idx_fuel_records_vehicle;
DROP INDEX IF EXISTS public.idx_fuel_records_tenant;
DROP INDEX IF EXISTS public.idx_contractor_requests_warehouse;
DROP INDEX IF EXISTS public.idx_contractor_requests_tenant;
DROP INDEX IF EXISTS public.idx_contractor_requests_status;
DROP INDEX IF EXISTS public.idx_contractor_memberships_tenant;
DROP INDEX IF EXISTS public.idx_contractor_memberships_status;
DROP INDEX IF EXISTS public.idx_contractor_memberships_contractor;
DROP INDEX IF EXISTS public.idx_audit_logs_tenant;
ALTER TABLE IF EXISTS ONLY public.warehouses DROP CONSTRAINT IF EXISTS warehouses_pkey;
ALTER TABLE IF EXISTS ONLY public.vehicles DROP CONSTRAINT IF EXISTS vehicles_tenant_plate_unique;
ALTER TABLE IF EXISTS ONLY public.vehicles DROP CONSTRAINT IF EXISTS vehicles_tenant_id_plate_number_key;
ALTER TABLE IF EXISTS ONLY public.vehicles DROP CONSTRAINT IF EXISTS vehicles_pkey;
ALTER TABLE IF EXISTS ONLY public.vehicle_driver_history DROP CONSTRAINT IF EXISTS vehicle_driver_history_pkey;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_tenant_username_unique;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_tenant_id_username_key;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_tenant_id_email_key;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_tenant_email_unique;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_pkey;
ALTER TABLE IF EXISTS ONLY public.units DROP CONSTRAINT IF EXISTS units_pkey;
ALTER TABLE IF EXISTS ONLY public.tenants DROP CONSTRAINT IF EXISTS tenants_slug_key;
ALTER TABLE IF EXISTS ONLY public.tenants DROP CONSTRAINT IF EXISTS tenants_pkey;
ALTER TABLE IF EXISTS ONLY public.tenants DROP CONSTRAINT IF EXISTS tenants_name_key;
ALTER TABLE IF EXISTS ONLY public.supply_requests DROP CONSTRAINT IF EXISTS supply_requests_pkey;
ALTER TABLE IF EXISTS ONLY public.shipments DROP CONSTRAINT IF EXISTS shipments_pkey;
ALTER TABLE IF EXISTS ONLY public.shipment_refuels DROP CONSTRAINT IF EXISTS shipment_refuels_pkey;
ALTER TABLE IF EXISTS ONLY public.shipment_items DROP CONSTRAINT IF EXISTS shipment_items_pkey;
ALTER TABLE IF EXISTS ONLY public.resources DROP CONSTRAINT IF EXISTS resources_pkey;
ALTER TABLE IF EXISTS ONLY public.resource_categories DROP CONSTRAINT IF EXISTS resource_categories_tenant_name_unique;
ALTER TABLE IF EXISTS ONLY public.resource_categories DROP CONSTRAINT IF EXISTS resource_categories_tenant_id_name_key;
ALTER TABLE IF EXISTS ONLY public.resource_categories DROP CONSTRAINT IF EXISTS resource_categories_pkey;
ALTER TABLE IF EXISTS ONLY public.resource_assignments DROP CONSTRAINT IF EXISTS resource_assignments_pkey;
ALTER TABLE IF EXISTS ONLY public.refresh_tokens DROP CONSTRAINT IF EXISTS refresh_tokens_pkey;
ALTER TABLE IF EXISTS ONLY public.notifications DROP CONSTRAINT IF EXISTS notifications_pkey;
ALTER TABLE IF EXISTS ONLY public.maintenance_records DROP CONSTRAINT IF EXISTS maintenance_records_pkey;
ALTER TABLE IF EXISTS ONLY public.invite_tokens DROP CONSTRAINT IF EXISTS invite_tokens_user_id_key;
ALTER TABLE IF EXISTS ONLY public.invite_tokens DROP CONSTRAINT IF EXISTS invite_tokens_pkey;
ALTER TABLE IF EXISTS ONLY public.inventory_checks DROP CONSTRAINT IF EXISTS inventory_checks_pkey;
ALTER TABLE IF EXISTS ONLY public.inventory_check_items DROP CONSTRAINT IF EXISTS inventory_check_items_pkey;
ALTER TABLE IF EXISTS ONLY public.gps_locations DROP CONSTRAINT IF EXISTS gps_locations_pkey;
ALTER TABLE IF EXISTS ONLY public.geofences DROP CONSTRAINT IF EXISTS geofences_pkey;
ALTER TABLE IF EXISTS ONLY public.geofence_alerts DROP CONSTRAINT IF EXISTS geofence_alerts_pkey;
ALTER TABLE IF EXISTS ONLY public.fuel_records DROP CONSTRAINT IF EXISTS fuel_records_pkey;
ALTER TABLE IF EXISTS ONLY public.contractor_requests DROP CONSTRAINT IF EXISTS contractor_requests_pkey;
ALTER TABLE IF EXISTS ONLY public.contractor_memberships DROP CONSTRAINT IF EXISTS contractor_memberships_pkey;
ALTER TABLE IF EXISTS ONLY public.contractor_memberships DROP CONSTRAINT IF EXISTS contractor_memberships_contractor_id_tenant_id_key;
ALTER TABLE IF EXISTS ONLY public.audit_logs DROP CONSTRAINT IF EXISTS audit_logs_pkey;
ALTER TABLE IF EXISTS public.units ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.gps_locations ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.geofences ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.geofence_alerts ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.audit_logs ALTER COLUMN id DROP DEFAULT;
DROP TABLE IF EXISTS public.warehouses;
DROP TABLE IF EXISTS public.vehicles;
DROP TABLE IF EXISTS public.vehicle_driver_history;
DROP TABLE IF EXISTS public.users;
DROP SEQUENCE IF EXISTS public.units_id_seq;
DROP TABLE IF EXISTS public.units;
DROP TABLE IF EXISTS public.tenants;
DROP TABLE IF EXISTS public.supply_requests;
DROP TABLE IF EXISTS public.shipments;
DROP TABLE IF EXISTS public.shipment_refuels;
DROP TABLE IF EXISTS public.shipment_items;
DROP TABLE IF EXISTS public.resources;
DROP TABLE IF EXISTS public.resource_categories;
DROP TABLE IF EXISTS public.resource_assignments;
DROP TABLE IF EXISTS public.refresh_tokens;
DROP TABLE IF EXISTS public.notifications;
DROP TABLE IF EXISTS public.maintenance_records;
DROP TABLE IF EXISTS public.invite_tokens;
DROP TABLE IF EXISTS public.inventory_checks;
DROP TABLE IF EXISTS public.inventory_check_items;
DROP SEQUENCE IF EXISTS public.gps_locations_id_seq;
DROP TABLE IF EXISTS public.gps_locations;
DROP SEQUENCE IF EXISTS public.geofences_id_seq;
DROP TABLE IF EXISTS public.geofences;
DROP SEQUENCE IF EXISTS public.geofence_alerts_id_seq;
DROP TABLE IF EXISTS public.geofence_alerts;
DROP TABLE IF EXISTS public.fuel_records;
DROP TABLE IF EXISTS public.contractor_requests;
DROP TABLE IF EXISTS public.contractor_memberships;
DROP SEQUENCE IF EXISTS public.audit_logs_id_seq;
DROP TABLE IF EXISTS public.audit_logs;
SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: audit_logs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.audit_logs (
    id integer NOT NULL,
    user_id uuid,
    action_type character varying(50) NOT NULL,
    entity_type character varying(50) NOT NULL,
    entity_id character varying(50),
    details text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    tenant_id uuid
);

ALTER TABLE ONLY public.audit_logs FORCE ROW LEVEL SECURITY;


--
-- Name: audit_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.audit_logs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: audit_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.audit_logs_id_seq OWNED BY public.audit_logs.id;


--
-- Name: contractor_memberships; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.contractor_memberships (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    contractor_id uuid NOT NULL,
    tenant_id uuid NOT NULL,
    status character varying(20) DEFAULT 'PENDING'::character varying NOT NULL,
    note text,
    requested_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    decided_at timestamp with time zone,
    decided_by uuid
);


--
-- Name: contractor_requests; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.contractor_requests (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    tenant_id uuid,
    created_by uuid NOT NULL,
    unit_id bigint,
    title character varying(255) NOT NULL,
    description text,
    status character varying(20) DEFAULT 'OPEN'::character varying NOT NULL,
    taken_by uuid,
    taken_at timestamp with time zone,
    completed_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deadline timestamp with time zone,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    target_warehouse_id uuid
);

ALTER TABLE ONLY public.contractor_requests FORCE ROW LEVEL SECURITY;


--
-- Name: fuel_records; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.fuel_records (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    tenant_id uuid,
    vehicle_id uuid NOT NULL,
    liters numeric(10,2) NOT NULL,
    odometer_km integer,
    record_type character varying(20) NOT NULL,
    created_by uuid,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    is_anomaly boolean DEFAULT false NOT NULL,
    anomaly_reason text,
    anomaly_excess_liters numeric(10,2) DEFAULT 0 NOT NULL
);

ALTER TABLE ONLY public.fuel_records FORCE ROW LEVEL SECURITY;


--
-- Name: geofence_alerts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.geofence_alerts (
    id bigint NOT NULL,
    vehicle_id uuid NOT NULL,
    geofence_id bigint NOT NULL,
    event_type character varying(20) NOT NULL,
    latitude double precision NOT NULL,
    longitude double precision NOT NULL,
    "timestamp" timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: geofence_alerts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.geofence_alerts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: geofence_alerts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.geofence_alerts_id_seq OWNED BY public.geofence_alerts.id;


--
-- Name: geofences; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.geofences (
    id bigint NOT NULL,
    unit_id bigint NOT NULL,
    name character varying(255) NOT NULL,
    latitude double precision NOT NULL,
    longitude double precision NOT NULL,
    radius double precision NOT NULL,
    type character varying(50) NOT NULL,
    active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    tenant_id uuid
);

ALTER TABLE ONLY public.geofences FORCE ROW LEVEL SECURITY;


--
-- Name: geofences_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.geofences_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: geofences_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.geofences_id_seq OWNED BY public.geofences.id;


--
-- Name: gps_locations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.gps_locations (
    id bigint NOT NULL,
    vehicle_id uuid NOT NULL,
    unit_id bigint NOT NULL,
    latitude double precision NOT NULL,
    longitude double precision NOT NULL,
    altitude double precision,
    speed double precision,
    heading double precision,
    accuracy double precision,
    "timestamp" timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    tenant_id uuid
);

ALTER TABLE ONLY public.gps_locations FORCE ROW LEVEL SECURITY;


--
-- Name: gps_locations_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.gps_locations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: gps_locations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.gps_locations_id_seq OWNED BY public.gps_locations.id;


--
-- Name: inventory_check_items; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.inventory_check_items (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    check_id uuid NOT NULL,
    resource_id uuid NOT NULL,
    book_quantity integer NOT NULL,
    actual_quantity integer,
    difference integer GENERATED ALWAYS AS ((actual_quantity - book_quantity)) STORED,
    verified_at timestamp with time zone
);


--
-- Name: inventory_checks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.inventory_checks (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    warehouse_id uuid NOT NULL,
    created_by uuid NOT NULL,
    status character varying(30) DEFAULT 'IN_PROGRESS'::character varying NOT NULL,
    started_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    completed_at timestamp with time zone,
    notes text
);


--
-- Name: invite_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.invite_tokens (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    token_hash character varying(255) NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    used_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: maintenance_records; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.maintenance_records (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    vehicle_id uuid NOT NULL,
    odometer_km integer NOT NULL,
    description text NOT NULL,
    performed_by character varying(255),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    cost_amount numeric(10,2) DEFAULT 0 NOT NULL,
    document_url text,
    driver_id uuid
);


--
-- Name: notifications; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.notifications (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    tenant_id uuid NOT NULL,
    type character varying(50) NOT NULL,
    title character varying(255) NOT NULL,
    message text NOT NULL,
    related_id uuid,
    is_read boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    read_at timestamp with time zone
);

ALTER TABLE ONLY public.notifications FORCE ROW LEVEL SECURITY;


--
-- Name: refresh_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.refresh_tokens (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    token_hash character varying(255) NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    revoked_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: resource_assignments; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.resource_assignments (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    resource_id uuid NOT NULL,
    user_id uuid NOT NULL,
    quantity integer DEFAULT 1 NOT NULL,
    status character varying(50) DEFAULT 'ACTIVE'::character varying NOT NULL,
    issued_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    returned_at timestamp without time zone,
    notes text
);


--
-- Name: resource_categories; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.resource_categories (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    tenant_id uuid,
    name character varying(100) NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: resources; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.resources (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    tenant_id uuid,
    category_id uuid NOT NULL,
    unit_id bigint,
    name character varying(255) NOT NULL,
    description text,
    quantity integer DEFAULT 0 NOT NULL,
    serial_number character varying(100),
    location character varying(255),
    condition character varying(20) DEFAULT 'NEW'::character varying NOT NULL,
    min_quantity integer DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    unit_type character varying(50) DEFAULT 'PCS'::character varying,
    warehouse_id uuid,
    weight_kg numeric(10,2) DEFAULT 1.00 NOT NULL,
    barcode character varying(255) DEFAULT ''::character varying,
    unit_price numeric(12,2) DEFAULT 0 NOT NULL
);

ALTER TABLE ONLY public.resources FORCE ROW LEVEL SECURITY;


--
-- Name: shipment_items; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.shipment_items (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    shipment_id uuid NOT NULL,
    resource_id uuid NOT NULL,
    quantity integer NOT NULL,
    request_id uuid,
    CONSTRAINT shipment_items_quantity_check CHECK ((quantity > 0))
);


--
-- Name: shipment_refuels; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.shipment_refuels (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    shipment_id uuid NOT NULL,
    vehicle_id uuid NOT NULL,
    liters numeric(10,2) NOT NULL,
    odometer_km integer,
    station_name character varying(200),
    cost_uah numeric(10,2),
    created_by uuid,
    tenant_id uuid,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT shipment_refuels_liters_check CHECK ((liters > (0)::numeric))
);

ALTER TABLE ONLY public.shipment_refuels FORCE ROW LEVEL SECURITY;


--
-- Name: shipments; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.shipments (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    from_warehouse_id uuid NOT NULL,
    to_warehouse_id uuid NOT NULL,
    vehicle_id uuid NOT NULL,
    priority character varying(20) DEFAULT 'NORMAL'::character varying NOT NULL,
    status character varying(30) DEFAULT 'DISPATCHED'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    tenant_id uuid,
    started_at timestamp with time zone,
    delivered_at timestamp with time zone,
    direction character varying(20),
    distance_km double precision DEFAULT 0,
    actual_km double precision DEFAULT 0
);

ALTER TABLE ONLY public.shipments FORCE ROW LEVEL SECURITY;


--
-- Name: supply_requests; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.supply_requests (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    tenant_id uuid,
    created_by uuid NOT NULL,
    resource_id uuid,
    quantity integer NOT NULL,
    status character varying(20) DEFAULT 'PENDING'::character varying NOT NULL,
    approved_by uuid,
    approved_at timestamp with time zone,
    comment text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    target_warehouse_id uuid,
    resource_name character varying(255),
    resource_category_id uuid
);

ALTER TABLE ONLY public.supply_requests FORCE ROW LEVEL SECURITY;


--
-- Name: tenants; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tenants (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(255) NOT NULL,
    slug character varying(100) NOT NULL,
    subscription_tier character varying(30) DEFAULT 'FREE'::character varying NOT NULL,
    subscription_expires_at timestamp with time zone,
    owner_email character varying(255),
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: units; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.units (
    id bigint NOT NULL,
    tenant_id uuid,
    parent_id bigint,
    name character varying(255) NOT NULL,
    unit_type character varying(20) NOT NULL,
    subscription_tier character varying(20) DEFAULT 'BASIC'::character varying
);

ALTER TABLE ONLY public.units FORCE ROW LEVEL SECURITY;


--
-- Name: units_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.units_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: units_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.units_id_seq OWNED BY public.units.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    tenant_id uuid,
    username character varying(100) NOT NULL,
    email character varying(255) NOT NULL,
    full_name character varying(255),
    phone character varying(50),
    password_hash text,
    role character varying(30) DEFAULT 'CONTRACTOR'::character varying NOT NULL,
    status character varying(20) DEFAULT 'PENDING'::character varying NOT NULL,
    unit_id bigint,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE ONLY public.users FORCE ROW LEVEL SECURITY;


--
-- Name: vehicle_driver_history; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.vehicle_driver_history (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    vehicle_id uuid,
    driver_id uuid,
    assigned_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: vehicles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.vehicles (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    tenant_id uuid,
    brand character varying(100) NOT NULL,
    model character varying(100),
    plate_number character varying(20) NOT NULL,
    status character varying(20) DEFAULT 'ACTIVE'::character varying NOT NULL,
    driver_id uuid,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    tank_capacity numeric(10,2) DEFAULT 0 NOT NULL,
    fuel_norm numeric(10,2) DEFAULT 0 NOT NULL,
    maintenance_interval_km integer DEFAULT 10000 NOT NULL,
    last_maintenance_odometer integer DEFAULT 0 NOT NULL,
    status_reason text,
    type character varying(50) DEFAULT 'VAN'::character varying,
    capacity_kg numeric(10,2) DEFAULT 1500.00 NOT NULL,
    current_warehouse_id uuid,
    home_warehouse_id uuid
);

ALTER TABLE ONLY public.vehicles FORCE ROW LEVEL SECURITY;


--
-- Name: warehouses; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.warehouses (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    unit_id bigint,
    name character varying(255) NOT NULL,
    location_type character varying(50) DEFAULT 'STATIONARY'::character varying,
    latitude double precision,
    longitude double precision,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    tenant_id uuid
);

ALTER TABLE ONLY public.warehouses FORCE ROW LEVEL SECURITY;


--
-- Name: audit_logs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.audit_logs ALTER COLUMN id SET DEFAULT nextval('public.audit_logs_id_seq'::regclass);


--
-- Name: geofence_alerts id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.geofence_alerts ALTER COLUMN id SET DEFAULT nextval('public.geofence_alerts_id_seq'::regclass);


--
-- Name: geofences id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.geofences ALTER COLUMN id SET DEFAULT nextval('public.geofences_id_seq'::regclass);


--
-- Name: gps_locations id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.gps_locations ALTER COLUMN id SET DEFAULT nextval('public.gps_locations_id_seq'::regclass);


--
-- Name: units id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.units ALTER COLUMN id SET DEFAULT nextval('public.units_id_seq'::regclass);


--
-- Data for Name: audit_logs; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.audit_logs (id, user_id, action_type, entity_type, entity_id, details, created_at, tenant_id) FROM stdin;
1	21833e68-101f-4378-a916-62c120a9f192	CREATE	TENANT	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	Створено нову організацію: Пійлівський Логістичний Центр	2026-04-23 00:11:12.212161+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
2	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-04-23 00:36:04.015106+00	d8729234-8e3b-41d9-83bc-9c725fe65838
3	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-23 10:22:32.327244+00	d8729234-8e3b-41d9-83bc-9c725fe65838
4	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-04-23 10:22:51.760412+00	d8729234-8e3b-41d9-83bc-9c725fe65838
5	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-04-23 10:34:09.949251+00	d8729234-8e3b-41d9-83bc-9c725fe65838
6	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-04-23 10:34:10.860489+00	d8729234-8e3b-41d9-83bc-9c725fe65838
7	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-23 11:00:36.72993+00	d8729234-8e3b-41d9-83bc-9c725fe65838
8	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-04-23 11:07:10.926804+00	d8729234-8e3b-41d9-83bc-9c725fe65838
9	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-04-23 11:07:12.853609+00	d8729234-8e3b-41d9-83bc-9c725fe65838
10	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-23 11:07:17.625776+00	d8729234-8e3b-41d9-83bc-9c725fe65838
11	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-23 11:08:05.125409+00	d8729234-8e3b-41d9-83bc-9c725fe65838
12	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-23 11:08:10.601289+00	d8729234-8e3b-41d9-83bc-9c725fe65838
13	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-04-23 11:08:55.986598+00	d8729234-8e3b-41d9-83bc-9c725fe65838
14	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-04-23 11:08:59.787311+00	d8729234-8e3b-41d9-83bc-9c725fe65838
15	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-04-23 11:09:24.258734+00	d8729234-8e3b-41d9-83bc-9c725fe65838
16	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-04-23 11:09:36.015843+00	d8729234-8e3b-41d9-83bc-9c725fe65838
17	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-23 11:29:21.919603+00	d8729234-8e3b-41d9-83bc-9c725fe65838
18	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-23 11:29:23.513523+00	d8729234-8e3b-41d9-83bc-9c725fe65838
19	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-23 11:29:33.653844+00	d8729234-8e3b-41d9-83bc-9c725fe65838
20	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-04-23 11:29:37.905841+00	d8729234-8e3b-41d9-83bc-9c725fe65838
21	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-04-23 11:29:39.15783+00	d8729234-8e3b-41d9-83bc-9c725fe65838
22	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-04-23 11:30:51.03176+00	d8729234-8e3b-41d9-83bc-9c725fe65838
23	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-04-23 11:30:52.168111+00	d8729234-8e3b-41d9-83bc-9c725fe65838
24	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-23 11:31:23.427755+00	d8729234-8e3b-41d9-83bc-9c725fe65838
25	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-04-23 11:32:04.618478+00	d8729234-8e3b-41d9-83bc-9c725fe65838
26	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-04-23 11:32:10.324985+00	d8729234-8e3b-41d9-83bc-9c725fe65838
27	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-04-23 11:57:16.202024+00	d8729234-8e3b-41d9-83bc-9c725fe65838
28	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-04-23 11:57:17.980978+00	d8729234-8e3b-41d9-83bc-9c725fe65838
29	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-23 11:58:37.399551+00	d8729234-8e3b-41d9-83bc-9c725fe65838
30	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-04-23 11:59:01.497325+00	d8729234-8e3b-41d9-83bc-9c725fe65838
31	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-04-23 11:59:13.714737+00	d8729234-8e3b-41d9-83bc-9c725fe65838
32	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-23 12:16:40.768349+00	d8729234-8e3b-41d9-83bc-9c725fe65838
33	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-23 12:16:43.497772+00	d8729234-8e3b-41d9-83bc-9c725fe65838
34	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-04-23 12:17:04.097573+00	d8729234-8e3b-41d9-83bc-9c725fe65838
35	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-04-23 12:17:05.227352+00	d8729234-8e3b-41d9-83bc-9c725fe65838
36	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-23 12:17:10.90909+00	d8729234-8e3b-41d9-83bc-9c725fe65838
37	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-23 12:18:48.92196+00	d8729234-8e3b-41d9-83bc-9c725fe65838
38	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-23 12:19:05.591974+00	d8729234-8e3b-41d9-83bc-9c725fe65838
39	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-23 12:19:14.272275+00	d8729234-8e3b-41d9-83bc-9c725fe65838
40	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-04-23 12:21:39.984271+00	d8729234-8e3b-41d9-83bc-9c725fe65838
41	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-04-23 12:21:43.159616+00	d8729234-8e3b-41d9-83bc-9c725fe65838
42	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-23 12:23:13.116949+00	d8729234-8e3b-41d9-83bc-9c725fe65838
43	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-23 12:32:45.879739+00	d8729234-8e3b-41d9-83bc-9c725fe65838
44	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-04-23 12:32:50.897792+00	d8729234-8e3b-41d9-83bc-9c725fe65838
45	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-04-23 12:32:52.043249+00	d8729234-8e3b-41d9-83bc-9c725fe65838
46	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-04-23 13:40:41.687269+00	d8729234-8e3b-41d9-83bc-9c725fe65838
47	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-04-23 13:56:35.792211+00	d8729234-8e3b-41d9-83bc-9c725fe65838
48	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-04-23 13:56:36.997145+00	d8729234-8e3b-41d9-83bc-9c725fe65838
49	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-23 13:57:04.290477+00	d8729234-8e3b-41d9-83bc-9c725fe65838
50	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-04-23 14:06:23.869231+00	d8729234-8e3b-41d9-83bc-9c725fe65838
51	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-04-23 14:06:25.22355+00	d8729234-8e3b-41d9-83bc-9c725fe65838
52	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-23 14:11:02.655803+00	d8729234-8e3b-41d9-83bc-9c725fe65838
53	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-04-23 14:29:07.127583+00	d8729234-8e3b-41d9-83bc-9c725fe65838
54	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-04-23 14:36:59.496729+00	d8729234-8e3b-41d9-83bc-9c725fe65838
55	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-04-23 14:37:00.394021+00	d8729234-8e3b-41d9-83bc-9c725fe65838
56	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-23 14:37:28.400991+00	d8729234-8e3b-41d9-83bc-9c725fe65838
57	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-23 14:52:40.922584+00	d8729234-8e3b-41d9-83bc-9c725fe65838
58	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-23 14:55:08.537025+00	d8729234-8e3b-41d9-83bc-9c725fe65838
59	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-04-23 15:03:52.559905+00	d8729234-8e3b-41d9-83bc-9c725fe65838
60	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-04-23 15:03:53.486606+00	d8729234-8e3b-41d9-83bc-9c725fe65838
61	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-04-23 15:06:46.564203+00	d8729234-8e3b-41d9-83bc-9c725fe65838
62	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-23 20:37:52.348885+00	d8729234-8e3b-41d9-83bc-9c725fe65838
63	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-04-23 20:38:12.519028+00	d8729234-8e3b-41d9-83bc-9c725fe65838
64	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-04-23 20:38:15.816477+00	d8729234-8e3b-41d9-83bc-9c725fe65838
65	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-23 20:46:30.152526+00	d8729234-8e3b-41d9-83bc-9c725fe65838
66	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-04-23 20:46:48.058595+00	d8729234-8e3b-41d9-83bc-9c725fe65838
67	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-04-23 20:46:49.488964+00	d8729234-8e3b-41d9-83bc-9c725fe65838
68	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-24 14:36:29.97881+00	d8729234-8e3b-41d9-83bc-9c725fe65838
69	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-04-24 14:37:17.243394+00	d8729234-8e3b-41d9-83bc-9c725fe65838
70	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-04-24 14:37:48.793669+00	d8729234-8e3b-41d9-83bc-9c725fe65838
71	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-04-24 14:37:56.765057+00	d8729234-8e3b-41d9-83bc-9c725fe65838
72	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-04-24 18:08:36.676532+00	d8729234-8e3b-41d9-83bc-9c725fe65838
73	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-04-24 18:08:50.245907+00	d8729234-8e3b-41d9-83bc-9c725fe65838
74	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-04-24 18:08:51.117435+00	d8729234-8e3b-41d9-83bc-9c725fe65838
75	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	1	Створено новий відділ: Центральний регіон управління постачанням	2026-04-25 16:01:23.422932+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
76	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	2	Створено новий відділ: Центральний регіон управління постачанням	2026-04-25 16:09:10.85915+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
77	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	3	Створено новий відділ: Київський головний розподільчий центр	2026-04-25 16:28:14.460338+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
78	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	4	Створено новий відділ: Відділ складської логістики	2026-04-25 16:28:47.478787+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
79	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	5	Створено новий відділ: Зміна денних комплектувальників	2026-04-25 16:29:12.715149+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
80	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	6	Створено новий відділ: Група операторів висотних штабелерів	2026-04-25 16:29:50.200157+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
81	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	7	Створено новий відділ: Диспетчерський департамент	2026-04-25 16:30:08.753118+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
82	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	8	Створено новий відділ: Команда моніторингу GPS	2026-04-25 16:30:32.109718+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
83	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	9	Створено новий відділ: Бригада водіїв великогабаритного транспорту	2026-04-25 16:30:52.628512+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
84	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	10	Створено новий відділ: Відділ зворотної логістики	2026-04-25 16:32:09.550754+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
85	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	11	Створено новий відділ: Група інспекції повернутих товарів	2026-04-25 16:32:26.260818+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
86	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	12	Створено новий відділ: Західний логістичний хаб	2026-04-25 16:32:42.574054+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
87	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	13	Створено новий відділ: Івано-Франківський транзитний термінал	2026-04-25 16:33:09.976533+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
88	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	14	Створено новий відділ: Відділ сортування та крос-докінгу	2026-04-25 16:35:56.728437+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
89	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	15	Створено новий відділ: Команда швидкого розвантаження	2026-04-25 16:36:28.332421+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
90	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	16	Створено новий відділ: Бригада пакувальників	2026-04-25 16:36:52.191256+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
91	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	17	Створено новий відділ: Відділ технічного забезпечення	2026-04-25 16:37:09.048695+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
92	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	18	Створено новий відділ: Мобільна ремонтна бригада автопарку	2026-04-25 16:37:42.845079+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
93	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	19	Створено новий відділ: Львівський митно-ліцензійний комплекс	2026-04-25 16:39:55.746905+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
94	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	20	Створено новий відділ: Відділ міжнародного транзиту	2026-04-25 16:40:13.170476+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
95	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	21	Створено новий відділ: Група оформлення супровідної документації	2026-04-25 16:40:32.131765+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
96	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	22	Створено новий відділ: Команда приймання імпорту	2026-04-25 16:40:52.446735+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
98	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	24	Створено новий відділ: Одеський портовий термінал	2026-04-25 16:41:31.312709+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
99	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	25	Створено новий відділ: Відділ контейнерних перевезень	2026-04-25 16:41:50.291812+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
101	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	27	Створено новий відділ: Група кріплення вантажів	2026-04-25 16:42:39.673438+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
104	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	30	Створено новий відділ: Миколаївський центр зберігання	2026-04-25 16:43:37.646942+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
105	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	31	Створено новий відділ: Відділ управління запасами	2026-04-25 16:44:09.502754+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
97	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	23	Створено новий відділ: Південний регіон	2026-04-25 16:41:17.158409+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
100	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	26	Створено новий відділ: Команда крановщиків	2026-04-25 16:42:22.092706+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
103	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	29	Створено новий відділ: Команда координації "море-залізниця"	2026-04-25 16:43:19.378713+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
106	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	32	Створено новий відділ: Команда інвентаризації	2026-04-25 16:45:08.251826+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
102	21833e68-101f-4378-a916-62c120a9f192	CREATE	UNIT	28	Створено новий відділ: Департамент мультимодальних перевезень	2026-04-25 16:42:59.220437+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
107	21833e68-101f-4378-a916-62c120a9f192	CREATE	USER	a2d66dd7-76b5-4baf-8818-bb130e29cf1a	Зареєстровано нового користувача: Ткаченко Олександр Іванович	2026-04-25 17:10:48.971771+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
108	21833e68-101f-4378-a916-62c120a9f192	CREATE	USER	acc38b6c-89a4-4403-8cb2-916107ff017a	Зареєстровано нового користувача: Коваленко Олена Вікторівна	2026-05-02 14:40:04.784708+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
109	21833e68-101f-4378-a916-62c120a9f192	CREATE	USER	78205adb-41af-4988-9cd4-6bee3337f215	Зареєстровано нового користувача: Бойко Наталія Сергіївна	2026-05-02 14:40:59.223102+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
110	21833e68-101f-4378-a916-62c120a9f192	CREATE	USER	bb96be7a-c867-48c2-a128-fb143bf6409a	Зареєстровано нового користувача: Кравченко Юлія Анатоліївна	2026-05-02 14:42:35.479137+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
111	21833e68-101f-4378-a916-62c120a9f192	CREATE	USER	744c25e6-ad82-4471-8cd6-8da1cb476d26	Зареєстровано нового користувача: Ткаченко Дмитро Олегович	2026-05-02 14:44:04.586732+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
112	21833e68-101f-4378-a916-62c120a9f192	CREATE	USER	44382f55-1ab8-496b-a63e-631171e55486	Зареєстровано нового користувача: Мельник Ігор Володимирович	2026-05-02 14:44:53.786321+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
113	21833e68-101f-4378-a916-62c120a9f192	CREATE	USER	1baed435-4416-444a-9939-724c8858ab78	Зареєстровано нового користувача: Гриценко Максим Іванович	2026-05-02 14:45:58.012879+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
114	21833e68-101f-4378-a916-62c120a9f192	CREATE	USER	8aa95ca9-5e34-4ab6-a033-be734bf87629	Зареєстровано нового користувача: Лисенко Олександр Петрович	2026-05-02 14:46:42.096073+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
115	21833e68-101f-4378-a916-62c120a9f192	CREATE	USER	6b43bc82-4a5b-46cc-af09-6ef4adecd95b	Зареєстровано нового користувача: Марченко Сергій Олександрович	2026-05-02 14:47:46.764943+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
116	21833e68-101f-4378-a916-62c120a9f192	CREATE	USER	b6635a10-c63a-4151-b9a4-ec34d35c99ee	Зареєстровано нового користувача: Руденко Віталій Миколайович	2026-05-02 14:48:33.001275+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
117	21833e68-101f-4378-a916-62c120a9f192	CREATE	USER	05eb7a5c-0413-4a6b-ad6a-5ad545e70a33	Зареєстровано нового користувача: Савченко Ірина Борисівна	2026-05-02 14:49:14.942435+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
118	21833e68-101f-4378-a916-62c120a9f192	CREATE	USER	c22f350c-6f18-4839-8d2d-79fd623e6a85	Зареєстровано нового користувача: Павленко Віктор Андрійович	2026-05-02 14:50:03.988322+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
119	21833e68-101f-4378-a916-62c120a9f192	CREATE	USER	39f7d374-63c3-4b75-b412-1a3e03a618b2	Зареєстровано нового користувача: Гаврилюк Анна Тарасівна	2026-05-02 14:50:48.937061+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
120	21833e68-101f-4378-a916-62c120a9f192	CREATE	USER	200577bb-7a24-4a4e-a4a8-7146fb112d37	Зареєстровано нового користувача: Романенко Тарас Юрійович	2026-05-02 14:59:29.473116+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
121	21833e68-101f-4378-a916-62c120a9f192	CREATE	USER	94579c72-8ab3-4faf-8d91-8d248206ebe6	Зареєстровано нового користувача: Сидоренко Олег Віталійович	2026-05-02 15:05:33.151737+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
122	21833e68-101f-4378-a916-62c120a9f192	CREATE	USER	26de8e1e-4ea0-46f8-907a-4c7a3390904b	Зареєстровано нового користувача: Поліщук Василь Степанович	2026-05-02 15:06:28.25025+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
123	21833e68-101f-4378-a916-62c120a9f192	CREATE	USER	0a4bc519-6dc3-4f19-9ec1-ecff82c91d03	Зареєстровано нового користувача: Козак Іван Павлович	2026-05-02 15:07:15.990204+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
124	21833e68-101f-4378-a916-62c120a9f192	CREATE	USER	b4bb3fff-a72f-437e-8d2f-48ed65d12eda	Зареєстровано нового користувача: Гончаренко Михайло Сергійович	2026-05-02 15:08:23.032367+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
125	21833e68-101f-4378-a916-62c120a9f192	CREATE	USER	dab88a27-f886-47b0-add8-8569e62942fa	Зареєстровано нового користувача: Литвин Андрій Анатолійович	2026-05-02 15:09:01.681809+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
126	21833e68-101f-4378-a916-62c120a9f192	CREATE	USER	4817a99c-f2c5-4157-8e3c-eedc2241fbc2	Зареєстровано нового користувача: Мороз Володимир Ігорович	2026-05-02 15:09:43.380737+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
127	21833e68-101f-4378-a916-62c120a9f192	CREATE	USER	2b080e9c-a604-4935-98c3-3bed5932986a	Зареєстровано нового користувача: Петренко Микола Васильович	2026-05-02 15:10:29.86306+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
128	21833e68-101f-4378-a916-62c120a9f192	CREATE	USER	3cd8d6bf-fad2-45f6-b6fb-cfce134df557	Зареєстровано нового користувача: Захарченко Ольга Павлівна	2026-05-02 15:11:27.378083+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
129	21833e68-101f-4378-a916-62c120a9f192	CREATE	USER	471f0149-166b-4f29-803b-3c910fa974c1	Зареєстровано нового користувача: Власенко Марина Юріївна	2026-05-02 15:12:09.127246+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
130	21833e68-101f-4378-a916-62c120a9f192	CREATE	WAREHOUSE	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	Створено склад	2026-05-02 15:20:16.918985+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
131	21833e68-101f-4378-a916-62c120a9f192	CREATE	WAREHOUSE	f64e8882-623d-41cd-8f72-2be30b783f8d	Створено склад	2026-05-02 15:20:46.472922+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
132	21833e68-101f-4378-a916-62c120a9f192	CREATE	WAREHOUSE	28f6a833-60fc-48bd-bb5c-30dfca2e3ace	Створено склад	2026-05-02 15:21:11.845873+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
133	21833e68-101f-4378-a916-62c120a9f192	CREATE	WAREHOUSE	b0981b66-2238-4044-ad6a-eef5917728f0	Створено склад	2026-05-02 15:21:43.472428+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
134	21833e68-101f-4378-a916-62c120a9f192	CREATE	WAREHOUSE	47bbf66e-ebd4-4c59-ad77-87314f7a302d	Створено склад	2026-05-02 15:22:18.53935+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
135	21833e68-101f-4378-a916-62c120a9f192	CREATE	WAREHOUSE	30efffba-f59c-4990-b1f6-e371a2dffca7	Створено склад	2026-05-02 15:23:07.889129+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
136	21833e68-101f-4378-a916-62c120a9f192	CREATE	WAREHOUSE	29918579-9237-4f65-9ce3-db58b825b86c	Створено склад	2026-05-02 15:23:35.534892+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
137	21833e68-101f-4378-a916-62c120a9f192	CREATE	WAREHOUSE	7a1639fb-d584-409a-89d9-e2b09fa379cd	Створено склад	2026-05-02 15:24:07.667547+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
138	21833e68-101f-4378-a916-62c120a9f192	CREATE	WAREHOUSE	b7de7c00-1ff9-4234-a909-40c7106ace7a	Створено склад	2026-05-02 15:24:41.943754+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
139	21833e68-101f-4378-a916-62c120a9f192	CREATE	WAREHOUSE	62ba10c5-14ee-4628-bfcc-cfb2f99d4e46	Створено склад	2026-05-02 15:25:05.621102+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
140	21833e68-101f-4378-a916-62c120a9f192	CREATE	VEHICLE	8355c89e-b072-4113-a6d2-adecc94a5a31	Додано новий автомобіль: MAN (KA1234AA)	2026-05-02 15:32:39.854582+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
141	21833e68-101f-4378-a916-62c120a9f192	CREATE	VEHICLE	5482a426-5bcf-48a2-82cd-d61441e76810	Додано новий автомобіль: Isuzu (AI4455EE)	2026-05-02 15:33:37.945336+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
142	21833e68-101f-4378-a916-62c120a9f192	CREATE	VEHICLE	dac43200-3bfe-47c8-8e60-17bf1ce6cf40	Додано новий автомобіль: Mercedes-Benz (BC5678CX)	2026-05-02 15:34:23.076803+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
143	21833e68-101f-4378-a916-62c120a9f192	CREATE	VEHICLE	48471a6b-ac6a-4327-b13c-dce6768ac4a0	Додано новий автомобіль: Renault (AB3322OO)	2026-05-02 15:34:56.409588+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
144	21833e68-101f-4378-a916-62c120a9f192	CREATE	VEHICLE	ea092cad-3147-4245-8ce6-ee7b8029806b	Додано новий автомобіль: Mitsubishi (CE9012MP)	2026-05-02 15:35:32.642153+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
301	21833e68-101f-4378-a916-62c120a9f192	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-05-07 14:09:41.94955+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
145	21833e68-101f-4378-a916-62c120a9f192	CREATE	CATEGORY	c69be86c-afeb-4892-91bb-043cffcb1487	Створено нову категорію: IT обладнання та зв'язок	2026-05-02 15:40:19.26213+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
146	21833e68-101f-4378-a916-62c120a9f192	CREATE	CATEGORY	6a273219-ea10-4ea7-aff4-9647696668f4	Створено нову категорію: Складський інвентар	2026-05-02 15:40:31.48327+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
147	21833e68-101f-4378-a916-62c120a9f192	CREATE	CATEGORY	733c1f11-0b6b-4f34-ab03-0c16946b9156	Створено нову категорію: Пакувальні та витратні матеріали	2026-05-02 15:40:43.474687+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
148	21833e68-101f-4378-a916-62c120a9f192	CREATE	CATEGORY	08a0ee32-49dc-4cec-aa93-1f95909488af	Створено нову категорію: Навігаційні пристрої	2026-05-02 15:40:58.690424+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
149	21833e68-101f-4378-a916-62c120a9f192	CREATE	RESOURCE	2a5ead75-9a60-4e0a-943c-bf8356e29539	Створено нову картку майна: Термінал збору даних (ТЗД) Zebra TC21	2026-05-02 15:42:09.750082+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
150	21833e68-101f-4378-a916-62c120a9f192	CREATE	RESOURCE	a231f49e-b368-4cfb-8d1c-6b7cd596f0f0	Створено нову картку майна: Гідравлічний візок (Рохля) 2.5т	2026-05-02 15:46:34.523687+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
151	21833e68-101f-4378-a916-62c120a9f192	CREATE	RESOURCE	a30253b6-4e38-43f9-9e93-2e8ab763d291	Створено нову картку майна: Стретч-плівка пакувальна (рулон 500мм, 20мкм)	2026-05-02 15:47:37.924386+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
152	21833e68-101f-4378-a916-62c120a9f192	CREATE	RESOURCE	974c8ddc-7326-44d7-aac1-aabc54aba0ce	Створено нову картку майна: GPS-трекер Teltonika FMB120	2026-05-02 15:48:34.720601+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
153	21833e68-101f-4378-a916-62c120a9f192	CREATE	WAREHOUSE	8b6d08d6-48bf-4a96-ba57-a4e6ec53fd55	Створено склад	2026-05-02 19:29:09.965228+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
154	21833e68-101f-4378-a916-62c120a9f192	DELETE	WAREHOUSE	8b6d08d6-48bf-4a96-ba57-a4e6ec53fd55	Видалено порожній склад	2026-05-02 20:03:49.948887+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
155	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	0c7e7fcd-23fe-4e08-9093-4cae1f867a8f	Створено нову заявку на забезпечення	2026-05-03 13:08:12.943673+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
156	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	8567c524-15d3-42ec-a50b-b61711722db9	Створено нову заявку на забезпечення	2026-05-03 13:22:29.6683+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
157	21833e68-101f-4378-a916-62c120a9f192	UPDATE	VEHICLE	ea092cad-3147-4245-8ce6-ee7b8029806b	Призначено нового водія на транспортний засіб	2026-05-03 13:23:41.080168+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
158	acc38b6c-89a4-4403-8cb2-916107ff017a	CREATE	USER	389997eb-96b7-4e37-a6ab-2db5692a6255	Зареєстровано нового користувача: Петренко Іван Миколайович	2026-05-03 13:30:14.323395+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
159	acc38b6c-89a4-4403-8cb2-916107ff017a	CREATE	USER	b1b097e6-452e-4b2a-b62f-ff8e543a59a0	Зареєстровано нового користувача: Ткаченко Олена Василівна	2026-05-03 13:31:18.115755+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
160	21833e68-101f-4378-a916-62c120a9f192	UPDATE	VEHICLE	48471a6b-ac6a-4327-b13c-dce6768ac4a0	Призначено нового водія на транспортний засіб	2026-05-03 13:41:00.66006+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
161	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-03 14:13:30.076988+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
162	21833e68-101f-4378-a916-62c120a9f192	CANCEL	SUPPLY_REQUEST	0c7e7fcd-23fe-4e08-9093-4cae1f867a8f	Скасовано власну заявку	2026-05-03 14:24:30.091868+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
163	21833e68-101f-4378-a916-62c120a9f192	READ	SHIPMENT	55a417c4-d5a8-4656-8a14-741f2c2c0273	Завантажено ТТН рейсу	2026-05-03 14:25:32.46272+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
164	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	55a417c4-d5a8-4656-8a14-741f2c2c0273	Вантаж прийнято на склад	2026-05-03 14:31:25.486957+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
165	21833e68-101f-4378-a916-62c120a9f192	READ	RESOURCE	a19cdafb-a938-4bf7-83ac-58f9df641b4f	Завантажено QR-код майна	2026-05-03 14:32:13.757755+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
166	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	c3e9e926-6d94-4d1a-aeba-592c3b2c5d19	Створено нову заявку на забезпечення	2026-05-03 15:02:10.218548+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
167	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	dcb85aa1-957f-4e44-8fa4-6ed940da4a5a	Створено нову заявку на забезпечення	2026-05-03 15:03:23.784453+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
168	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	706b07cc-1c61-43c6-a635-b241d4a968a8	Створено нову заявку на забезпечення	2026-05-03 15:05:17.679163+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
169	21833e68-101f-4378-a916-62c120a9f192	UPDATE	VEHICLE	48471a6b-ac6a-4327-b13c-dce6768ac4a0	Знято водія з транспортного засобу	2026-05-03 15:06:30.753257+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
170	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	527944b9-13c9-41f2-ba5b-8f1a03a2e404	Створено нову заявку на забезпечення	2026-05-03 15:06:43.138655+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
171	21833e68-101f-4378-a916-62c120a9f192	UPDATE	VEHICLE	48471a6b-ac6a-4327-b13c-dce6768ac4a0	Призначено нового водія на транспортний засіб	2026-05-03 15:07:31.059929+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
172	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-03 15:40:42.337816+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
173	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	8839bbd3-ea60-4e00-b25f-823449d5e004	Вантаж прийнято на склад	2026-05-03 15:49:26.045825+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
174	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	de0a3d0f-e7e3-4d96-ba02-2a5329367b3a	Створено нову заявку на забезпечення	2026-05-03 16:05:11.729831+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
175	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-03 16:05:32.152993+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
176	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	916f6d90-76f8-4cb1-b4a6-327cb2985d79	Створено нову заявку на забезпечення	2026-05-03 16:08:37.650481+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
177	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	5d03ed08-289e-4825-af7c-d9134ecd4107	Вантаж прийнято на склад	2026-05-03 16:09:39.352822+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
178	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-03 16:20:44.692945+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
179	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	a5504524-8653-44aa-8c38-a86832a5963e	Вантаж прийнято на склад	2026-05-03 16:24:14.656312+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
180	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	06f2fdf5-29f4-4beb-867d-47a46008d223	Створено нову заявку на забезпечення	2026-05-03 16:24:38.061397+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
181	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-03 16:25:00.835401+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
182	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	2f71fc41-c742-48b4-8d2f-9c7b6263c3e0	Вантаж прийнято на склад	2026-05-03 16:38:21.840005+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
183	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	eabaf774-962d-47c2-836d-466296a900fc	Створено нову заявку на забезпечення	2026-05-03 16:38:41.672117+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
184	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-03 16:38:53.577662+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
185	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	20bb818e-2a8c-4058-b865-5f52388a30ea	Вантаж прийнято на склад	2026-05-03 16:44:15.780641+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
186	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	ec0485d8-b5e4-4cc3-b75c-0b2f77aaaf39	Створено нову заявку на забезпечення	2026-05-03 16:44:29.182335+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
187	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-03 16:44:40.570261+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
188	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	0c54467b-ad81-4d0b-8554-2a2e44e46b22	Вантаж прийнято на склад	2026-05-03 16:58:00.976597+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
189	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	ad573105-8f14-4191-93d1-30b0153a79a7	Створено нову заявку на забезпечення	2026-05-03 16:58:23.936881+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
190	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-03 16:58:33.48684+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
191	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	ddb8e40b-3f73-419b-a82f-9cac667af6a9	Створено нову заявку на забезпечення	2026-05-03 17:12:05.181357+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
192	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-03 17:19:42.498128+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
193	389997eb-96b7-4e37-a6ab-2db5692a6255	UPDATE	SHIPMENT	798d6084-4572-441c-9f52-3f05892de225	Рейс розпочато (виїзд підтверджено)	2026-05-03 17:39:58.813008+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
194	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	798d6084-4572-441c-9f52-3f05892de225	Вантаж прийнято на склад	2026-05-03 17:45:33.301779+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
195	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	08441a0c-ad2c-48a8-bb4e-27dff31af672	Створено нову заявку на забезпечення	2026-05-03 17:49:18.005054+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
196	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-03 17:49:31.441684+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
197	389997eb-96b7-4e37-a6ab-2db5692a6255	UPDATE	SHIPMENT	7e631304-8a13-4a8b-9605-e83fbd90e834	Рейс розпочато (виїзд підтверджено)	2026-05-03 17:50:06.03447+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
198	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	7e631304-8a13-4a8b-9605-e83fbd90e834	Вантаж прийнято на склад	2026-05-03 17:50:16.327347+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
199	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	734691af-8adc-471c-9caa-c56eb6d52e1c	Створено нову заявку на забезпечення	2026-05-03 17:57:50.135403+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
200	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	910bb041-cc4d-4ffa-9543-859c30ee8429	Створено нову заявку на забезпечення	2026-05-03 17:58:00.874058+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
201	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-03 17:58:11.78287+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
202	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	9cd279eb-1626-4932-938c-63cd7f4829e2	Рейс розпочато (виїзд підтверджено)	2026-05-03 18:12:15.475402+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
203	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	9cd279eb-1626-4932-938c-63cd7f4829e2	Вантаж прийнято на склад	2026-05-03 18:12:18.686418+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
204	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	49bdeabb-c3b2-4d5d-a928-aa9bba50f4e9	Створено нову заявку на забезпечення	2026-05-03 19:16:25.525781+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
205	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	0f3ccaca-b253-4be2-8082-04a71e357b51	Створено нову заявку на забезпечення	2026-05-03 19:21:55.861438+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
206	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	41ba70af-d560-44b9-94cd-b56d29265920	Створено нову заявку на забезпечення	2026-05-03 19:28:54.309236+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
207	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	2bf81fd0-182d-42b6-b507-345fd0b2e16f	Створено нову заявку на забезпечення	2026-05-03 19:35:32.10284+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
208	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-03 19:49:15.094041+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
209	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	b05dd831-264e-4bed-a7a9-e3e3d06c1393	Рейс розпочато (виїзд підтверджено)	2026-05-03 19:52:17.868112+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
210	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	b05dd831-264e-4bed-a7a9-e3e3d06c1393	Вантаж прийнято на склад	2026-05-03 19:52:19.474518+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
211	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-03 20:01:53.852898+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
212	389997eb-96b7-4e37-a6ab-2db5692a6255	UPDATE	SHIPMENT	8e1ca805-ed59-4b6f-bded-84ad66ae0138	Рейс розпочато (виїзд підтверджено)	2026-05-03 20:02:12.40762+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
213	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	09413e76-b72e-4749-b11c-9547e16e6dee	Створено нову заявку на забезпечення	2026-05-03 20:05:28.763399+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
214	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	c5f0423f-c748-4ee8-90e6-4b6cf113ffc0	Створено нову заявку на забезпечення	2026-05-03 20:13:06.107804+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
215	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	8e1ca805-ed59-4b6f-bded-84ad66ae0138	Вантаж прийнято на склад	2026-05-03 20:28:14.91686+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
216	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	40947ff2-d081-4520-913d-a89527332190	Створено нову заявку на забезпечення	2026-05-03 21:09:38.34251+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
217	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	617e59a8-38a6-4cf7-8c5d-12c13a09ae2d	Створено нову заявку на забезпечення	2026-05-03 21:21:13.848649+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
218	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-03 21:38:16.994788+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
219	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	e948e8a6-1151-4086-8a32-41f1efd5350e	Створено нову заявку на забезпечення	2026-05-03 21:44:25.406095+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
220	21833e68-101f-4378-a916-62c120a9f192	SLA_VIOLATION	REQUEST	617e59a8-38a6-4cf7-8c5d-12c13a09ae2d	Заявка очікує підтвердження понад 24 годин. Автоматична ескалація.	2026-05-05 06:21:38.09617+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
221	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	44878184-a03c-4e93-9b39-bd9ba1b650e6	Створено нову заявку на забезпечення	2026-05-05 12:33:10.104924+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
222	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	34af91e3-95bb-4633-a8ef-3ce9a7336400	Рейс розпочато (виїзд підтверджено)	2026-05-05 12:33:32.201594+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
223	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	34af91e3-95bb-4633-a8ef-3ce9a7336400	Вантаж прийнято на склад	2026-05-05 12:33:33.093674+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
224	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-05 12:33:45.659479+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
225	389997eb-96b7-4e37-a6ab-2db5692a6255	UPDATE	SHIPMENT	704ce237-6254-415d-aadb-7089f825e48b	Рейс розпочато (виїзд підтверджено)	2026-05-05 12:33:57.546291+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
226	389997eb-96b7-4e37-a6ab-2db5692a6255	UPDATE	SHIPMENT	704ce237-6254-415d-aadb-7089f825e48b	Вантаж прийнято на склад	2026-05-05 12:34:04.993935+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
227	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-05-05 12:37:48.240177+00	d8729234-8e3b-41d9-83bc-9c725fe65838
228	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-05-05 12:37:49.548037+00	d8729234-8e3b-41d9-83bc-9c725fe65838
229	33a91fae-503d-45d2-aa42-4fdf78fcc983	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-05-05 12:37:51.063903+00	d8729234-8e3b-41d9-83bc-9c725fe65838
230	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-05 18:32:10.682972+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
231	21833e68-101f-4378-a916-62c120a9f192	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-05-06 15:55:50.361987+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
232	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	230457f2-442e-4eb3-80b0-7feb166ec89d	Створено нову заявку на забезпечення	2026-05-06 17:45:40.183868+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
233	389997eb-96b7-4e37-a6ab-2db5692a6255	UPDATE	SHIPMENT	8ff169fe-d4d5-41d5-ab14-355008cc9a0e	Рейс розпочато (виїзд підтверджено)	2026-05-06 17:46:49.354766+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
234	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	8ff169fe-d4d5-41d5-ab14-355008cc9a0e	Вантаж прийнято на склад	2026-05-06 21:04:26.406123+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
235	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-06 21:04:40.195161+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
236	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	51af7bd7-5340-4970-ac2a-61bba55cf810	Рейс розпочато (виїзд підтверджено)	2026-05-06 21:05:19.113504+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
485	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-30 15:15:16.458037+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
237	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	51af7bd7-5340-4970-ac2a-61bba55cf810	Вантаж прийнято на склад	2026-05-06 21:05:19.84114+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
238	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	8bcb06cf-5b5e-4e2c-b8b9-a766449e409a	Створено нову заявку на забезпечення	2026-05-06 21:05:36.442174+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
239	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-06 21:07:03.695459+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
240	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	570a39e6-bc00-4985-aa2c-bf24ed2d4375	Рейс розпочато (виїзд підтверджено)	2026-05-06 21:13:06.213938+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
241	\N	SLA_VIOLATION	REQUEST	4db21c95-d272-47db-bf82-2d2951ec94e8	Заявка очікує підтвердження понад 24 годин. Автоматична ескалація.	2026-05-07 10:33:13.797069+00	d8729234-8e3b-41d9-83bc-9c725fe65838
242	\N	SLA_VIOLATION	REQUEST	615ce782-705f-438d-bbd5-2dc0b8ab2197	Заявка очікує підтвердження понад 24 годин. Автоматична ескалація.	2026-05-07 10:33:13.797842+00	d8729234-8e3b-41d9-83bc-9c725fe65838
243	\N	SLA_VIOLATION	REQUEST	2c147665-3a3a-4382-b527-b1d390987f8a	Заявка очікує підтвердження понад 24 годин. Автоматична ескалація.	2026-05-07 10:33:13.800579+00	d8729234-8e3b-41d9-83bc-9c725fe65838
244	\N	SLA_VIOLATION	REQUEST	ac87ffc4-f314-41f2-bd62-9d565a5f3245	Заявка очікує підтвердження понад 24 годин. Автоматична ескалація.	2026-05-07 10:33:13.802393+00	d8729234-8e3b-41d9-83bc-9c725fe65838
245	\N	SLA_VIOLATION	REQUEST	b12516c5-3629-4ad7-99b5-4870bcc66049	Заявка очікує підтвердження понад 24 годин. Автоматична ескалація.	2026-05-07 10:33:13.802385+00	d8729234-8e3b-41d9-83bc-9c725fe65838
246	\N	SLA_VIOLATION	REQUEST	d47a2dc8-f76d-4ac8-a59a-fd0a4c095fa1	Заявка очікує підтвердження понад 24 годин. Автоматична ескалація.	2026-05-07 10:33:13.803302+00	d8729234-8e3b-41d9-83bc-9c725fe65838
247	\N	SLA_VIOLATION	REQUEST	ffe46dff-1709-4d2c-83db-97f2c39848a2	Заявка очікує підтвердження понад 24 годин. Автоматична ескалація.	2026-05-07 10:33:13.807941+00	d8729234-8e3b-41d9-83bc-9c725fe65838
248	\N	SLA_VIOLATION	REQUEST	cd0e794a-b14b-4596-bad2-c5e17679a95f	Заявка очікує підтвердження понад 24 годин. Автоматична ескалація.	2026-05-07 10:33:13.809282+00	d8729234-8e3b-41d9-83bc-9c725fe65838
249	\N	SLA_VIOLATION	REQUEST	4d7a61fd-c66a-4c7f-a9e4-838ccebc1434	Заявка очікує підтвердження понад 24 годин. Автоматична ескалація.	2026-05-07 10:33:13.809348+00	d8729234-8e3b-41d9-83bc-9c725fe65838
250	\N	SLA_VIOLATION	REQUEST	45878390-d17e-478a-962b-eede6daacb09	Заявка очікує підтвердження понад 24 годин. Автоматична ескалація.	2026-05-07 10:33:13.810632+00	d8729234-8e3b-41d9-83bc-9c725fe65838
251	3e59277f-0c1a-4aa3-973c-cf3db1e19497	SLA_VIOLATION	REQUEST	8cca3a73-5564-4f47-a697-4539ba62045d	Заявка очікує підтвердження понад 24 годин. Автоматична ескалація.	2026-05-07 10:34:13.791478+00	0997831c-654f-471b-934f-cedafbc54ea5
253	601e6d16-947b-4eb0-a67c-5304edafedd8	SLA_VIOLATION	REQUEST	b967d3a4-1450-41d5-8d3b-144e75530da8	Заявка очікує підтвердження понад 24 годин. Автоматична ескалація.	2026-05-07 10:34:13.804026+00	0997831c-654f-471b-934f-cedafbc54ea5
252	601e6d16-947b-4eb0-a67c-5304edafedd8	SLA_VIOLATION	REQUEST	3b670dbb-11af-4eeb-9e09-fd743b30f14b	Заявка очікує підтвердження понад 24 годин. Автоматична ескалація.	2026-05-07 10:34:13.803942+00	0997831c-654f-471b-934f-cedafbc54ea5
262	21833e68-101f-4378-a916-62c120a9f192	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-05-07 10:40:21.374142+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
254	612e483e-e25c-4d11-89b7-8327b8257c4f	SLA_VIOLATION	REQUEST	2bf1df38-c5d2-4a55-84de-d1f1a8b064c7	Заявка очікує підтвердження понад 24 годин. Автоматична ескалація.	2026-05-07 10:34:13.803814+00	0997831c-654f-471b-934f-cedafbc54ea5
255	135fe7a0-2fc0-4b17-9e22-4d1c32d9658b	SLA_VIOLATION	REQUEST	aa978e04-66ca-417f-aff4-5b1a92334e11	Заявка очікує підтвердження понад 24 годин. Автоматична ескалація.	2026-05-07 10:34:13.803762+00	0997831c-654f-471b-934f-cedafbc54ea5
256	601e6d16-947b-4eb0-a67c-5304edafedd8	SLA_VIOLATION	REQUEST	592d1bb0-aca0-4696-ba3a-afaf0a91ee38	Заявка очікує підтвердження понад 24 годин. Автоматична ескалація.	2026-05-07 10:34:13.803997+00	0997831c-654f-471b-934f-cedafbc54ea5
257	612e483e-e25c-4d11-89b7-8327b8257c4f	SLA_VIOLATION	REQUEST	154ed374-3d5b-405a-9404-41b867503f2f	Заявка очікує підтвердження понад 24 годин. Автоматична ескалація.	2026-05-07 10:34:13.812197+00	0997831c-654f-471b-934f-cedafbc54ea5
258	3e59277f-0c1a-4aa3-973c-cf3db1e19497	SLA_VIOLATION	REQUEST	33fa3eb0-15e1-4bf4-aec8-276d818742bd	Заявка очікує підтвердження понад 24 годин. Автоматична ескалація.	2026-05-07 10:34:13.812633+00	0997831c-654f-471b-934f-cedafbc54ea5
259	e169450b-acde-48e4-b894-fc4538fa26a9	SLA_VIOLATION	REQUEST	f7d3c5d0-45bd-4d23-88e5-e7bfa3765b65	Заявка очікує підтвердження понад 24 годин. Автоматична ескалація.	2026-05-07 10:34:13.813808+00	0997831c-654f-471b-934f-cedafbc54ea5
260	135fe7a0-2fc0-4b17-9e22-4d1c32d9658b	SLA_VIOLATION	REQUEST	f92a5020-2bf8-4a34-bfeb-a82d047bd09d	Заявка очікує підтвердження понад 24 годин. Автоматична ескалація.	2026-05-07 10:34:13.813797+00	0997831c-654f-471b-934f-cedafbc54ea5
261	8e4524be-8fb5-46a9-abb3-9d132f4fcda9	CREATE	SUPPLY_REQUEST	96929ab1-f11f-431a-91b9-e0ead5785672	Створено нову заявку на забезпечення	2026-05-07 10:39:03.65594+00	0997831c-654f-471b-934f-cedafbc54ea5
263	21833e68-101f-4378-a916-62c120a9f192	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-05-07 10:42:05.382706+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
264	21833e68-101f-4378-a916-62c120a9f192	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-05-07 10:43:09.759754+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
265	21833e68-101f-4378-a916-62c120a9f192	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-05-07 10:43:24.758788+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
266	21833e68-101f-4378-a916-62c120a9f192	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-05-07 10:43:25.874503+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
267	21833e68-101f-4378-a916-62c120a9f192	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-05-07 10:43:44.052892+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
268	21833e68-101f-4378-a916-62c120a9f192	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-05-07 11:58:29.78534+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
269	21833e68-101f-4378-a916-62c120a9f192	CREATE	SMART_REPLENISH		Сформовано автоматичні заявки на поповнення (1 шт)	2026-05-07 11:58:45.602235+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
270	21833e68-101f-4378-a916-62c120a9f192	CREATE	RESOURCE	e3ab8171-e065-4eb0-8f06-280de9694a57	Створено нову картку майна: Провід для заряджання	2026-05-07 12:00:07.642449+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
271	21833e68-101f-4378-a916-62c120a9f192	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-05-07 12:00:14.492946+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
272	21833e68-101f-4378-a916-62c120a9f192	CREATE	SMART_REPLENISH		Сформовано автоматичні заявки на поповнення (1 шт)	2026-05-07 12:00:21.98903+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
273	21833e68-101f-4378-a916-62c120a9f192	EXPORT	FUEL		Експорт звіту витрат пального (Excel)	2026-05-07 12:02:32.763176+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
274	21833e68-101f-4378-a916-62c120a9f192	EXPORT	INVENTORY		Експорт звіту залишків на складах (Excel)	2026-05-07 12:02:36.127502+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
275	21833e68-101f-4378-a916-62c120a9f192	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-05-07 12:04:16.461464+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
276	21833e68-101f-4378-a916-62c120a9f192	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-05-07 12:04:21.201753+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
277	21833e68-101f-4378-a916-62c120a9f192	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-05-07 12:04:26.797098+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
278	21833e68-101f-4378-a916-62c120a9f192	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-05-07 12:06:10.993651+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
279	21833e68-101f-4378-a916-62c120a9f192	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-05-07 12:06:13.40724+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
280	21833e68-101f-4378-a916-62c120a9f192	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-05-07 12:06:15.27614+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
281	21833e68-101f-4378-a916-62c120a9f192	CREATE	RESOURCE	1cd798c4-9488-4365-be01-ca8e904ad6c9	Створено нову картку майна: Електросамокат	2026-05-07 12:29:19.880796+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
282	21833e68-101f-4378-a916-62c120a9f192	CREATE	SMART_REPLENISH		Сформовано автоматичні заявки на поповнення (1 шт)	2026-05-07 12:29:37.093802+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
283	21833e68-101f-4378-a916-62c120a9f192	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-05-07 12:30:21.299048+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
284	21833e68-101f-4378-a916-62c120a9f192	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-05-07 12:30:30.420684+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
285	21833e68-101f-4378-a916-62c120a9f192	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-05-07 12:30:38.528188+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
286	21833e68-101f-4378-a916-62c120a9f192	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-05-07 12:30:49.76217+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
287	21833e68-101f-4378-a916-62c120a9f192	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-05-07 12:47:42.136856+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
288	21833e68-101f-4378-a916-62c120a9f192	WRITE_OFF	RESOURCE	a30253b6-4e38-43f9-9e93-2e8ab763d291	Списання майна зі складу	2026-05-07 12:47:51.856133+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
289	21833e68-101f-4378-a916-62c120a9f192	CREATE	RESOURCE	9dd794cb-2a09-48f9-91e9-90822118a266	Створено нову картку майна: Ноут	2026-05-07 12:49:30.791048+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
290	21833e68-101f-4378-a916-62c120a9f192	CREATE	SMART_REPLENISH		Сформовано автоматичні заявки на поповнення (1 шт)	2026-05-07 12:49:37.975516+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
291	21833e68-101f-4378-a916-62c120a9f192	WRITE_OFF	RESOURCE	a30253b6-4e38-43f9-9e93-2e8ab763d291	Списання майна зі складу	2026-05-07 12:50:44.640243+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
292	21833e68-101f-4378-a916-62c120a9f192	WRITE_OFF	RESOURCE	a30253b6-4e38-43f9-9e93-2e8ab763d291	Списання майна зі складу	2026-05-07 12:51:03.10098+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
293	21833e68-101f-4378-a916-62c120a9f192	WRITE_OFF	RESOURCE	a30253b6-4e38-43f9-9e93-2e8ab763d291	Списання майна зі складу	2026-05-07 12:55:24.775323+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
294	21833e68-101f-4378-a916-62c120a9f192	WRITE_OFF	RESOURCE	a30253b6-4e38-43f9-9e93-2e8ab763d291	Списання майна зі складу	2026-05-07 13:09:47.60585+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
295	21833e68-101f-4378-a916-62c120a9f192	WRITE_OFF	RESOURCE	a30253b6-4e38-43f9-9e93-2e8ab763d291	Списання майна зі складу	2026-05-07 13:10:19.634215+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
296	21833e68-101f-4378-a916-62c120a9f192	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-05-07 13:21:22.347701+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
297	21833e68-101f-4378-a916-62c120a9f192	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-05-07 13:21:32.466361+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
298	21833e68-101f-4378-a916-62c120a9f192	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-05-07 13:25:02.263926+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
299	21833e68-101f-4378-a916-62c120a9f192	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-05-07 13:25:13.826444+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
300	21833e68-101f-4378-a916-62c120a9f192	VIEW	KPI		Переглянув розширену аналітику (KPI)	2026-05-07 13:25:18.005963+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
302	21833e68-101f-4378-a916-62c120a9f192	VIEW	MAINTENANCE_SCHEDULE		Переглянув прогноз обслуговування машин	2026-05-07 14:09:56.332638+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
303	21833e68-101f-4378-a916-62c120a9f192	VIEW	FUEL_ANOMALIES		Переглянув аналіз аномалій у витраті палива	2026-05-07 14:09:57.374909+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
304	21833e68-101f-4378-a916-62c120a9f192	WRITE_OFF	RESOURCE	a30253b6-4e38-43f9-9e93-2e8ab763d291	Списання майна зі складу	2026-05-08 13:23:41.830834+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
305	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	f713b5c9-223b-45ad-94de-1122170b287f	Створено нову заявку на забезпечення	2026-05-17 00:34:53.85642+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
306	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	570a39e6-bc00-4985-aa2c-bf24ed2d4375	Вантаж прийнято на склад	2026-05-17 00:35:14.217378+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
307	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-17 00:42:45.71509+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
308	389997eb-96b7-4e37-a6ab-2db5692a6255	UPDATE	SHIPMENT	59069426-28fa-45c7-8d5b-54503c66165c	Рейс розпочато (виїзд підтверджено)	2026-05-17 00:43:50.176865+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
309	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	4b6531da-0415-46b1-a90c-ab6bf6a48c9e	Створено нову заявку на забезпечення	2026-05-17 11:44:42.363547+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
310	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	53db2ad3-fd81-4f75-974d-8a633b8a86f5	Створено нову заявку на забезпечення	2026-05-22 08:42:51.302548+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
311	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	59069426-28fa-45c7-8d5b-54503c66165c	Вантаж прийнято на склад	2026-05-27 21:14:32.613011+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
312	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-27 21:21:41.249891+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
313	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	af0bd3a8-de7d-4f6f-ac96-db08cd1fb0c4	Рейс розпочато (виїзд підтверджено)	2026-05-27 21:21:54.311307+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
314	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	af0bd3a8-de7d-4f6f-ac96-db08cd1fb0c4	Вантаж прийнято на склад	2026-05-27 21:21:55.887003+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
315	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	a1d6a6a0-3749-4303-a799-30d09b346245	Створено нову заявку на забезпечення	2026-05-27 21:48:57.399653+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
316	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-27 21:49:18.573574+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
317	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	9e717db0-b052-4834-a3f9-f8ec4279e19f	Рейс розпочато (виїзд підтверджено)	2026-05-27 21:49:24.753623+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
318	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	9e717db0-b052-4834-a3f9-f8ec4279e19f	Вантаж прийнято на склад	2026-05-27 21:49:34.580311+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
319	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	6d4c35df-6eeb-4a42-a59c-fa023ba79096	Створено нову заявку на забезпечення	2026-05-27 21:51:24.49904+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
320	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-27 21:52:09.690924+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
321	21833e68-101f-4378-a916-62c120a9f192	CREATE	VEHICLE	1565aa02-d4ff-464c-bc5b-c96f63fd9be1	Додано новий автомобіль: Ford (AT1234BH)	2026-05-28 07:44:53.450135+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
322	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	9bbd6df5-403d-4b06-828f-42b3979ec189	Рейс розпочато (виїзд підтверджено)	2026-05-28 07:45:35.97943+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
323	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	9bbd6df5-403d-4b06-828f-42b3979ec189	Вантаж прийнято на склад	2026-05-28 07:45:38.4161+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
324	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	7a304c19-ae39-4673-802b-e7019c06ed81	Створено нову заявку на забезпечення	2026-05-28 07:47:02.757729+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
325	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT	smart-batch	Smart Розподіл: створено 1 рейсів	2026-05-28 08:12:02.159067+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
326	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	cfb9dcda-8034-49a7-82f1-8744022ee386	Рейс розпочато (виїзд підтверджено)	2026-05-28 08:14:43.635836+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
327	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	cfb9dcda-8034-49a7-82f1-8744022ee386	Вантаж прийнято на склад	2026-05-28 08:27:33.77023+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
328	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	92d41b0a-acda-4500-90c4-cf6d697d0dfc	Створено нову заявку на забезпечення	2026-05-28 08:28:57.593861+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
329	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT	smart-batch	Smart Розподіл: створено 1 рейсів	2026-05-28 08:29:12.952812+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
330	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	2c639beb-e3c1-4d18-a672-651deca05a3b	Створено нову заявку на забезпечення	2026-05-28 08:42:25.266998+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
331	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	87cd4448-e066-4c10-a0f8-1e1faa25b499	Створено нову заявку на забезпечення	2026-05-28 12:11:29.752345+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
332	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT	smart-batch	Smart Розподіл: створено 1 рейсів	2026-05-28 12:11:44.030946+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
333	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	a1b425ea-0777-4d0d-965b-5082a819db50	Рейс розпочато (виїзд підтверджено)	2026-05-28 12:12:08.041763+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
334	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	103a7304-ed77-4f89-9ba5-1746383c6cba	Рейс розпочато (виїзд підтверджено)	2026-05-28 12:12:10.574877+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
335	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	a1b425ea-0777-4d0d-965b-5082a819db50	Вантаж прийнято на склад	2026-05-28 12:12:11.351409+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
336	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	103a7304-ed77-4f89-9ba5-1746383c6cba	Вантаж прийнято на склад	2026-05-28 12:12:12.124152+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
337	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	19031419-655e-44bf-9863-c68c70b42ad0	Створено нову заявку на забезпечення	2026-05-28 12:12:41.207893+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
338	21833e68-101f-4378-a916-62c120a9f192	UPDATE	VEHICLE	dac43200-3bfe-47c8-8e60-17bf1ce6cf40	Призначено нового водія на транспортний засіб	2026-05-28 12:21:28.443516+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
339	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT	smart-batch	Smart Розподіл: створено 1 рейсів	2026-05-28 12:23:37.877222+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
340	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	a190d65f-0532-4580-945a-a1efbedb5d4a	Рейс розпочато (виїзд підтверджено)	2026-05-28 12:23:46.037125+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
341	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	a190d65f-0532-4580-945a-a1efbedb5d4a	Вантаж прийнято на склад	2026-05-28 12:23:46.906867+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
342	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	75f07327-df11-46d1-bbe3-08f089791c8b	Створено нову заявку на забезпечення	2026-05-28 12:24:02.652211+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
343	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-28 12:24:09.017775+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
344	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	ad4a7da0-ca3b-4e94-a08a-160c9357931d	Рейс розпочато (виїзд підтверджено)	2026-05-28 12:24:15.012087+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
345	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	ad4a7da0-ca3b-4e94-a08a-160c9357931d	Вантаж прийнято на склад	2026-05-28 12:24:15.639448+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
346	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	c3de2651-f9f3-4b48-8ed7-e6474b5bf1d3	Створено нову заявку на забезпечення	2026-05-28 13:07:22.591914+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
347	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT	smart-batch	Smart Розподіл: створено 1 рейсів	2026-05-28 13:07:34.033274+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
348	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	7aee9783-f7fc-41e5-b369-b295796f2366	Рейс розпочато (виїзд підтверджено)	2026-05-28 13:11:26.962804+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
349	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	7cf6f81a-664b-47b6-949e-35385a3d3270	Створено нову заявку на забезпечення	2026-05-28 13:12:00.6734+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
350	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT	smart-batch	Smart Розподіл: створено 1 рейсів	2026-05-28 13:12:06.06385+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
351	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	7a0d4a54-d7b7-4fde-8942-55a04d5577d8	Створено нову заявку на забезпечення	2026-05-28 13:12:35.279472+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
352	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT	smart-batch	Smart Розподіл: створено 1 рейсів	2026-05-28 13:12:40.475451+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
353	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	2f064ad8-9087-4fc4-8b64-73afe2c25925	Створено нову заявку на забезпечення	2026-05-28 13:16:13.15371+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
354	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	ebd81f67-bf05-4329-8837-6a6c19dcd454	Рейс розпочато (виїзд підтверджено)	2026-05-28 14:31:05.376731+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
355	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	2925f911-f566-47fc-9ffe-a9488127b7fc	Рейс розпочато (виїзд підтверджено)	2026-05-28 14:31:06.410174+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
356	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	7aee9783-f7fc-41e5-b369-b295796f2366	Вантаж прийнято на склад	2026-05-28 14:31:06.949992+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
357	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	2925f911-f566-47fc-9ffe-a9488127b7fc	Вантаж прийнято на склад	2026-05-28 14:31:07.297475+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
358	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	ebd81f67-bf05-4329-8837-6a6c19dcd454	Вантаж прийнято на склад	2026-05-28 14:31:08.17686+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
359	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	2fb378e7-d3be-41ba-bfae-fe84e7c985f2	Створено нову заявку на забезпечення	2026-05-28 14:31:36.572949+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
360	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT	smart-batch	Smart Розподіл: створено 1 рейсів	2026-05-28 14:31:54.588438+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
361	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	c60622f1-fea1-4c60-b80c-4a83ed2be749	Створено нову заявку на забезпечення	2026-05-28 14:32:21.082543+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
362	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT	smart-batch	Smart Розподіл: створено 1 рейсів	2026-05-28 14:32:28.74039+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
363	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	298b0158-fd71-4224-ba82-67f7d753e108	Рейс розпочато (виїзд підтверджено)	2026-05-28 14:36:52.433691+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
364	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	298b0158-fd71-4224-ba82-67f7d753e108	Вантаж прийнято на склад	2026-05-28 14:40:29.493474+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
365	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	d7f8a04d-1d50-40cf-9be4-1cabc516f15b	Рейс розпочато (виїзд підтверджено)	2026-05-28 14:40:31.328966+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
366	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	d7f8a04d-1d50-40cf-9be4-1cabc516f15b	Вантаж прийнято на склад	2026-05-28 14:40:32.024334+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
367	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	a1484e78-3114-4773-9cd1-fd7dfdb3fe54	Створено нову заявку на забезпечення	2026-05-28 15:30:57.285558+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
368	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT	smart-batch	Smart Розподіл: створено 1 рейсів	2026-05-28 15:44:49.607514+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
369	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	9722f33c-6e4d-45ea-a1e8-346b318c1dda	Створено нову заявку на забезпечення	2026-05-28 15:46:03.106096+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
370	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	ac220e0a-de5e-45ed-8635-8a0eaae954aa	Створено нову заявку на забезпечення	2026-05-28 15:46:14.627008+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
371	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT	smart-batch	Smart Розподіл: створено 1 рейсів	2026-05-28 15:46:22.785247+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
372	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT	smart-batch	Smart Розподіл: створено 1 рейсів	2026-05-28 15:46:27.869907+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
373	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	850d5366-3794-474d-b863-bf6beaeb2d36	Рейс розпочато (виїзд підтверджено)	2026-05-28 15:47:27.874488+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
374	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	850d5366-3794-474d-b863-bf6beaeb2d36	Вантаж прийнято на склад	2026-05-28 15:47:43.667318+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
375	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	0b3ab0c2-0ef3-42a3-a6ab-9bd4bcc9cd2f	Створено нову заявку на забезпечення	2026-05-28 15:48:18.616274+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
376	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT	smart-batch	Smart Розподіл: створено 1 рейсів	2026-05-28 15:48:26.02302+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
377	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	a55619ee-2d84-4572-8a92-f93f506fa5e7	Створено нову заявку на забезпечення	2026-05-28 16:14:24.159069+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
378	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT	smart-batch	Smart Розподіл: створено 1 рейсів	2026-05-28 16:14:32.914044+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
379	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	aedbe475-34da-4a74-801e-a3752873ef8f	Рейс розпочато (виїзд підтверджено)	2026-05-28 16:15:30.532629+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
380	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	afc7246a-4a83-44d2-81d7-278ebab69699	Рейс розпочато (виїзд підтверджено)	2026-05-28 16:15:31.071996+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
381	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	d07d5ba7-5520-4782-bd60-a296058a22e4	Рейс розпочато (виїзд підтверджено)	2026-05-28 16:15:31.921751+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
382	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	d5c49bb6-339a-4a58-9ff0-c360087993b0	Рейс розпочато (виїзд підтверджено)	2026-05-28 16:15:32.981451+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
383	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	d5c49bb6-339a-4a58-9ff0-c360087993b0	Вантаж прийнято на склад	2026-05-28 16:15:33.499735+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
384	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	d07d5ba7-5520-4782-bd60-a296058a22e4	Вантаж прийнято на склад	2026-05-28 16:15:34.151135+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
385	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	afc7246a-4a83-44d2-81d7-278ebab69699	Вантаж прийнято на склад	2026-05-28 16:15:34.760244+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
386	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	aedbe475-34da-4a74-801e-a3752873ef8f	Вантаж прийнято на склад	2026-05-28 16:15:35.325053+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
387	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	d83944c7-221b-4499-96ea-d22dd53b9e4e	Створено нову заявку на забезпечення	2026-05-28 16:16:45.107646+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
388	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	e4bbb1a5-8e7b-46b3-be05-0c103fa57a30	Створено нову заявку на забезпечення	2026-05-28 20:30:28.973322+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
389	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	960f3ecb-6ce2-4224-88b9-abf8241825f6	Створено нову заявку на забезпечення	2026-05-28 20:30:38.639876+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
390	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT	smart-batch	Smart Розподіл: створено 1 рейсів	2026-05-28 20:58:22.185705+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
391	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT	smart-batch	Smart Розподіл: створено 1 рейсів	2026-05-28 20:58:32.593477+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
392	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT	smart-batch	Smart Розподіл: створено 1 рейсів	2026-05-28 20:58:42.268232+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
393	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	a58fc4df-9977-4289-ac15-b2defd44a741	Створено нову заявку на забезпечення	2026-05-28 20:59:11.167057+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
394	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT	smart-batch	Smart Розподіл: створено 1 рейсів	2026-05-28 21:29:50.258721+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
395	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	627c4181-535a-466b-8af3-c896dbe98eb5	Створено нову заявку на забезпечення	2026-05-28 21:30:03.850829+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
396	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT	smart-batch	Smart Розподіл: створено 1 рейсів	2026-05-28 21:30:08.980852+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
397	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	59f2d6d8-3472-46e7-8879-502ade5630c4	Створено нову заявку на забезпечення	2026-05-28 21:30:28.060473+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
398	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	b44f2f01-9882-48ab-bc7c-8ea84cc30f9c	Рейс розпочато (виїзд підтверджено)	2026-05-28 21:51:56.667676+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
399	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	b44f2f01-9882-48ab-bc7c-8ea84cc30f9c	Вантаж прийнято на склад	2026-05-29 07:31:30.566848+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
400	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	4a06e118-4ae2-4ede-99fa-128dc83cd510	Створено нову заявку на забезпечення	2026-05-29 07:31:52.897943+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
401	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT	smart-batch	Smart Розподіл: створено 1 рейсів	2026-05-29 07:33:02.553814+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
402	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	1d2b8d5a-2259-4641-a0d3-275a1dcd3647	Створено нову заявку на забезпечення	2026-05-29 07:33:38.94399+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
403	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	3feb3240-4a1b-4f0d-834c-037b62a8efa9	Створено нову заявку на забезпечення	2026-05-29 07:51:51.007281+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
404	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	7160f8f6-77a7-4864-8d43-1ef6fbdde1b5	Рейс розпочато (виїзд підтверджено)	2026-05-29 07:55:24.010227+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
405	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	7160f8f6-77a7-4864-8d43-1ef6fbdde1b5	Вантаж прийнято на склад	2026-05-29 07:55:25.996453+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
406	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	f7e1401b-8706-4858-9fa8-e4aa4cbba26e	Рейс розпочато (виїзд підтверджено)	2026-05-29 07:55:26.572129+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
407	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	f7e1401b-8706-4858-9fa8-e4aa4cbba26e	Вантаж прийнято на склад	2026-05-29 07:55:28.782402+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
408	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	b1f16ff5-0fb9-4380-b573-d5ffde39ba59	Рейс розпочато (виїзд підтверджено)	2026-05-29 07:55:29.608114+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
409	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	b1f16ff5-0fb9-4380-b573-d5ffde39ba59	Вантаж прийнято на склад	2026-05-29 07:55:32.742993+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
410	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	59153e98-ae5d-45a0-8e0f-f3fd2fe0435d	Рейс розпочато (виїзд підтверджено)	2026-05-29 07:55:33.954138+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
411	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	59153e98-ae5d-45a0-8e0f-f3fd2fe0435d	Вантаж прийнято на склад	2026-05-29 07:55:36.198033+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
412	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	a1093fd8-3282-473d-8713-bf8500ad2c5b	Рейс розпочато (виїзд підтверджено)	2026-05-29 07:55:37.374422+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
413	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	a1093fd8-3282-473d-8713-bf8500ad2c5b	Вантаж прийнято на склад	2026-05-29 07:55:39.677129+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
414	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	d40732b3-c42c-495b-8b33-fe4d1a7d80e1	Створено нову заявку на забезпечення	2026-05-29 07:56:07.278471+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
415	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	79e917e5-660f-4d4c-8ba8-32ef9d37c229	Створено нову заявку на забезпечення	2026-05-29 07:56:32.925798+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
416	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	4a8f8ca5-995a-4253-9536-23333720b7c4	Створено нову заявку на забезпечення	2026-05-29 07:56:50.212087+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
417	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT	smart-batch	Smart Розподіл: створено 1 рейсів	2026-05-29 07:56:59.708069+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
418	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT	smart-batch	Smart Розподіл: створено 1 рейсів	2026-05-29 07:57:08.431956+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
419	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT	smart-batch	Smart Розподіл: створено 1 рейсів	2026-05-29 07:57:30.519811+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
420	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	0a5ba720-77d5-4fc6-8ee1-5ade154ef999	Створено нову заявку на забезпечення	2026-05-29 07:59:12.19656+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
421	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	1c444910-7c78-409c-852a-11d186c31578	Створено нову заявку на забезпечення	2026-05-29 08:14:48.374765+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
422	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT	smart-batch	Smart Розподіл: створено 1 рейсів	2026-05-29 08:35:34.732916+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
423	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT	smart-batch	Smart Розподіл: створено 1 рейсів	2026-05-29 09:18:19.055367+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
424	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	ecc2d04a-15bb-439f-a7fc-9f070cff31ff	Створено нову заявку на забезпечення	2026-05-29 09:33:12.197918+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
425	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	84a00a91-1625-43db-83ac-78a9d92c0962	Рейс розпочато (виїзд підтверджено)	2026-05-29 09:36:47.36336+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
426	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	84a00a91-1625-43db-83ac-78a9d92c0962	Вантаж прийнято на склад	2026-05-29 09:36:48.897641+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
427	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	d3ed0394-d72e-4d4b-8023-72fd2cf72b91	Рейс розпочато (виїзд підтверджено)	2026-05-29 09:36:50.045586+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
428	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	d3ed0394-d72e-4d4b-8023-72fd2cf72b91	Вантаж прийнято на склад	2026-05-29 09:36:51.885244+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
429	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	4d1d8708-f792-45fe-9932-21de489833bf	Рейс розпочато (виїзд підтверджено)	2026-05-29 09:36:52.660725+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
430	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	4d1d8708-f792-45fe-9932-21de489833bf	Вантаж прийнято на склад	2026-05-29 09:36:54.599594+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
431	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	a15772fc-002a-4e30-b90a-52c850191858	Рейс розпочато (виїзд підтверджено)	2026-05-29 09:36:55.573585+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
432	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	a15772fc-002a-4e30-b90a-52c850191858	Вантаж прийнято на склад	2026-05-29 09:36:57.569634+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
433	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	9329b2e7-5c15-447d-805b-b1f8d6d3bdf4	Рейс розпочато (виїзд підтверджено)	2026-05-29 09:36:58.360675+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
434	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	9329b2e7-5c15-447d-805b-b1f8d6d3bdf4	Вантаж прийнято на склад	2026-05-29 09:37:00.111233+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
435	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	67f8edf9-03ed-4387-b855-2f14a2a55035	Створено нову заявку на забезпечення	2026-05-29 09:37:30.420804+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
436	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	b3ef8912-7968-4d8f-8057-4a7b5109ee4d	Створено нову заявку на забезпечення	2026-05-29 12:13:24.251686+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
437	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	3482bf32-472b-48af-9db0-c40e10294bb8	Створено нову заявку на забезпечення	2026-05-29 12:15:03.6609+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
438	21833e68-101f-4378-a916-62c120a9f192	UPDATE	VEHICLE	48471a6b-ac6a-4327-b13c-dce6768ac4a0	Знято водія з транспортного засобу	2026-05-29 12:47:31.410717+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
439	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-29 14:59:26.62308+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
440	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	6c320299-71b2-4811-9d1a-0210fa33d864	Рейс розпочато (виїзд підтверджено)	2026-05-29 14:59:39.613712+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
441	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	6c320299-71b2-4811-9d1a-0210fa33d864	Вантаж прийнято на склад	2026-05-29 21:31:58.194855+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
442	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	4ac04bc3-7209-4b80-bf41-de836563be7b	Створено нову заявку на забезпечення	2026-05-29 21:33:29.592835+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
443	21833e68-101f-4378-a916-62c120a9f192	CREATE	FUEL_RECORD	ea092cad-3147-4245-8ce6-ee7b8029806b	Додано запис про пальне (Тип: REFUEL, Літри: 50.00)	2026-05-29 21:34:25.665461+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
444	21833e68-101f-4378-a916-62c120a9f192	CREATE	FUEL_RECORD	1565aa02-d4ff-464c-bc5b-c96f63fd9be1	Додано запис про пальне (Тип: REFUEL, Літри: 30.00)	2026-05-29 21:35:35.128444+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
445	21833e68-101f-4378-a916-62c120a9f192	CREATE	FUEL_RECORD	1565aa02-d4ff-464c-bc5b-c96f63fd9be1	Додано запис про пальне (Тип: REFUEL, Літри: 10.20)	2026-05-29 21:42:26.753495+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
446	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-29 21:50:55.314609+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
447	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	772a04e3-29ca-470c-b8f9-a54221d88de1	Рейс розпочато (виїзд підтверджено)	2026-05-29 21:51:14.030029+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
448	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	fe351879-2b3a-49e0-80f7-e7b9ffb8203b	Створено нову заявку на забезпечення	2026-05-29 22:05:03.824832+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
449	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-29 22:05:13.473383+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
450	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	d8e0dc0f-a992-4782-8fae-1b50b90185e1	Рейс розпочато (виїзд підтверджено)	2026-05-29 22:05:45.411071+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
451	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	d8e0dc0f-a992-4782-8fae-1b50b90185e1	Вантаж прийнято на склад	2026-05-29 22:33:47.219404+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
452	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	772a04e3-29ca-470c-b8f9-a54221d88de1	Вантаж прийнято на склад	2026-05-29 22:33:50.864359+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
453	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	2c5f4d28-7147-4d40-9638-78b4d6e4347b	Створено нову заявку на забезпечення	2026-05-29 22:34:14.943267+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
454	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-29 22:34:39.186216+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
455	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	c189fea9-e9ac-4722-8f9a-799037a4035c	Рейс розпочато (виїзд підтверджено)	2026-05-29 22:40:44.731223+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
456	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	c189fea9-e9ac-4722-8f9a-799037a4035c	Вантаж прийнято на склад	2026-05-29 22:48:29.490786+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
457	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	79f1e0c2-1e2c-4418-9c3f-c61698e3b81e	Створено нову заявку на забезпечення	2026-05-29 22:48:44.116965+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
458	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-29 22:48:54.598233+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
459	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	d0e4a50b-3888-4c08-85b7-bbc9b2bec7d9	Рейс розпочато (виїзд підтверджено)	2026-05-29 22:49:17.027158+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
460	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	d0e4a50b-3888-4c08-85b7-bbc9b2bec7d9	Вантаж прийнято на склад	2026-05-29 23:07:38.839976+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
461	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	4fb347a1-a5b1-426a-9aee-80122f4b719e	Створено нову заявку на забезпечення	2026-05-29 23:07:49.670692+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
462	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-29 23:07:58.060607+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
463	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	a1d3a9d8-3a6d-411d-a1a1-c8dfbe0e71ad	Рейс розпочато (виїзд підтверджено)	2026-05-29 23:08:14.848424+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
464	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	e7b7c7b0-1407-46f8-accf-c178408bbf17	Створено нову заявку на забезпечення	2026-05-29 23:24:05.936722+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
465	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-29 23:24:20.793291+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
466	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	9bd63b6e-2560-48cf-bf88-54b89904271b	Рейс розпочато (виїзд підтверджено)	2026-05-29 23:24:25.896881+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
467	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	a7742526-a89e-4b0a-93fc-3925d487ca3f	Створено нову заявку на забезпечення	2026-05-30 10:16:42.917913+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
468	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	9bd63b6e-2560-48cf-bf88-54b89904271b	Вантаж прийнято на склад	2026-05-30 10:17:25.327603+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
469	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-30 10:17:43.779017+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
470	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	00eaa206-0780-4c97-89c7-5245593a747e	Рейс розпочато (виїзд підтверджено)	2026-05-30 10:17:50.925381+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
471	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	407cf304-0c12-4fb7-beaf-b1164247e366	Створено нову заявку на забезпечення	2026-05-30 12:47:19.879808+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
472	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	00eaa206-0780-4c97-89c7-5245593a747e	Вантаж прийнято. Пробіг: 8 км (GPS не записувався)	2026-05-30 12:48:02.855107+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
473	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-30 12:48:21.563285+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
474	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	fb8f5382-6ba2-44ed-9170-c0ad9ead2d86	Рейс розпочато (виїзд підтверджено)	2026-05-30 12:48:31.02495+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
475	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	fb8f5382-6ba2-44ed-9170-c0ad9ead2d86	Вантаж прийнято. Пробіг: 8 км (GPS не записувався)	2026-05-30 12:52:08.817297+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
476	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	46cb862f-6806-41be-ae82-8a003a2c53c7	Створено нову заявку на забезпечення	2026-05-30 12:53:52.995648+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
477	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-30 12:54:09.878569+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
478	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	84ac7e4e-76c9-4d76-a84e-87fa5adc65b8	Рейс розпочато (виїзд підтверджено)	2026-05-30 12:54:41.393772+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
479	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	84ac7e4e-76c9-4d76-a84e-87fa5adc65b8	Вантаж прийнято. Пробіг: 16 км (GPS не записувався)	2026-05-30 13:12:25.826027+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
480	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	570b02c0-74b7-4c64-b541-33f15c12e0bc	Створено нову заявку на забезпечення	2026-05-30 13:12:55.670292+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
481	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-30 13:13:26.050326+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
482	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	57663c7f-c1e6-4e4c-b12a-b7273eb2f2be	Рейс розпочато (виїзд підтверджено)	2026-05-30 13:15:06.631711+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
483	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	a6853850-640e-464c-8f9e-a31544524197	Створено нову заявку на забезпечення	2026-05-30 15:12:33.835246+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
484	21833e68-101f-4378-a916-62c120a9f192	CREATE	FUEL_RECORD	dac43200-3bfe-47c8-8e60-17bf1ce6cf40	Додано запис про пальне (Тип: REFUEL, Літри: 40.00)	2026-05-30 15:14:53.163244+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
486	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	c61746a4-50f6-48d6-907e-7bfffdd17776	Рейс розпочато (виїзд підтверджено)	2026-05-30 15:16:26.768639+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
487	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	c61746a4-50f6-48d6-907e-7bfffdd17776	Вантаж прийнято. Пробіг: 16.0 км (GPS не записувався)	2026-05-30 15:16:51.548709+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
488	21833e68-101f-4378-a916-62c120a9f192	CREATE	FUEL_RECORD	5482a426-5bcf-48a2-82cd-d61441e76810	Додано запис про пальне (Тип: REFUEL, Літри: 40.00)	2026-05-30 21:05:11.517015+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
489	21833e68-101f-4378-a916-62c120a9f192	CREATE	SUPPLY_REQUEST	e26f128f-43d2-4bc0-a1d5-cf7748986488	Створено нову заявку на забезпечення	2026-05-30 21:30:23.738507+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
490	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT		Сформовано новий рейс	2026-05-30 21:30:32.025746+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
491	21833e68-101f-4378-a916-62c120a9f192	CREATE	SHIPMENT_REFUEL	a1d3a9d8-3a6d-411d-a1a1-c8dfbe0e71ad	Дозаправка в дорозі: 5.0 л	2026-05-30 22:02:20.752163+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
492	21833e68-101f-4378-a916-62c120a9f192	UPDATE	SHIPMENT	a1d3a9d8-3a6d-411d-a1a1-c8dfbe0e71ad	Вантаж прийнято. Пробіг: 16.0 км (GPS не записувався)	2026-05-31 22:53:35.221521+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
493	21833e68-101f-4378-a916-62c120a9f192	WRITE_OFF	RESOURCE	a30253b6-4e38-43f9-9e93-2e8ab763d291	Списання майна зі складу	2026-05-31 22:57:26.081837+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
494	21833e68-101f-4378-a916-62c120a9f192	UPDATE	RESOURCE		Майно видано користувачу	2026-05-31 23:41:24.176477+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
495	21833e68-101f-4378-a916-62c120a9f192	REPORT	RESOURCE_ASSIGNMENT	d226fa15-98d1-4d58-9f6f-a76f8ea0a47c	Подано рапорт щодо майна: BROKEN	2026-06-01 00:13:52.406706+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
496	21833e68-101f-4378-a916-62c120a9f192	CREATE	CONTRACTOR_REQUEST	2c26917d-a600-4fe4-bf36-eda3fc3eec98	Створено новий зовнішній запит	2026-06-08 09:50:16.419718+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
497	21833e68-101f-4378-a916-62c120a9f192	CREATE	CONTRACTOR_REQUEST	c59e8f03-b667-45a3-9673-775c037ddff0	Створено новий зовнішній запит	2026-06-08 14:41:25.637242+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
498	dd8a8b66-805e-4b3f-92e5-ea158a88f421	UPDATE	CONTRACTOR_MEMBERSHIP	a3bf70c9-c5d0-4f43-9025-fd97cf8278c7	Підрядника схвалено для співпраці	2026-06-08 15:53:31.864025+00	0997831c-654f-471b-934f-cedafbc54ea5
499	\N	UPDATE	CONTRACTOR_REQUEST	263d8252-a281-441a-ba33-bd40758cd092	Заявку взято в роботу	2026-06-08 15:53:32.128757+00	d8729234-8e3b-41d9-83bc-9c725fe65838
502	dd8a8b66-805e-4b3f-92e5-ea158a88f421	UPDATE	CONTRACTOR_MEMBERSHIP	95fbd39b-daf8-4532-9265-63cbfd74d0c8	Підрядника схвалено для співпраці	2026-06-08 16:37:28.62512+00	0997831c-654f-471b-934f-cedafbc54ea5
500	\N	CREATE	CONTRACTOR_MEMBERSHIP	0997831c-654f-471b-934f-cedafbc54ea5	Підрядник надіслав заявку на співпрацю	2026-06-08 16:36:37.760873+00	d8729234-8e3b-41d9-83bc-9c725fe65838
501	\N	CREATE	CONTRACTOR_MEMBERSHIP	0997831c-654f-471b-934f-cedafbc54ea5	Підрядник надіслав заявку на співпрацю	2026-06-08 16:36:37.762746+00	d8729234-8e3b-41d9-83bc-9c725fe65838
503	\N	CREATE	CONTRACTOR_MEMBERSHIP	0997831c-654f-471b-934f-cedafbc54ea5	Підрядник надіслав заявку на співпрацю	2026-06-08 16:37:28.795115+00	d8729234-8e3b-41d9-83bc-9c725fe65838
505	21833e68-101f-4378-a916-62c120a9f192	UPDATE	CONTRACTOR_MEMBERSHIP	4f4d03e5-74e1-4272-b64b-10d9da71740c	Підрядника схвалено для співпраці	2026-06-08 16:43:05.91602+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
504	c301c093-ebc6-491a-9c40-eb425e472857	CREATE	CONTRACTOR_MEMBERSHIP	0997831c-654f-471b-934f-cedafbc54ea5	Підрядник надіслав заявку на співпрацю	2026-06-08 16:41:59.13537+00	d8729234-8e3b-41d9-83bc-9c725fe65838
506	c301c093-ebc6-491a-9c40-eb425e472857	UPDATE	CONTRACTOR_REQUEST	c59e8f03-b667-45a3-9673-775c037ddff0	Заявку взято в роботу	2026-06-08 16:43:20.520094+00	d8729234-8e3b-41d9-83bc-9c725fe65838
507	c301c093-ebc6-491a-9c40-eb425e472857	UPDATE	CONTRACTOR_REQUEST	c59e8f03-b667-45a3-9673-775c037ddff0	Статус змінено на 'Доставлено'	2026-06-08 16:44:41.944364+00	d8729234-8e3b-41d9-83bc-9c725fe65838
508	21833e68-101f-4378-a916-62c120a9f192	EXPORT	FUEL		Експорт звіту витрат пального (Excel)	2026-06-09 01:45:19.536567+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
509	21833e68-101f-4378-a916-62c120a9f192	UPDATE	CONTRACTOR_REQUEST	c59e8f03-b667-45a3-9673-775c037ddff0	Майно успішно прийнято на баланс	2026-06-09 17:57:37.447327+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
\.


--
-- Data for Name: contractor_memberships; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.contractor_memberships (id, contractor_id, tenant_id, status, note, requested_at, decided_at, decided_by) FROM stdin;
57b59559-6f6b-4ff1-9951-58e924d3c21f	c301c093-ebc6-491a-9c40-eb425e472857	0997831c-654f-471b-934f-cedafbc54ea5	PENDING	\N	2026-06-08 16:41:59.131249+00	\N	\N
4f4d03e5-74e1-4272-b64b-10d9da71740c	c301c093-ebc6-491a-9c40-eb425e472857	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	APPROVED	\N	2026-06-08 15:59:06.454186+00	2026-06-08 16:43:05.904795+00	21833e68-101f-4378-a916-62c120a9f192
\.


--
-- Data for Name: contractor_requests; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.contractor_requests (id, tenant_id, created_by, unit_id, title, description, status, taken_by, taken_at, completed_at, created_at, deadline, updated_at, target_warehouse_id) FROM stdin;
43d5914c-aa5d-4cfa-8d52-b5340277d1a3	0997831c-654f-471b-934f-cedafbc54ea5	3e59277f-0c1a-4aa3-973c-cf3db1e19497	75	Запчастини для Ford Transit	Заявка створена автоматично seed-скриптом для тестування.	IN_PROGRESS	e4e50963-10f3-46c5-86cd-988e5799c3fe	2026-05-01 10:33:21.274364+00	\N	2026-04-27 10:33:21.041491+00	\N	2026-05-27 12:49:14.50152+00	\N
e98e02a6-186d-41e2-a346-c8853fc715ff	0997831c-654f-471b-934f-cedafbc54ea5	3e59277f-0c1a-4aa3-973c-cf3db1e19497	75	Акумулятори LiFePO4 (2 шт.)	Заявка створена автоматично seed-скриптом для тестування.	OPEN	\N	\N	\N	2026-04-27 10:33:21.041491+00	\N	2026-05-27 12:49:14.50152+00	\N
f034b209-856b-41cb-ac62-a0f5ee2e85e1	0997831c-654f-471b-934f-cedafbc54ea5	3e59277f-0c1a-4aa3-973c-cf3db1e19497	75	Ноутбуки — 2 шт. для офісу	Заявка створена автоматично seed-скриптом для тестування.	OPEN	\N	\N	\N	2026-04-27 10:33:21.041491+00	\N	2026-05-27 12:49:14.50152+00	\N
fb96aae3-b5f3-4f6e-9539-a49d796eb4e5	0997831c-654f-471b-934f-cedafbc54ea5	3e59277f-0c1a-4aa3-973c-cf3db1e19497	75	Потрібен генератор 5 кВт	Заявка створена автоматично seed-скриптом для тестування.	COMPLETED	a2fbe6df-f0aa-432e-a5df-6759d8b5dff6	2026-05-04 10:33:21.277031+00	2026-05-05 10:33:21.277031+00	2026-04-27 10:33:21.041491+00	\N	2026-05-27 12:49:14.50152+00	\N
10557735-c178-42f2-9bf0-af943ab7fad9	0997831c-654f-471b-934f-cedafbc54ea5	223e9c26-98f3-4f00-8af8-dc066e2de9f4	81	Запчастини для Ford Transit	Заявка створена автоматично seed-скриптом для тестування.	COMPLETED	a2fbe6df-f0aa-432e-a5df-6759d8b5dff6	2026-05-05 10:33:21.278529+00	2026-05-08 10:33:21.278529+00	2026-04-27 10:33:21.041491+00	\N	2026-05-27 12:49:14.50152+00	\N
b25d587a-42b7-42bd-8dcf-409ba30ca721	0997831c-654f-471b-934f-cedafbc54ea5	223e9c26-98f3-4f00-8af8-dc066e2de9f4	81	Комплект рацій Motorola (5 шт.)	Заявка створена автоматично seed-скриптом для тестування.	OPEN	\N	\N	\N	2026-04-27 10:33:21.041491+00	\N	2026-05-27 12:49:14.50152+00	\N
1dbe5a18-b4f9-4654-bd15-5ef877fa5ad2	0997831c-654f-471b-934f-cedafbc54ea5	223e9c26-98f3-4f00-8af8-dc066e2de9f4	81	Спальні мішки для команди (6 шт.)	Заявка створена автоматично seed-скриптом для тестування.	COMPLETED	a2fbe6df-f0aa-432e-a5df-6759d8b5dff6	2026-04-30 10:33:21.280157+00	2026-05-04 10:33:21.280157+00	2026-04-27 10:33:21.041491+00	\N	2026-05-27 12:49:14.50152+00	\N
7de291aa-5ded-468f-8fd0-9b85e6a89c61	0997831c-654f-471b-934f-cedafbc54ea5	81810dfd-0d94-402c-a398-fa6d5559686e	84	Потрібні аптечки IFAK (10 комплектів)	Заявка створена автоматично seed-скриптом для тестування.	OPEN	\N	\N	\N	2026-04-27 10:33:21.041491+00	\N	2026-05-27 12:49:14.50152+00	\N
cc37c50a-8d8e-41d2-b87b-a86fb5a02a73	0997831c-654f-471b-934f-cedafbc54ea5	81810dfd-0d94-402c-a398-fa6d5559686e	84	Потрібен генератор 5 кВт	Заявка створена автоматично seed-скриптом для тестування.	OPEN	\N	\N	\N	2026-04-27 10:33:21.041491+00	\N	2026-05-27 12:49:14.50152+00	\N
72403559-47f6-4b51-b3f8-6d32478915d1	0997831c-654f-471b-934f-cedafbc54ea5	03ced0df-4ca8-4e03-9115-a390abdf9e15	86	Акумулятори LiFePO4 (2 шт.)	Заявка створена автоматично seed-скриптом для тестування.	OPEN	\N	\N	\N	2026-04-27 10:33:21.041491+00	\N	2026-05-27 12:49:14.50152+00	\N
b3d3361d-0143-4c1f-a765-31641812c4a3	0997831c-654f-471b-934f-cedafbc54ea5	03ced0df-4ca8-4e03-9115-a390abdf9e15	86	Ноутбуки — 2 шт. для офісу	Заявка створена автоматично seed-скриптом для тестування.	COMPLETED	e4e50963-10f3-46c5-86cd-988e5799c3fe	2026-05-03 10:33:21.282987+00	2026-05-06 10:33:21.282987+00	2026-04-27 10:33:21.041491+00	\N	2026-05-27 12:49:14.50152+00	\N
c9a790a3-097f-4585-855b-0380cba68cfb	0997831c-654f-471b-934f-cedafbc54ea5	03ced0df-4ca8-4e03-9115-a390abdf9e15	86	Ноутбуки — 2 шт. для офісу	Заявка створена автоматично seed-скриптом для тестування.	IN_PROGRESS	e4e50963-10f3-46c5-86cd-988e5799c3fe	2026-04-30 10:33:21.283738+00	\N	2026-04-27 10:33:21.041491+00	\N	2026-05-27 12:49:14.50152+00	\N
2c26917d-a600-4fe4-bf36-eda3fc3eec98	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	Термоукладка	Тест	OPEN	\N	\N	\N	2026-06-08 09:50:16.406939+00	\N	2026-06-08 09:50:16.406939+00	\N
263d8252-a281-441a-ba33-bd40758cd092	0997831c-654f-471b-934f-cedafbc54ea5	223e9c26-98f3-4f00-8af8-dc066e2de9f4	81	Спальні мішки для команди (6 шт.)	Заявка створена автоматично seed-скриптом для тестування.	OPEN	\N	\N	\N	2026-04-27 10:33:21.041491+00	\N	2026-05-27 12:49:14.50152+00	\N
c59e8f03-b667-45a3-9673-775c037ddff0	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	Гідропідсилювач	Потрібно 5 штук	ACCEPTED	c301c093-ebc6-491a-9c40-eb425e472857	2026-06-08 16:43:20.510304+00	2026-06-09 17:57:37.432268+00	2026-06-08 14:41:25.622409+00	\N	2026-06-08 14:41:25.622409+00	f64e8882-623d-41cd-8f72-2be30b783f8d
\.


--
-- Data for Name: fuel_records; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.fuel_records (id, tenant_id, vehicle_id, liters, odometer_km, record_type, created_by, created_at, is_anomaly, anomaly_reason, anomaly_excess_liters) FROM stdin;
23bd3383-85a9-4eae-a258-e9ef1c564114	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	48471a6b-ac6a-4327-b13c-dce6768ac4a0	0.11	1	EXPENSE	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 07:31:30.539102+00	f	\N	0.00
b3dc76fb-0c4a-4c01-9943-f47fb2e0aeed	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	48471a6b-ac6a-4327-b13c-dce6768ac4a0	0.11	2	EXPENSE	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 07:55:25.962499+00	f	\N	0.00
a8b79e44-dab9-4190-a79e-3caa320e5bf9	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	8355c89e-b072-4113-a6d2-adecc94a5a31	0.33	2	EXPENSE	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 09:36:48.874212+00	f	\N	0.00
c25668cf-8b91-416c-bef4-3d32c7b93caa	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	ea092cad-3147-4245-8ce6-ee7b8029806b	0.10	3	EXPENSE	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 21:31:58.173825+00	f	\N	0.00
bcce8e66-efdc-4cad-9739-046198ea68c9	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	ea092cad-3147-4245-8ce6-ee7b8029806b	50.00	3	REFUEL	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 21:34:25.653175+00	f	\N	0.00
f18a03ac-376e-436c-88ee-2dd0cb799e49	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	1565aa02-d4ff-464c-bc5b-c96f63fd9be1	10.20	2	REFUEL	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 21:42:26.736485+00	f	\N	0.00
17a4bc91-c364-494d-ab12-b7308706f842	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	1565aa02-d4ff-464c-bc5b-c96f63fd9be1	0.10	3	EXPENSE	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 22:33:47.197123+00	f	\N	0.00
238fc951-5e44-4de2-9bee-f9952a8f71e6	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	1565aa02-d4ff-464c-bc5b-c96f63fd9be1	0.64	11	EXPENSE	21833e68-101f-4378-a916-62c120a9f192	2026-05-30 10:17:25.310516+00	f	\N	0.00
5e7122a1-87e4-407c-9a25-882172001bbc	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	1565aa02-d4ff-464c-bc5b-c96f63fd9be1	0.64	19	EXPENSE	21833e68-101f-4378-a916-62c120a9f192	2026-05-30 12:48:02.825064+00	f	\N	0.00
8c9a3f33-80b1-4cb4-b524-62f652a8d5fa	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	1565aa02-d4ff-464c-bc5b-c96f63fd9be1	0.64	27	EXPENSE	21833e68-101f-4378-a916-62c120a9f192	2026-05-30 12:52:08.804323+00	f	\N	0.00
be1662ff-30a1-48ad-b8ef-d7db1ffa8195	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	ea092cad-3147-4245-8ce6-ee7b8029806b	5.00	2	REFUEL	21833e68-101f-4378-a916-62c120a9f192	2026-05-30 22:02:20.735099+00	f	\N	0.00
724751cb-a1b5-4961-8eff-579727aef814	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	ea092cad-3147-4245-8ce6-ee7b8029806b	1.52	22	EXPENSE	21833e68-101f-4378-a916-62c120a9f192	2026-05-31 22:53:35.193414+00	f	\N	0.00
d0b9a3ff-19dc-4d81-ace9-79f0ab9c680e	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	8355c89e-b072-4113-a6d2-adecc94a5a31	0.33	1	EXPENSE	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 07:55:28.763902+00	f	\N	0.00
68863b63-115a-40b6-a2c3-6f9834712b10	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	dac43200-3bfe-47c8-8e60-17bf1ce6cf40	0.11	1	EXPENSE	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 07:55:32.737308+00	f	\N	0.00
fe30041e-bddd-4dc4-8275-472438483256	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	dac43200-3bfe-47c8-8e60-17bf1ce6cf40	0.11	2	EXPENSE	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 09:36:51.861145+00	f	\N	0.00
771373fd-c57c-4121-8178-37d1520fb1e8	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	1565aa02-d4ff-464c-bc5b-c96f63fd9be1	30.00	2	REFUEL	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 21:35:35.108501+00	f	\N	0.00
7ddeb19c-f42a-430d-8379-0ce5c74eb3e9	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	ea092cad-3147-4245-8ce6-ee7b8029806b	0.10	4	EXPENSE	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 22:33:50.827435+00	f	\N	0.00
56abddac-a2a7-4d2e-b4d5-5f85a48dd3e7	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	1565aa02-d4ff-464c-bc5b-c96f63fd9be1	1.28	43	EXPENSE	21833e68-101f-4378-a916-62c120a9f192	2026-05-30 13:12:25.804356+00	f	\N	0.00
2844e06e-6cf8-4da0-9aa5-e90f636093a5	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	1565aa02-d4ff-464c-bc5b-c96f63fd9be1	0.10	1	EXPENSE	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 07:55:36.162222+00	f	\N	0.00
8ba485b8-5811-452b-b775-7cb900e83348	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	1565aa02-d4ff-464c-bc5b-c96f63fd9be1	0.10	2	EXPENSE	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 09:36:54.577963+00	f	\N	0.00
4f4fa89d-5af8-4d3b-af82-17f2e757b179	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	48471a6b-ac6a-4327-b13c-dce6768ac4a0	0.11	3	EXPENSE	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 09:36:57.56492+00	f	\N	0.00
87d4376b-ef48-4d5f-ad89-7eb660e93782	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	ea092cad-3147-4245-8ce6-ee7b8029806b	0.10	5	EXPENSE	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 22:48:29.468411+00	f	\N	0.00
349e319f-1318-405f-8c4d-499d102fd8d0	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	dac43200-3bfe-47c8-8e60-17bf1ce6cf40	40.00	2	REFUEL	21833e68-101f-4378-a916-62c120a9f192	2026-05-30 15:14:53.152351+00	f	\N	0.00
9e0fa8e0-3189-41bf-8689-516ca94b9322	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	ea092cad-3147-4245-8ce6-ee7b8029806b	0.10	1	EXPENSE	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 07:55:39.657918+00	f	\N	0.00
c0e00a3c-9160-46ef-b8e0-63433b9a8190	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	ea092cad-3147-4245-8ce6-ee7b8029806b	0.10	2	EXPENSE	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 09:37:00.091532+00	f	\N	0.00
a5007d03-e297-4aa2-aa9b-f4bb59a8b56b	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	ea092cad-3147-4245-8ce6-ee7b8029806b	0.10	6	EXPENSE	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 23:07:38.798591+00	f	\N	0.00
5f44eb3c-317e-4c45-bb62-8ff691943f93	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	dac43200-3bfe-47c8-8e60-17bf1ce6cf40	1.79	18	EXPENSE	21833e68-101f-4378-a916-62c120a9f192	2026-05-30 15:16:51.525543+00	f	\N	0.00
66a2138c-44c5-4fca-9921-b7452f0b5f08	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	5482a426-5bcf-48a2-82cd-d61441e76810	40.00	0	REFUEL	21833e68-101f-4378-a916-62c120a9f192	2026-05-30 21:05:11.507248+00	f	\N	0.00
edb35478-d649-45f8-845c-765cee16da19	0997831c-654f-471b-934f-cedafbc54ea5	1fe8f920-1f68-4f98-ac28-06256ac21c55	59.23	35315	REFUEL	\N	2026-03-23 10:33:21.192073+00	f	\N	0.00
864b2852-8f25-4925-8bae-1079b05958d5	0997831c-654f-471b-934f-cedafbc54ea5	1fe8f920-1f68-4f98-ac28-06256ac21c55	20.90	35515	REFUEL	\N	2026-03-28 10:33:21.193275+00	f	\N	0.00
b4f50888-414b-4d00-96f8-46f00c55c264	0997831c-654f-471b-934f-cedafbc54ea5	1fe8f920-1f68-4f98-ac28-06256ac21c55	32.35	35757	REFUEL	\N	2026-04-02 10:33:21.193827+00	f	\N	0.00
96554d9d-76e9-4d0c-b445-faec7d45fa0c	0997831c-654f-471b-934f-cedafbc54ea5	1fe8f920-1f68-4f98-ac28-06256ac21c55	26.29	36176	REFUEL	\N	2026-04-07 10:33:21.194348+00	f	\N	0.00
4362b4ae-e282-46c1-9e56-96f067278013	0997831c-654f-471b-934f-cedafbc54ea5	1fe8f920-1f68-4f98-ac28-06256ac21c55	74.19	36290	REFUEL	\N	2026-04-12 10:33:21.194908+00	f	\N	0.00
6009d4af-84aa-45ab-97f6-63556709c7e6	0997831c-654f-471b-934f-cedafbc54ea5	1fe8f920-1f68-4f98-ac28-06256ac21c55	49.47	36854	REFUEL	\N	2026-04-17 10:33:21.195497+00	f	\N	0.00
e8098e4b-beb5-4a8f-b459-3f01faa87f98	0997831c-654f-471b-934f-cedafbc54ea5	1fe8f920-1f68-4f98-ac28-06256ac21c55	27.74	37119	REFUEL	\N	2026-04-22 10:33:21.196223+00	f	\N	0.00
dbd983a6-d598-4b13-ad1e-23011427eb95	0997831c-654f-471b-934f-cedafbc54ea5	1fe8f920-1f68-4f98-ac28-06256ac21c55	57.50	37239	REFUEL	\N	2026-04-27 10:33:21.196698+00	f	\N	0.00
2c0d8375-a0f6-477d-bd2e-700fa374f9c6	0997831c-654f-471b-934f-cedafbc54ea5	1fe8f920-1f68-4f98-ac28-06256ac21c55	70.63	37795	REFUEL	\N	2026-05-02 10:33:21.19706+00	f	\N	0.00
d5643d39-167d-42e1-ae00-ac167897652a	0997831c-654f-471b-934f-cedafbc54ea5	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	58.95	102314	REFUEL	\N	2026-03-23 10:33:21.197388+00	f	\N	0.00
80ded903-3940-47bc-a36a-1853d5fef6ca	0997831c-654f-471b-934f-cedafbc54ea5	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	68.18	102807	REFUEL	\N	2026-03-26 10:33:21.19768+00	f	\N	0.00
5ac79dbe-4b5e-41d7-9160-e8c0403afc6d	0997831c-654f-471b-934f-cedafbc54ea5	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	70.76	102927	REFUEL	\N	2026-03-29 10:33:21.19801+00	f	\N	0.00
596cdcc3-a59f-46f6-be6f-1729ea67d291	0997831c-654f-471b-934f-cedafbc54ea5	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	36.94	103451	REFUEL	\N	2026-04-01 10:33:21.198412+00	f	\N	0.00
a5bffe2b-d81f-4be7-8128-0241d9f58aea	0997831c-654f-471b-934f-cedafbc54ea5	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	43.86	103919	REFUEL	\N	2026-04-04 10:33:21.198725+00	f	\N	0.00
610338a0-6ba3-45e5-8ce4-ed46c297c74d	0997831c-654f-471b-934f-cedafbc54ea5	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	39.98	104311	REFUEL	\N	2026-04-07 10:33:21.198982+00	f	\N	0.00
208ee127-461d-4e7e-a2ce-16710a15e3a8	0997831c-654f-471b-934f-cedafbc54ea5	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	21.19	104702	REFUEL	\N	2026-04-10 10:33:21.199232+00	f	\N	0.00
409ddccb-ff27-4ba6-b4e6-97a2cc0582aa	0997831c-654f-471b-934f-cedafbc54ea5	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	70.02	104942	REFUEL	\N	2026-04-13 10:33:21.199501+00	f	\N	0.00
1e432386-bbad-4a9f-957f-4ff352b90c49	0997831c-654f-471b-934f-cedafbc54ea5	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	77.49	105256	REFUEL	\N	2026-04-16 10:33:21.199787+00	f	\N	0.00
371d988f-4677-4957-9def-bcdacfaac1e5	0997831c-654f-471b-934f-cedafbc54ea5	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	58.76	105711	REFUEL	\N	2026-04-19 10:33:21.2001+00	f	\N	0.00
777e83f3-7c3f-4d08-8aba-9f8bc0f443bd	0997831c-654f-471b-934f-cedafbc54ea5	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	55.04	106300	REFUEL	\N	2026-04-22 10:33:21.200425+00	f	\N	0.00
099318d1-5d82-4991-ab65-ee68cc202816	0997831c-654f-471b-934f-cedafbc54ea5	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	55.83	106406	REFUEL	\N	2026-04-25 10:33:21.200773+00	f	\N	0.00
a3168be0-0b63-48b1-915e-41152508d787	0997831c-654f-471b-934f-cedafbc54ea5	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	22.91	106527	REFUEL	\N	2026-04-28 10:33:21.201115+00	f	\N	0.00
4c9ec573-4f6a-4e11-9136-9d8978f00a73	0997831c-654f-471b-934f-cedafbc54ea5	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	69.68	106815	REFUEL	\N	2026-05-01 10:33:21.201436+00	f	\N	0.00
05255bf4-daf2-4c47-91b2-4f650f0595b3	0997831c-654f-471b-934f-cedafbc54ea5	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	51.87	107360	REFUEL	\N	2026-05-04 10:33:21.201802+00	f	\N	0.00
5c6132b4-4ad5-44e8-9c99-f864afadee14	0997831c-654f-471b-934f-cedafbc54ea5	f678ceb7-8984-4534-abac-41d262aaea81	74.40	47675	REFUEL	\N	2026-03-23 10:33:21.202111+00	f	\N	0.00
4a400482-b7a1-4aab-a6db-4c22fabce30d	0997831c-654f-471b-934f-cedafbc54ea5	f678ceb7-8984-4534-abac-41d262aaea81	49.32	47903	REFUEL	\N	2026-03-25 10:33:21.202394+00	f	\N	0.00
dbe75855-4937-4557-aaf9-bde93aaa93c9	0997831c-654f-471b-934f-cedafbc54ea5	f678ceb7-8984-4534-abac-41d262aaea81	39.21	48396	REFUEL	\N	2026-03-27 10:33:21.202698+00	f	\N	0.00
0c4c0f9a-83d7-4fab-99da-dab3d3ea70fc	0997831c-654f-471b-934f-cedafbc54ea5	f678ceb7-8984-4534-abac-41d262aaea81	65.11	48794	REFUEL	\N	2026-03-29 10:33:21.203047+00	f	\N	0.00
175eb38b-2801-4e95-9897-528c5aa7485d	0997831c-654f-471b-934f-cedafbc54ea5	f678ceb7-8984-4534-abac-41d262aaea81	74.57	49060	REFUEL	\N	2026-03-31 10:33:21.203369+00	f	\N	0.00
b9599d78-17d8-4e2c-a93e-534e1e3f891d	0997831c-654f-471b-934f-cedafbc54ea5	f678ceb7-8984-4534-abac-41d262aaea81	75.79	49453	REFUEL	\N	2026-04-02 10:33:21.203696+00	f	\N	0.00
e3fe4167-a456-4198-8edc-73c863207ed7	0997831c-654f-471b-934f-cedafbc54ea5	f678ceb7-8984-4534-abac-41d262aaea81	37.48	49787	REFUEL	\N	2026-04-04 10:33:21.204267+00	f	\N	0.00
53fc4975-5d13-440e-b250-5676a65c3a22	0997831c-654f-471b-934f-cedafbc54ea5	f678ceb7-8984-4534-abac-41d262aaea81	53.81	49947	REFUEL	\N	2026-04-06 10:33:21.204605+00	f	\N	0.00
d6e27b79-7c8b-4799-94b4-ea8515068f49	0997831c-654f-471b-934f-cedafbc54ea5	f678ceb7-8984-4534-abac-41d262aaea81	76.55	50148	REFUEL	\N	2026-04-08 10:33:21.204906+00	f	\N	0.00
6b45192d-bd51-4145-8fdd-f896ad39a8b9	0997831c-654f-471b-934f-cedafbc54ea5	f678ceb7-8984-4534-abac-41d262aaea81	29.56	50390	REFUEL	\N	2026-04-10 10:33:21.205191+00	f	\N	0.00
93766f6b-6a89-496c-9c3d-964fee29f023	0997831c-654f-471b-934f-cedafbc54ea5	f678ceb7-8984-4534-abac-41d262aaea81	64.86	50883	REFUEL	\N	2026-04-12 10:33:21.205488+00	f	\N	0.00
b0d8287d-1d95-458e-affd-f4c0debd8c08	0997831c-654f-471b-934f-cedafbc54ea5	f678ceb7-8984-4534-abac-41d262aaea81	30.02	51250	REFUEL	\N	2026-04-14 10:33:21.205866+00	f	\N	0.00
632d6ef9-2dac-4a35-93f4-7e1edac448a9	0997831c-654f-471b-934f-cedafbc54ea5	f678ceb7-8984-4534-abac-41d262aaea81	68.65	51532	REFUEL	\N	2026-04-16 10:33:21.206145+00	f	\N	0.00
08e5abc9-8f52-44eb-9134-2f4c5a242c7d	0997831c-654f-471b-934f-cedafbc54ea5	f678ceb7-8984-4534-abac-41d262aaea81	73.61	51673	REFUEL	\N	2026-04-18 10:33:21.206403+00	f	\N	0.00
c92cc622-19e2-4414-9a93-b859e53d87d0	0997831c-654f-471b-934f-cedafbc54ea5	f678ceb7-8984-4534-abac-41d262aaea81	23.35	51801	REFUEL	\N	2026-04-20 10:33:21.206663+00	f	\N	0.00
c176e693-0a17-44cf-9e2e-5f1b94725210	0997831c-654f-471b-934f-cedafbc54ea5	f678ceb7-8984-4534-abac-41d262aaea81	56.20	51944	REFUEL	\N	2026-04-22 10:33:21.206905+00	f	\N	0.00
eed29a63-b65c-467e-ac54-a97592b82652	0997831c-654f-471b-934f-cedafbc54ea5	f678ceb7-8984-4534-abac-41d262aaea81	37.05	52053	REFUEL	\N	2026-04-24 10:33:21.20715+00	f	\N	0.00
8650b68b-5a22-4746-8cfa-0b9e18a991f1	0997831c-654f-471b-934f-cedafbc54ea5	f678ceb7-8984-4534-abac-41d262aaea81	66.69	52521	REFUEL	\N	2026-04-26 10:33:21.207408+00	f	\N	0.00
17738ded-d374-46ea-b07a-07bb57192294	0997831c-654f-471b-934f-cedafbc54ea5	f678ceb7-8984-4534-abac-41d262aaea81	66.14	52925	REFUEL	\N	2026-04-28 10:33:21.207665+00	f	\N	0.00
a105311f-1fe5-4214-91e7-d82d140796fe	0997831c-654f-471b-934f-cedafbc54ea5	f678ceb7-8984-4534-abac-41d262aaea81	67.26	53123	REFUEL	\N	2026-04-30 10:33:21.207928+00	f	\N	0.00
a6588a53-c294-48a3-ae23-e70a9f0a53ca	0997831c-654f-471b-934f-cedafbc54ea5	f678ceb7-8984-4534-abac-41d262aaea81	67.03	53608	REFUEL	\N	2026-05-02 10:33:21.208209+00	f	\N	0.00
5c159a5d-91f6-414d-8610-ae6d1927ef4f	0997831c-654f-471b-934f-cedafbc54ea5	f678ceb7-8984-4534-abac-41d262aaea81	78.39	54062	REFUEL	\N	2026-05-04 10:33:21.208489+00	f	\N	0.00
3d6402f4-cda0-462c-943a-d0ad72aa98da	0997831c-654f-471b-934f-cedafbc54ea5	f678ceb7-8984-4534-abac-41d262aaea81	57.43	54417	REFUEL	\N	2026-05-06 10:33:21.208798+00	f	\N	0.00
5deefae1-ca01-4018-bf3a-8c66e4cd5ad4	0997831c-654f-471b-934f-cedafbc54ea5	6a4d89a0-f357-4ccf-bfaf-8f0791f22192	36.47	123598	REFUEL	\N	2026-03-23 10:33:21.209173+00	f	\N	0.00
1a7c2f5e-31bf-4bb8-b173-e1f0eca57457	0997831c-654f-471b-934f-cedafbc54ea5	6a4d89a0-f357-4ccf-bfaf-8f0791f22192	22.46	123926	REFUEL	\N	2026-03-28 10:33:21.209728+00	f	\N	0.00
2ed8f962-2c71-49bf-a72a-776980bb3e51	0997831c-654f-471b-934f-cedafbc54ea5	6a4d89a0-f357-4ccf-bfaf-8f0791f22192	69.73	124439	REFUEL	\N	2026-04-02 10:33:21.2101+00	f	\N	0.00
945f0905-533d-40b3-834d-f7be9bba8962	0997831c-654f-471b-934f-cedafbc54ea5	6a4d89a0-f357-4ccf-bfaf-8f0791f22192	44.46	124714	REFUEL	\N	2026-04-07 10:33:21.210379+00	f	\N	0.00
89bbfbc4-688e-4e66-b02e-51a15f063ffe	0997831c-654f-471b-934f-cedafbc54ea5	6a4d89a0-f357-4ccf-bfaf-8f0791f22192	77.24	124914	REFUEL	\N	2026-04-12 10:33:21.210626+00	f	\N	0.00
8d3b3fd5-f27f-4cd6-9fba-7f3e4363400b	0997831c-654f-471b-934f-cedafbc54ea5	6a4d89a0-f357-4ccf-bfaf-8f0791f22192	62.11	125380	REFUEL	\N	2026-04-17 10:33:21.210864+00	f	\N	0.00
39cc61d9-b537-46e2-8b42-1e914115e871	0997831c-654f-471b-934f-cedafbc54ea5	6a4d89a0-f357-4ccf-bfaf-8f0791f22192	63.72	125782	REFUEL	\N	2026-04-22 10:33:21.211088+00	f	\N	0.00
0c7fc90a-1231-47eb-a0c9-97b185eeab7a	0997831c-654f-471b-934f-cedafbc54ea5	6a4d89a0-f357-4ccf-bfaf-8f0791f22192	65.41	126226	REFUEL	\N	2026-04-27 10:33:21.211405+00	f	\N	0.00
fa6ec054-7cc2-48a7-b149-6c164083f776	0997831c-654f-471b-934f-cedafbc54ea5	6a4d89a0-f357-4ccf-bfaf-8f0791f22192	64.26	126641	REFUEL	\N	2026-05-02 10:33:21.211659+00	f	\N	0.00
ee1faab8-8039-4bab-a000-db29d02dfcab	0997831c-654f-471b-934f-cedafbc54ea5	9e479722-a3d7-4cc7-ad11-0f26f67be07b	69.59	36564	REFUEL	\N	2026-03-23 10:33:21.211939+00	f	\N	0.00
20b2de1f-a8ab-41fd-ac1c-c4eea172c8f3	0997831c-654f-471b-934f-cedafbc54ea5	9e479722-a3d7-4cc7-ad11-0f26f67be07b	33.73	36944	REFUEL	\N	2026-03-28 10:33:21.212207+00	f	\N	0.00
ae776bb7-3387-4df5-bb0d-c5f72c0f8ee1	0997831c-654f-471b-934f-cedafbc54ea5	9e479722-a3d7-4cc7-ad11-0f26f67be07b	50.06	37203	REFUEL	\N	2026-04-02 10:33:21.212509+00	f	\N	0.00
3ce07896-9ca3-4702-9b51-f534d3177b6c	0997831c-654f-471b-934f-cedafbc54ea5	9e479722-a3d7-4cc7-ad11-0f26f67be07b	62.39	37340	REFUEL	\N	2026-04-07 10:33:21.212763+00	f	\N	0.00
709f0cc0-23ac-4b78-a839-ec2d443d6b5b	0997831c-654f-471b-934f-cedafbc54ea5	9e479722-a3d7-4cc7-ad11-0f26f67be07b	65.80	37470	REFUEL	\N	2026-04-12 10:33:21.213012+00	f	\N	0.00
d368ebce-7276-4472-a1fc-efa3e65d8749	0997831c-654f-471b-934f-cedafbc54ea5	9e479722-a3d7-4cc7-ad11-0f26f67be07b	77.60	37887	REFUEL	\N	2026-04-17 10:33:21.21326+00	f	\N	0.00
cb2a3648-fef4-4a52-b00d-08dfb49525c5	0997831c-654f-471b-934f-cedafbc54ea5	9e479722-a3d7-4cc7-ad11-0f26f67be07b	55.03	38035	REFUEL	\N	2026-04-22 10:33:21.213519+00	f	\N	0.00
7a172917-3060-42ec-9fa9-a0380ac3d32c	0997831c-654f-471b-934f-cedafbc54ea5	9e479722-a3d7-4cc7-ad11-0f26f67be07b	61.45	38215	REFUEL	\N	2026-04-27 10:33:21.213952+00	f	\N	0.00
e4c67144-6ba6-4f46-8442-c0f58091e252	0997831c-654f-471b-934f-cedafbc54ea5	9e479722-a3d7-4cc7-ad11-0f26f67be07b	72.79	38488	REFUEL	\N	2026-05-02 10:33:21.214277+00	f	\N	0.00
1587296a-b8b7-4fce-a5a0-76a72598c07b	0997831c-654f-471b-934f-cedafbc54ea5	4775a9bb-25f8-49ef-a65b-1002543d468b	65.76	78311	REFUEL	\N	2026-03-23 10:33:21.214526+00	f	\N	0.00
5d82fa66-de7b-4fcd-8e20-bc31901aa6ff	0997831c-654f-471b-934f-cedafbc54ea5	4775a9bb-25f8-49ef-a65b-1002543d468b	77.58	78413	REFUEL	\N	2026-03-25 10:33:21.214773+00	f	\N	0.00
7ca433ac-f320-4859-91f7-d5dc516609db	0997831c-654f-471b-934f-cedafbc54ea5	4775a9bb-25f8-49ef-a65b-1002543d468b	69.47	78863	REFUEL	\N	2026-03-27 10:33:21.215017+00	f	\N	0.00
dbb78a44-6a30-4308-a850-55d45cf80249	0997831c-654f-471b-934f-cedafbc54ea5	4775a9bb-25f8-49ef-a65b-1002543d468b	80.00	79110	REFUEL	\N	2026-03-29 10:33:21.215427+00	f	\N	0.00
487b79f6-ce56-4807-8f4a-d0ea8e8c9191	0997831c-654f-471b-934f-cedafbc54ea5	4775a9bb-25f8-49ef-a65b-1002543d468b	52.56	79549	REFUEL	\N	2026-03-31 10:33:21.215675+00	f	\N	0.00
29f7d4f7-514f-4e3b-8b9f-87e25564cca8	0997831c-654f-471b-934f-cedafbc54ea5	4775a9bb-25f8-49ef-a65b-1002543d468b	60.05	79791	REFUEL	\N	2026-04-02 10:33:21.215924+00	f	\N	0.00
edc5872c-3050-4237-be6a-f093133b8173	0997831c-654f-471b-934f-cedafbc54ea5	4775a9bb-25f8-49ef-a65b-1002543d468b	57.41	80204	REFUEL	\N	2026-04-04 10:33:21.216166+00	f	\N	0.00
936c84ae-38bc-49a0-887c-fbd8ee4a1de2	0997831c-654f-471b-934f-cedafbc54ea5	4775a9bb-25f8-49ef-a65b-1002543d468b	69.91	80577	REFUEL	\N	2026-04-06 10:33:21.216404+00	f	\N	0.00
a604c458-a012-4307-9d0b-4ee8ea3fd2dc	0997831c-654f-471b-934f-cedafbc54ea5	4775a9bb-25f8-49ef-a65b-1002543d468b	66.94	80917	REFUEL	\N	2026-04-08 10:33:21.216637+00	f	\N	0.00
fa3a8c98-18c5-47dd-8124-eee1eb624d08	0997831c-654f-471b-934f-cedafbc54ea5	4775a9bb-25f8-49ef-a65b-1002543d468b	66.41	81290	REFUEL	\N	2026-04-10 10:33:21.216874+00	f	\N	0.00
94c379de-7652-42dd-9cc7-f0453c78b574	0997831c-654f-471b-934f-cedafbc54ea5	4775a9bb-25f8-49ef-a65b-1002543d468b	77.16	81512	REFUEL	\N	2026-04-12 10:33:21.217106+00	f	\N	0.00
0b2f0d63-ad6d-4053-8d38-7abd0fb44dfb	0997831c-654f-471b-934f-cedafbc54ea5	4775a9bb-25f8-49ef-a65b-1002543d468b	33.39	81695	REFUEL	\N	2026-04-14 10:33:21.217406+00	f	\N	0.00
bc8e8abe-8b0a-4158-9e33-d26a400c0539	0997831c-654f-471b-934f-cedafbc54ea5	4775a9bb-25f8-49ef-a65b-1002543d468b	65.93	82112	REFUEL	\N	2026-04-16 10:33:21.217666+00	f	\N	0.00
a0d8c157-b05c-4570-8bf8-385439017cf2	0997831c-654f-471b-934f-cedafbc54ea5	4775a9bb-25f8-49ef-a65b-1002543d468b	47.05	82316	REFUEL	\N	2026-04-18 10:33:21.217922+00	f	\N	0.00
5771feba-e52d-4571-bb74-c8ea2fe85786	0997831c-654f-471b-934f-cedafbc54ea5	4775a9bb-25f8-49ef-a65b-1002543d468b	48.96	82762	REFUEL	\N	2026-04-20 10:33:21.21816+00	f	\N	0.00
edfdb960-8d46-45de-be4d-658166c2fc34	0997831c-654f-471b-934f-cedafbc54ea5	4775a9bb-25f8-49ef-a65b-1002543d468b	66.62	83054	REFUEL	\N	2026-04-22 10:33:21.218387+00	f	\N	0.00
af158130-6865-4fdf-88cb-512aa3aa288a	0997831c-654f-471b-934f-cedafbc54ea5	4775a9bb-25f8-49ef-a65b-1002543d468b	20.54	83435	REFUEL	\N	2026-04-24 10:33:21.218604+00	f	\N	0.00
767b953a-a122-4361-a017-b9f01b6d8591	0997831c-654f-471b-934f-cedafbc54ea5	4775a9bb-25f8-49ef-a65b-1002543d468b	49.06	83946	REFUEL	\N	2026-04-26 10:33:21.218835+00	f	\N	0.00
32ff16e2-cc0f-40ee-9d44-84c252d6796f	0997831c-654f-471b-934f-cedafbc54ea5	4775a9bb-25f8-49ef-a65b-1002543d468b	76.60	84103	REFUEL	\N	2026-04-28 10:33:21.219061+00	f	\N	0.00
022bd8c2-f758-4746-add3-bb22a97d0a75	0997831c-654f-471b-934f-cedafbc54ea5	4775a9bb-25f8-49ef-a65b-1002543d468b	36.49	84410	REFUEL	\N	2026-04-30 10:33:21.219293+00	f	\N	0.00
3e73b8cc-d51b-4845-911a-947d7681d30f	0997831c-654f-471b-934f-cedafbc54ea5	4775a9bb-25f8-49ef-a65b-1002543d468b	29.19	84909	REFUEL	\N	2026-05-02 10:33:21.21952+00	f	\N	0.00
421f72b2-b4b5-4c6d-8ab5-48e87a8fce03	0997831c-654f-471b-934f-cedafbc54ea5	4775a9bb-25f8-49ef-a65b-1002543d468b	26.49	85122	REFUEL	\N	2026-05-04 10:33:21.219747+00	f	\N	0.00
b662fa5c-8d6d-4033-8e63-e01d63a92809	0997831c-654f-471b-934f-cedafbc54ea5	4775a9bb-25f8-49ef-a65b-1002543d468b	51.04	85231	REFUEL	\N	2026-05-06 10:33:21.219973+00	f	\N	0.00
aafdc908-138d-4840-a01a-df3b42039464	0997831c-654f-471b-934f-cedafbc54ea5	75da1df1-ee2f-41b8-90da-385bfaa4d080	70.35	77027	REFUEL	\N	2026-03-23 10:33:21.220229+00	f	\N	0.00
8bc70e8c-4427-4441-b429-1ddaa0f9fb9e	0997831c-654f-471b-934f-cedafbc54ea5	75da1df1-ee2f-41b8-90da-385bfaa4d080	39.25	77471	REFUEL	\N	2026-03-28 10:33:21.220464+00	f	\N	0.00
e32f3e06-8844-4aad-a85c-eba4da120cb1	0997831c-654f-471b-934f-cedafbc54ea5	75da1df1-ee2f-41b8-90da-385bfaa4d080	60.81	77737	REFUEL	\N	2026-04-02 10:33:21.220734+00	f	\N	0.00
0f1fb494-ea79-4cf9-bbe1-e78a7b5fcef8	0997831c-654f-471b-934f-cedafbc54ea5	75da1df1-ee2f-41b8-90da-385bfaa4d080	52.47	78300	REFUEL	\N	2026-04-07 10:33:21.221017+00	f	\N	0.00
66116acb-8b72-46c7-bb4f-ca4559dc0951	0997831c-654f-471b-934f-cedafbc54ea5	75da1df1-ee2f-41b8-90da-385bfaa4d080	46.61	78766	REFUEL	\N	2026-04-12 10:33:21.221252+00	f	\N	0.00
6a2ea3a2-f5c8-442a-ac8f-afe7741f344c	0997831c-654f-471b-934f-cedafbc54ea5	75da1df1-ee2f-41b8-90da-385bfaa4d080	53.41	78988	REFUEL	\N	2026-04-17 10:33:21.22146+00	f	\N	0.00
75db201d-fcfc-4a51-a94b-6977ec3e5686	0997831c-654f-471b-934f-cedafbc54ea5	75da1df1-ee2f-41b8-90da-385bfaa4d080	75.32	79398	REFUEL	\N	2026-04-22 10:33:21.221669+00	f	\N	0.00
81cc9acc-ee51-4e0f-af00-4d10b91aef6f	0997831c-654f-471b-934f-cedafbc54ea5	75da1df1-ee2f-41b8-90da-385bfaa4d080	62.99	79640	REFUEL	\N	2026-04-27 10:33:21.221863+00	f	\N	0.00
ec6a6dbb-25ad-4330-9a87-57757ff0998b	0997831c-654f-471b-934f-cedafbc54ea5	75da1df1-ee2f-41b8-90da-385bfaa4d080	29.49	80050	REFUEL	\N	2026-05-02 10:33:21.222047+00	f	\N	0.00
cec15a00-83e4-4851-8415-3d45519895ec	0997831c-654f-471b-934f-cedafbc54ea5	3863dff1-509f-4519-b6fd-eb60da807648	35.85	29860	REFUEL	\N	2026-03-23 10:33:21.222232+00	f	\N	0.00
d7ebbab7-ea79-418c-8ce8-23e3c4102bee	0997831c-654f-471b-934f-cedafbc54ea5	3863dff1-509f-4519-b6fd-eb60da807648	58.84	30155	REFUEL	\N	2026-03-25 10:33:21.222583+00	f	\N	0.00
9a3a3b84-f6aa-417f-a1c7-d3280486355e	0997831c-654f-471b-934f-cedafbc54ea5	3863dff1-509f-4519-b6fd-eb60da807648	72.74	30618	REFUEL	\N	2026-03-27 10:33:21.222838+00	f	\N	0.00
b4c3294f-e00f-41ea-a844-2c323bbedf86	0997831c-654f-471b-934f-cedafbc54ea5	3863dff1-509f-4519-b6fd-eb60da807648	79.82	30943	REFUEL	\N	2026-03-29 10:33:21.223082+00	f	\N	0.00
5136b04a-d021-410c-a3a9-d4f5de5b68c0	0997831c-654f-471b-934f-cedafbc54ea5	3863dff1-509f-4519-b6fd-eb60da807648	40.41	31291	REFUEL	\N	2026-03-31 10:33:21.223325+00	f	\N	0.00
881593c5-d0b9-4184-9340-0f4d82fbd86d	0997831c-654f-471b-934f-cedafbc54ea5	3863dff1-509f-4519-b6fd-eb60da807648	38.13	31429	REFUEL	\N	2026-04-02 10:33:21.223571+00	f	\N	0.00
24e2b630-b3b9-495e-b232-49ede8d73e4b	0997831c-654f-471b-934f-cedafbc54ea5	3863dff1-509f-4519-b6fd-eb60da807648	56.48	31534	REFUEL	\N	2026-04-04 10:33:21.223932+00	f	\N	0.00
d8028689-7fa2-40af-a0bb-e9888fc761fc	0997831c-654f-471b-934f-cedafbc54ea5	3863dff1-509f-4519-b6fd-eb60da807648	51.67	32020	REFUEL	\N	2026-04-06 10:33:21.224152+00	f	\N	0.00
f05653d7-87bf-4656-8c80-3100cf994208	0997831c-654f-471b-934f-cedafbc54ea5	3863dff1-509f-4519-b6fd-eb60da807648	56.67	32435	REFUEL	\N	2026-04-08 10:33:21.224361+00	f	\N	0.00
cd071434-359d-4b8c-a539-1985584b36fb	0997831c-654f-471b-934f-cedafbc54ea5	3863dff1-509f-4519-b6fd-eb60da807648	39.73	32549	REFUEL	\N	2026-04-10 10:33:21.224567+00	f	\N	0.00
e8a8bfc9-06c8-45e4-8b13-e542c2b072c2	0997831c-654f-471b-934f-cedafbc54ea5	3863dff1-509f-4519-b6fd-eb60da807648	51.01	32828	REFUEL	\N	2026-04-12 10:33:21.224779+00	f	\N	0.00
cade36ea-d50b-45d2-a9a7-fcd0a998c2f7	0997831c-654f-471b-934f-cedafbc54ea5	3863dff1-509f-4519-b6fd-eb60da807648	51.18	33180	REFUEL	\N	2026-04-14 10:33:21.224969+00	f	\N	0.00
590b8103-40e2-470b-9071-30495740cd1d	0997831c-654f-471b-934f-cedafbc54ea5	3863dff1-509f-4519-b6fd-eb60da807648	43.62	33505	REFUEL	\N	2026-04-16 10:33:21.22515+00	f	\N	0.00
44f7c104-721b-49fa-86d2-5d307a299bf8	0997831c-654f-471b-934f-cedafbc54ea5	3863dff1-509f-4519-b6fd-eb60da807648	51.12	33723	REFUEL	\N	2026-04-18 10:33:21.225331+00	f	\N	0.00
d1cfa718-4715-428c-993e-0393a7251553	0997831c-654f-471b-934f-cedafbc54ea5	3863dff1-509f-4519-b6fd-eb60da807648	60.61	34182	REFUEL	\N	2026-04-20 10:33:21.225514+00	f	\N	0.00
29bbb1b7-63d6-4698-b7d2-2b36c0c37291	0997831c-654f-471b-934f-cedafbc54ea5	3863dff1-509f-4519-b6fd-eb60da807648	55.44	34292	REFUEL	\N	2026-04-22 10:33:21.225726+00	f	\N	0.00
54db395e-c5fc-4605-a11f-e36a90707e59	0997831c-654f-471b-934f-cedafbc54ea5	3863dff1-509f-4519-b6fd-eb60da807648	51.66	34591	REFUEL	\N	2026-04-24 10:33:21.22596+00	f	\N	0.00
127a7daf-3306-46d6-a4e3-2f0c5168ef13	0997831c-654f-471b-934f-cedafbc54ea5	3863dff1-509f-4519-b6fd-eb60da807648	36.85	34899	REFUEL	\N	2026-04-26 10:33:21.226205+00	f	\N	0.00
15c71ed4-6ddc-4bdd-8ed1-a9624ab1db0f	0997831c-654f-471b-934f-cedafbc54ea5	3863dff1-509f-4519-b6fd-eb60da807648	23.43	35037	REFUEL	\N	2026-04-28 10:33:21.226429+00	f	\N	0.00
5a47ee4e-af64-4388-856c-cbfbf13a6641	0997831c-654f-471b-934f-cedafbc54ea5	3863dff1-509f-4519-b6fd-eb60da807648	47.54	35332	REFUEL	\N	2026-04-30 10:33:21.22662+00	f	\N	0.00
923a6565-a26b-4ffe-b9ab-7f2c39bba535	0997831c-654f-471b-934f-cedafbc54ea5	3863dff1-509f-4519-b6fd-eb60da807648	34.90	35922	REFUEL	\N	2026-05-02 10:33:21.22682+00	f	\N	0.00
677f0cc4-96d0-4c42-9c39-b276f91e3780	0997831c-654f-471b-934f-cedafbc54ea5	3863dff1-509f-4519-b6fd-eb60da807648	21.88	36380	REFUEL	\N	2026-05-04 10:33:21.227006+00	f	\N	0.00
62260e27-551b-4098-8c19-59ff2918ad24	0997831c-654f-471b-934f-cedafbc54ea5	3863dff1-509f-4519-b6fd-eb60da807648	24.94	36511	REFUEL	\N	2026-05-06 10:33:21.227193+00	f	\N	0.00
6707d219-b15f-4032-bde7-1aabebf3f1e7	0997831c-654f-471b-934f-cedafbc54ea5	9e508594-5002-4e3e-bb16-32d710685d17	27.11	36655	REFUEL	\N	2026-03-23 10:33:21.227436+00	f	\N	0.00
41eb2a4e-87c8-4160-8df0-7aff57642639	0997831c-654f-471b-934f-cedafbc54ea5	9e508594-5002-4e3e-bb16-32d710685d17	30.30	36898	REFUEL	\N	2026-03-25 10:33:21.227691+00	f	\N	0.00
1ba37dfc-ae84-487e-86a1-bb140981a2ee	0997831c-654f-471b-934f-cedafbc54ea5	9e508594-5002-4e3e-bb16-32d710685d17	36.81	37351	REFUEL	\N	2026-03-27 10:33:21.227935+00	f	\N	0.00
b3f8c994-1687-4af4-9eaa-55d0184d5184	0997831c-654f-471b-934f-cedafbc54ea5	9e508594-5002-4e3e-bb16-32d710685d17	39.07	37681	REFUEL	\N	2026-03-29 10:33:21.228176+00	f	\N	0.00
eeebf065-e9c5-4b6e-b63f-66c2e3c67d93	0997831c-654f-471b-934f-cedafbc54ea5	9e508594-5002-4e3e-bb16-32d710685d17	72.57	38092	REFUEL	\N	2026-03-31 10:33:21.22841+00	f	\N	0.00
46db02c3-7718-40db-9f86-9f49f04516a6	0997831c-654f-471b-934f-cedafbc54ea5	9e508594-5002-4e3e-bb16-32d710685d17	41.59	38352	REFUEL	\N	2026-04-02 10:33:21.228644+00	f	\N	0.00
e35afdd9-bb97-4f93-a0eb-1b7c679000a3	0997831c-654f-471b-934f-cedafbc54ea5	9e508594-5002-4e3e-bb16-32d710685d17	42.59	38510	REFUEL	\N	2026-04-04 10:33:21.228875+00	f	\N	0.00
a217f1cd-7b8d-48df-b3c2-dc38ee900809	0997831c-654f-471b-934f-cedafbc54ea5	9e508594-5002-4e3e-bb16-32d710685d17	70.76	38913	REFUEL	\N	2026-04-06 10:33:21.229148+00	f	\N	0.00
f7dfda1a-b1e2-4abb-bb4f-9dddd05b0cd3	0997831c-654f-471b-934f-cedafbc54ea5	9e508594-5002-4e3e-bb16-32d710685d17	77.35	39323	REFUEL	\N	2026-04-08 10:33:21.229392+00	f	\N	0.00
3a6d7a71-b94f-466c-842b-446a801550c6	0997831c-654f-471b-934f-cedafbc54ea5	9e508594-5002-4e3e-bb16-32d710685d17	45.52	39874	REFUEL	\N	2026-04-10 10:33:21.2297+00	f	\N	0.00
75639562-5465-45e3-980c-6298af5c8810	0997831c-654f-471b-934f-cedafbc54ea5	9e508594-5002-4e3e-bb16-32d710685d17	75.48	40145	REFUEL	\N	2026-04-12 10:33:21.229913+00	f	\N	0.00
8bac38a6-20da-439c-b4e8-65c20bb91c9b	0997831c-654f-471b-934f-cedafbc54ea5	9e508594-5002-4e3e-bb16-32d710685d17	57.78	40728	REFUEL	\N	2026-04-14 10:33:21.23011+00	f	\N	0.00
a871ff78-8821-45e6-b8a0-720aff65a93f	0997831c-654f-471b-934f-cedafbc54ea5	9e508594-5002-4e3e-bb16-32d710685d17	42.96	41246	REFUEL	\N	2026-04-16 10:33:21.230329+00	f	\N	0.00
eb04cb2c-7b4d-4c84-aefd-790c2ca59410	0997831c-654f-471b-934f-cedafbc54ea5	9e508594-5002-4e3e-bb16-32d710685d17	37.91	41450	REFUEL	\N	2026-04-18 10:33:21.23057+00	f	\N	0.00
5ab76234-14b7-466a-9d0e-8119d4a9bac1	0997831c-654f-471b-934f-cedafbc54ea5	9e508594-5002-4e3e-bb16-32d710685d17	55.85	42017	REFUEL	\N	2026-04-20 10:33:21.230804+00	f	\N	0.00
4375b786-c80c-4cdb-8478-8b82fe7c185d	0997831c-654f-471b-934f-cedafbc54ea5	9e508594-5002-4e3e-bb16-32d710685d17	75.87	42307	REFUEL	\N	2026-04-22 10:33:21.231047+00	f	\N	0.00
0ae8bc7b-690c-45dd-a088-9acba0849ad8	0997831c-654f-471b-934f-cedafbc54ea5	9e508594-5002-4e3e-bb16-32d710685d17	53.09	42884	REFUEL	\N	2026-04-24 10:33:21.231259+00	f	\N	0.00
1cffce40-316a-435f-b14e-a9a6324a8b16	0997831c-654f-471b-934f-cedafbc54ea5	9e508594-5002-4e3e-bb16-32d710685d17	45.75	43349	REFUEL	\N	2026-04-26 10:33:21.231428+00	f	\N	0.00
f3db1776-b423-4782-90bf-c46aac1a5fab	0997831c-654f-471b-934f-cedafbc54ea5	9e508594-5002-4e3e-bb16-32d710685d17	27.21	43850	REFUEL	\N	2026-04-28 10:33:21.23164+00	f	\N	0.00
f7160c2a-dcc8-461f-8597-691e9f98c86f	0997831c-654f-471b-934f-cedafbc54ea5	9e508594-5002-4e3e-bb16-32d710685d17	60.41	44090	REFUEL	\N	2026-04-30 10:33:21.231813+00	f	\N	0.00
6b27175f-6dfe-49b8-949e-3f8219d80b90	0997831c-654f-471b-934f-cedafbc54ea5	9e508594-5002-4e3e-bb16-32d710685d17	61.73	44542	REFUEL	\N	2026-05-02 10:33:21.231979+00	f	\N	0.00
1998dec3-7631-4a71-ac16-ee6d52850fb2	0997831c-654f-471b-934f-cedafbc54ea5	9e508594-5002-4e3e-bb16-32d710685d17	55.69	44787	REFUEL	\N	2026-05-04 10:33:21.232148+00	f	\N	0.00
4a8e1b34-62d9-49fc-a282-0b9608a26564	0997831c-654f-471b-934f-cedafbc54ea5	9e508594-5002-4e3e-bb16-32d710685d17	71.42	44939	REFUEL	\N	2026-05-06 10:33:21.232321+00	f	\N	0.00
5a3a018f-d91b-48d3-9f17-a56af12acbe1	0997831c-654f-471b-934f-cedafbc54ea5	23fbee13-4b35-48cb-81fe-d0c496119765	24.31	81843	REFUEL	\N	2026-03-23 10:33:21.232498+00	f	\N	0.00
69fb6bc2-a8cc-4d0d-91b7-e3f37ab904c8	0997831c-654f-471b-934f-cedafbc54ea5	23fbee13-4b35-48cb-81fe-d0c496119765	40.10	82307	REFUEL	\N	2026-03-26 10:33:21.23267+00	f	\N	0.00
0719fc2a-8a79-4395-aa77-8e7e665c83b5	0997831c-654f-471b-934f-cedafbc54ea5	23fbee13-4b35-48cb-81fe-d0c496119765	29.58	82864	REFUEL	\N	2026-03-29 10:33:21.232849+00	f	\N	0.00
d8d32b4a-fdbb-4cfd-bc8a-efe808fbbc45	0997831c-654f-471b-934f-cedafbc54ea5	23fbee13-4b35-48cb-81fe-d0c496119765	25.41	83155	REFUEL	\N	2026-04-01 10:33:21.233064+00	f	\N	0.00
279ee57d-63f2-42e6-9b24-af3763b0ff5d	0997831c-654f-471b-934f-cedafbc54ea5	23fbee13-4b35-48cb-81fe-d0c496119765	40.36	83367	REFUEL	\N	2026-04-04 10:33:21.233296+00	f	\N	0.00
daa0421c-d7f3-4c7f-8fb9-6d2f7d3576b9	0997831c-654f-471b-934f-cedafbc54ea5	23fbee13-4b35-48cb-81fe-d0c496119765	36.64	83882	REFUEL	\N	2026-04-07 10:33:21.233543+00	f	\N	0.00
f8be4ad1-8bcf-4c35-92e2-006385a3e37d	0997831c-654f-471b-934f-cedafbc54ea5	23fbee13-4b35-48cb-81fe-d0c496119765	49.47	84322	REFUEL	\N	2026-04-10 10:33:21.233798+00	f	\N	0.00
e8ff300d-6175-4ca5-abcd-0fd563e31eaa	0997831c-654f-471b-934f-cedafbc54ea5	23fbee13-4b35-48cb-81fe-d0c496119765	34.32	84806	REFUEL	\N	2026-04-13 10:33:21.234005+00	f	\N	0.00
cdafa57a-ed67-4be2-8555-6523a44af839	0997831c-654f-471b-934f-cedafbc54ea5	23fbee13-4b35-48cb-81fe-d0c496119765	56.34	85173	REFUEL	\N	2026-04-16 10:33:21.234226+00	f	\N	0.00
6dd069d9-bc34-4835-be54-3429bb990e1c	0997831c-654f-471b-934f-cedafbc54ea5	23fbee13-4b35-48cb-81fe-d0c496119765	56.78	85344	REFUEL	\N	2026-04-19 10:33:21.234457+00	f	\N	0.00
3a9a2bd7-21f0-41b0-9842-730bd3e874b8	0997831c-654f-471b-934f-cedafbc54ea5	23fbee13-4b35-48cb-81fe-d0c496119765	21.93	85780	REFUEL	\N	2026-04-22 10:33:21.234686+00	f	\N	0.00
f9fac2f4-5f53-4d78-ad45-34175f499f50	0997831c-654f-471b-934f-cedafbc54ea5	23fbee13-4b35-48cb-81fe-d0c496119765	53.24	86106	REFUEL	\N	2026-04-25 10:33:21.234907+00	f	\N	0.00
95d444ef-ed2c-4227-b5d8-cc89f1676891	0997831c-654f-471b-934f-cedafbc54ea5	23fbee13-4b35-48cb-81fe-d0c496119765	51.79	86390	REFUEL	\N	2026-04-28 10:33:21.235131+00	f	\N	0.00
7d127400-5b90-4995-98e0-f908c4a51eab	0997831c-654f-471b-934f-cedafbc54ea5	23fbee13-4b35-48cb-81fe-d0c496119765	73.58	86812	REFUEL	\N	2026-05-01 10:33:21.235352+00	f	\N	0.00
336f6075-2594-476a-b132-8e6bab3f3b23	0997831c-654f-471b-934f-cedafbc54ea5	23fbee13-4b35-48cb-81fe-d0c496119765	45.11	87277	REFUEL	\N	2026-05-04 10:33:21.235569+00	f	\N	0.00
8e50c6a2-58d0-408e-9029-35e4e596af80	0997831c-654f-471b-934f-cedafbc54ea5	7c691151-27c5-420c-80e8-cf4bedb56aa9	44.39	67496	REFUEL	\N	2026-03-23 10:33:21.235793+00	f	\N	0.00
5f09eec9-b916-4f38-bba9-09b892307f38	0997831c-654f-471b-934f-cedafbc54ea5	7c691151-27c5-420c-80e8-cf4bedb56aa9	74.99	67654	REFUEL	\N	2026-03-26 10:33:21.236014+00	f	\N	0.00
4a7f3bc8-aa89-4c84-bb0b-ae0cd4ce43ca	0997831c-654f-471b-934f-cedafbc54ea5	7c691151-27c5-420c-80e8-cf4bedb56aa9	53.48	67945	REFUEL	\N	2026-03-29 10:33:21.236239+00	f	\N	0.00
bb198328-51a6-4428-b1cb-82f6192759a6	0997831c-654f-471b-934f-cedafbc54ea5	7c691151-27c5-420c-80e8-cf4bedb56aa9	35.17	68336	REFUEL	\N	2026-04-01 10:33:21.23646+00	f	\N	0.00
fbdcb88a-b5fe-44da-b10d-1d4d9eb20d78	0997831c-654f-471b-934f-cedafbc54ea5	7c691151-27c5-420c-80e8-cf4bedb56aa9	48.08	68639	REFUEL	\N	2026-04-04 10:33:21.236677+00	f	\N	0.00
e59766ff-e8a7-472a-8b3c-4dfde4013c55	0997831c-654f-471b-934f-cedafbc54ea5	7c691151-27c5-420c-80e8-cf4bedb56aa9	58.09	68846	REFUEL	\N	2026-04-07 10:33:21.236895+00	f	\N	0.00
aef43385-9270-41be-a3da-47a1c3d7a80b	0997831c-654f-471b-934f-cedafbc54ea5	7c691151-27c5-420c-80e8-cf4bedb56aa9	20.05	69137	REFUEL	\N	2026-04-10 10:33:21.237111+00	f	\N	0.00
65aa4ca7-feb0-4bf7-85c6-47ef8be06657	0997831c-654f-471b-934f-cedafbc54ea5	7c691151-27c5-420c-80e8-cf4bedb56aa9	38.09	69295	REFUEL	\N	2026-04-13 10:33:21.237326+00	f	\N	0.00
9cf52403-5576-458e-bb90-fad04d60726b	0997831c-654f-471b-934f-cedafbc54ea5	7c691151-27c5-420c-80e8-cf4bedb56aa9	34.81	69507	REFUEL	\N	2026-04-16 10:33:21.237539+00	f	\N	0.00
8702062f-6396-4be4-83d5-1dbf1d68757e	0997831c-654f-471b-934f-cedafbc54ea5	7c691151-27c5-420c-80e8-cf4bedb56aa9	32.56	70032	REFUEL	\N	2026-04-19 10:33:21.237782+00	f	\N	0.00
845784e7-5253-415f-860c-b3e140689dd9	0997831c-654f-471b-934f-cedafbc54ea5	7c691151-27c5-420c-80e8-cf4bedb56aa9	63.41	70397	REFUEL	\N	2026-04-22 10:33:21.237999+00	f	\N	0.00
d9bfcbc0-03af-4306-8dd8-2bc471df3da6	0997831c-654f-471b-934f-cedafbc54ea5	7c691151-27c5-420c-80e8-cf4bedb56aa9	46.27	70923	REFUEL	\N	2026-04-25 10:33:21.238207+00	f	\N	0.00
a3c3f953-4e6c-488f-801a-f66b72a58b08	0997831c-654f-471b-934f-cedafbc54ea5	7c691151-27c5-420c-80e8-cf4bedb56aa9	65.16	71225	REFUEL	\N	2026-04-28 10:33:21.238411+00	f	\N	0.00
b5dd6ed5-95d5-4f53-ad5a-366430362bd6	0997831c-654f-471b-934f-cedafbc54ea5	7c691151-27c5-420c-80e8-cf4bedb56aa9	59.71	71708	REFUEL	\N	2026-05-01 10:33:21.238615+00	f	\N	0.00
cb22133a-e83e-46b4-8875-a21d76c1283c	0997831c-654f-471b-934f-cedafbc54ea5	7c691151-27c5-420c-80e8-cf4bedb56aa9	53.66	72214	REFUEL	\N	2026-05-04 10:33:21.238819+00	f	\N	0.00
8180b16f-bff4-42db-82e0-59caf8c7a641	0997831c-654f-471b-934f-cedafbc54ea5	131f4e80-eafd-458f-a070-9c51bb8605b7	60.28	83455	REFUEL	\N	2026-03-23 10:33:21.23902+00	f	\N	0.00
e434927d-c3dc-44c4-9674-c38af373733b	0997831c-654f-471b-934f-cedafbc54ea5	131f4e80-eafd-458f-a070-9c51bb8605b7	35.24	83891	REFUEL	\N	2026-03-26 10:33:21.239215+00	f	\N	0.00
0e680b59-2856-4689-ba4c-17586152f629	0997831c-654f-471b-934f-cedafbc54ea5	131f4e80-eafd-458f-a070-9c51bb8605b7	62.33	84375	REFUEL	\N	2026-03-29 10:33:21.239417+00	f	\N	0.00
f81cd33d-9a3d-448a-b499-fd298b1064d1	0997831c-654f-471b-934f-cedafbc54ea5	131f4e80-eafd-458f-a070-9c51bb8605b7	32.89	84895	REFUEL	\N	2026-04-01 10:33:21.239617+00	f	\N	0.00
f07cd9e6-9ef0-40dc-986d-725b5723b07d	0997831c-654f-471b-934f-cedafbc54ea5	131f4e80-eafd-458f-a070-9c51bb8605b7	22.61	85216	REFUEL	\N	2026-04-04 10:33:21.239825+00	f	\N	0.00
d941a51a-e855-4ff6-8a3c-245560d37aab	0997831c-654f-471b-934f-cedafbc54ea5	131f4e80-eafd-458f-a070-9c51bb8605b7	73.30	85535	REFUEL	\N	2026-04-07 10:33:21.240068+00	f	\N	0.00
edb71233-2c7b-4c31-8080-ae43043e1287	0997831c-654f-471b-934f-cedafbc54ea5	131f4e80-eafd-458f-a070-9c51bb8605b7	54.76	86015	REFUEL	\N	2026-04-10 10:33:21.240307+00	f	\N	0.00
83de843f-29ef-4f2d-9e6d-e88006df099b	0997831c-654f-471b-934f-cedafbc54ea5	131f4e80-eafd-458f-a070-9c51bb8605b7	43.37	86256	REFUEL	\N	2026-04-13 10:33:21.240537+00	f	\N	0.00
1e5b811a-dbb6-477f-b5da-3ed3241a19ac	0997831c-654f-471b-934f-cedafbc54ea5	131f4e80-eafd-458f-a070-9c51bb8605b7	63.37	86404	REFUEL	\N	2026-04-16 10:33:21.240808+00	f	\N	0.00
bbe0a1ec-1b74-4d51-b3f0-ff636fc96e53	0997831c-654f-471b-934f-cedafbc54ea5	131f4e80-eafd-458f-a070-9c51bb8605b7	65.21	86869	REFUEL	\N	2026-04-19 10:33:21.24104+00	f	\N	0.00
0e22a349-9999-456d-868c-1bcafc151533	0997831c-654f-471b-934f-cedafbc54ea5	131f4e80-eafd-458f-a070-9c51bb8605b7	26.17	87185	REFUEL	\N	2026-04-22 10:33:21.241275+00	f	\N	0.00
04ff4210-e532-48da-90a0-a6abfc06e04f	0997831c-654f-471b-934f-cedafbc54ea5	131f4e80-eafd-458f-a070-9c51bb8605b7	56.83	87535	REFUEL	\N	2026-04-25 10:33:21.241527+00	f	\N	0.00
9a65fc44-594b-40a0-9682-cdbd83f95ce1	0997831c-654f-471b-934f-cedafbc54ea5	131f4e80-eafd-458f-a070-9c51bb8605b7	35.47	87863	REFUEL	\N	2026-04-28 10:33:21.241792+00	f	\N	0.00
7dbdb636-5d54-478a-920d-e65b01a97980	0997831c-654f-471b-934f-cedafbc54ea5	131f4e80-eafd-458f-a070-9c51bb8605b7	28.29	88046	REFUEL	\N	2026-05-01 10:33:21.242104+00	f	\N	0.00
e1df4083-fefe-42df-94b2-0708fe6d7811	0997831c-654f-471b-934f-cedafbc54ea5	131f4e80-eafd-458f-a070-9c51bb8605b7	52.93	88447	REFUEL	\N	2026-05-04 10:33:21.242361+00	f	\N	0.00
eafa5794-f482-44ac-a198-9090c09be52e	0997831c-654f-471b-934f-cedafbc54ea5	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	45.02	94705	REFUEL	\N	2026-03-23 10:33:21.242607+00	f	\N	0.00
c4235787-7c8c-440c-9e63-5166906facc0	0997831c-654f-471b-934f-cedafbc54ea5	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	36.63	94988	REFUEL	\N	2026-03-26 10:33:21.242848+00	f	\N	0.00
6e49fb56-7686-4e28-b767-2f68ef6438e9	0997831c-654f-471b-934f-cedafbc54ea5	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	30.24	95164	REFUEL	\N	2026-03-29 10:33:21.243084+00	f	\N	0.00
a9d2ea28-67b2-4756-af6d-18e3e7a98ec0	0997831c-654f-471b-934f-cedafbc54ea5	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	40.36	95584	REFUEL	\N	2026-04-01 10:33:21.243323+00	f	\N	0.00
10b9cfc1-114e-4ad1-88e0-a3e1733ba910	0997831c-654f-471b-934f-cedafbc54ea5	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	29.51	95962	REFUEL	\N	2026-04-04 10:33:21.243555+00	f	\N	0.00
2869f1d8-8d9d-4c1b-b4a4-f2105648df64	0997831c-654f-471b-934f-cedafbc54ea5	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	54.56	96292	REFUEL	\N	2026-04-07 10:33:21.243789+00	f	\N	0.00
e43d9e0d-57df-4519-97d6-81efc351d853	0997831c-654f-471b-934f-cedafbc54ea5	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	55.50	96837	REFUEL	\N	2026-04-10 10:33:21.244085+00	f	\N	0.00
52efb7ef-fc8e-422e-851c-3603b0ee5095	0997831c-654f-471b-934f-cedafbc54ea5	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	56.23	97360	REFUEL	\N	2026-04-13 10:33:21.244315+00	f	\N	0.00
83b131be-d099-4ed6-a036-61d24df65425	0997831c-654f-471b-934f-cedafbc54ea5	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	26.57	97918	REFUEL	\N	2026-04-16 10:33:21.244522+00	f	\N	0.00
9b341867-4b19-4c94-9f95-a12211028405	0997831c-654f-471b-934f-cedafbc54ea5	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	56.85	98019	REFUEL	\N	2026-04-19 10:33:21.244727+00	f	\N	0.00
bec6442c-a1d8-47e3-a03f-5eba81e90d44	0997831c-654f-471b-934f-cedafbc54ea5	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	46.09	98322	REFUEL	\N	2026-04-22 10:33:21.244929+00	f	\N	0.00
f6c5cc0c-cc20-45ae-a001-a105b048ff90	0997831c-654f-471b-934f-cedafbc54ea5	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	45.23	98473	REFUEL	\N	2026-04-25 10:33:21.245137+00	f	\N	0.00
328f67d6-3bce-4da0-bbf8-f09cda47a3ee	0997831c-654f-471b-934f-cedafbc54ea5	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	71.77	98783	REFUEL	\N	2026-04-28 10:33:21.245342+00	f	\N	0.00
4ae8cdd5-50b9-4b97-a4aa-957cb17d544a	0997831c-654f-471b-934f-cedafbc54ea5	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	68.46	98986	REFUEL	\N	2026-05-01 10:33:21.24564+00	f	\N	0.00
a8cebfa6-5e08-4f49-ae0d-c86ff0c9b3d2	0997831c-654f-471b-934f-cedafbc54ea5	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	31.19	99086	REFUEL	\N	2026-05-04 10:33:21.245939+00	f	\N	0.00
bebbe79f-6991-4bf1-abcb-c609a96bbd69	0997831c-654f-471b-934f-cedafbc54ea5	d5605709-91a8-4772-92c4-d4bd92b65717	39.84	128949	REFUEL	\N	2026-03-23 10:33:21.246213+00	f	\N	0.00
ef00cd16-b55e-483a-99fe-cca54c35211c	0997831c-654f-471b-934f-cedafbc54ea5	d5605709-91a8-4772-92c4-d4bd92b65717	63.99	129544	REFUEL	\N	2026-03-28 10:33:21.246484+00	f	\N	0.00
dbea03c2-ad89-4273-853f-538f79883bd2	0997831c-654f-471b-934f-cedafbc54ea5	d5605709-91a8-4772-92c4-d4bd92b65717	55.55	130111	REFUEL	\N	2026-04-02 10:33:21.246765+00	f	\N	0.00
be25acea-d93b-43fa-9887-5dec14d5fbb6	0997831c-654f-471b-934f-cedafbc54ea5	d5605709-91a8-4772-92c4-d4bd92b65717	52.76	130356	REFUEL	\N	2026-04-07 10:33:21.247097+00	f	\N	0.00
ebbc1a70-5bd3-4dee-a915-1f03c02642bc	0997831c-654f-471b-934f-cedafbc54ea5	d5605709-91a8-4772-92c4-d4bd92b65717	35.90	130552	REFUEL	\N	2026-04-12 10:33:21.247363+00	f	\N	0.00
fb300e98-ccc5-4afe-abad-7bf6c29caaa2	0997831c-654f-471b-934f-cedafbc54ea5	d5605709-91a8-4772-92c4-d4bd92b65717	71.74	130899	REFUEL	\N	2026-04-17 10:33:21.247579+00	f	\N	0.00
6f9c4ff6-2825-4792-9a8a-79fdade83baa	0997831c-654f-471b-934f-cedafbc54ea5	d5605709-91a8-4772-92c4-d4bd92b65717	56.22	131253	REFUEL	\N	2026-04-22 10:33:21.247787+00	f	\N	0.00
f228dd47-4323-4d51-85e2-7c68c0de1068	0997831c-654f-471b-934f-cedafbc54ea5	d5605709-91a8-4772-92c4-d4bd92b65717	45.30	131584	REFUEL	\N	2026-04-27 10:33:21.247996+00	f	\N	0.00
8ccedad5-99a5-487f-b477-12bfd46c20e5	0997831c-654f-471b-934f-cedafbc54ea5	d5605709-91a8-4772-92c4-d4bd92b65717	79.36	132173	REFUEL	\N	2026-05-02 10:33:21.248234+00	f	\N	0.00
0edf228a-1643-407b-9432-dcc9233cd31a	0997831c-654f-471b-934f-cedafbc54ea5	eff2cb75-547b-42eb-89f5-233a879ac549	57.72	56612	REFUEL	\N	2026-03-23 10:33:21.248528+00	f	\N	0.00
81e52921-364f-446e-b0d2-39920c7b7d49	0997831c-654f-471b-934f-cedafbc54ea5	eff2cb75-547b-42eb-89f5-233a879ac549	49.87	56721	REFUEL	\N	2026-03-28 10:33:21.248876+00	f	\N	0.00
e52a2140-ba0a-4d1b-b749-885cfad87bc0	0997831c-654f-471b-934f-cedafbc54ea5	eff2cb75-547b-42eb-89f5-233a879ac549	56.71	57106	REFUEL	\N	2026-04-02 10:33:21.24917+00	f	\N	0.00
db82ee58-9547-4150-bb7a-f857f4f8b171	0997831c-654f-471b-934f-cedafbc54ea5	eff2cb75-547b-42eb-89f5-233a879ac549	34.11	57662	REFUEL	\N	2026-04-07 10:33:21.24944+00	f	\N	0.00
ed4e2183-9f61-492a-85c8-ed25a4ae620a	0997831c-654f-471b-934f-cedafbc54ea5	eff2cb75-547b-42eb-89f5-233a879ac549	37.78	58038	REFUEL	\N	2026-04-12 10:33:21.249707+00	f	\N	0.00
03c94680-aa1d-40b2-8670-9f5ab71601db	0997831c-654f-471b-934f-cedafbc54ea5	eff2cb75-547b-42eb-89f5-233a879ac549	64.45	58342	REFUEL	\N	2026-04-17 10:33:21.249945+00	f	\N	0.00
6cc51f2f-3ae2-401f-bdf5-a3d648873a22	0997831c-654f-471b-934f-cedafbc54ea5	eff2cb75-547b-42eb-89f5-233a879ac549	63.78	58523	REFUEL	\N	2026-04-22 10:33:21.25017+00	f	\N	0.00
898d913f-54c0-48ab-bb4f-46119a851403	0997831c-654f-471b-934f-cedafbc54ea5	eff2cb75-547b-42eb-89f5-233a879ac549	61.56	58857	REFUEL	\N	2026-04-27 10:33:21.250406+00	f	\N	0.00
0b97b7ae-535a-4766-a53b-27069075aff8	0997831c-654f-471b-934f-cedafbc54ea5	eff2cb75-547b-42eb-89f5-233a879ac549	55.37	59372	REFUEL	\N	2026-05-02 10:33:21.250625+00	f	\N	0.00
6de53431-4c4c-4c64-afd6-370d28e295c5	0997831c-654f-471b-934f-cedafbc54ea5	d3c1b403-57b5-4f49-8828-7b31b580ddd0	54.24	70341	REFUEL	\N	2026-03-23 10:33:21.25085+00	f	\N	0.00
0fa7690d-5474-4fb9-8275-9c14053214c0	0997831c-654f-471b-934f-cedafbc54ea5	d3c1b403-57b5-4f49-8828-7b31b580ddd0	30.78	70558	REFUEL	\N	2026-03-25 10:33:21.251065+00	f	\N	0.00
b8be6283-300c-4980-964c-73e14f82fcd8	0997831c-654f-471b-934f-cedafbc54ea5	d3c1b403-57b5-4f49-8828-7b31b580ddd0	27.83	70893	REFUEL	\N	2026-03-27 10:33:21.251282+00	f	\N	0.00
6e625440-0820-4839-beb3-7b27db2949ac	0997831c-654f-471b-934f-cedafbc54ea5	d3c1b403-57b5-4f49-8828-7b31b580ddd0	36.10	71242	REFUEL	\N	2026-03-29 10:33:21.25147+00	f	\N	0.00
e97ccaa3-67eb-478a-b8f6-2df66144f7dc	0997831c-654f-471b-934f-cedafbc54ea5	d3c1b403-57b5-4f49-8828-7b31b580ddd0	70.05	71802	REFUEL	\N	2026-03-31 10:33:21.25166+00	f	\N	0.00
01f2a852-7fbb-4a6b-b179-9e2805ab26c1	0997831c-654f-471b-934f-cedafbc54ea5	d3c1b403-57b5-4f49-8828-7b31b580ddd0	24.01	72125	REFUEL	\N	2026-04-02 10:33:21.25185+00	f	\N	0.00
43db8124-e75c-42cb-b1d9-c02dfa4bad57	0997831c-654f-471b-934f-cedafbc54ea5	d3c1b403-57b5-4f49-8828-7b31b580ddd0	52.58	72311	REFUEL	\N	2026-04-04 10:33:21.252039+00	f	\N	0.00
669a043a-0b1f-4354-8fbe-377fb5953a57	0997831c-654f-471b-934f-cedafbc54ea5	d3c1b403-57b5-4f49-8828-7b31b580ddd0	43.75	72810	REFUEL	\N	2026-04-06 10:33:21.252237+00	f	\N	0.00
7bff3439-995f-49ab-bc9d-69ec6d460789	0997831c-654f-471b-934f-cedafbc54ea5	d3c1b403-57b5-4f49-8828-7b31b580ddd0	38.34	73357	REFUEL	\N	2026-04-08 10:33:21.252462+00	f	\N	0.00
1b36cb6d-96ad-43b3-8db8-8bf8b0d82abd	0997831c-654f-471b-934f-cedafbc54ea5	d3c1b403-57b5-4f49-8828-7b31b580ddd0	70.06	73588	REFUEL	\N	2026-04-10 10:33:21.252648+00	f	\N	0.00
38c5d0ec-71c7-4191-81b5-256214f489d2	0997831c-654f-471b-934f-cedafbc54ea5	d3c1b403-57b5-4f49-8828-7b31b580ddd0	33.81	73937	REFUEL	\N	2026-04-12 10:33:21.252833+00	f	\N	0.00
d0c45e33-752d-4ee7-8baa-1634e01a4d31	0997831c-654f-471b-934f-cedafbc54ea5	d3c1b403-57b5-4f49-8828-7b31b580ddd0	24.75	74392	REFUEL	\N	2026-04-14 10:33:21.253015+00	f	\N	0.00
1e243e71-83e2-48d0-9a2e-d61e12d5809f	0997831c-654f-471b-934f-cedafbc54ea5	d3c1b403-57b5-4f49-8828-7b31b580ddd0	32.36	74775	REFUEL	\N	2026-04-16 10:33:21.2532+00	f	\N	0.00
0716c951-2170-41d5-b3b8-b403c01e99fe	0997831c-654f-471b-934f-cedafbc54ea5	d3c1b403-57b5-4f49-8828-7b31b580ddd0	45.44	74915	REFUEL	\N	2026-04-18 10:33:21.253405+00	f	\N	0.00
47f3cc5b-1c25-47db-986f-7190ee5f1b92	0997831c-654f-471b-934f-cedafbc54ea5	d3c1b403-57b5-4f49-8828-7b31b580ddd0	74.20	75245	REFUEL	\N	2026-04-20 10:33:21.253597+00	f	\N	0.00
4ad45335-e3ba-4a29-b033-f9ee3d54e6d9	0997831c-654f-471b-934f-cedafbc54ea5	d3c1b403-57b5-4f49-8828-7b31b580ddd0	47.85	75788	REFUEL	\N	2026-04-22 10:33:21.253875+00	f	\N	0.00
b83fbf08-7dc1-4e2d-bd32-9d49a633faf6	0997831c-654f-471b-934f-cedafbc54ea5	d3c1b403-57b5-4f49-8828-7b31b580ddd0	20.41	76067	REFUEL	\N	2026-04-24 10:33:21.254115+00	f	\N	0.00
da680d2c-3358-476c-8474-abb307cb2043	0997831c-654f-471b-934f-cedafbc54ea5	d3c1b403-57b5-4f49-8828-7b31b580ddd0	46.61	76466	REFUEL	\N	2026-04-26 10:33:21.254375+00	f	\N	0.00
a8f0f5c2-1274-407c-be8f-076c78b272b1	0997831c-654f-471b-934f-cedafbc54ea5	d3c1b403-57b5-4f49-8828-7b31b580ddd0	30.04	76766	REFUEL	\N	2026-04-28 10:33:21.254614+00	f	\N	0.00
7f63486b-63cc-429e-912b-7c3a03c18167	0997831c-654f-471b-934f-cedafbc54ea5	d3c1b403-57b5-4f49-8828-7b31b580ddd0	33.40	77121	REFUEL	\N	2026-04-30 10:33:21.254835+00	f	\N	0.00
ae7f65ba-ff32-41cb-af9f-491fcd9468ef	0997831c-654f-471b-934f-cedafbc54ea5	d3c1b403-57b5-4f49-8828-7b31b580ddd0	54.65	77486	REFUEL	\N	2026-05-02 10:33:21.255061+00	f	\N	0.00
c8ef7851-4af3-435c-8850-51de547036ef	0997831c-654f-471b-934f-cedafbc54ea5	d3c1b403-57b5-4f49-8828-7b31b580ddd0	45.18	78040	REFUEL	\N	2026-05-04 10:33:21.255278+00	f	\N	0.00
c6455d74-60b5-4cab-b658-4a71ddc58345	0997831c-654f-471b-934f-cedafbc54ea5	d3c1b403-57b5-4f49-8828-7b31b580ddd0	68.98	78614	REFUEL	\N	2026-05-06 10:33:21.255531+00	f	\N	0.00
c01e74c4-55bb-4379-81a8-f8969059108f	0997831c-654f-471b-934f-cedafbc54ea5	3a1de34e-544b-4b67-aa92-20e6cfb09e72	33.60	89184	REFUEL	\N	2026-03-23 10:33:21.255782+00	f	\N	0.00
f65ee5bb-d952-49d3-9eb1-95c7d1cd9298	0997831c-654f-471b-934f-cedafbc54ea5	3a1de34e-544b-4b67-aa92-20e6cfb09e72	66.21	89769	REFUEL	\N	2026-03-28 10:33:21.256032+00	f	\N	0.00
82ae6bb9-9865-4049-b2a7-393d9a1bd67a	0997831c-654f-471b-934f-cedafbc54ea5	3a1de34e-544b-4b67-aa92-20e6cfb09e72	30.00	90326	REFUEL	\N	2026-04-02 10:33:21.25628+00	f	\N	0.00
403d616f-919c-4948-ba1a-56b268617df7	0997831c-654f-471b-934f-cedafbc54ea5	3a1de34e-544b-4b67-aa92-20e6cfb09e72	24.46	90636	REFUEL	\N	2026-04-07 10:33:21.256524+00	f	\N	0.00
db35e420-d09c-471a-95d6-bee072fd5632	0997831c-654f-471b-934f-cedafbc54ea5	3a1de34e-544b-4b67-aa92-20e6cfb09e72	45.37	91179	REFUEL	\N	2026-04-12 10:33:21.256765+00	f	\N	0.00
9ac3e099-6d97-4da0-9bf3-53d9e0921a4b	0997831c-654f-471b-934f-cedafbc54ea5	3a1de34e-544b-4b67-aa92-20e6cfb09e72	69.98	91375	REFUEL	\N	2026-04-17 10:33:21.257008+00	f	\N	0.00
eaff171c-1e28-45b7-a291-1a943cf29f85	0997831c-654f-471b-934f-cedafbc54ea5	3a1de34e-544b-4b67-aa92-20e6cfb09e72	68.03	91822	REFUEL	\N	2026-04-22 10:33:21.257254+00	f	\N	0.00
d73bf389-7b53-4d2d-baf2-b7d534b526d3	0997831c-654f-471b-934f-cedafbc54ea5	3a1de34e-544b-4b67-aa92-20e6cfb09e72	67.58	92170	REFUEL	\N	2026-04-27 10:33:21.257489+00	f	\N	0.00
04c6d4d1-518f-47f3-9cd0-f3513a54ca04	0997831c-654f-471b-934f-cedafbc54ea5	3a1de34e-544b-4b67-aa92-20e6cfb09e72	47.92	92537	REFUEL	\N	2026-05-02 10:33:21.257739+00	f	\N	0.00
3fd20044-ad8f-4d36-aaaa-2c063840de7c	0997831c-654f-471b-934f-cedafbc54ea5	e51f66b4-9055-42e2-8790-5063fe231566	38.26	111510	REFUEL	\N	2026-03-23 10:33:21.257977+00	f	\N	0.00
e1d2ea5e-5583-46fd-b85d-b0085cca295a	0997831c-654f-471b-934f-cedafbc54ea5	e51f66b4-9055-42e2-8790-5063fe231566	77.37	111940	REFUEL	\N	2026-03-28 10:33:21.258207+00	f	\N	0.00
61e59f64-63ce-41de-8eb0-5f5ac5258474	0997831c-654f-471b-934f-cedafbc54ea5	e51f66b4-9055-42e2-8790-5063fe231566	77.74	112180	REFUEL	\N	2026-04-02 10:33:21.258436+00	f	\N	0.00
af793da5-5e93-4ea0-863f-47b47fc64f6e	0997831c-654f-471b-934f-cedafbc54ea5	e51f66b4-9055-42e2-8790-5063fe231566	22.68	112652	REFUEL	\N	2026-04-07 10:33:21.258639+00	f	\N	0.00
3b1bb4b6-ef9b-4f1f-8f51-2177c4207c35	0997831c-654f-471b-934f-cedafbc54ea5	e51f66b4-9055-42e2-8790-5063fe231566	48.60	113064	REFUEL	\N	2026-04-12 10:33:21.258845+00	f	\N	0.00
d7c30049-2639-4669-906e-417379cb255b	0997831c-654f-471b-934f-cedafbc54ea5	e51f66b4-9055-42e2-8790-5063fe231566	75.75	113256	REFUEL	\N	2026-04-17 10:33:21.259047+00	f	\N	0.00
5c5850b6-0281-44d2-8d27-820494b93a9d	0997831c-654f-471b-934f-cedafbc54ea5	e51f66b4-9055-42e2-8790-5063fe231566	47.52	113825	REFUEL	\N	2026-04-22 10:33:21.25932+00	f	\N	0.00
6fbdaaa5-6443-4a6f-93bc-9bcb044ea429	0997831c-654f-471b-934f-cedafbc54ea5	e51f66b4-9055-42e2-8790-5063fe231566	34.04	114369	REFUEL	\N	2026-04-27 10:33:21.259548+00	f	\N	0.00
f79bdc3e-2fb8-47b8-8c4c-c01b335b8f6a	0997831c-654f-471b-934f-cedafbc54ea5	e51f66b4-9055-42e2-8790-5063fe231566	57.81	114625	REFUEL	\N	2026-05-02 10:33:21.259759+00	f	\N	0.00
\.


--
-- Data for Name: geofence_alerts; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.geofence_alerts (id, vehicle_id, geofence_id, event_type, latitude, longitude, "timestamp", created_at) FROM stdin;
\.


--
-- Data for Name: geofences; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.geofences (id, unit_id, name, latitude, longitude, radius, type, active, created_at, updated_at, tenant_id) FROM stdin;
13	76	Безпечна зона штабу Чернігів	51.4982	31.2893	400	SAFE	t	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
14	81	Портова зона Одеса	46.4774	30.7326	600	SAFE	t	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
15	85	Логістичний центр Київ	50.4501	30.5234	500	SAFE	t	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
16	75	Заборонена зона	51.6	31.5	800	FORBIDDEN	t	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
\.


--
-- Data for Name: gps_locations; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.gps_locations (id, vehicle_id, unit_id, latitude, longitude, altitude, speed, heading, accuracy, "timestamp", created_at, tenant_id) FROM stdin;
1	48471a6b-ac6a-4327-b13c-dce6768ac4a0	2	50.4501	30.5234	0	0	0	150	2026-05-06 18:09:18.331267+00	2026-05-06 18:09:18.332886+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
2	48471a6b-ac6a-4327-b13c-dce6768ac4a0	2	50.4501	30.5234	0	0	0	150	2026-05-06 18:09:28.323366+00	2026-05-06 18:09:28.324462+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
3	48471a6b-ac6a-4327-b13c-dce6768ac4a0	2	50.4501	30.5234	0	0	0	150	2026-05-06 18:09:38.319967+00	2026-05-06 18:09:38.320155+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
4	48471a6b-ac6a-4327-b13c-dce6768ac4a0	2	50.4501	30.5234	0	0	0	150	2026-05-06 18:09:48.321603+00	2026-05-06 18:09:48.321855+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
5	48471a6b-ac6a-4327-b13c-dce6768ac4a0	2	50.4501	30.5234	0	0	0	150	2026-05-06 18:09:58.322146+00	2026-05-06 18:09:58.323023+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
6	48471a6b-ac6a-4327-b13c-dce6768ac4a0	2	50.4501	30.5234	0	0	0	150	2026-05-06 18:10:08.323575+00	2026-05-06 18:10:08.323876+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
7	48471a6b-ac6a-4327-b13c-dce6768ac4a0	2	50.4501	30.5234	0	0	0	150	2026-05-06 18:13:46.274855+00	2026-05-06 18:13:46.275111+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
8	48471a6b-ac6a-4327-b13c-dce6768ac4a0	2	50.4501	30.5234	0	0	0	150	2026-05-06 18:13:56.269701+00	2026-05-06 18:13:56.269951+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
9	48471a6b-ac6a-4327-b13c-dce6768ac4a0	2	50.4501	30.5234	0	0	0	150	2026-05-06 18:14:06.26718+00	2026-05-06 18:14:06.267492+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
10	48471a6b-ac6a-4327-b13c-dce6768ac4a0	2	50.4501	30.5234	0	0	0	150	2026-05-06 18:14:16.526953+00	2026-05-06 18:14:16.527279+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
11	48471a6b-ac6a-4327-b13c-dce6768ac4a0	2	50.4501	30.5234	0	0	0	150	2026-05-06 18:14:26.267374+00	2026-05-06 18:14:26.267681+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
12	48471a6b-ac6a-4327-b13c-dce6768ac4a0	2	50.4501	30.5234	0	0	0	150	2026-05-06 18:14:36.267013+00	2026-05-06 18:14:36.267397+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
13	48471a6b-ac6a-4327-b13c-dce6768ac4a0	2	50.4501	30.5234	0	0	0	150	2026-05-06 18:14:46.267935+00	2026-05-06 18:14:46.268217+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
14	48471a6b-ac6a-4327-b13c-dce6768ac4a0	2	50.4501	30.5234	0	0	0	150	2026-05-06 18:23:44.103333+00	2026-05-06 18:23:44.104346+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
15	48471a6b-ac6a-4327-b13c-dce6768ac4a0	2	50.4501	30.5234	0	0	0	150	2026-05-06 18:23:54.102019+00	2026-05-06 18:23:54.10316+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
16	48471a6b-ac6a-4327-b13c-dce6768ac4a0	2	50.4501	30.5234	0	0	0	150	2026-05-06 18:24:04.107755+00	2026-05-06 18:24:04.10888+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
17	48471a6b-ac6a-4327-b13c-dce6768ac4a0	2	50.4501	30.5234	0	0	0	150	2026-05-06 18:24:14.105403+00	2026-05-06 18:24:14.106564+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
18	48471a6b-ac6a-4327-b13c-dce6768ac4a0	2	50.4501	30.5234	0	0	0	150	2026-05-06 19:04:01.53726+00	2026-05-06 19:04:01.538091+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
1099	48471a6b-ac6a-4327-b13c-dce6768ac4a0	2	49.0242048	24.3531776	0	0	0	78179.3534530165	2026-05-17 00:43:50.488983+00	2026-05-17 00:43:50.489496+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
829	1fe8f920-1f68-4f98-ac28-06256ac21c55	75	51.47237131902229	31.30550616086932	\N	24.67	29.51	6.78	2026-05-07 10:18:21.297048+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
830	1fe8f920-1f68-4f98-ac28-06256ac21c55	75	51.478599870506336	31.294632144356367	\N	13.32	212.52	14.75	2026-05-07 10:19:21.298753+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
831	1fe8f920-1f68-4f98-ac28-06256ac21c55	75	51.46212547759531	31.26411825770035	\N	38.86	27.64	12.73	2026-05-07 10:20:21.299504+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
832	1fe8f920-1f68-4f98-ac28-06256ac21c55	75	51.484644882790256	31.31448511501475	\N	17.75	228.2	8.94	2026-05-07 10:21:21.300198+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
833	1fe8f920-1f68-4f98-ac28-06256ac21c55	75	51.47879447659899	31.25578262901432	\N	50.73	49.12	8.63	2026-05-07 10:22:21.300854+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
834	1fe8f920-1f68-4f98-ac28-06256ac21c55	75	51.46215437040987	31.25592533788396	\N	60.27	257.01	3.78	2026-05-07 10:23:21.301481+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
835	1fe8f920-1f68-4f98-ac28-06256ac21c55	75	51.50602732582542	31.255946423536354	\N	84.74	71.49	8.52	2026-05-07 10:24:21.30235+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
836	1fe8f920-1f68-4f98-ac28-06256ac21c55	75	51.506377300539455	31.325958023709187	\N	84.55	51.03	8	2026-05-07 10:25:21.302988+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
837	1fe8f920-1f68-4f98-ac28-06256ac21c55	75	51.506354497722796	31.329149863907876	\N	79.6	156.15	12.83	2026-05-07 10:26:21.30339+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
838	1fe8f920-1f68-4f98-ac28-06256ac21c55	75	51.51789843744194	31.274442820953038	\N	15.45	226.79	9.39	2026-05-07 10:27:21.303773+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
839	1fe8f920-1f68-4f98-ac28-06256ac21c55	75	51.491839398533834	31.277913971641173	\N	59.63	164.48	10.82	2026-05-07 10:28:21.30414+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
840	1fe8f920-1f68-4f98-ac28-06256ac21c55	75	51.49039338978379	31.28331024728163	\N	11.59	265.85	5.16	2026-05-07 10:29:21.304486+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
841	1fe8f920-1f68-4f98-ac28-06256ac21c55	75	51.49864905033063	31.312602976375505	\N	11.08	78.44	14.41	2026-05-07 10:30:21.304843+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
842	1fe8f920-1f68-4f98-ac28-06256ac21c55	75	51.47539499970865	31.253067328971937	\N	89.67	179.72	7.05	2026-05-07 10:31:21.305234+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
843	1fe8f920-1f68-4f98-ac28-06256ac21c55	75	51.488802753805366	31.32596998116129	\N	33.25	275.19	3.1	2026-05-07 10:32:21.305654+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
844	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	75	51.46095320646165	31.278742563331075	\N	54.51	94.37	7.11	2026-05-07 10:18:21.306107+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
845	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	75	51.50002234964669	31.26728855636418	\N	86.16	350.59	6.53	2026-05-07 10:19:21.306558+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
846	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	75	51.486735154677405	31.28865030470154	\N	14.56	168.89	9.84	2026-05-07 10:20:21.306985+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
847	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	75	51.49615454991329	31.301465513403254	\N	44.84	238.71	6.85	2026-05-07 10:21:21.307344+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
848	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	75	51.4616850847657	31.262109310105053	\N	30.76	81.05	6.52	2026-05-07 10:22:21.307792+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
849	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	75	51.536139433070645	31.319845314225677	\N	38.9	20.07	10.23	2026-05-07 10:23:21.308345+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
850	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	75	51.51078827962529	31.307731681755598	\N	8.63	13.81	7.25	2026-05-07 10:24:21.308767+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
851	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	75	51.47681037438053	31.271839219766015	\N	65.42	95.65	2.69	2026-05-07 10:25:21.309066+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
852	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	75	51.471536601842764	31.278185817519176	\N	11.21	45.68	11.51	2026-05-07 10:26:21.309368+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
853	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	75	51.50913928006826	31.249924471920313	\N	78	59.88	12.36	2026-05-07 10:27:21.309722+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
854	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	75	51.513293062640706	31.319689120511025	\N	72.88	118.35	2.12	2026-05-07 10:28:21.310012+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
855	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	75	51.498104816611786	31.297196426296694	\N	60.09	125.97	9.68	2026-05-07 10:29:21.310324+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
856	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	75	51.46672839235979	31.273660672895627	\N	69.46	107.2	6.28	2026-05-07 10:30:21.310762+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
857	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	75	51.497900617098566	31.306364806443334	\N	45.37	100.93	6.2	2026-05-07 10:31:21.311178+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
858	1d9b95fd-005c-4ad4-ba63-4d468ce9e913	75	51.47867561097017	31.28836841211016	\N	46.64	284.21	9.36	2026-05-07 10:32:21.311657+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
859	f678ceb7-8984-4534-abac-41d262aaea81	75	51.529008236769926	31.31046474650964	\N	67.08	19.5	10.34	2026-05-07 10:18:21.312001+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
860	f678ceb7-8984-4534-abac-41d262aaea81	75	51.50136918089196	31.265099798610866	\N	64.71	359.97	4.81	2026-05-07 10:19:21.312336+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
861	f678ceb7-8984-4534-abac-41d262aaea81	75	51.52307904024376	31.303115666327567	\N	34.67	354.84	4.52	2026-05-07 10:20:21.312629+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
862	f678ceb7-8984-4534-abac-41d262aaea81	75	51.49019931922841	31.30420717775779	\N	25.04	345.11	12.94	2026-05-07 10:21:21.312896+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
863	f678ceb7-8984-4534-abac-41d262aaea81	75	51.46434232647667	31.294701523564505	\N	0.14	116.45	9.58	2026-05-07 10:22:21.313184+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
864	f678ceb7-8984-4534-abac-41d262aaea81	75	51.53257983636267	31.24963700260395	\N	25.06	51.77	11.75	2026-05-07 10:23:21.3135+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
865	f678ceb7-8984-4534-abac-41d262aaea81	75	51.49526543824485	31.287577376409068	\N	44.16	48.56	5.45	2026-05-07 10:24:21.313857+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
866	f678ceb7-8984-4534-abac-41d262aaea81	75	51.506479286764296	31.311863224324384	\N	56.98	104.94	14.79	2026-05-07 10:25:21.314205+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
867	f678ceb7-8984-4534-abac-41d262aaea81	75	51.4710142921481	31.321061085399002	\N	59.14	90.79	10.59	2026-05-07 10:26:21.314521+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
868	f678ceb7-8984-4534-abac-41d262aaea81	75	51.475256178171406	31.326304118222243	\N	5.11	298.34	7.92	2026-05-07 10:27:21.314803+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
869	f678ceb7-8984-4534-abac-41d262aaea81	75	51.53816896502666	31.322159942017787	\N	47.55	147.36	13.45	2026-05-07 10:28:21.315085+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
870	f678ceb7-8984-4534-abac-41d262aaea81	75	51.536467915584616	31.261638938008964	\N	19.41	268.98	11.38	2026-05-07 10:29:21.315373+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
871	f678ceb7-8984-4534-abac-41d262aaea81	75	51.5373541507189	31.260423437361037	\N	89.32	232.39	7.33	2026-05-07 10:30:21.3157+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
872	f678ceb7-8984-4534-abac-41d262aaea81	75	51.50722292938677	31.29209696390674	\N	87.62	49.55	3.73	2026-05-07 10:31:21.31602+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
873	f678ceb7-8984-4534-abac-41d262aaea81	75	51.532656785239325	31.296371359216888	\N	10.58	102.19	11.69	2026-05-07 10:32:21.316314+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
874	6a4d89a0-f357-4ccf-bfaf-8f0791f22192	75	51.49651771351044	31.302118429445102	\N	8.59	209.38	10.81	2026-05-07 10:18:21.31661+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
875	6a4d89a0-f357-4ccf-bfaf-8f0791f22192	75	51.520082922870756	31.311646457191973	\N	67.87	67.44	11	2026-05-07 10:19:21.316891+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
876	6a4d89a0-f357-4ccf-bfaf-8f0791f22192	75	51.4673298137176	31.297162235416966	\N	44.27	45.38	10.32	2026-05-07 10:20:21.317184+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
877	6a4d89a0-f357-4ccf-bfaf-8f0791f22192	75	51.47634623425886	31.31907625643342	\N	30.9	254.35	11.15	2026-05-07 10:21:21.317468+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
878	6a4d89a0-f357-4ccf-bfaf-8f0791f22192	75	51.511797608634694	31.258365777624302	\N	55.09	92.94	14.28	2026-05-07 10:22:21.317833+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
879	6a4d89a0-f357-4ccf-bfaf-8f0791f22192	75	51.499954116360335	31.284525755219445	\N	74.1	76.64	8.95	2026-05-07 10:23:21.318194+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
880	6a4d89a0-f357-4ccf-bfaf-8f0791f22192	75	51.495590395016386	31.285861501542115	\N	86.46	200.06	10.76	2026-05-07 10:24:21.318544+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
881	6a4d89a0-f357-4ccf-bfaf-8f0791f22192	75	51.51459569984053	31.300964021968177	\N	45.47	254.26	9.22	2026-05-07 10:25:21.318929+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
882	6a4d89a0-f357-4ccf-bfaf-8f0791f22192	75	51.48270256068529	31.31657488847735	\N	73.77	299.52	6.47	2026-05-07 10:26:21.31935+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
883	6a4d89a0-f357-4ccf-bfaf-8f0791f22192	75	51.49673854052677	31.280926271604063	\N	30.73	164.56	2.7	2026-05-07 10:27:21.319849+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
884	6a4d89a0-f357-4ccf-bfaf-8f0791f22192	75	51.489529145813876	31.268836986929568	\N	75.26	327.26	13.21	2026-05-07 10:28:21.320221+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
885	6a4d89a0-f357-4ccf-bfaf-8f0791f22192	75	51.47490521153049	31.253171601965498	\N	61.6	223.95	14.91	2026-05-07 10:29:21.320574+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
886	6a4d89a0-f357-4ccf-bfaf-8f0791f22192	75	51.46703084729166	31.285187055422547	\N	46.65	338.58	6.9	2026-05-07 10:30:21.32106+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
887	6a4d89a0-f357-4ccf-bfaf-8f0791f22192	75	51.51439030111261	31.268177625138083	\N	0.99	7.06	14.06	2026-05-07 10:31:21.321655+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
888	6a4d89a0-f357-4ccf-bfaf-8f0791f22192	75	51.5102566611013	31.253464278437118	\N	79.63	204.46	4.3	2026-05-07 10:32:21.322029+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
889	9e479722-a3d7-4cc7-ad11-0f26f67be07b	75	51.5258140550085	31.28819774120323	\N	46.87	22.3	2.43	2026-05-07 10:18:21.322339+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
890	9e479722-a3d7-4cc7-ad11-0f26f67be07b	75	51.49138190547018	31.325998146735486	\N	6.7	115.24	13.82	2026-05-07 10:19:21.322692+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
891	9e479722-a3d7-4cc7-ad11-0f26f67be07b	75	51.47663425780686	31.32696536822027	\N	57.61	47.09	4.57	2026-05-07 10:20:21.322989+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
892	9e479722-a3d7-4cc7-ad11-0f26f67be07b	75	51.458246565713594	31.281289656867802	\N	4.83	228.8	3.28	2026-05-07 10:21:21.32329+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
893	9e479722-a3d7-4cc7-ad11-0f26f67be07b	75	51.524506123602436	31.297890893300657	\N	32.17	264.52	8.68	2026-05-07 10:22:21.323584+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
894	9e479722-a3d7-4cc7-ad11-0f26f67be07b	75	51.536016859239815	31.278333519561144	\N	12.26	219.77	13.44	2026-05-07 10:23:21.323877+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
895	9e479722-a3d7-4cc7-ad11-0f26f67be07b	75	51.48409730992211	31.30791394665224	\N	28.76	85.15	6.54	2026-05-07 10:24:21.324184+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
896	9e479722-a3d7-4cc7-ad11-0f26f67be07b	75	51.53817116513294	31.285425110092852	\N	34.93	281.65	12.49	2026-05-07 10:25:21.324485+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
897	9e479722-a3d7-4cc7-ad11-0f26f67be07b	75	51.47100883308724	31.280763859491447	\N	62.76	238.96	2.31	2026-05-07 10:26:21.324774+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
898	9e479722-a3d7-4cc7-ad11-0f26f67be07b	75	51.533228463081954	31.28599478511122	\N	0.4	84.69	14.46	2026-05-07 10:27:21.325077+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
899	9e479722-a3d7-4cc7-ad11-0f26f67be07b	75	51.50512875405231	31.253760122147515	\N	38.46	269.87	10.49	2026-05-07 10:28:21.325367+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
900	9e479722-a3d7-4cc7-ad11-0f26f67be07b	75	51.52667429959758	31.26287884403653	\N	5.46	209.86	5.97	2026-05-07 10:29:21.325682+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
901	9e479722-a3d7-4cc7-ad11-0f26f67be07b	75	51.53462272000486	31.318783464160372	\N	13.91	192.35	4.97	2026-05-07 10:30:21.325979+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
902	9e479722-a3d7-4cc7-ad11-0f26f67be07b	75	51.525275802668645	31.266782681758627	\N	45.34	9.44	3.09	2026-05-07 10:31:21.326271+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
903	9e479722-a3d7-4cc7-ad11-0f26f67be07b	75	51.46339499996616	31.283157518937863	\N	13	226.09	7.04	2026-05-07 10:32:21.326552+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
904	4775a9bb-25f8-49ef-a65b-1002543d468b	75	51.53114472748923	31.319635496946518	\N	76.69	280.76	13.54	2026-05-07 10:18:21.326835+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
905	4775a9bb-25f8-49ef-a65b-1002543d468b	75	51.50478102791771	31.269550375574955	\N	39.13	194.66	8.58	2026-05-07 10:19:21.327114+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
906	4775a9bb-25f8-49ef-a65b-1002543d468b	75	51.5287028030911	31.32788918988898	\N	29.72	304.99	13.81	2026-05-07 10:20:21.327392+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
907	4775a9bb-25f8-49ef-a65b-1002543d468b	75	51.46156428786019	31.310544772240565	\N	11.95	240.09	12.45	2026-05-07 10:21:21.327701+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
908	4775a9bb-25f8-49ef-a65b-1002543d468b	75	51.49499135854018	31.30009957897967	\N	27.25	246.1	7.86	2026-05-07 10:22:21.328005+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
909	4775a9bb-25f8-49ef-a65b-1002543d468b	75	51.51465011711586	31.30070766589838	\N	14.83	221.35	10.49	2026-05-07 10:23:21.328289+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
910	4775a9bb-25f8-49ef-a65b-1002543d468b	75	51.53432353583728	31.30651853667485	\N	45.33	294.89	13.59	2026-05-07 10:24:21.328571+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
911	4775a9bb-25f8-49ef-a65b-1002543d468b	75	51.508358257643614	31.258957169060857	\N	28.87	8.23	7.85	2026-05-07 10:25:21.328827+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
912	4775a9bb-25f8-49ef-a65b-1002543d468b	75	51.52701799340116	31.25401780380772	\N	19.48	11.68	14.46	2026-05-07 10:26:21.329094+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
913	4775a9bb-25f8-49ef-a65b-1002543d468b	75	51.51030550572842	31.284674820601285	\N	11.1	12.91	3.48	2026-05-07 10:27:21.329348+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
914	4775a9bb-25f8-49ef-a65b-1002543d468b	75	51.517041652589924	31.25392971087865	\N	46.34	334.99	5.77	2026-05-07 10:28:21.329599+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
915	4775a9bb-25f8-49ef-a65b-1002543d468b	75	51.46603518621249	31.312268673802933	\N	12.18	137.97	13.13	2026-05-07 10:29:21.32988+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
916	4775a9bb-25f8-49ef-a65b-1002543d468b	75	51.49076325655909	31.261919355082306	\N	71.98	274.95	12.41	2026-05-07 10:30:21.330141+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
917	4775a9bb-25f8-49ef-a65b-1002543d468b	75	51.53201937515049	31.289673561048872	\N	77.98	101.08	14.32	2026-05-07 10:31:21.330399+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
918	4775a9bb-25f8-49ef-a65b-1002543d468b	75	51.482439223098694	31.27738377529527	\N	86.06	48.32	7.93	2026-05-07 10:32:21.330656+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
919	75da1df1-ee2f-41b8-90da-385bfaa4d080	80	46.45815330013051	30.71574939506796	\N	10.63	47.82	12.72	2026-05-07 10:18:21.330916+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
920	75da1df1-ee2f-41b8-90da-385bfaa4d080	80	46.45401591937172	30.771034584666342	\N	25.21	104.88	10.11	2026-05-07 10:19:21.331195+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
921	75da1df1-ee2f-41b8-90da-385bfaa4d080	80	46.48292198561023	30.709479324691223	\N	32.35	225.9	8.1	2026-05-07 10:20:21.331491+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
922	75da1df1-ee2f-41b8-90da-385bfaa4d080	80	46.45036796862385	30.719200152712478	\N	84.87	300.2	12.05	2026-05-07 10:21:21.331849+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
923	75da1df1-ee2f-41b8-90da-385bfaa4d080	80	46.45593886312092	30.74289883907378	\N	65.54	325.35	9.26	2026-05-07 10:22:21.332175+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
924	75da1df1-ee2f-41b8-90da-385bfaa4d080	80	46.50980639582102	30.698859929821374	\N	2.37	176.99	5.12	2026-05-07 10:23:21.332497+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
925	75da1df1-ee2f-41b8-90da-385bfaa4d080	80	46.46863030183542	30.75747613016187	\N	0.72	107.56	4.32	2026-05-07 10:24:21.332805+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
926	75da1df1-ee2f-41b8-90da-385bfaa4d080	80	46.5066652420185	30.72933313527147	\N	56.66	271.97	7.14	2026-05-07 10:25:21.333115+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
927	75da1df1-ee2f-41b8-90da-385bfaa4d080	80	46.514347046007465	30.749941956463637	\N	59.82	52.6	13.82	2026-05-07 10:26:21.333417+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
928	75da1df1-ee2f-41b8-90da-385bfaa4d080	80	46.4936718843695	30.712573904560205	\N	89.75	91.83	6.16	2026-05-07 10:27:21.333735+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
929	75da1df1-ee2f-41b8-90da-385bfaa4d080	80	46.50046957888074	30.766357327885917	\N	37.41	294.47	10.68	2026-05-07 10:28:21.334036+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
930	75da1df1-ee2f-41b8-90da-385bfaa4d080	80	46.493098301850466	30.7491266693271	\N	83.19	28.89	11.14	2026-05-07 10:29:21.334337+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
931	75da1df1-ee2f-41b8-90da-385bfaa4d080	80	46.47257256356727	30.718444746313452	\N	29.81	257.77	9.8	2026-05-07 10:30:21.334645+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
932	75da1df1-ee2f-41b8-90da-385bfaa4d080	80	46.47944364363614	30.701891419926604	\N	41.98	171.83	6.98	2026-05-07 10:31:21.334938+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
933	75da1df1-ee2f-41b8-90da-385bfaa4d080	80	46.51084560021937	30.723950464199707	\N	28.88	99.88	6.31	2026-05-07 10:32:21.335226+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
934	3863dff1-509f-4519-b6fd-eb60da807648	80	46.45024311152535	30.716127361965423	\N	4.43	298.42	3.22	2026-05-07 10:18:21.335514+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
935	3863dff1-509f-4519-b6fd-eb60da807648	80	46.49063663383368	30.705017756758036	\N	5.84	10.88	6.27	2026-05-07 10:19:21.335778+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
936	3863dff1-509f-4519-b6fd-eb60da807648	80	46.477636480629265	30.705162236088835	\N	66.48	173.94	4.13	2026-05-07 10:20:21.336066+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
937	3863dff1-509f-4519-b6fd-eb60da807648	80	46.47083038904206	30.74783082695136	\N	47.82	319.72	2.73	2026-05-07 10:21:21.336333+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
938	3863dff1-509f-4519-b6fd-eb60da807648	80	46.512431353528555	30.700547865790657	\N	63.75	57.98	5.44	2026-05-07 10:22:21.3366+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
939	3863dff1-509f-4519-b6fd-eb60da807648	80	46.444241969541366	30.758134457992565	\N	17.92	99.34	14.34	2026-05-07 10:23:21.336857+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
940	3863dff1-509f-4519-b6fd-eb60da807648	80	46.47149057764746	30.720626416786732	\N	22.54	348.06	7.78	2026-05-07 10:24:21.337119+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
941	3863dff1-509f-4519-b6fd-eb60da807648	80	46.48085080239236	30.697385024083395	\N	35.23	204.35	14.56	2026-05-07 10:25:21.337396+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
942	3863dff1-509f-4519-b6fd-eb60da807648	80	46.448503688635114	30.720686057200336	\N	50.02	184.34	13.85	2026-05-07 10:26:21.337745+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
943	3863dff1-509f-4519-b6fd-eb60da807648	80	46.45779776447887	30.722927099507896	\N	27.14	204.05	3.42	2026-05-07 10:27:21.338098+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
944	3863dff1-509f-4519-b6fd-eb60da807648	80	46.51312759820277	30.76953855041903	\N	86.97	231.93	11.38	2026-05-07 10:28:21.338397+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
945	3863dff1-509f-4519-b6fd-eb60da807648	80	46.45832459438794	30.73681520320846	\N	45.73	297.28	11.75	2026-05-07 10:29:21.33871+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
946	3863dff1-509f-4519-b6fd-eb60da807648	80	46.45867991443266	30.720856171079145	\N	51.59	254.73	10.87	2026-05-07 10:30:21.339019+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
947	3863dff1-509f-4519-b6fd-eb60da807648	80	46.45428520976498	30.758654986029367	\N	89.32	301.31	11.79	2026-05-07 10:31:21.339322+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
948	3863dff1-509f-4519-b6fd-eb60da807648	80	46.47267166389806	30.751275721195036	\N	58.24	107.16	9.71	2026-05-07 10:32:21.33972+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
949	9e508594-5002-4e3e-bb16-32d710685d17	80	46.50099668895327	30.74171848523117	\N	61.3	107.77	7.41	2026-05-07 10:18:21.340178+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
950	9e508594-5002-4e3e-bb16-32d710685d17	80	46.48849936639696	30.7664222502153	\N	44.81	14.78	2.23	2026-05-07 10:19:21.340505+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
951	9e508594-5002-4e3e-bb16-32d710685d17	80	46.48531930695559	30.716335609085288	\N	87	189.43	7.83	2026-05-07 10:20:21.340792+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
952	9e508594-5002-4e3e-bb16-32d710685d17	80	46.43749168802783	30.75542828477928	\N	43.61	197.63	4.43	2026-05-07 10:21:21.341084+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
953	9e508594-5002-4e3e-bb16-32d710685d17	80	46.510840419837386	30.743832924993043	\N	77.39	35.42	5.72	2026-05-07 10:22:21.341374+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
954	9e508594-5002-4e3e-bb16-32d710685d17	80	46.511124216161456	30.70832287085862	\N	78.6	233.29	12.81	2026-05-07 10:23:21.341676+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
955	9e508594-5002-4e3e-bb16-32d710685d17	80	46.490587451660446	30.723250988907832	\N	12.81	225.62	5.14	2026-05-07 10:24:21.341974+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
956	9e508594-5002-4e3e-bb16-32d710685d17	80	46.45264468765084	30.70731358001109	\N	12.75	186.07	10.66	2026-05-07 10:25:21.342264+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
957	9e508594-5002-4e3e-bb16-32d710685d17	80	46.45323211967226	30.765773624100248	\N	37.65	52.23	12.07	2026-05-07 10:26:21.34254+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
958	9e508594-5002-4e3e-bb16-32d710685d17	80	46.44922992417473	30.702252712279545	\N	67.72	20.97	12.87	2026-05-07 10:27:21.342815+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
959	9e508594-5002-4e3e-bb16-32d710685d17	80	46.48630222312934	30.749226952293448	\N	79.11	279.36	7.91	2026-05-07 10:28:21.343094+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
960	9e508594-5002-4e3e-bb16-32d710685d17	80	46.47607969500415	30.75673087748595	\N	60.86	71.03	14.24	2026-05-07 10:29:21.343359+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
961	9e508594-5002-4e3e-bb16-32d710685d17	80	46.47645379124998	30.702630259345334	\N	3.58	171.92	11.8	2026-05-07 10:30:21.343644+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
962	9e508594-5002-4e3e-bb16-32d710685d17	80	46.440070273364064	30.730895780641056	\N	73.03	177.45	10.99	2026-05-07 10:31:21.343922+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
963	9e508594-5002-4e3e-bb16-32d710685d17	80	46.48875252335967	30.72490808660355	\N	18.11	293.67	3.84	2026-05-07 10:32:21.344385+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
964	23fbee13-4b35-48cb-81fe-d0c496119765	80	46.4451216655657	30.73952022484536	\N	36.88	75.01	14.55	2026-05-07 10:18:21.344798+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
965	23fbee13-4b35-48cb-81fe-d0c496119765	80	46.441744678501394	30.75202262548179	\N	71.68	157.85	5.31	2026-05-07 10:19:21.34514+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
966	23fbee13-4b35-48cb-81fe-d0c496119765	80	46.48603153445014	30.71069762306682	\N	76.65	329.48	4.48	2026-05-07 10:20:21.345461+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
967	23fbee13-4b35-48cb-81fe-d0c496119765	80	46.46557688428404	30.704819245319907	\N	9.58	211.17	4.63	2026-05-07 10:21:21.345964+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
968	23fbee13-4b35-48cb-81fe-d0c496119765	80	46.505959481467926	30.753493638917007	\N	44.29	14.3	8.63	2026-05-07 10:22:21.346284+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
969	23fbee13-4b35-48cb-81fe-d0c496119765	80	46.50087915221281	30.71214690233562	\N	51.3	119.24	4.85	2026-05-07 10:23:21.346584+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
970	23fbee13-4b35-48cb-81fe-d0c496119765	80	46.51277479943898	30.757044266697527	\N	24.98	104.39	12.49	2026-05-07 10:24:21.346881+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
971	23fbee13-4b35-48cb-81fe-d0c496119765	80	46.43832312058932	30.749140313307713	\N	24.81	95.96	5.7	2026-05-07 10:25:21.34717+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
972	23fbee13-4b35-48cb-81fe-d0c496119765	80	46.49376549531266	30.750367371285403	\N	68.98	175.15	2.88	2026-05-07 10:26:21.347462+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
973	23fbee13-4b35-48cb-81fe-d0c496119765	80	46.506957759968	30.710484757680547	\N	27.56	1.42	13.99	2026-05-07 10:27:21.347747+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
974	23fbee13-4b35-48cb-81fe-d0c496119765	80	46.50344287097673	30.753222518917163	\N	16.84	196.97	4.69	2026-05-07 10:28:21.348041+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
975	23fbee13-4b35-48cb-81fe-d0c496119765	80	46.47036029579946	30.707482928209817	\N	51.82	57.92	12.91	2026-05-07 10:29:21.348334+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
976	23fbee13-4b35-48cb-81fe-d0c496119765	80	46.46021781609564	30.728980388978123	\N	17.59	91.63	9.66	2026-05-07 10:30:21.348607+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
977	23fbee13-4b35-48cb-81fe-d0c496119765	80	46.4731870464333	30.765198711951648	\N	68.36	123.15	6.37	2026-05-07 10:31:21.348878+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
978	23fbee13-4b35-48cb-81fe-d0c496119765	80	46.4610940806616	30.71675447310808	\N	2.81	117.32	11.94	2026-05-07 10:32:21.349139+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
979	7c691151-27c5-420c-80e8-cf4bedb56aa9	80	46.51492410628908	30.77166177167932	\N	31.67	26.96	7.09	2026-05-07 10:18:21.3494+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
980	7c691151-27c5-420c-80e8-cf4bedb56aa9	80	46.484872469953025	30.717106562535513	\N	1.95	280.85	6.47	2026-05-07 10:19:21.349667+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
981	7c691151-27c5-420c-80e8-cf4bedb56aa9	80	46.49383459064971	30.766605990107504	\N	51.39	348.04	12.61	2026-05-07 10:20:21.349938+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
982	7c691151-27c5-420c-80e8-cf4bedb56aa9	80	46.45132859691749	30.70194717654244	\N	65.65	342.55	4.73	2026-05-07 10:21:21.350202+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
983	7c691151-27c5-420c-80e8-cf4bedb56aa9	80	46.48151645060349	30.767751738505027	\N	19.15	94.06	10.65	2026-05-07 10:22:21.350457+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
984	7c691151-27c5-420c-80e8-cf4bedb56aa9	80	46.478793156669845	30.721211828324417	\N	8.43	118.63	6.66	2026-05-07 10:23:21.350721+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
985	7c691151-27c5-420c-80e8-cf4bedb56aa9	80	46.51620135328344	30.758939565205115	\N	53.34	280.87	13.5	2026-05-07 10:24:21.350982+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
986	7c691151-27c5-420c-80e8-cf4bedb56aa9	80	46.50895644254853	30.77026643582795	\N	20.06	22.69	10.57	2026-05-07 10:25:21.351246+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
987	7c691151-27c5-420c-80e8-cf4bedb56aa9	80	46.47445897564521	30.728141839455343	\N	38.88	172.69	7.89	2026-05-07 10:26:21.351524+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
988	7c691151-27c5-420c-80e8-cf4bedb56aa9	80	46.452172992288375	30.769349735165203	\N	36.91	25.83	5.06	2026-05-07 10:27:21.351795+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
989	7c691151-27c5-420c-80e8-cf4bedb56aa9	80	46.51253678802117	30.74817262691352	\N	5.67	9.16	8.14	2026-05-07 10:28:21.352068+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
990	7c691151-27c5-420c-80e8-cf4bedb56aa9	80	46.46150806495614	30.753857491818557	\N	62.96	108.96	6.22	2026-05-07 10:29:21.352534+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
991	7c691151-27c5-420c-80e8-cf4bedb56aa9	80	46.49550278176595	30.74999661555722	\N	4.66	26.26	6.9	2026-05-07 10:30:21.352892+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
992	7c691151-27c5-420c-80e8-cf4bedb56aa9	80	46.43778814553617	30.69476862679187	\N	10.3	12.17	14.66	2026-05-07 10:31:21.353202+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
993	7c691151-27c5-420c-80e8-cf4bedb56aa9	80	46.51726878932327	30.712141749202036	\N	86	137.51	10.08	2026-05-07 10:32:21.353504+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
994	131f4e80-eafd-458f-a070-9c51bb8605b7	84	50.472377861038176	30.555825126782256	\N	62.94	85.21	3.31	2026-05-07 10:18:21.353853+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
995	131f4e80-eafd-458f-a070-9c51bb8605b7	84	50.4249336839058	30.525698080167295	\N	67.98	325.33	4.53	2026-05-07 10:19:21.354166+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
996	131f4e80-eafd-458f-a070-9c51bb8605b7	84	50.439976844545335	30.53708592584964	\N	45.95	305.84	14.34	2026-05-07 10:20:21.354469+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
997	131f4e80-eafd-458f-a070-9c51bb8605b7	84	50.48457378198318	30.487468526246246	\N	51.49	69.96	9.97	2026-05-07 10:21:21.354762+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
998	131f4e80-eafd-458f-a070-9c51bb8605b7	84	50.46382500683807	30.548053427028695	\N	86.46	136.08	8.99	2026-05-07 10:22:21.355048+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
999	131f4e80-eafd-458f-a070-9c51bb8605b7	84	50.418321096114525	30.56147125707222	\N	71.14	234.98	6.84	2026-05-07 10:23:21.355346+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1000	131f4e80-eafd-458f-a070-9c51bb8605b7	84	50.467144364041566	30.561992176121837	\N	56.56	315.46	8.99	2026-05-07 10:24:21.355647+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1001	131f4e80-eafd-458f-a070-9c51bb8605b7	84	50.465776004181905	30.56055293443987	\N	58.3	45.46	3.05	2026-05-07 10:25:21.355928+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1002	131f4e80-eafd-458f-a070-9c51bb8605b7	84	50.417929859110565	30.527442225775715	\N	55.11	148.28	8.83	2026-05-07 10:26:21.356205+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1003	131f4e80-eafd-458f-a070-9c51bb8605b7	84	50.44151199809549	30.553159558282164	\N	61.6	43.83	12.04	2026-05-07 10:27:21.356529+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1004	131f4e80-eafd-458f-a070-9c51bb8605b7	84	50.486144510636294	30.507319312668194	\N	45.62	173.46	13.1	2026-05-07 10:28:21.356916+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1005	131f4e80-eafd-458f-a070-9c51bb8605b7	84	50.45925146430054	30.54800150864036	\N	11.26	289.89	7.51	2026-05-07 10:29:21.357244+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1006	131f4e80-eafd-458f-a070-9c51bb8605b7	84	50.42446315589546	30.52087581332972	\N	77.68	211.95	12.69	2026-05-07 10:30:21.357665+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1007	131f4e80-eafd-458f-a070-9c51bb8605b7	84	50.443775885476605	30.502810592478454	\N	27.25	203.63	2.33	2026-05-07 10:31:21.358009+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1008	131f4e80-eafd-458f-a070-9c51bb8605b7	84	50.441718371447365	30.517992337870844	\N	54.23	62.8	9.62	2026-05-07 10:32:21.358445+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1009	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	84	50.42716008910867	30.499141914924184	\N	38.92	233.9	13.28	2026-05-07 10:18:21.358948+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1010	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	84	50.47797875443614	30.53361686276209	\N	76.23	341.3	4.95	2026-05-07 10:19:21.359382+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1011	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	84	50.43690459230066	30.504624481344067	\N	88.38	218.09	11.19	2026-05-07 10:20:21.35975+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1012	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	84	50.41437621402453	30.546540617751393	\N	74.24	37.56	7.2	2026-05-07 10:21:21.360087+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1013	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	84	50.489211611319874	30.50570238492905	\N	30.56	298.61	8.8	2026-05-07 10:22:21.360425+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1014	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	84	50.48291643474541	30.529787847837287	\N	40.9	2.52	10.15	2026-05-07 10:23:21.360767+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1015	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	84	50.45485412982441	30.535864780226138	\N	7.49	167.85	8.09	2026-05-07 10:24:21.361087+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1016	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	84	50.48114178067523	30.538160905734813	\N	13.94	193.12	2.62	2026-05-07 10:25:21.361425+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1017	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	84	50.428460360666996	30.534042225665825	\N	4.62	194.89	9.73	2026-05-07 10:26:21.361772+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1018	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	84	50.46780442402491	30.529088232120618	\N	81.74	122.1	4.35	2026-05-07 10:27:21.362144+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1019	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	84	50.48984202289835	30.492874601389108	\N	1.41	97.66	12.21	2026-05-07 10:28:21.362464+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1020	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	84	50.41107864518632	30.526762763381583	\N	68.73	41.61	10.59	2026-05-07 10:29:21.362763+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1021	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	84	50.43336425317978	30.497958989459484	\N	89.34	100.85	12.18	2026-05-07 10:30:21.363073+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1022	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	84	50.434883185327045	30.501964090704103	\N	6.78	350.59	9.94	2026-05-07 10:31:21.363446+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1023	ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	84	50.466659207603776	30.54569651067212	\N	32.67	271.02	11.46	2026-05-07 10:32:21.363833+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1024	d5605709-91a8-4772-92c4-d4bd92b65717	84	50.4533209903675	30.488717339082157	\N	49.09	204.34	2	2026-05-07 10:18:21.364157+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1025	d5605709-91a8-4772-92c4-d4bd92b65717	84	50.47353056541257	30.56322708799318	\N	28.35	146.72	11.58	2026-05-07 10:19:21.364474+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1026	d5605709-91a8-4772-92c4-d4bd92b65717	84	50.44092484998177	30.52196935100041	\N	40.98	321.5	9.42	2026-05-07 10:20:21.364798+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1027	d5605709-91a8-4772-92c4-d4bd92b65717	84	50.464533096772136	30.496299400390203	\N	52.44	226.02	11.74	2026-05-07 10:21:21.365119+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1028	d5605709-91a8-4772-92c4-d4bd92b65717	84	50.436587999018194	30.520698370920666	\N	70.4	95.52	9.09	2026-05-07 10:22:21.365438+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1029	d5605709-91a8-4772-92c4-d4bd92b65717	84	50.44561289206663	30.540114378333534	\N	68.48	207.88	8.3	2026-05-07 10:23:21.365823+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1030	d5605709-91a8-4772-92c4-d4bd92b65717	84	50.464783846666734	30.501278653384723	\N	46.13	358.85	7.57	2026-05-07 10:24:21.366206+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1031	d5605709-91a8-4772-92c4-d4bd92b65717	84	50.42448288470094	30.507966103948185	\N	70.79	23.64	4.74	2026-05-07 10:25:21.366707+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1032	d5605709-91a8-4772-92c4-d4bd92b65717	84	50.48036061799658	30.55378292886168	\N	27.83	223.57	10.49	2026-05-07 10:26:21.367145+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1033	d5605709-91a8-4772-92c4-d4bd92b65717	84	50.424960752564026	30.553913088085416	\N	65.34	118.71	7.81	2026-05-07 10:27:21.367568+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1034	d5605709-91a8-4772-92c4-d4bd92b65717	84	50.42304948932024	30.516284860096395	\N	80.66	5.69	12.77	2026-05-07 10:28:21.367957+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1035	d5605709-91a8-4772-92c4-d4bd92b65717	84	50.427406073994504	30.505990785594346	\N	73.11	220.46	14.98	2026-05-07 10:29:21.368442+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1036	d5605709-91a8-4772-92c4-d4bd92b65717	84	50.4112778490994	30.519842175616883	\N	34.63	191.42	6.35	2026-05-07 10:30:21.368854+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1037	d5605709-91a8-4772-92c4-d4bd92b65717	84	50.4497899514911	30.54928763085185	\N	32.76	171.21	14.9	2026-05-07 10:31:21.369178+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1038	d5605709-91a8-4772-92c4-d4bd92b65717	84	50.45776971763945	30.55564465977283	\N	86.23	172.41	5.66	2026-05-07 10:32:21.369493+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1039	eff2cb75-547b-42eb-89f5-233a879ac549	84	50.47340850938081	30.488685925325495	\N	44.42	292.31	6.82	2026-05-07 10:18:21.369838+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1040	eff2cb75-547b-42eb-89f5-233a879ac549	84	50.43362286136174	30.5066300393066	\N	23.62	237.95	8.3	2026-05-07 10:19:21.370142+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1041	eff2cb75-547b-42eb-89f5-233a879ac549	84	50.461760236322995	30.523891927405213	\N	67.09	274.09	6.11	2026-05-07 10:20:21.370439+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1042	eff2cb75-547b-42eb-89f5-233a879ac549	84	50.47434221602715	30.495192537044385	\N	26.95	164.12	14.38	2026-05-07 10:21:21.370745+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1043	eff2cb75-547b-42eb-89f5-233a879ac549	84	50.43764150433935	30.512056879983092	\N	11.49	316.57	9.99	2026-05-07 10:22:21.371044+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1044	eff2cb75-547b-42eb-89f5-233a879ac549	84	50.45993944803349	30.518610203276072	\N	38.42	220.5	3.87	2026-05-07 10:23:21.371342+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1045	eff2cb75-547b-42eb-89f5-233a879ac549	84	50.41074636352072	30.50120044279019	\N	20.03	37.22	4.34	2026-05-07 10:24:21.371643+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1046	eff2cb75-547b-42eb-89f5-233a879ac549	84	50.45351417224355	30.55366081571944	\N	78.62	62.25	10.8	2026-05-07 10:25:21.371957+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1047	eff2cb75-547b-42eb-89f5-233a879ac549	84	50.44002139170573	30.51788602842199	\N	23.1	241.26	13.96	2026-05-07 10:26:21.372272+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1048	eff2cb75-547b-42eb-89f5-233a879ac549	84	50.41742768383859	30.520419162224194	\N	50.69	336.99	9.86	2026-05-07 10:27:21.372575+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1049	eff2cb75-547b-42eb-89f5-233a879ac549	84	50.46279570092339	30.525423514711132	\N	35.74	92.83	11.71	2026-05-07 10:28:21.372875+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1050	eff2cb75-547b-42eb-89f5-233a879ac549	84	50.43042642412057	30.53754861635339	\N	5.71	183.38	8.45	2026-05-07 10:29:21.373199+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1051	eff2cb75-547b-42eb-89f5-233a879ac549	84	50.413608270967494	30.492388424073628	\N	60.91	352.27	11.02	2026-05-07 10:30:21.373509+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1052	eff2cb75-547b-42eb-89f5-233a879ac549	84	50.45474706055126	30.520093895157014	\N	68.65	7.51	12.06	2026-05-07 10:31:21.373855+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1053	eff2cb75-547b-42eb-89f5-233a879ac549	84	50.46276073192271	30.548461497054962	\N	87.11	4.26	8.36	2026-05-07 10:32:21.374171+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1054	d3c1b403-57b5-4f49-8828-7b31b580ddd0	84	50.414261262165056	30.523548147043382	\N	71.79	85.21	13.7	2026-05-07 10:18:21.374471+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1055	d3c1b403-57b5-4f49-8828-7b31b580ddd0	84	50.416007685708806	30.48669111892246	\N	56.98	99.04	12.66	2026-05-07 10:19:21.374765+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1056	d3c1b403-57b5-4f49-8828-7b31b580ddd0	84	50.41968196698941	30.54098273248665	\N	64.35	321.41	14.36	2026-05-07 10:20:21.375059+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1057	d3c1b403-57b5-4f49-8828-7b31b580ddd0	84	50.44098135472447	30.505908502284246	\N	65.46	352.2	5.38	2026-05-07 10:21:21.375348+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1058	d3c1b403-57b5-4f49-8828-7b31b580ddd0	84	50.466162119750294	30.49902129066092	\N	62.88	71.05	13.66	2026-05-07 10:22:21.375644+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1059	d3c1b403-57b5-4f49-8828-7b31b580ddd0	84	50.41758564501844	30.533253707816694	\N	60.21	160.68	14.6	2026-05-07 10:23:21.375945+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1060	d3c1b403-57b5-4f49-8828-7b31b580ddd0	84	50.46249202429023	30.526525230980535	\N	7.45	341.41	14.11	2026-05-07 10:24:21.376277+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1061	d3c1b403-57b5-4f49-8828-7b31b580ddd0	84	50.45732786735803	30.493755671701415	\N	76.33	204.95	5.26	2026-05-07 10:25:21.376575+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1062	d3c1b403-57b5-4f49-8828-7b31b580ddd0	84	50.41702878581548	30.556063103856193	\N	73.69	305.33	12.9	2026-05-07 10:26:21.376867+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1063	d3c1b403-57b5-4f49-8828-7b31b580ddd0	84	50.43052771497902	30.542220361664896	\N	44.57	358.96	13.22	2026-05-07 10:27:21.377156+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1064	d3c1b403-57b5-4f49-8828-7b31b580ddd0	84	50.449121262178906	30.519135357097493	\N	55.41	342.05	3.87	2026-05-07 10:28:21.37751+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1065	d3c1b403-57b5-4f49-8828-7b31b580ddd0	84	50.469439359254835	30.551408854827528	\N	18.5	333.53	10.36	2026-05-07 10:29:21.377879+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1066	d3c1b403-57b5-4f49-8828-7b31b580ddd0	84	50.471923354295356	30.520735586291238	\N	22.62	219.05	11.48	2026-05-07 10:30:21.378181+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1067	d3c1b403-57b5-4f49-8828-7b31b580ddd0	84	50.45180218548184	30.517075592313603	\N	46.38	94.15	5.63	2026-05-07 10:31:21.378494+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1068	d3c1b403-57b5-4f49-8828-7b31b580ddd0	84	50.46924895351487	30.511275822563775	\N	80.06	355.07	14.13	2026-05-07 10:32:21.378796+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1069	3a1de34e-544b-4b67-aa92-20e6cfb09e72	84	50.44651744796918	30.501516659596096	\N	42.85	349.08	7.51	2026-05-07 10:18:21.379103+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1070	3a1de34e-544b-4b67-aa92-20e6cfb09e72	84	50.431594294493486	30.49266576987462	\N	24.55	65.02	12.78	2026-05-07 10:19:21.379416+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1071	3a1de34e-544b-4b67-aa92-20e6cfb09e72	84	50.44485107496745	30.561609087685753	\N	85.25	96.7	14.34	2026-05-07 10:20:21.379733+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1072	3a1de34e-544b-4b67-aa92-20e6cfb09e72	84	50.4523559255467	30.54846001060286	\N	47.59	231.91	13.24	2026-05-07 10:21:21.380059+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1073	3a1de34e-544b-4b67-aa92-20e6cfb09e72	84	50.425972657594514	30.50172183685575	\N	43.22	309.44	2.37	2026-05-07 10:22:21.380381+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1074	3a1de34e-544b-4b67-aa92-20e6cfb09e72	84	50.48288124653617	30.51980847120045	\N	35.2	161.57	8	2026-05-07 10:23:21.380724+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1075	3a1de34e-544b-4b67-aa92-20e6cfb09e72	84	50.42431485125788	30.54424387198789	\N	61.3	45.26	13.77	2026-05-07 10:24:21.381039+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1076	3a1de34e-544b-4b67-aa92-20e6cfb09e72	84	50.45089478005843	30.50879481155703	\N	85.58	139.07	6.28	2026-05-07 10:25:21.381351+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1077	3a1de34e-544b-4b67-aa92-20e6cfb09e72	84	50.42407219744715	30.516636484302445	\N	43.44	103.8	2.22	2026-05-07 10:26:21.38168+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1078	3a1de34e-544b-4b67-aa92-20e6cfb09e72	84	50.4514348023165	30.492737820241675	\N	11.04	332.1	13.75	2026-05-07 10:27:21.381997+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1079	3a1de34e-544b-4b67-aa92-20e6cfb09e72	84	50.428766094458716	30.517680308479278	\N	28.92	318.16	9	2026-05-07 10:28:21.382304+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1080	3a1de34e-544b-4b67-aa92-20e6cfb09e72	84	50.44873661503489	30.553055360227393	\N	11.94	300.44	6.33	2026-05-07 10:29:21.382609+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1081	3a1de34e-544b-4b67-aa92-20e6cfb09e72	84	50.443674432383816	30.530234000947008	\N	65.53	194.36	9.31	2026-05-07 10:30:21.382916+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1082	3a1de34e-544b-4b67-aa92-20e6cfb09e72	84	50.41784495858691	30.50433265131629	\N	86.29	36.49	14.13	2026-05-07 10:31:21.383216+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1083	3a1de34e-544b-4b67-aa92-20e6cfb09e72	84	50.48473459709598	30.495564107320494	\N	12.64	191.74	8.93	2026-05-07 10:32:21.383512+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1084	e51f66b4-9055-42e2-8790-5063fe231566	84	50.48850829819036	30.54799436099062	\N	38.65	39.98	14.65	2026-05-07 10:18:21.383827+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1085	e51f66b4-9055-42e2-8790-5063fe231566	84	50.455201334790964	30.502625960288885	\N	78.76	152.43	5.04	2026-05-07 10:19:21.384129+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1086	e51f66b4-9055-42e2-8790-5063fe231566	84	50.47527487003642	30.515893116058407	\N	31.83	52.19	7.24	2026-05-07 10:20:21.384452+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1087	e51f66b4-9055-42e2-8790-5063fe231566	84	50.42454763933974	30.52981771703449	\N	31.57	20	2.11	2026-05-07 10:21:21.384754+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1088	e51f66b4-9055-42e2-8790-5063fe231566	84	50.433760999756046	30.561771574404986	\N	58.02	251.91	13.07	2026-05-07 10:22:21.38505+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1089	e51f66b4-9055-42e2-8790-5063fe231566	84	50.45457985673594	30.49975234836297	\N	6.8	28.36	8.81	2026-05-07 10:23:21.38534+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1090	e51f66b4-9055-42e2-8790-5063fe231566	84	50.435761778102055	30.505059792816272	\N	21.54	243.4	14.06	2026-05-07 10:24:21.385701+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1091	e51f66b4-9055-42e2-8790-5063fe231566	84	50.470392265958296	30.500891004656616	\N	82.07	131.33	11.63	2026-05-07 10:25:21.386023+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1092	e51f66b4-9055-42e2-8790-5063fe231566	84	50.438891530759946	30.527238696334926	\N	68.28	184.11	13.62	2026-05-07 10:26:21.386329+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1093	e51f66b4-9055-42e2-8790-5063fe231566	84	50.47839379297103	30.525975711984984	\N	51.13	40.1	9.21	2026-05-07 10:27:21.386645+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1094	e51f66b4-9055-42e2-8790-5063fe231566	84	50.477402737760556	30.555010786322374	\N	63.61	327.83	11.59	2026-05-07 10:28:21.386993+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1095	e51f66b4-9055-42e2-8790-5063fe231566	84	50.43109181552138	30.497224726561036	\N	0.44	315.99	9.22	2026-05-07 10:29:21.387334+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1096	e51f66b4-9055-42e2-8790-5063fe231566	84	50.436755406592084	30.48573477342765	\N	38.19	111.9	13.34	2026-05-07 10:30:21.387661+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1097	e51f66b4-9055-42e2-8790-5063fe231566	84	50.46453236229362	30.538532471372704	\N	50.19	6.16	4.32	2026-05-07 10:31:21.387975+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
1098	e51f66b4-9055-42e2-8790-5063fe231566	84	50.42692550274427	30.485871735136648	\N	43.4	181.03	14.47	2026-05-07 10:32:21.388291+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
\.


--
-- Data for Name: inventory_check_items; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.inventory_check_items (id, check_id, resource_id, book_quantity, actual_quantity, verified_at) FROM stdin;
\.


--
-- Data for Name: inventory_checks; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.inventory_checks (id, warehouse_id, created_by, status, started_at, completed_at, notes) FROM stdin;
\.


--
-- Data for Name: invite_tokens; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.invite_tokens (id, user_id, token_hash, expires_at, used_at, created_at) FROM stdin;
ef81a8e4-314d-484b-9d8f-8a0ac145d2f4	acc38b6c-89a4-4403-8cb2-916107ff017a	871128e40be0d9473da37803ccfbeb356997562a515e8afa8ca22e9b5c173243	2026-05-03 14:40:04.770944+00	\N	2026-05-02 14:40:04.770947+00
ca10c3cf-2168-4ddf-bdaf-23016c6b2303	78205adb-41af-4988-9cd4-6bee3337f215	63c796e15239a0f8850d10f2bd832546694e36ec5f129ee1f4ef2f13c19ea18c	2026-05-03 14:40:59.211567+00	\N	2026-05-02 14:40:59.211571+00
14aafae4-4d68-4579-aa4d-51135c1083be	bb96be7a-c867-48c2-a128-fb143bf6409a	83df395bb6ad4134ec5c4fdb1d82c201d816e4e8a98a3d76cffb14ca21dc369e	2026-05-03 14:42:35.46645+00	\N	2026-05-02 14:42:35.466455+00
59c6b419-ba1f-48df-ae09-74c98528ef38	744c25e6-ad82-4471-8cd6-8da1cb476d26	89d827ecadc6f76ad0c63ddf9f9aa429798f39304e229f6c949cda60152c2c7d	2026-05-03 14:44:04.576543+00	\N	2026-05-02 14:44:04.576547+00
f0a92dea-1928-4cd5-992b-d99a8b04a0c4	44382f55-1ab8-496b-a63e-631171e55486	e1c66a0f82bb1fd2b95dc82e9c4db44e1b113fa37cc3788d005ba4c14b125fa5	2026-05-03 14:44:53.773505+00	\N	2026-05-02 14:44:53.773509+00
52c47e40-6118-4897-9d3e-735399c53712	1baed435-4416-444a-9939-724c8858ab78	76574c946f8c2d19fd0aacae27e25a10fbc511c95af0f35a50f0f70a174cd952	2026-05-03 14:45:58.009172+00	\N	2026-05-02 14:45:58.009175+00
de291e12-f90d-401c-af75-86d45eedc89e	8aa95ca9-5e34-4ab6-a033-be734bf87629	1df786ffe48cb6e43292ae8910d370e0abb15088818dbc22f8f1b7f17d2e2858	2026-05-03 14:46:42.083569+00	\N	2026-05-02 14:46:42.083574+00
8d9810a3-a272-4157-928c-f644b45f9888	6b43bc82-4a5b-46cc-af09-6ef4adecd95b	2c86b5b8076377903d2ef5584b3419e4d8f7d5c2640c4a68e73c8acb469eccc5	2026-05-03 14:47:46.762618+00	\N	2026-05-02 14:47:46.762623+00
468afd5e-0cef-4945-a72f-4f34f5ccf4ee	b6635a10-c63a-4151-b9a4-ec34d35c99ee	7e805ebf324deecd4be900ed57f4c826c3b8d810ce89b390dd9dd4024b0dc5d4	2026-05-03 14:48:32.989562+00	\N	2026-05-02 14:48:32.989566+00
124e5d80-a39a-4cd1-83d2-e67ab3a5e766	05eb7a5c-0413-4a6b-ad6a-5ad545e70a33	f553fa3bd655e49549ca1cf389cb97bbe04f18f87389bdca0ffabb154e409c30	2026-05-03 14:49:14.932265+00	\N	2026-05-02 14:49:14.932268+00
b09055f2-9ba5-4103-bc2d-ab722215679f	c22f350c-6f18-4839-8d2d-79fd623e6a85	e2e5caad3b5c9911bfbb7ec65358a3aa8366b7b11ae235a3daa87421cdf188cb	2026-05-03 14:50:03.978286+00	\N	2026-05-02 14:50:03.978291+00
25fe5bb4-35b8-40cd-97ca-1e31534fe2b7	39f7d374-63c3-4b75-b412-1a3e03a618b2	b8810f4a72d0a67e43a9a2539b21e818c00d7709b5c205d9372fb1096add6e53	2026-05-03 14:50:48.926624+00	\N	2026-05-02 14:50:48.926629+00
476f2748-a528-49c3-9b8a-995d6037e5fa	200577bb-7a24-4a4e-a4a8-7146fb112d37	fec2c44e4400fb775fd8a71d5caa9be394b9c2eb97baefd9784c99c02c4809f3	2026-05-03 14:59:29.458677+00	\N	2026-05-02 14:59:29.45868+00
fa0c31ca-0cd5-404f-85db-e4e08ed26608	94579c72-8ab3-4faf-8d91-8d248206ebe6	6d8f60bf40330fe9c4b18538207dfec8957ad68448e756b895b9e53d7453281b	2026-05-03 15:05:33.139133+00	\N	2026-05-02 15:05:33.139136+00
8856873a-c40f-442c-a920-8364887ebc67	26de8e1e-4ea0-46f8-907a-4c7a3390904b	3484393def03814b801ea1444253a3a2f174315df6143ba138779e3c48fecaf2	2026-05-03 15:06:28.237129+00	\N	2026-05-02 15:06:28.237133+00
336ccbd4-b3e1-4d47-bdb1-84d8561a1a46	0a4bc519-6dc3-4f19-9ec1-ecff82c91d03	4e8e6821fc118d7d498bb3e8a36428ff893401a99ef45481d350aae7a65af262	2026-05-03 15:07:15.986339+00	\N	2026-05-02 15:07:15.986344+00
32f484e8-a309-48ed-8449-26e0a2583654	b4bb3fff-a72f-437e-8d2f-48ed65d12eda	b7a43eae174b2677c7b69935ef5676f1c9dfbc69eb1b3d980804c8f3a1086f6f	2026-05-03 15:08:23.020645+00	\N	2026-05-02 15:08:23.02065+00
f63838da-7d6e-48a4-a36c-530b813e6177	dab88a27-f886-47b0-add8-8569e62942fa	87ab0657713999d5299748af5721d310a5df3a9678d00ad24259c827a0267716	2026-05-03 15:09:01.670246+00	\N	2026-05-02 15:09:01.67025+00
1bfc0363-116e-48bd-a3a9-00b43ea89838	4817a99c-f2c5-4157-8e3c-eedc2241fbc2	4d76f725a6a683ddc369c217b93afacbe1d5908ad56b26f480ac4cffb2605942	2026-05-03 15:09:43.370984+00	\N	2026-05-02 15:09:43.370988+00
dc9af046-48ba-469e-b357-53377dfd6a12	2b080e9c-a604-4935-98c3-3bed5932986a	3d51e4d2a4a4620baeb2da49aed8f35295fbf649d7dafa4f73c8f1882f23c215	2026-05-03 15:10:29.852359+00	\N	2026-05-02 15:10:29.852363+00
63950f01-97cc-4f30-9fa9-5b0f4fa44c98	3cd8d6bf-fad2-45f6-b6fb-cfce134df557	8bce7a5d1fd379455152a012c57b15bade128600e85315b62112aa73b9e1a977	2026-05-03 15:11:27.367154+00	\N	2026-05-02 15:11:27.367159+00
1626507f-7f3c-49ef-bc2c-b71c808e5db0	471f0149-166b-4f29-803b-3c910fa974c1	1a303751676ce70fa9e3c97dc4437889fbe689bdbd09608f975fa07deef303b9	2026-05-03 15:12:09.115434+00	\N	2026-05-02 15:12:09.115439+00
92912018-15fa-49fa-a98a-3d3404556c46	389997eb-96b7-4e37-a6ab-2db5692a6255	63978f04da595a2dc828b8ec1bf3a83aeeaa96897ad62c159b77e76f3099dfb7	2026-05-04 13:30:14.310123+00	\N	2026-05-03 13:30:14.310126+00
45d9033f-7e0b-4764-a0a9-efef7cf2ee88	b1b097e6-452e-4b2a-b62f-ff8e543a59a0	97d83741a8d137bb6dd673094e76a1da0baae7710305b8b7099fbc599797c69f	2026-05-04 13:31:18.10543+00	\N	2026-05-03 13:31:18.105432+00
65e92bb4-9377-4a92-91f5-5da5d820c9d5	21833e68-101f-4378-a916-62c120a9f192	520c9f81e87c0550cb1790b1408ef7e755f6632eebaa18cb61d3584e4795bd28	2026-06-09 00:53:55.737727+00	\N	2026-06-08 23:53:55.737728+00
\.


--
-- Data for Name: maintenance_records; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.maintenance_records (id, vehicle_id, odometer_km, description, performed_by, created_at, cost_amount, document_url, driver_id) FROM stdin;
\.


--
-- Data for Name: notifications; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.notifications (id, user_id, tenant_id, type, title, message, related_id, is_read, created_at, read_at) FROM stdin;
75d93849-2b39-4025-92e5-e0c482212581	389997eb-96b7-4e37-a6ab-2db5692a6255	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	798d6084-4572-441c-9f52-3f05892de225	t	2026-05-03 17:19:42.491222+00	2026-05-03 17:27:42.200543+00
01a3da64-1178-462c-86be-fd3f9b5d2291	389997eb-96b7-4e37-a6ab-2db5692a6255	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	edb7fed9-ddc8-4360-915f-23b639e5bac2	t	2026-05-03 16:58:33.484792+00	2026-05-03 20:02:12.413909+00
49db8e86-0b86-418b-8afa-e5cb19443aa2	389997eb-96b7-4e37-a6ab-2db5692a6255	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	7e631304-8a13-4a8b-9605-e83fbd90e834	t	2026-05-03 17:49:31.43885+00	2026-05-03 20:02:12.413909+00
2d05cb4c-8f2e-4194-b118-627d44a950a6	389997eb-96b7-4e37-a6ab-2db5692a6255	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	9cd279eb-1626-4932-938c-63cd7f4829e2	t	2026-05-03 17:58:11.779184+00	2026-05-03 20:02:12.413909+00
93ea33a6-79c5-4675-9432-3840e88c798e	389997eb-96b7-4e37-a6ab-2db5692a6255	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Розподільчий центр "Київ-Північ" → Головний транзитний хаб "Центр". Перевірте деталі в розділі «Транспорт».	b05dd831-264e-4bed-a7a9-e3e3d06c1393	t	2026-05-03 19:49:15.090894+00	2026-05-03 20:02:12.413909+00
d85f15ff-fa0f-41ae-aab3-1c2503dea2ac	389997eb-96b7-4e37-a6ab-2db5692a6255	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Розподільчий центр "Київ-Північ" → Головний транзитний хаб "Центр". Перевірте деталі в розділі «Транспорт».	8e1ca805-ed59-4b6f-bded-84ad66ae0138	t	2026-05-03 20:01:53.84861+00	2026-05-03 20:02:12.413909+00
778aad4b-91a7-4895-b8cf-5fc982aeff4f	389997eb-96b7-4e37-a6ab-2db5692a6255	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Розподільчий центр "Київ-Північ" → Головний транзитний хаб "Центр". Перевірте деталі в розділі «Транспорт».	34af91e3-95bb-4633-a8ef-3ce9a7336400	t	2026-05-03 21:38:16.992706+00	2026-05-05 12:33:57.551485+00
7ead8090-d275-4d55-85f3-b1895863e5ce	389997eb-96b7-4e37-a6ab-2db5692a6255	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	704ce237-6254-415d-aadb-7089f825e48b	t	2026-05-05 12:33:45.654849+00	2026-05-05 12:33:57.551485+00
d1a46fc5-1787-496e-b359-03399ee8f584	389997eb-96b7-4e37-a6ab-2db5692a6255	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Розподільчий центр "Київ-Північ" → Головний транзитний хаб "Центр". Перевірте деталі в розділі «Транспорт».	8ff169fe-d4d5-41d5-ab14-355008cc9a0e	f	2026-05-05 18:32:10.679696+00	\N
94287c02-dae6-4a1e-a888-ca29fc7425ad	389997eb-96b7-4e37-a6ab-2db5692a6255	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	51af7bd7-5340-4970-ac2a-61bba55cf810	f	2026-05-06 21:04:40.189703+00	\N
bb4d0f79-2813-435c-b9a9-0df3ff7f45b2	389997eb-96b7-4e37-a6ab-2db5692a6255	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	570a39e6-bc00-4985-aa2c-bf24ed2d4375	f	2026-05-06 21:07:03.691891+00	\N
459f1de1-9264-4693-8c82-7da93e8b0df9	3e59277f-0c1a-4aa3-973c-cf3db1e19497	0997831c-654f-471b-934f-cedafbc54ea5	VEHICLE_MAINTENANCE	Тестове сповіщення	Це автоматично згенероване сповіщення від seed-скрипта.	\N	t	2026-05-07 03:33:21.041491+00	\N
118f9271-0388-4ff5-9965-794596fcb811	3e59277f-0c1a-4aa3-973c-cf3db1e19497	0997831c-654f-471b-934f-cedafbc54ea5	SUPPLY_APPROVED	Тестове сповіщення	Це автоматично згенероване сповіщення від seed-скрипта.	\N	t	2026-05-07 03:33:21.041491+00	\N
a23458a0-a875-476b-af8f-74f9123601f6	3e59277f-0c1a-4aa3-973c-cf3db1e19497	0997831c-654f-471b-934f-cedafbc54ea5	LOW_STOCK	Тестове сповіщення	Це автоматично згенероване сповіщення від seed-скрипта.	\N	t	2026-05-07 03:33:21.041491+00	\N
d7ce7015-fa2d-440c-aa88-be803411eb27	3e59277f-0c1a-4aa3-973c-cf3db1e19497	0997831c-654f-471b-934f-cedafbc54ea5	SUPPLY_APPROVED	Тестове сповіщення	Це автоматично згенероване сповіщення від seed-скрипта.	\N	t	2026-05-07 03:33:21.041491+00	\N
bdd4bdb9-5944-4931-aafa-9b16883b0cfe	3e59277f-0c1a-4aa3-973c-cf3db1e19497	0997831c-654f-471b-934f-cedafbc54ea5	NEW_USER	Тестове сповіщення	Це автоматично згенероване сповіщення від seed-скрипта.	\N	t	2026-05-07 03:33:21.041491+00	\N
d6dcb699-6bf5-4325-92c5-aa3fc1c203b4	3e59277f-0c1a-4aa3-973c-cf3db1e19497	0997831c-654f-471b-934f-cedafbc54ea5	LOW_STOCK	Тестове сповіщення	Це автоматично згенероване сповіщення від seed-скрипта.	\N	f	2026-05-07 03:33:21.041491+00	\N
37fd0940-c95d-4a3f-9021-72945f37d158	3e59277f-0c1a-4aa3-973c-cf3db1e19497	0997831c-654f-471b-934f-cedafbc54ea5	SUPPLY_APPROVED	Тестове сповіщення	Це автоматично згенероване сповіщення від seed-скрипта.	\N	t	2026-05-07 03:33:21.041491+00	\N
d7d9e98e-3346-44f6-9c8d-cd36320eda74	f03a133e-617b-4c6f-9bab-3560b041c358	0997831c-654f-471b-934f-cedafbc54ea5	SUPPLY_APPROVED	Тестове сповіщення	Це автоматично згенероване сповіщення від seed-скрипта.	\N	t	2026-05-07 03:33:21.041491+00	\N
759f2208-6fa4-4ced-980c-fb4f93118e18	f03a133e-617b-4c6f-9bab-3560b041c358	0997831c-654f-471b-934f-cedafbc54ea5	LOW_STOCK	Тестове сповіщення	Це автоматично згенероване сповіщення від seed-скрипта.	\N	f	2026-05-07 03:33:21.041491+00	\N
154d0c1c-b872-4bc1-b217-2b634aba37d2	f03a133e-617b-4c6f-9bab-3560b041c358	0997831c-654f-471b-934f-cedafbc54ea5	SUPPLY_APPROVED	Тестове сповіщення	Це автоматично згенероване сповіщення від seed-скрипта.	\N	f	2026-05-07 03:33:21.041491+00	\N
5e3d71a8-10be-4990-9d0f-6f544590c7b1	f03a133e-617b-4c6f-9bab-3560b041c358	0997831c-654f-471b-934f-cedafbc54ea5	NEW_USER	Тестове сповіщення	Це автоматично згенероване сповіщення від seed-скрипта.	\N	f	2026-05-07 03:33:21.041491+00	\N
4c5000eb-aa43-45bf-a02c-4cf05a6cd138	81810dfd-0d94-402c-a398-fa6d5559686e	0997831c-654f-471b-934f-cedafbc54ea5	LOW_STOCK	Тестове сповіщення	Це автоматично згенероване сповіщення від seed-скрипта.	\N	t	2026-05-07 03:33:21.041491+00	\N
85ec2f4b-7165-4f1d-94cb-d5d6abed8132	81810dfd-0d94-402c-a398-fa6d5559686e	0997831c-654f-471b-934f-cedafbc54ea5	VEHICLE_MAINTENANCE	Тестове сповіщення	Це автоматично згенероване сповіщення від seed-скрипта.	\N	t	2026-05-07 03:33:21.041491+00	\N
0b5f6acd-a7d4-41fd-bcc3-1929842ba466	81810dfd-0d94-402c-a398-fa6d5559686e	0997831c-654f-471b-934f-cedafbc54ea5	VEHICLE_MAINTENANCE	Тестове сповіщення	Це автоматично згенероване сповіщення від seed-скрипта.	\N	f	2026-05-07 03:33:21.041491+00	\N
9f88ab3d-0b31-46a5-8242-ad3c539efbdd	81810dfd-0d94-402c-a398-fa6d5559686e	0997831c-654f-471b-934f-cedafbc54ea5	LOW_STOCK	Тестове сповіщення	Це автоматично згенероване сповіщення від seed-скрипта.	\N	f	2026-05-07 03:33:21.041491+00	\N
2a23bfcf-a3fc-404c-a824-7bb01c37b46e	81810dfd-0d94-402c-a398-fa6d5559686e	0997831c-654f-471b-934f-cedafbc54ea5	SUPPLY_APPROVED	Тестове сповіщення	Це автоматично згенероване сповіщення від seed-скрипта.	\N	t	2026-05-07 03:33:21.041491+00	\N
7b10f862-d36b-469a-b231-f3d8f4024e61	81810dfd-0d94-402c-a398-fa6d5559686e	0997831c-654f-471b-934f-cedafbc54ea5	LOW_STOCK	Тестове сповіщення	Це автоматично згенероване сповіщення від seed-скрипта.	\N	t	2026-05-07 03:33:21.041491+00	\N
22ed7a18-1e70-40ab-90c2-7fbc68cbda3a	81810dfd-0d94-402c-a398-fa6d5559686e	0997831c-654f-471b-934f-cedafbc54ea5	VEHICLE_MAINTENANCE	Тестове сповіщення	Це автоматично згенероване сповіщення від seed-скрипта.	\N	f	2026-05-07 03:33:21.041491+00	\N
db4ff8ff-14bf-4b1f-828b-23d265af1080	dd8a8b66-805e-4b3f-92e5-ea158a88f421	0997831c-654f-471b-934f-cedafbc54ea5	NEW_USER	Тестове сповіщення	Це автоматично згенероване сповіщення від seed-скрипта.	\N	f	2026-05-07 03:33:21.041491+00	\N
0110ffd9-4deb-4ba6-9580-abad2eddacc9	dd8a8b66-805e-4b3f-92e5-ea158a88f421	0997831c-654f-471b-934f-cedafbc54ea5	SUPPLY_REJECTED	Тестове сповіщення	Це автоматично згенероване сповіщення від seed-скрипта.	\N	f	2026-05-07 03:33:21.041491+00	\N
4a28ccb2-090f-4024-b49b-7f0b6af00f91	dd8a8b66-805e-4b3f-92e5-ea158a88f421	0997831c-654f-471b-934f-cedafbc54ea5	SUPPLY_REJECTED	Тестове сповіщення	Це автоматично згенероване сповіщення від seed-скрипта.	\N	f	2026-05-07 03:33:21.041491+00	\N
5fb5d323-276b-4c70-88c0-1e77f9d4bb6b	dd8a8b66-805e-4b3f-92e5-ea158a88f421	0997831c-654f-471b-934f-cedafbc54ea5	VEHICLE_MAINTENANCE	Тестове сповіщення	Це автоматично згенероване сповіщення від seed-скрипта.	\N	f	2026-05-07 03:33:21.041491+00	\N
4ed01878-641f-43ca-9fc8-6bbee99c40fa	dd8a8b66-805e-4b3f-92e5-ea158a88f421	0997831c-654f-471b-934f-cedafbc54ea5	NEW_USER	Тестове сповіщення	Це автоматично згенероване сповіщення від seed-скрипта.	\N	t	2026-05-07 03:33:21.041491+00	\N
dc91f85c-4cb0-407c-9077-debfcf4f96bf	dd8a8b66-805e-4b3f-92e5-ea158a88f421	0997831c-654f-471b-934f-cedafbc54ea5	VEHICLE_MAINTENANCE	Тестове сповіщення	Це автоматично згенероване сповіщення від seed-скрипта.	\N	t	2026-05-07 03:33:21.041491+00	\N
6253fe75-b3d2-4a69-85e5-091fe48a5cb3	dd8a8b66-805e-4b3f-92e5-ea158a88f421	0997831c-654f-471b-934f-cedafbc54ea5	VEHICLE_MAINTENANCE	Тестове сповіщення	Це автоматично згенероване сповіщення від seed-скрипта.	\N	f	2026-05-07 03:33:21.041491+00	\N
43ed120f-267c-469f-bbd3-4b7220d83ebf	389997eb-96b7-4e37-a6ab-2db5692a6255	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	59069426-28fa-45c7-8d5b-54503c66165c	f	2026-05-17 00:42:45.712388+00	\N
7cf465bc-2471-4d6a-b07a-5aaea91fc996	389997eb-96b7-4e37-a6ab-2db5692a6255	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	af0bd3a8-de7d-4f6f-ac96-db08cd1fb0c4	f	2026-05-27 21:21:41.246078+00	\N
1505289f-a0c2-4977-9c90-81357d26d43d	389997eb-96b7-4e37-a6ab-2db5692a6255	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	9e717db0-b052-4834-a3f9-f8ec4279e19f	f	2026-05-27 21:49:18.557949+00	\N
5d54ab62-cb34-4854-ab98-e32027e6ce79	389997eb-96b7-4e37-a6ab-2db5692a6255	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	9bbd6df5-403d-4b06-828f-42b3979ec189	f	2026-05-27 21:52:09.686287+00	\N
bf902700-628e-4c2d-8b83-a15cce2142b0	4817a99c-f2c5-4157-8e3c-eedc2241fbc2	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	cfb9dcda-8034-49a7-82f1-8744022ee386	f	2026-05-28 08:12:02.156725+00	\N
f8ebfb98-9fe6-4397-a7d1-484043abcdd8	4817a99c-f2c5-4157-8e3c-eedc2241fbc2	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	103a7304-ed77-4f89-9ba5-1746383c6cba	f	2026-05-28 08:29:12.948487+00	\N
f02caac5-1813-4dfe-9576-15b6c4696c9d	0a4bc519-6dc3-4f19-9ec1-ecff82c91d03	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	a1b425ea-0777-4d0d-965b-5082a819db50	f	2026-05-28 12:11:44.017249+00	\N
ba8e38c8-5183-4277-9076-926c763dd9e2	0a4bc519-6dc3-4f19-9ec1-ecff82c91d03	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	a190d65f-0532-4580-945a-a1efbedb5d4a	f	2026-05-28 12:23:37.874905+00	\N
6c89e9f0-20e9-4f88-a7f1-d5234fb86546	b1b097e6-452e-4b2a-b62f-ff8e543a59a0	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	ad4a7da0-ca3b-4e94-a08a-160c9357931d	f	2026-05-28 12:24:09.015812+00	\N
3ea52f8a-b110-48c3-ba22-231d02f9f63a	0a4bc519-6dc3-4f19-9ec1-ecff82c91d03	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	7aee9783-f7fc-41e5-b369-b295796f2366	f	2026-05-28 13:07:34.024006+00	\N
d4f9c5be-e4b5-44e6-8251-720cbe2c281e	94579c72-8ab3-4faf-8d91-8d248206ebe6	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	2925f911-f566-47fc-9ffe-a9488127b7fc	f	2026-05-28 13:12:06.061187+00	\N
1e53bc8b-98bd-4bae-9460-d34f839e1bc9	389997eb-96b7-4e37-a6ab-2db5692a6255	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	ebd81f67-bf05-4329-8837-6a6c19dcd454	f	2026-05-28 13:12:40.469145+00	\N
b73dc4c3-e040-49eb-b852-079446022832	0a4bc519-6dc3-4f19-9ec1-ecff82c91d03	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	298b0158-fd71-4224-ba82-67f7d753e108	f	2026-05-28 14:31:54.584986+00	\N
2c7579ee-a4b9-462c-9a17-e1c50d087ceb	94579c72-8ab3-4faf-8d91-8d248206ebe6	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	d7f8a04d-1d50-40cf-9be4-1cabc516f15b	f	2026-05-28 14:32:28.728003+00	\N
ced7a002-55b2-4a97-acb4-51405f41c8d2	0a4bc519-6dc3-4f19-9ec1-ecff82c91d03	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	850d5366-3794-474d-b863-bf6beaeb2d36	f	2026-05-28 15:44:49.603349+00	\N
5cacf881-05b8-49f4-99c6-53c780d8b668	389997eb-96b7-4e37-a6ab-2db5692a6255	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	aedbe475-34da-4a74-801e-a3752873ef8f	f	2026-05-28 15:46:22.781766+00	\N
9baf64bf-9ab8-486c-9caf-915b8611496d	94579c72-8ab3-4faf-8d91-8d248206ebe6	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	afc7246a-4a83-44d2-81d7-278ebab69699	f	2026-05-28 15:46:27.867798+00	\N
b18a6495-9f39-418f-8a80-ee4acab85769	0a4bc519-6dc3-4f19-9ec1-ecff82c91d03	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	d07d5ba7-5520-4782-bd60-a296058a22e4	f	2026-05-28 15:48:26.020562+00	\N
3f9af95d-be91-49d8-abf1-4d3f5a218a8a	b1b097e6-452e-4b2a-b62f-ff8e543a59a0	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	d5c49bb6-339a-4a58-9ff0-c360087993b0	f	2026-05-28 16:14:32.909996+00	\N
426daea9-6751-45f4-927a-77ec9fbb44ed	0a4bc519-6dc3-4f19-9ec1-ecff82c91d03	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	a1093fd8-3282-473d-8713-bf8500ad2c5b	f	2026-05-28 20:58:22.180162+00	\N
50de9e68-07f0-4057-bbd8-e47e3bceaa6b	94579c72-8ab3-4faf-8d91-8d248206ebe6	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	59153e98-ae5d-45a0-8e0f-f3fd2fe0435d	f	2026-05-28 20:58:32.591834+00	\N
276ee7d5-bc2b-4cf0-a6ca-4dd3b473a53a	389997eb-96b7-4e37-a6ab-2db5692a6255	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	b44f2f01-9882-48ab-bc7c-8ea84cc30f9c	f	2026-05-28 20:58:42.263893+00	\N
bac13b4b-c632-4c61-8aff-1ffccae03d5f	b1b097e6-452e-4b2a-b62f-ff8e543a59a0	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	b1f16ff5-0fb9-4380-b573-d5ffde39ba59	f	2026-05-28 21:29:50.256171+00	\N
a3bfe5ae-d25a-464f-8a68-667a675c667b	4817a99c-f2c5-4157-8e3c-eedc2241fbc2	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	f7e1401b-8706-4858-9fa8-e4aa4cbba26e	f	2026-05-28 21:30:08.972535+00	\N
fa522cd8-affb-426b-ac7c-131ee36a4bf0	389997eb-96b7-4e37-a6ab-2db5692a6255	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	7160f8f6-77a7-4864-8d43-1ef6fbdde1b5	f	2026-05-29 07:33:02.551207+00	\N
23e9ec47-e8c8-499f-a2e1-cd0953e3e5f9	0a4bc519-6dc3-4f19-9ec1-ecff82c91d03	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	9329b2e7-5c15-447d-805b-b1f8d6d3bdf4	f	2026-05-29 07:56:59.705504+00	\N
8717d6fc-0b2d-45a0-94e6-81ca47125deb	389997eb-96b7-4e37-a6ab-2db5692a6255	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	a15772fc-002a-4e30-b90a-52c850191858	f	2026-05-29 07:57:08.427342+00	\N
220879e3-6982-45db-8ec0-017e7dd4b2f7	94579c72-8ab3-4faf-8d91-8d248206ebe6	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	4d1d8708-f792-45fe-9932-21de489833bf	f	2026-05-29 07:57:30.517276+00	\N
073f742a-fc07-4675-92ff-37337f493b8c	b1b097e6-452e-4b2a-b62f-ff8e543a59a0	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	d3ed0394-d72e-4d4b-8023-72fd2cf72b91	f	2026-05-29 08:35:34.727404+00	\N
d1be5b8c-c865-4dfd-9941-cd508272d68c	4817a99c-f2c5-4157-8e3c-eedc2241fbc2	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	84a00a91-1625-43db-83ac-78a9d92c0962	f	2026-05-29 09:18:19.051879+00	\N
9d74d279-0f6f-4594-b63a-5775c67ea9e1	0a4bc519-6dc3-4f19-9ec1-ecff82c91d03	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	6c320299-71b2-4811-9d1a-0210fa33d864	f	2026-05-29 14:59:26.620108+00	\N
c310b7b9-63e6-4de5-944c-16f99a35ca72	0a4bc519-6dc3-4f19-9ec1-ecff82c91d03	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	772a04e3-29ca-470c-b8f9-a54221d88de1	f	2026-05-29 21:50:55.311686+00	\N
29576ba0-1486-4bcb-874f-18ff5e4b9e05	94579c72-8ab3-4faf-8d91-8d248206ebe6	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	d8e0dc0f-a992-4782-8fae-1b50b90185e1	f	2026-05-29 22:05:13.469859+00	\N
26971150-4f9d-47f5-82d7-5a8bc3d91da7	0a4bc519-6dc3-4f19-9ec1-ecff82c91d03	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	c189fea9-e9ac-4722-8f9a-799037a4035c	f	2026-05-29 22:34:39.181581+00	\N
4da2c619-37f4-4d4f-81f8-029220c2fe56	0a4bc519-6dc3-4f19-9ec1-ecff82c91d03	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	d0e4a50b-3888-4c08-85b7-bbc9b2bec7d9	f	2026-05-29 22:48:54.595806+00	\N
83d1ea11-18df-494b-a154-9e4f7c8d48d2	0a4bc519-6dc3-4f19-9ec1-ecff82c91d03	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	a1d3a9d8-3a6d-411d-a1a1-c8dfbe0e71ad	f	2026-05-29 23:07:58.057623+00	\N
e5bfd72a-e10e-4b02-b1e2-ab7e850680a2	94579c72-8ab3-4faf-8d91-8d248206ebe6	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	9bd63b6e-2560-48cf-bf88-54b89904271b	f	2026-05-29 23:24:20.785537+00	\N
611efefd-225f-47be-bec2-61d9a2ce7f51	94579c72-8ab3-4faf-8d91-8d248206ebe6	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	00eaa206-0780-4c97-89c7-5245593a747e	f	2026-05-30 10:17:43.77586+00	\N
bc0d7a03-6d12-42f5-86b0-1d1538705953	94579c72-8ab3-4faf-8d91-8d248206ebe6	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	fb8f5382-6ba2-44ed-9170-c0ad9ead2d86	f	2026-05-30 12:48:21.560785+00	\N
ac7e1679-4857-49ac-92e9-38a159ca2e98	94579c72-8ab3-4faf-8d91-8d248206ebe6	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	84ac7e4e-76c9-4d76-a84e-87fa5adc65b8	f	2026-05-30 12:54:09.875111+00	\N
641682e5-c3e8-4f94-856a-c92540dd2c93	94579c72-8ab3-4faf-8d91-8d248206ebe6	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	57663c7f-c1e6-4e4c-b12a-b7273eb2f2be	f	2026-05-30 13:13:26.040957+00	\N
672633cc-ada4-4654-b48e-fca851424859	b1b097e6-452e-4b2a-b62f-ff8e543a59a0	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	c61746a4-50f6-48d6-907e-7bfffdd17776	f	2026-05-30 15:15:16.45424+00	\N
9153adb3-33ce-4933-b1a3-1b24956d6131	b1b097e6-452e-4b2a-b62f-ff8e543a59a0	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	SHIPMENT_ASSIGNED	🚚 Новий рейс призначено	Вам призначено рейс: Головний транзитний хаб "Центр" → Розподільчий центр "Київ-Північ". Перевірте деталі в розділі «Транспорт».	0dc954b2-401a-4e53-8b14-8d2e7a3eea16	f	2026-05-30 21:30:32.02243+00	\N
31e064cc-ff8e-4f1c-8a76-3961ae98fb83	744c25e6-ad82-4471-8cd6-8da1cb476d26	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	EQUIPMENT_REPORT	⚠️ Рапорт щодо майна: Електросамокат	Струтинський Марко Валерійович подав(ла) рапорт на списання: Електросамокат. Причина: Зламано / Пошкоджено.	d226fa15-98d1-4d58-9f6f-a76f8ea0a47c	f	2026-06-01 00:13:52.393568+00	\N
b29fd8f7-e5fc-474b-a269-0e85211ea78f	44382f55-1ab8-496b-a63e-631171e55486	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	EQUIPMENT_REPORT	⚠️ Рапорт щодо майна: Електросамокат	Струтинський Марко Валерійович подав(ла) рапорт на списання: Електросамокат. Причина: Зламано / Пошкоджено.	d226fa15-98d1-4d58-9f6f-a76f8ea0a47c	f	2026-06-01 00:13:52.399939+00	\N
75c0424c-717b-417f-884e-424845b470be	1baed435-4416-444a-9939-724c8858ab78	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	EQUIPMENT_REPORT	⚠️ Рапорт щодо майна: Електросамокат	Струтинський Марко Валерійович подав(ла) рапорт на списання: Електросамокат. Причина: Зламано / Пошкоджено.	d226fa15-98d1-4d58-9f6f-a76f8ea0a47c	f	2026-06-01 00:13:52.40275+00	\N
434bddc1-5d70-4811-b4a5-11ddca52fc32	8aa95ca9-5e34-4ab6-a033-be734bf87629	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	EQUIPMENT_REPORT	⚠️ Рапорт щодо майна: Електросамокат	Струтинський Марко Валерійович подав(ла) рапорт на списання: Електросамокат. Причина: Зламано / Пошкоджено.	d226fa15-98d1-4d58-9f6f-a76f8ea0a47c	f	2026-06-01 00:13:52.404304+00	\N
\.


--
-- Data for Name: refresh_tokens; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at) FROM stdin;
24ea21fb-4d72-4efb-8635-bd1d39991bc8	21833e68-101f-4378-a916-62c120a9f192	125c52411c351a3a920485462837dbf5255bc8b134b18f6c442f2634ef274602	2026-04-30 00:11:34.257682+00	\N	2026-04-23 00:11:34.257686+00
ff738017-5c1e-475d-8b28-d4f87e231075	33a91fae-503d-45d2-aa42-4fdf78fcc983	76bb20ac81c03f9cc2e12f40d128dd24d6ae901308da5c73c635c9f9785cfcdc	2026-04-30 00:35:35.212699+00	\N	2026-04-23 00:35:35.212703+00
df5cda5a-1dd5-4c39-b9c0-e661ab6ad323	33a91fae-503d-45d2-aa42-4fdf78fcc983	a66c14fbdfada5591af7654bb86605c0c33c7c5a61c75ab0fecd667cb6555cb9	2026-04-30 10:34:33.638207+00	\N	2026-04-23 10:34:33.638211+00
0ba33977-6018-4df1-9c23-31e4df8b0712	21833e68-101f-4378-a916-62c120a9f192	ff6625c9e37fb779df4f827be783df705ef79b7233a91cf8525c5eb328bc96ab	2026-04-30 15:08:21.929854+00	\N	2026-04-23 15:08:21.929857+00
b12d84b1-812c-4b38-841c-e5cf174b6b8a	33a91fae-503d-45d2-aa42-4fdf78fcc983	ac40cf0f597783ef314e917530aaf16b4c4ed075f432633c36304fcbb4025537	2026-04-30 16:01:02.964365+00	\N	2026-04-23 16:01:02.96437+00
07dbdb37-cb00-4522-b24c-ed2696a3df36	33a91fae-503d-45d2-aa42-4fdf78fcc983	96d1f4ec951cf7f799c9f324e258f12ddaf8b7a671921438d4acdd4e71abb2f9	2026-04-30 16:10:28.089104+00	\N	2026-04-23 16:10:28.089109+00
a1cd6ad1-6cee-43d0-924f-0ad19082694c	33a91fae-503d-45d2-aa42-4fdf78fcc983	879a4063fb109d60ce992df714c942fc222172482cc814087d6597403e837f30	2026-04-30 17:22:58.473872+00	\N	2026-04-23 17:22:58.473875+00
fb77d8e6-6adb-47e9-8270-48eace1dffcf	33a91fae-503d-45d2-aa42-4fdf78fcc983	155483945e6e6b1e639370765a1522fea1566b83d9157ae789835c49718a2567	2026-04-30 18:24:19.679966+00	\N	2026-04-23 18:24:19.679972+00
af9aeb0a-15c2-4a57-8c53-517bc435605c	21833e68-101f-4378-a916-62c120a9f192	887bdf2a52940f44a46ac5c4b2adcd9b43dbd0cc865e90933f2772547bfbc875	2026-04-30 20:50:20.769447+00	\N	2026-04-23 20:50:20.769451+00
0e18433b-b6ac-45d2-81b0-18799e160105	33a91fae-503d-45d2-aa42-4fdf78fcc983	6db645470fd887fc3bcb9bf4cf4fe5316f598eaf8f8ee803652961fd4be5d78e	2026-04-30 22:02:19.806651+00	\N	2026-04-23 22:02:19.806655+00
28c5f723-803a-42db-82e1-eaae664f42c4	21833e68-101f-4378-a916-62c120a9f192	5b50aa579887dc3d53eda33573508cf2b7a5f58db2b1a4582d12c6fb8eb591db	2026-04-30 22:06:57.080067+00	\N	2026-04-23 22:06:57.08007+00
2cd14166-cec9-484b-acbc-0f22cab447b2	33a91fae-503d-45d2-aa42-4fdf78fcc983	1df0887cf9ecb9ab38e0ce8d8dae5f4cf9b1c7e47cfde1d11ddfbcb344be72f3	2026-05-01 14:35:49.989868+00	\N	2026-04-24 14:35:49.989872+00
214d63eb-043d-4e1f-926e-f7cca0ae8370	33a91fae-503d-45d2-aa42-4fdf78fcc983	a32fdbdef6c22c913785076187bef8ea1d64dff7294b29cda94ce693aebf0b4d	2026-05-01 18:15:45.414402+00	\N	2026-04-24 18:15:45.414405+00
259e8af4-a89a-4517-8024-a111026f8331	21833e68-101f-4378-a916-62c120a9f192	823e18305086237ac47fbbacc827d51138f8c70ce326dc49e7e02b33a4e439f7	2026-05-01 18:16:18.662652+00	\N	2026-04-24 18:16:18.662656+00
f26d928f-9a61-4531-8ae8-71fc32db5f7a	21833e68-101f-4378-a916-62c120a9f192	5906060bbdf5dbc0600d639325b7aa23afad4dc8f1aac5cddad7523314a41590	2026-05-01 18:27:47.628797+00	\N	2026-04-24 18:27:47.628802+00
d99cf64f-0999-481d-b9ea-83fc7677e6ee	33a91fae-503d-45d2-aa42-4fdf78fcc983	a82c2d0e907891bf54b4ec66f8931132f0af1ce8e38f9a4b57dfae02e9eba35b	2026-05-01 18:49:23.674794+00	\N	2026-04-24 18:49:23.6748+00
2c43002b-31e3-4d20-af64-90125a6228aa	21833e68-101f-4378-a916-62c120a9f192	8ac373554a4c0b0067d22d817a2a6526027ecf8751da41ec803319d7ba6349f4	2026-05-01 18:49:52.533881+00	\N	2026-04-24 18:49:52.533888+00
993f88b0-28a6-4679-ba1d-9539cca08d68	21833e68-101f-4378-a916-62c120a9f192	cdafcfdecdd8970efaa7c9bba10a6c9f41bec1d4bb42b0f76b2fc896d2ede72f	2026-05-01 18:56:19.336357+00	\N	2026-04-24 18:56:19.336363+00
6c6c8b9f-c34c-4456-8819-7ad1142f4fb5	21833e68-101f-4378-a916-62c120a9f192	a2605470122261a05cd6b61003eb14c7e35b09b243f428ab3c82a538d6438542	2026-05-02 15:23:24.350518+00	\N	2026-04-25 15:23:24.350522+00
a88fc6c8-cd93-4e2c-aa90-8e88b94d2925	33a91fae-503d-45d2-aa42-4fdf78fcc983	ad113e8fd17ae11f7c2a65b9e74f1fe28aa70b3949e64d8657fb36e526adcff5	2026-05-02 15:23:44.859041+00	\N	2026-04-25 15:23:44.859044+00
471a883c-d9cf-4ae4-ab0c-1c9721719632	21833e68-101f-4378-a916-62c120a9f192	d09d830538959073cb8fb4c683be2cf2df6784a0c63e45d32c23744afa779a7d	2026-05-02 15:24:34.138132+00	\N	2026-04-25 15:24:34.138135+00
c19ba20c-13e5-4216-b8ab-2c057bdbed63	21833e68-101f-4378-a916-62c120a9f192	767ea8339f6f4f731949b53e9545fb7afacc1e3cbf202317f6b408ace730a5ed	2026-05-02 16:21:31.143576+00	\N	2026-04-25 16:21:31.143581+00
cf771624-0c7a-4fa3-9702-3ebe1c6e5ee5	33a91fae-503d-45d2-aa42-4fdf78fcc983	62d4bc120057d3ddc8fc9f4effe959b0c2403174187638d72c697eedbb267e17	2026-05-03 20:20:55.516165+00	\N	2026-04-26 20:20:55.516169+00
f6fac676-cd49-429d-8360-e9d941c9a2f1	33a91fae-503d-45d2-aa42-4fdf78fcc983	b8d34d216039216e5435dc934396c3d7bd7898c3621ce97e60021ad0f4b2749b	2026-05-05 12:34:53.113152+00	\N	2026-04-28 12:34:53.113155+00
db051779-068a-4be2-8dc1-467b572e0434	21833e68-101f-4378-a916-62c120a9f192	2e6d4b3c43a3ebbe4314b74f268dc28111e9351a63c7610f5c24dfa77256b61a	2026-05-06 22:12:01.371051+00	\N	2026-04-29 22:12:01.371055+00
04b0aa1b-1e7c-4139-967c-4792e353e642	21833e68-101f-4378-a916-62c120a9f192	cfa42615bad352ada4a6a95ec52e71bf703d340bc21cfec9a9fde74e16867b6b	2026-05-09 14:24:07.899543+00	\N	2026-05-02 14:24:07.899549+00
6e50c6e3-0491-42e1-8cb7-1dea89d84e5e	21833e68-101f-4378-a916-62c120a9f192	a83947ef3134743670876d23449b7115e688470241c4c6c59fbb071a5a0b4fba	2026-05-10 13:04:11.982531+00	\N	2026-05-03 13:04:11.982534+00
fa9ee6fd-d89b-4b51-aa00-205fabab20fa	acc38b6c-89a4-4403-8cb2-916107ff017a	19404cf7aa80df5f2ebd210ba094050f576f9d936ede934aa073016dc4b5c2a6	2026-05-10 13:09:59.736222+00	\N	2026-05-03 13:09:59.736226+00
101f964c-ef3f-49f9-9bcb-1a2100243a50	acc38b6c-89a4-4403-8cb2-916107ff017a	4d53d3f7a479e8e2f06639563a461260b2b2137ff64d6116096a851a2de286c4	2026-05-10 13:29:19.676872+00	\N	2026-05-03 13:29:19.676875+00
972bc8fa-efb9-4fcb-926a-08f473ccd2a5	21833e68-101f-4378-a916-62c120a9f192	01fd7c96ca6893d55927d79ea8d953ad36510df05e025e2ea4b2d8e5bc416f1a	2026-05-10 13:32:44.898815+00	\N	2026-05-03 13:32:44.898818+00
736e057f-4157-4a49-a494-5f4916730139	389997eb-96b7-4e37-a6ab-2db5692a6255	953764437290953789d180bafb94d90d9dd0e68b85fb8a7957c37f28b679efb8	2026-05-10 15:41:41.742206+00	\N	2026-05-03 15:41:41.742209+00
c98d30ef-295a-4fd9-82ab-d33a0e2674f5	acc38b6c-89a4-4403-8cb2-916107ff017a	adb06665b1467f70129c43a98c4ac8f10b0602c45d286e5620635b8218b00dc4	2026-05-10 15:42:05.918882+00	\N	2026-05-03 15:42:05.918886+00
4413dbd3-a74c-495f-874d-7cc36fde5019	21833e68-101f-4378-a916-62c120a9f192	2b65d15681c42b71f26c3fadb01d890d7a69c93ed10ba44be436ace3655e8f0c	2026-05-10 15:42:14.120068+00	\N	2026-05-03 15:42:14.120071+00
0dcdd9ca-1aee-4d57-9936-bb159bca568e	389997eb-96b7-4e37-a6ab-2db5692a6255	0c350485715b9ebad1918dfb2a20947850288c3535aa10b83734317c075e651a	2026-05-10 16:05:40.839941+00	\N	2026-05-03 16:05:40.839944+00
c91675bb-b891-46a8-951c-6208f0fff048	21833e68-101f-4378-a916-62c120a9f192	d449245fb817f248aa6a166c559a82b2bd82dd9bd18f45bbd6ec8f9f7781039d	2026-05-10 16:06:03.469521+00	\N	2026-05-03 16:06:03.469525+00
931ee2b1-d310-4819-8d82-b17366ebf5af	21833e68-101f-4378-a916-62c120a9f192	f596a474075833bd2179e775d186e879eaf1d64fa79d8b980c0c949db5014f38	2026-05-10 16:19:31.076587+00	\N	2026-05-03 16:19:31.07659+00
370a9b20-63f6-4843-8a50-5c39ea9b7ff9	389997eb-96b7-4e37-a6ab-2db5692a6255	bb6767ea212f4e160ccb23d9d9a0c0756d0d3c3e0032e900cac856a4cf1c282f	2026-05-10 16:20:55.509122+00	\N	2026-05-03 16:20:55.509125+00
bdee1710-1d69-48e8-b380-a21867ce75d3	21833e68-101f-4378-a916-62c120a9f192	dfe3558aa622154c747af15238af446a5735204e6844852f8e80fdfc25c8c48d	2026-05-10 16:23:59.757203+00	\N	2026-05-03 16:23:59.757207+00
ea0835bb-e554-46a2-b0a5-423856ab2656	389997eb-96b7-4e37-a6ab-2db5692a6255	ac672f23d5e27443818936d6a4c25f2ac52a9354365f35ba4df269c2be9a56ef	2026-05-10 16:25:10.706255+00	\N	2026-05-03 16:25:10.706259+00
7cc51e4a-4df7-4ce2-a953-885f59eb47b6	21833e68-101f-4378-a916-62c120a9f192	916cd97dc4663895a46c16ff0df85fbf96b9d952f31c08f4fdfa66d73728f5a9	2026-05-10 16:38:14.246549+00	\N	2026-05-03 16:38:14.246552+00
d3baf820-5683-4db9-80b0-670c5f0f231d	389997eb-96b7-4e37-a6ab-2db5692a6255	96f18ee2a56b8b1eea15142378a38b19adcc5148a891ef20cd0b7b0ec3bc654a	2026-05-10 16:39:00.806922+00	\N	2026-05-03 16:39:00.806926+00
cfcb2109-136b-4b94-b794-f445023bc077	21833e68-101f-4378-a916-62c120a9f192	713f5b8bc7dccc3262fbd26d9e1e25aaaf03d1f5c8d6109f3da1e27771058254	2026-05-10 16:43:59.93495+00	\N	2026-05-03 16:43:59.934953+00
d2933975-b5c2-4406-ac8a-42e48e026641	389997eb-96b7-4e37-a6ab-2db5692a6255	d4ed33dbea228729609ba9a5cf60be87a036e9946b0c99ee61b0243f9289eac2	2026-05-10 16:44:52.790427+00	\N	2026-05-03 16:44:52.790431+00
83ded671-480a-40c0-a5cb-62459a41a0c3	21833e68-101f-4378-a916-62c120a9f192	5eab041cf590230aade19e0f9aaf6ee0d8278e838a8b175ba80f150991972a0f	2026-05-10 16:57:55.416385+00	\N	2026-05-03 16:57:55.416389+00
5f1c1e54-2c20-49d6-b7b8-b2b5481b87d1	389997eb-96b7-4e37-a6ab-2db5692a6255	162f1d54f75decd4104bdfb32de17403dfde509e6b7d30428a5cdf24fc91fd2f	2026-05-10 16:58:40.950964+00	\N	2026-05-03 16:58:40.950979+00
d61115b0-00a5-4ea8-9b39-edd7b7261e6e	21833e68-101f-4378-a916-62c120a9f192	3c9dfda7839c71e3c4d4944bf885ad422aae79f899c811823a18ceee6a283700	2026-05-10 17:08:34.886241+00	\N	2026-05-03 17:08:34.886245+00
d7fff4b4-36bc-417f-9bf7-4007929a4cea	389997eb-96b7-4e37-a6ab-2db5692a6255	c4ed29e34f245071d58fea66bec45aafc2f0c3361ba2ae2ced4b8224d937af5c	2026-05-10 17:19:53.192867+00	\N	2026-05-03 17:19:53.19287+00
72ab9aad-4b13-4c0b-8e24-d111d133f0dc	389997eb-96b7-4e37-a6ab-2db5692a6255	ef1615f1ecc758848b3244ccade7e539c4d813aaf2b511eadd3ec9fb698f2aa1	2026-05-10 17:40:08.469492+00	\N	2026-05-03 17:40:08.469495+00
3626c37e-49b6-4c64-b164-e8e2bbe86dd2	21833e68-101f-4378-a916-62c120a9f192	1a93a42878707cbaa915216ee283148bc724fbeb272b64d33b7e379bdcc7df4b	2026-05-10 17:42:01.071965+00	\N	2026-05-03 17:42:01.071968+00
a3ccfb52-d56c-44eb-8db7-9a8fcf534e5f	389997eb-96b7-4e37-a6ab-2db5692a6255	caf970f8c79a3b406978383c356a7d94ab06d5accc6181399414ffea86057b15	2026-05-10 17:49:37.896009+00	\N	2026-05-03 17:49:37.896012+00
5ad4df2f-f224-4c16-ac14-a88e8651fa5f	21833e68-101f-4378-a916-62c120a9f192	a6e88b677f70e97ae8167ad89fd1b8742023545f131d7ebcf7c755b282eccd4e	2026-05-10 17:50:12.071337+00	\N	2026-05-03 17:50:12.071341+00
83afbe06-c791-4c39-81b2-44651331e8e0	389997eb-96b7-4e37-a6ab-2db5692a6255	a5c90fc516b5742f5794927de8a76f9bb1b01973f61582c901cbc8807b1166b1	2026-05-10 17:58:24.550518+00	\N	2026-05-03 17:58:24.550522+00
75b592d6-f35f-4923-ae92-12b4a2fc6d03	21833e68-101f-4378-a916-62c120a9f192	146a716c8be32fd0f1a8ca1317553a08fd456fbabf73c2398556e5d4f3f82c79	2026-05-10 18:05:58.590397+00	\N	2026-05-03 18:05:58.5904+00
f93c175e-b1ca-4dec-876d-7b5aecbe0d16	389997eb-96b7-4e37-a6ab-2db5692a6255	9b8af01b655060be3e11929fdc1e7bffffce6ffe3b73c722140d084404fde303	2026-05-10 18:06:28.558233+00	\N	2026-05-03 18:06:28.558236+00
5cec65ec-641c-4588-8498-d19235fee829	21833e68-101f-4378-a916-62c120a9f192	a9993b2d32af7be74bd3ff88ccc3c583c2fc300a2f997898c32832ccebdf2237	2026-05-10 18:12:00.535609+00	\N	2026-05-03 18:12:00.535612+00
3dc9ee92-eded-4cbc-b709-cdf2c333e3e4	21833e68-101f-4378-a916-62c120a9f192	5d452fb0d3d6932e7dff012b2b9350a39c751e7b053893e0eb80aeb648127911	2026-05-10 19:39:45.50779+00	\N	2026-05-03 19:39:45.50781+00
5cbfe823-6bc6-437e-a4df-fcb4a2cf48f0	389997eb-96b7-4e37-a6ab-2db5692a6255	c85684c4a76b79ff8f1d52b157cdea8c1eda9163ea576039d8f68164e4b5b550	2026-05-10 20:02:02.830645+00	\N	2026-05-03 20:02:02.83065+00
76ba70d2-7afc-4ef4-81dd-3fb1403bd80c	21833e68-101f-4378-a916-62c120a9f192	dc657253ff001e25866e8663c9815e6cf8f9fe6bd6a35cd10060d24c9095aa97	2026-05-10 20:02:24.218957+00	\N	2026-05-03 20:02:24.21896+00
242f17be-547a-4603-b672-0d4c3500dc25	21833e68-101f-4378-a916-62c120a9f192	8aca9614c962e869c7338f082e979fdb6d23a9558306ba734acc868901c8c690	2026-05-10 20:50:46.929532+00	\N	2026-05-03 20:50:46.929535+00
4c13fb38-9868-4766-8c4d-e0f4fc2d05d0	389997eb-96b7-4e37-a6ab-2db5692a6255	63da9aa25df3afe863470a0ce7beccdfcb0f2599a466238018a16001a0de5a01	2026-05-11 07:24:01.972973+00	\N	2026-05-04 07:24:01.972978+00
380eab9e-a402-4a0e-95ef-ffbc8e6b0b67	21833e68-101f-4378-a916-62c120a9f192	567d9bf3c3b340c5f6fd96623ee46b22b0f422ab3ec4c11ffe775dcbf4f2c728	2026-05-11 07:24:28.945324+00	\N	2026-05-04 07:24:28.945328+00
6a8fb36b-fb3b-430d-9539-6cedba5ecde6	389997eb-96b7-4e37-a6ab-2db5692a6255	88ca5223ca62206a3ac91ab6ee126da2bdb7edbfbc2c7af216e07fa9ba03d924	2026-05-12 12:33:53.243975+00	\N	2026-05-05 12:33:53.243978+00
0320b68d-e03b-43b5-975d-4f45806bda7b	33a91fae-503d-45d2-aa42-4fdf78fcc983	b8976ad2b1206c028a3dbec65f15662641dfd8de520426130905f457b1fd7545	2026-05-12 12:37:15.48936+00	\N	2026-05-05 12:37:15.489364+00
c24930bb-da2a-4712-8b78-fac617fa6a82	21833e68-101f-4378-a916-62c120a9f192	8193d0fe1f3e22eb3e60fb0e9828dd172a7871316e00d85e4d68ce0b0c5c9eea	2026-05-12 12:39:20.901042+00	\N	2026-05-05 12:39:20.901046+00
10a3abce-36b1-45ec-bce3-69c82c8ad40e	21833e68-101f-4378-a916-62c120a9f192	6a88adc79e192abe8ff0fca2ccc632ced82190d7898d8e13ac2d54341e6aae2b	2026-05-12 20:45:46.640841+00	\N	2026-05-05 20:45:46.640851+00
2f3671ae-6a2e-46a1-b6d4-caabdb5932a9	33a91fae-503d-45d2-aa42-4fdf78fcc983	375da9860464760597e3b8a8aaf5b19966d45839ef0c042c816512dd3d66a089	2026-05-13 15:18:54.235872+00	\N	2026-05-06 15:18:54.235876+00
cd7115e1-79db-4613-8812-b1c70793e0da	21833e68-101f-4378-a916-62c120a9f192	5643fd60c9a6dfbc25b0d01313ed6d067f53323b789398e805afa853be4be514	2026-05-13 15:19:08.624612+00	\N	2026-05-06 15:19:08.624615+00
cdd17b36-d5ac-4c26-bc90-e477aa8c6793	33a91fae-503d-45d2-aa42-4fdf78fcc983	3b0e30489255426fe302c7d53b5213e21137f55cab892794af1c491cf23bfd3e	2026-05-13 15:19:23.025168+00	\N	2026-05-06 15:19:23.025171+00
3f37bf8a-9626-4449-b4ca-7d03ac69b751	21833e68-101f-4378-a916-62c120a9f192	a9f90544969f526e82c1a777ae56e17469235de709d72c24ae1e55c8193ceeb8	2026-05-13 15:19:41.152442+00	\N	2026-05-06 15:19:41.152446+00
352af1ce-ecc1-4282-a31b-6df562418c47	21833e68-101f-4378-a916-62c120a9f192	c6560565aa3a6d702c3d359fef98db78363e87cbc7ae5c2f41560ca82eeaf277	2026-05-13 15:31:46.055595+00	\N	2026-05-06 15:31:46.055598+00
6e9bb5fc-ffd4-47d0-8faa-a9bc4149bc11	33a91fae-503d-45d2-aa42-4fdf78fcc983	2887b6f098d3ce3571560ed092edfe8b57d5b52ddb47af70d67093415eb2eeb7	2026-05-13 15:56:25.632693+00	\N	2026-05-06 15:56:25.632697+00
8675d4c3-3491-4783-ba8c-1a24cea2951e	21833e68-101f-4378-a916-62c120a9f192	aaef82ae57d9244d76f7cb6632b6f79b5fcfcfb095fee5eee6303aa714d77083	2026-05-13 15:56:43.335054+00	\N	2026-05-06 15:56:43.335057+00
31e6704c-20f7-4ccc-82e3-b85ced611d71	33a91fae-503d-45d2-aa42-4fdf78fcc983	191052ff2af3fbc2bc5e33343244185e1cefd1803610cb940435354659e56229	2026-05-13 16:44:46.953664+00	\N	2026-05-06 16:44:46.953669+00
60946960-6503-42b3-bbfb-ddafbec0fabc	21833e68-101f-4378-a916-62c120a9f192	868630cf89992cf08a3e7015ad008beedf1c2788ab0946aec65d8b5a85763b93	2026-05-13 16:45:05.062908+00	\N	2026-05-06 16:45:05.062912+00
6d58ff17-8239-40d1-b203-9c237b76a156	33a91fae-503d-45d2-aa42-4fdf78fcc983	f83973744abcd59bb4349e8b1e0cb1eb60a0b85521d46d47f8c042b0f8199321	2026-05-13 16:45:34.550483+00	\N	2026-05-06 16:45:34.550487+00
7484f06a-7d97-42ec-acbe-814e88cd9b6b	21833e68-101f-4378-a916-62c120a9f192	c2c58cc930fe80f324a3c90ee31de79570c1adec92f542b0f5821066d455028b	2026-05-13 16:45:48.66012+00	\N	2026-05-06 16:45:48.660123+00
9bbf3f40-5f77-439f-b73a-8ed294c6ac9c	21833e68-101f-4378-a916-62c120a9f192	89bbf11076eb2ab39d029f5976717c335f11b9df7ee9c187660740d211b4aecc	2026-05-13 16:46:02.540516+00	\N	2026-05-06 16:46:02.54052+00
00453f4d-ded6-42e0-873d-6f33e92378fd	33a91fae-503d-45d2-aa42-4fdf78fcc983	95fb04560e41bb03da28e149ab1ad189ea4cd3e73807606522ef2bd1aa12fc19	2026-05-13 16:46:16.230556+00	\N	2026-05-06 16:46:16.230559+00
47308f5d-1fc4-4256-a6be-a2885e0dc717	21833e68-101f-4378-a916-62c120a9f192	91b6a79698a67f97159bc02702ee439f7f2961ee50d4822a1d3d422e0886b23a	2026-05-13 16:46:51.037341+00	\N	2026-05-06 16:46:51.037344+00
1987849e-a6d0-4f54-8cf4-3a31ca3b42e4	33a91fae-503d-45d2-aa42-4fdf78fcc983	d9a8391ac002cc5a21de337dfab1cd6c02546fbf71a6d3c2b12f283ba3afa215	2026-05-13 16:53:24.785911+00	\N	2026-05-06 16:53:24.785914+00
925b9993-2623-4b17-ab3e-f33a681faa87	21833e68-101f-4378-a916-62c120a9f192	fb7ed4f8a5660dfab066d616a9d0b853ad630c9d3b5b27300089c86853f70106	2026-05-13 16:53:47.286488+00	\N	2026-05-06 16:53:47.286491+00
992338ef-d327-4d7f-9658-2224a9a7cae1	21833e68-101f-4378-a916-62c120a9f192	86cdca1e64f5d819732fafb00ff5f7bc7b8bee5fe8ae6c577b6e8f98ce7a790d	2026-05-13 16:57:07.886759+00	\N	2026-05-06 16:57:07.886762+00
0fa2225e-2785-4ebd-b476-68edb5859e68	33a91fae-503d-45d2-aa42-4fdf78fcc983	14d5290386b78675f85a58a16095b06b767800169bd3657ebb9db54842096be9	2026-05-13 16:57:18.178668+00	\N	2026-05-06 16:57:18.178673+00
d7f23215-6354-46a4-8eed-1b3b1a575e82	21833e68-101f-4378-a916-62c120a9f192	c8275a3f4fb8ecc820b68bef804c50f90e2653952b681b4acf7d86bea874892e	2026-05-13 16:57:31.185341+00	\N	2026-05-06 16:57:31.185344+00
7048fe17-1b90-40f3-bd1b-23102b3a5a9c	33a91fae-503d-45d2-aa42-4fdf78fcc983	08bd33371cf3f7868b288c3ce719221e762cd2698df187178d8eaa91d6ecf1f4	2026-05-13 17:00:55.199992+00	\N	2026-05-06 17:00:55.199996+00
c183ecdf-bce7-4765-8209-2c1c5cebf459	21833e68-101f-4378-a916-62c120a9f192	14313d86d49d2679490d592b33c64f731913742b16a730870b017eab3fbb8c87	2026-05-13 17:01:41.860105+00	\N	2026-05-06 17:01:41.860116+00
5a65eb5c-d009-446d-a360-7dd418075ced	389997eb-96b7-4e37-a6ab-2db5692a6255	cb98b5278d966982df30d6a7e8b87691bdba6877cd931be09180765a8398fcb8	2026-05-13 17:45:56.508523+00	\N	2026-05-06 17:45:56.508529+00
6b5cc6f0-39a3-49aa-9240-7d7c221b0c04	21833e68-101f-4378-a916-62c120a9f192	04304127aa87c21d83c5e33e11bda6bc5a3c1e42dfa61357adc6ff6b812c7217	2026-05-13 18:14:58.902495+00	\N	2026-05-06 18:14:58.9025+00
bbabc2c4-c14e-453e-8185-4f090460a78b	389997eb-96b7-4e37-a6ab-2db5692a6255	6ccb79498d87f4e6458aa33e6f58356415eba3c7acebab7846ea41a37354dd25	2026-05-13 18:23:42.026043+00	\N	2026-05-06 18:23:42.026048+00
49bdea83-03af-4531-9e4b-27d59fa5133d	21833e68-101f-4378-a916-62c120a9f192	d5b88fa1579d9e9834e165c9bb17a21feba0ea874fb0e8fc91a8264a968725c7	2026-05-13 18:24:21.533593+00	\N	2026-05-06 18:24:21.533599+00
0c766d20-389b-4377-a41c-acf03942fb6f	21833e68-101f-4378-a916-62c120a9f192	6c4b297735116142af8c2c4ebfe46ab67643edd31976f59636c59557e0e317e7	2026-05-13 19:02:12.246017+00	\N	2026-05-06 19:02:12.246021+00
f01da621-a8ea-4edb-8e26-985db55b049e	389997eb-96b7-4e37-a6ab-2db5692a6255	3f35cd9c9469f10028be805ce9ac193b0caade04b34d70f5c8db572766a38c7c	2026-05-13 19:04:01.489024+00	\N	2026-05-06 19:04:01.489027+00
f25000d2-11c3-49f6-84cf-76504448d826	21833e68-101f-4378-a916-62c120a9f192	885a76cc45bc104c6282a88d334200cb73e3338743f9fcdcff4f6ce5ff89e5a4	2026-05-13 19:04:14.077742+00	\N	2026-05-06 19:04:14.077745+00
2091ab48-4abc-4c99-9032-f06fb404c3ba	8e4524be-8fb5-46a9-abb3-9d132f4fcda9	7401ebc9a9201e784c46c95ae6307933aea596427dc1fe5d9fd533fc74f8f08f	2026-05-14 10:36:23.610893+00	\N	2026-05-07 10:36:23.610898+00
5aab544d-8115-499f-bfc6-f6025c67dfe7	21833e68-101f-4378-a916-62c120a9f192	c67238eb868e21b37bcdf00698273d490886d82d16d1db24576ebac5b27ddc43	2026-05-14 10:39:32.536689+00	\N	2026-05-07 10:39:32.536692+00
0745a975-f838-48fb-8152-6ddd527cfbb0	21833e68-101f-4378-a916-62c120a9f192	e683fedffa09edb7e93a0cc17969d24d90d587ac88c5a728b9b4474e079ee221	2026-05-14 10:41:53.215142+00	\N	2026-05-07 10:41:53.215146+00
aec49176-166e-4b36-b061-9a5b1773e436	21833e68-101f-4378-a916-62c120a9f192	df3e1497c570e8a5cad5c6c7c1b9839298a6d955f86dce69ec7b6621932dc44b	2026-05-15 07:49:57.125932+00	\N	2026-05-08 07:49:57.125942+00
5cb06032-2eae-498d-ae46-2515e692b666	21833e68-101f-4378-a916-62c120a9f192	e5659a0a7d68ecc4bc6523eed82538fca9e9f5886bcfaedad76c379c2908a82d	2026-05-20 07:01:15.980559+00	\N	2026-05-13 07:01:15.980563+00
7c0df388-1d7b-409f-9b37-15946397ad5b	21833e68-101f-4378-a916-62c120a9f192	0c116b85f2f705fdb95bcc85b79f850e190e8ed9285e60619292155f1efbe8cb	2026-05-24 00:33:54.060338+00	\N	2026-05-17 00:33:54.060341+00
4c1d9748-edc8-4d7c-bca4-0fe5bbb7aca9	21833e68-101f-4378-a916-62c120a9f192	5fd3670cf507d97687e91023d7d2fd153755a4296d16723d642651b965f5a75a	2026-05-24 00:41:52.937567+00	\N	2026-05-17 00:41:52.937571+00
70fd402d-ec01-4637-94da-619c5f7be883	389997eb-96b7-4e37-a6ab-2db5692a6255	fab4b33a12a4df0910031981819be1a151740b7d8dde990f125dc4a60732dbd6	2026-05-24 00:43:44.39731+00	\N	2026-05-17 00:43:44.397313+00
a75c6c36-6ad1-4eaf-9448-636f3a20f494	21833e68-101f-4378-a916-62c120a9f192	ab81b35ecd3a2f585ebc07280f5ce424c5367128f3a11b239713ff98e5878a80	2026-05-24 00:43:58.420399+00	\N	2026-05-17 00:43:58.420401+00
4d77cf89-c4c2-48bf-bc67-c43a590dd36b	21833e68-101f-4378-a916-62c120a9f192	112fc0f395bb1b9ab1be067a1d3294cbe5c3f5866a6ac038ffb6afe58840d7a9	2026-05-24 00:46:30.32529+00	\N	2026-05-17 00:46:30.325294+00
e2740299-d8e7-4b2a-8b83-cfd6bf0c9971	21833e68-101f-4378-a916-62c120a9f192	8f3239f2b05cb7f9ff19bb99df1369fdfeea27f3af5b5abe10887ed9a640e42e	2026-05-24 11:43:16.439005+00	\N	2026-05-17 11:43:16.439009+00
3620fd8f-c30f-4475-84a8-b71c233027fc	21833e68-101f-4378-a916-62c120a9f192	5c095b64758bbb534e71318c3ba0953f7f98b7804a10cb889a3095cc9453a882	2026-06-03 12:49:32.329215+00	\N	2026-05-27 12:49:32.329216+00
8f145147-7515-4311-ad3d-51e3abf4f195	21833e68-101f-4378-a916-62c120a9f192	2e6ae7868bb68c4637683add4010ae42b47b94b676be2581f0aebe281b14938c	2026-06-04 07:42:57.426079+00	\N	2026-05-28 07:42:57.42608+00
e2f9a3a7-4300-416c-8a18-432e5c491427	21833e68-101f-4378-a916-62c120a9f192	b23c295c65df21efd0b594cccad20fc9363fb9ef61ec348d2635ca460af129b8	2026-06-04 13:05:54.270735+00	\N	2026-05-28 13:05:54.270736+00
98b48dc5-10aa-4736-ab56-effab88515db	21833e68-101f-4378-a916-62c120a9f192	49f921d1bb9802351eb0b9b664f6c1128771a3d5266af6abe43b48cd7ab74ecf	2026-06-04 14:28:06.68347+00	\N	2026-05-28 14:28:06.683472+00
9ee69628-0ffa-4317-93e4-c91a7bddb069	21833e68-101f-4378-a916-62c120a9f192	fcbec16dcfacd295acaa5587af439d490c3af5cee49eab0e8bb2a2683b4d3984	2026-06-06 15:10:30.709862+00	\N	2026-05-30 15:10:30.709864+00
8e3266b2-ab4e-4555-9f33-08de1b521153	21833e68-101f-4378-a916-62c120a9f192	dd0d2139221dad3d4475516aa8361c53ca5b0539ef5a9a481cfa1b2e4a94caa7	2026-06-15 09:20:30.002564+00	\N	2026-06-08 09:20:30.002565+00
01230038-ecb8-472e-a27d-753c28bd0b80	c301c093-ebc6-491a-9c40-eb425e472857	278ba0c31960dc55cb1fc373ed5cf6efa038b203e05dfb835b86b5370cc4258f	2026-06-15 10:18:11.09368+00	\N	2026-06-08 10:18:11.093682+00
4d471582-5029-434d-a8f2-e5a71c818418	21833e68-101f-4378-a916-62c120a9f192	ffd89247e58fb3640199a7c2938d9566793838ee468acb135bdac226f9d8728d	2026-06-15 12:31:17.26335+00	\N	2026-06-08 12:31:17.263351+00
5827b879-8cdd-43a8-8e43-a893a114e769	c301c093-ebc6-491a-9c40-eb425e472857	5ac40d53230b2e5f8270a27830b609455e563d943efd1e5bcbf0523ec6a0553e	2026-06-15 12:31:36.695739+00	\N	2026-06-08 12:31:36.69574+00
9738f7ea-8b5e-4959-8650-0aa022a5dc3f	21833e68-101f-4378-a916-62c120a9f192	349cc866b997de1d5cdaac58f807210c7219608940d239b281a52f0cc4545515	2026-06-15 12:32:46.537295+00	\N	2026-06-08 12:32:46.537296+00
0bca4495-2f7d-43c5-8bea-bc8111f4b3a5	c301c093-ebc6-491a-9c40-eb425e472857	e18297926a6c49af82d147f579b244989814f5f6e17848b0ddebcbd9c69ea00f	2026-06-15 13:04:25.893338+00	\N	2026-06-08 13:04:25.89334+00
d9def762-dd55-49c7-9704-319c5f9093bb	21833e68-101f-4378-a916-62c120a9f192	db8243e849feaca5a18e761426eb63a5053be1b65a35cafba77a2f8f6a366eac	2026-06-15 13:46:40.711647+00	\N	2026-06-08 13:46:40.711648+00
23d61679-5bc2-4023-92b9-23f039a5ff80	c301c093-ebc6-491a-9c40-eb425e472857	a87307c9a05815a220d6f862520fd132cbc027e6ef488354ad6b415d837c88a6	2026-06-15 13:48:17.598607+00	\N	2026-06-08 13:48:17.598608+00
dcc44c2d-9257-4d10-acd9-b25b78640186	21833e68-101f-4378-a916-62c120a9f192	b8c35ec9cc0064c763cfdf740e7b39389e7f3f477a19bc99914c935ee116a835	2026-06-15 13:48:50.387766+00	\N	2026-06-08 13:48:50.387767+00
272ffaa6-be4b-4265-bdc6-040f3ffd4c7b	c301c093-ebc6-491a-9c40-eb425e472857	7fb7d0c790b1499cb635f33901a14cc24d4cc6650a0486684dce1081165ccb36	2026-06-15 13:49:59.111062+00	\N	2026-06-08 13:49:59.111063+00
4cf41904-efee-42a5-a3fc-c1dc63e1a856	21833e68-101f-4378-a916-62c120a9f192	796a0b71c206e58f66574b0c2dd777ba9aa3af35e30a3c2e099059d9ca5a6832	2026-06-15 14:40:27.501845+00	\N	2026-06-08 14:40:27.501846+00
9d7a4651-0049-4d47-a26f-f71a6e1df5af	c301c093-ebc6-491a-9c40-eb425e472857	4e3637288eec86fedd0ec742ad53875408b68a3889697179679ecdfe087ddf64	2026-06-15 14:41:58.947795+00	\N	2026-06-08 14:41:58.947796+00
60c9ebf7-83e8-4320-848e-94135e1431bc	21833e68-101f-4378-a916-62c120a9f192	cf31c3d1889acb9eb9e3f038aeef7ff9439b680f1e7f238f644aaae99c0a559c	2026-06-15 14:42:10.291642+00	\N	2026-06-08 14:42:10.291643+00
7541ec80-313e-4690-ab83-82471efa2df4	a2fbe6df-f0aa-432e-a5df-6759d8b5dff6	d554cee209314ca11457b5d459587d03000cfe9455676b28302aa31f75253369	2026-06-15 14:53:59.862322+00	\N	2026-06-08 14:53:59.862323+00
58b5ebe4-f65f-4a69-9981-a01fd21b192b	a2fbe6df-f0aa-432e-a5df-6759d8b5dff6	e744998707faecd5b92badc189179f1a5f480a524256374ce1a9ff051a82b8bd	2026-06-15 14:55:03.842357+00	\N	2026-06-08 14:55:03.842359+00
66aea99f-81d1-483a-9f40-0744a2071973	a2fbe6df-f0aa-432e-a5df-6759d8b5dff6	288cef4f3a13410750046911a201929df58ab27d6cf0877e30bac0ec8b4da619	2026-06-15 14:56:30.980655+00	\N	2026-06-08 14:56:30.980656+00
a71a7cbf-f24d-4875-acc4-da00927b5f2c	a2fbe6df-f0aa-432e-a5df-6759d8b5dff6	4c5fa555277b64f5e90282623c49304a09249cac91a9a63919e59961d88370ac	2026-06-15 14:59:22.837634+00	\N	2026-06-08 14:59:22.837639+00
67c9bd93-aebf-42ee-a5e0-6b6ab8a8b3cf	dd8a8b66-805e-4b3f-92e5-ea158a88f421	229a2d97550b9cdfb72347662726bb9b9e0ed32329b1aad17fd5c8f0b7c33884	2026-06-15 15:34:07.817306+00	\N	2026-06-08 15:34:07.817307+00
a9d4f8be-2172-4b2b-8a76-e91ae2a38ca7	dd8a8b66-805e-4b3f-92e5-ea158a88f421	c17af8559d6626197c59fa514f6250d6d14645d437db4ca8100fae5a7910e8ad	2026-06-15 15:35:04.294051+00	\N	2026-06-08 15:35:04.294052+00
5991b94c-bef4-4c1b-99f8-2360cb8f225d	dd8a8b66-805e-4b3f-92e5-ea158a88f421	5c0c7fe76a0e198f6105cdc77bd1764a95243b6c86bb9738a28c4ddeba42cbd1	2026-06-15 15:40:38.883743+00	\N	2026-06-08 15:40:38.883744+00
5beb98c2-1770-4566-99d2-be1238c70a4d	dd8a8b66-805e-4b3f-92e5-ea158a88f421	0dcc1efebf6a97079bf63336ae07a2f00d2caeebbcd3823093417c195ad001b0	2026-06-15 15:53:06.937149+00	\N	2026-06-08 15:53:06.937151+00
b912afe7-9ec0-45fb-8204-ad9ee54adca5	dd8a8b66-805e-4b3f-92e5-ea158a88f421	705a8532ff21900572e9d50800e75fd27a47f123c52b2fa17abdad93eb1870dd	2026-06-15 15:53:31.770744+00	\N	2026-06-08 15:53:31.770745+00
8a79a4c9-db30-427a-a349-981b6410824f	c301c093-ebc6-491a-9c40-eb425e472857	51a7522050c61f3b6c297d8860c473e636e2dcd42d454d41e89c3202ca26269d	2026-06-15 15:58:49.532086+00	\N	2026-06-08 15:58:49.532087+00
40bb47a4-63f9-4f1d-8e48-16cb860efe4b	21833e68-101f-4378-a916-62c120a9f192	87b2c5168688d9e788d82c5ecceb74d181a56dfb37612dfa5acfba5d4d78e855	2026-06-15 16:00:59.392079+00	\N	2026-06-08 16:00:59.39208+00
0a8e7457-ec45-449c-812f-e9a372f37e69	dd8a8b66-805e-4b3f-92e5-ea158a88f421	265006ea6a046c815f17dd931031946409d859f5ca382276092b8f3d76ac534f	2026-06-15 16:37:28.422019+00	\N	2026-06-08 16:37:28.422021+00
1f798f2b-d835-43de-aa70-812298065c21	c301c093-ebc6-491a-9c40-eb425e472857	7ea9774cd50a6409835fa3103839139d7bbf74d2fd1c3d68e872546df7d411f9	2026-06-15 16:41:36.082582+00	\N	2026-06-08 16:41:36.082583+00
9152fe3b-bfa0-4c04-b6d0-4f72b796a608	c301c093-ebc6-491a-9c40-eb425e472857	7259c68bc8de6a397d56b56a41ac8d196cd4368382f044dae75ba24c20a123e5	2026-06-15 16:41:56.021867+00	\N	2026-06-08 16:41:56.021868+00
0c49096c-4745-4dac-8d07-0cecb536cb72	21833e68-101f-4378-a916-62c120a9f192	19e1ec1e8fc61c3b24b083972439e56640ae9ae1aa863f08790b5dcc5825e18d	2026-06-15 16:42:08.406009+00	\N	2026-06-08 16:42:08.406011+00
5b56ba0e-42f4-4791-9ec2-df3d87dfe02a	c301c093-ebc6-491a-9c40-eb425e472857	45656aaadf29fe3f88de0afdfeb23516a59d53cd5716a35b9645625d30ba7bb6	2026-06-15 16:43:14.36879+00	\N	2026-06-08 16:43:14.368791+00
d82138eb-66a9-4807-ac48-cff86eaa87f0	21833e68-101f-4378-a916-62c120a9f192	07c7fc050be5e5d94191be88ab211543810a1343aff5343fbbb6bf2bfb8dd6a3	2026-06-15 16:43:39.236515+00	\N	2026-06-08 16:43:39.236516+00
502a72b0-0913-4f21-83b4-5c9e5693b605	c301c093-ebc6-491a-9c40-eb425e472857	1fd61cb39c12d51ca760a60a00bca0659ca710ef314289689e1dd9631de42a5b	2026-06-15 16:44:39.392487+00	\N	2026-06-08 16:44:39.392488+00
726eb8c8-8a3f-46e5-a5cd-77805d1dee16	21833e68-101f-4378-a916-62c120a9f192	973b6ed59344845c6a3279044f8555380b9b77789ac94fea205a5c44341027a7	2026-06-15 16:44:49.658025+00	\N	2026-06-08 16:44:49.658026+00
29118d35-2291-43e6-a3cb-6b1d030f42d7	c301c093-ebc6-491a-9c40-eb425e472857	4d46d543a205a30e6dcb246b1649edca3af03196c7099f72673fb5aa3ce7a0ae	2026-06-15 16:45:57.690026+00	\N	2026-06-08 16:45:57.690027+00
8cff69cb-acfd-446b-a848-6e3ae53ac64a	21833e68-101f-4378-a916-62c120a9f192	c3799fb21661808be6ebed47a8e9901905efb9ee94611fd16d855cee6d2e715b	2026-06-15 16:55:24.374694+00	\N	2026-06-08 16:55:24.374695+00
78c751e3-048a-4393-b7d9-76af0844501f	c301c093-ebc6-491a-9c40-eb425e472857	66e3b11bf5c039910f5db2b26f82629237ccb9a2d7b4e2d86df1b0d665c1ca3b	2026-06-16 00:04:45.059682+00	\N	2026-06-09 00:04:45.059683+00
ec6b6cd4-32b6-45ca-a0b6-76dd294dc2cd	c301c093-ebc6-491a-9c40-eb425e472857	00ac2f1022b3c1bdafe3637f3dfa51e4080306bb5677fd60dd4f537b3acf4aed	2026-06-16 00:06:13.167025+00	\N	2026-06-09 00:06:13.167026+00
58ff203a-d173-49b0-bf90-6e0ace16f826	21833e68-101f-4378-a916-62c120a9f192	607827b9511f5d8b97cf81d480d9bef3ce5d87b51bded37ab93a644d3280d407	2026-06-16 00:08:50.116514+00	\N	2026-06-09 00:08:50.116515+00
af0f7da7-bb11-4f17-b455-0641a8555b2d	21833e68-101f-4378-a916-62c120a9f192	518e604d0c60d114c0ba9854ea5d654e468c70d8a53b38caf7508f1030b8631c	2026-06-16 01:43:55.07615+00	\N	2026-06-09 01:43:55.076152+00
2d62e869-56bc-4d61-a6a2-b99f604fe183	21833e68-101f-4378-a916-62c120a9f192	1ceb49cd3336100e54dd6404cc3949069408b13a7ec0f0fae4f90ac722f59ddd	2026-06-16 01:48:47.798314+00	\N	2026-06-09 01:48:47.798316+00
f7e03d0b-19b5-48e5-bae7-e47357b26fe6	21833e68-101f-4378-a916-62c120a9f192	ca5f4dcba6139630432777a525064fd632977dfbf47e69e64fefecf87667e840	2026-06-16 02:08:06.303554+00	\N	2026-06-09 02:08:06.303555+00
1eb26c21-4bb1-4127-b96d-eba08cd98c3b	21833e68-101f-4378-a916-62c120a9f192	e67c2ba4c3de0d55410d18f381b050732e68a3b100bc3658d609a20e83d6b2f0	2026-06-16 02:19:15.048613+00	\N	2026-06-09 02:19:15.048614+00
\.


--
-- Data for Name: resource_assignments; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.resource_assignments (id, resource_id, user_id, quantity, status, issued_at, returned_at, notes) FROM stdin;
d226fa15-98d1-4d58-9f6f-a76f8ea0a47c	1cd798c4-9488-4365-be01-ca8e904ad6c9	21833e68-101f-4378-a916-62c120a9f192	1	ACTIVE	2026-05-31 23:41:24.163344	\N	
\.


--
-- Data for Name: resource_categories; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.resource_categories (id, tenant_id, name, description, created_at) FROM stdin;
c69be86c-afeb-4892-91bb-043cffcb1487	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	IT обладнання та зв'язок	Комп'ютерна техніка, планшети, термінали збору даних (ТЗД) та мережеве обладнання для забезпечення роботи складів та персоналу.	2026-05-02 15:40:19.25067+00
6a273219-ea10-4ea7-aff4-9647696668f4	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	Складський інвентар	Ручне та напівмеханізоване обладнання для переміщення вантажів (візки, рохлі), а також інструменти для комплектації.	2026-05-02 15:40:31.471438+00
733c1f11-0b6b-4f34-ab03-0c16946b9156	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	Пакувальні та витратні матеріали	Матеріали для пакування, фіксації та маркування вантажів (стретч-плівка, скотч, картонні короби, етикетки).	2026-05-02 15:40:43.465774+00
08a0ee32-49dc-4cec-aa93-1f95909488af	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	Навігаційні пристрої	Трекери, маячки та датчики для моніторингу транспорту та цілісності вантажів.	2026-05-02 15:40:58.680456+00
d424ffe3-b576-47f0-b4bd-d671fb37e5e3	0997831c-654f-471b-934f-cedafbc54ea5	Канцелярія	Ручки, папір, офісне приладдя	2026-05-07 10:33:21.041491+00
3b6bf3aa-0eb5-43f5-bad7-ca0a1fd595ee	0997831c-654f-471b-934f-cedafbc54ea5	Електроніка	Ноутбуки, телефони, зарядки	2026-05-07 10:33:21.041491+00
51a3c25e-42cf-4164-a5ab-91125b65a6e8	0997831c-654f-471b-934f-cedafbc54ea5	Інструмент	Ручний та електричний інструмент	2026-05-07 10:33:21.041491+00
6c39e1a3-edbe-48bd-ae22-a08251726603	0997831c-654f-471b-934f-cedafbc54ea5	Медикаменти	Аптечки, ліки, перев'язувальні матеріали	2026-05-07 10:33:21.041491+00
134c5fe0-51e6-43c4-ac46-670f00460450	0997831c-654f-471b-934f-cedafbc54ea5	Продукти	Консерви та продукти довготривалого зберігання	2026-05-07 10:33:21.041491+00
d427d9e2-60a7-4c7a-8aa9-8abbb69f5159	0997831c-654f-471b-934f-cedafbc54ea5	Одяг	Форма, спецодяг, ЗІЗ	2026-05-07 10:33:21.041491+00
4701ed53-f70a-4a5d-85d4-dfdb8b63bb76	0997831c-654f-471b-934f-cedafbc54ea5	Паливно-мастильні	Оливи, фільтри, технічні рідини	2026-05-07 10:33:21.041491+00
c6312649-caa3-4bb6-9a8d-66bf91c33b7d	0997831c-654f-471b-934f-cedafbc54ea5	Запчастини	Авто-запчастини та комплектуючі	2026-05-07 10:33:21.041491+00
3ffe4388-5816-40df-aa72-955d3c15d6e5	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	Автозапчастини		2026-06-09 17:57:37.432268+00
\.


--
-- Data for Name: resources; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.resources (id, tenant_id, category_id, unit_id, name, description, quantity, serial_number, location, condition, min_quantity, created_at, updated_at, unit_type, warehouse_id, weight_kg, barcode, unit_price) FROM stdin;
a19cdafb-a938-4bf7-83ac-58f9df641b4f	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	733c1f11-0b6b-4f34-ab03-0c16946b9156	3	Стретч-плівка пакувальна (рулон 500мм, 20мкм)		50	\N	\N	NEW	100	2026-05-03 14:31:25.470134+00	2026-05-31 22:53:35.193414+00	PCS	f64e8882-623d-41cd-8f72-2be30b783f8d	2.20		280.00
a30253b6-4e38-43f9-9e93-2e8ab763d291	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	733c1f11-0b6b-4f34-ab03-0c16946b9156	2	Стретч-плівка пакувальна (рулон 500мм, 20мкм)		438		\N	NEW	100	2026-05-02 15:47:37.912707+00	2026-05-31 22:57:26.077431+00	PCS	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	2.20		280.00
ab1dcaae-e910-41de-92f9-feb21f872456	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	733c1f11-0b6b-4f34-ab03-0c16946b9156	2	Стретч-плівка пакувальна (рулон 500мм, 20мкм)		9		\N	WRITTEN_OFF	100	2026-05-07 12:47:51.849245+00	2026-05-31 22:57:26.07948+00	PCS	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	2.20		280.00
1cd798c4-9488-4365-be01-ca8e904ad6c9	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	c69be86c-afeb-4892-91bb-043cffcb1487	2	Електросамокат		99		\N	NEW	120	2026-05-07 12:29:19.868246+00	2026-05-31 23:41:24.163344+00	PCS	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	40.00		15999.00
adac998c-affa-4dce-ad2f-2ae7fadcb4cd	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	3ffe4388-5816-40df-aa72-955d3c15d6e5	3	Гідропідсилювач	\N	10	\N	\N	NEW	0	2026-06-09 17:57:37.432268+00	2026-06-09 17:57:37.432268+00	PCS	f64e8882-623d-41cd-8f72-2be30b783f8d	1.00		2000.00
2a5ead75-9a60-4e0a-943c-bf8356e29539	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	c69be86c-afeb-4892-91bb-043cffcb1487	3	Термінал збору даних (ТЗД) Zebra TC21		25		\N	NEW	5	2026-05-02 15:42:09.740527+00	2026-05-02 15:42:09.740527+00	PCS	f64e8882-623d-41cd-8f72-2be30b783f8d	0.30		18500.00
a231f49e-b368-4cfb-8d1c-6b7cd596f0f0	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	6a273219-ea10-4ea7-aff4-9647696668f4	4	Гідравлічний візок (Рохля) 2.5т		12		\N	NEW	2	2026-05-02 15:46:34.510881+00	2026-05-02 15:46:34.510881+00	PCS	28f6a833-60fc-48bd-bb5c-30dfca2e3ace	65.00		12200.00
974c8ddc-7326-44d7-aac1-aabc54aba0ce	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	08a0ee32-49dc-4cec-aa93-1f95909488af	8	GPS-трекер Teltonika FMB120		45		\N	NEW	10	2026-05-02 15:48:34.709737+00	2026-05-02 15:48:34.709737+00	PCS	7a1639fb-d584-409a-89d9-e2b09fa379cd	0.20		1850.00
e3ab8171-e065-4eb0-8f06-280de9694a57	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	c69be86c-afeb-4892-91bb-043cffcb1487	2	Провід для заряджання		100		\N	NEW	120	2026-05-07 12:00:07.637046+00	2026-05-07 12:00:07.637046+00	PCS	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	0.50		350.00
9dd794cb-2a09-48f9-91e9-90822118a266	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	c69be86c-afeb-4892-91bb-043cffcb1487	2	Ноут		10		\N	NEW	20	2026-05-07 12:49:30.779274+00	2026-05-07 12:49:30.779274+00	PCS	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	2.20		280.00
00737572-6cf5-4e7c-81be-aa295129cc72	0997831c-654f-471b-934f-cedafbc54ea5	c6312649-caa3-4bb6-9a8d-66bf91c33b7d	75	Мультитул Leatherman Wave #184	Тестовий ресурс #1 для підрозділу north	45	SN-C82A20AF	Склад Чернігів, стелаж B	GOOD	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
7cb330a2-bbad-496f-ad1b-92b33adbd531	0997831c-654f-471b-934f-cedafbc54ea5	3b6bf3aa-0eb5-43f5-bad7-ca0a1fd595ee	75	Олива моторна 5W-40 (4 л) #337	Тестовий ресурс #2 для підрозділу north	98	SN-978C8DA3	Склад Чернігів, стелаж A	NEW	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
910cb880-efc2-4e96-8c09-61f277906eab	0997831c-654f-471b-934f-cedafbc54ea5	6c39e1a3-edbe-48bd-ae22-a08251726603	75	Набір викруток 40 шт. #885	Тестовий ресурс #3 для підрозділу north	97	SN-3E385CD8	Склад Чернігів, стелаж A	GOOD	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
38df8f83-898c-4bf2-8b4b-08721c9a3bc7	0997831c-654f-471b-934f-cedafbc54ea5	3b6bf3aa-0eb5-43f5-bad7-ca0a1fd595ee	75	Фільтр повітряний Mann-Filter #128	Тестовий ресурс #4 для підрозділу north	200	SN-4C677EA0	Склад Чернігів, стелаж A	NEW	5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
d902e3a4-c9ed-4dca-a8f9-ea0bd6b82b55	0997831c-654f-471b-934f-cedafbc54ea5	4701ed53-f70a-4a5d-85d4-dfdb8b63bb76	75	Насос ручний для шин #333	Тестовий ресурс #5 для підрозділу north	118	SN-2DA15772	Склад Чернігів, стелаж A	FAIR	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
27a32e73-d7f3-415f-b6bb-84d6bb5611bf	0997831c-654f-471b-934f-cedafbc54ea5	6c39e1a3-edbe-48bd-ae22-a08251726603	75	Фільтр повітряний Mann-Filter #213	Тестовий ресурс #6 для підрозділу north	54	SN-20082140	Склад Суми, зона 1	NEW	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
e53a2eb2-7e96-4807-9f17-51d0d8456cbe	0997831c-654f-471b-934f-cedafbc54ea5	51a3c25e-42cf-4164-a5ab-91125b65a6e8	75	Планшет Samsung Galaxy Tab #218	Тестовий ресурс #7 для підрозділу north	142	SN-8992CFE7	Склад Суми, зона 1	NEW	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
615d5262-b41b-488f-b355-c71185ebdcb0	0997831c-654f-471b-934f-cedafbc54ea5	134c5fe0-51e6-43c4-ac46-670f00460450	75	Спальний мішок -10°C #883	Тестовий ресурс #8 для підрозділу north	114	SN-36B66F86	Склад Чернігів, стелаж B	DAMAGED	5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
ab810d5d-ca53-4668-a44b-91f438e3db19	0997831c-654f-471b-934f-cedafbc54ea5	c6312649-caa3-4bb6-9a8d-66bf91c33b7d	75	Тент армований 6×8 м #557	Тестовий ресурс #9 для підрозділу north	103	SN-3E484630	Склад Суми, зона 1	DAMAGED	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
a6af48eb-c59d-4d89-a3d9-dadbeae15cdc	0997831c-654f-471b-934f-cedafbc54ea5	d424ffe3-b576-47f0-b4bd-d671fb37e5e3	75	Рюкзак тактичний 60 л #770	Тестовий ресурс #10 для підрозділу north	195	SN-96DAC282	Склад Чернігів, стелаж A	FAIR	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
54b243f1-bd9d-489e-afbf-9a44b3bf0a28	0997831c-654f-471b-934f-cedafbc54ea5	134c5fe0-51e6-43c4-ac46-670f00460450	75	Дрон DJI Mini 3 #332	Тестовий ресурс #11 для підрозділу north	42	SN-56B53257	Склад Чернігів, стелаж B	GOOD	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
acb3e99a-6901-4d7a-b28f-df4719cf5f8c	0997831c-654f-471b-934f-cedafbc54ea5	51a3c25e-42cf-4164-a5ab-91125b65a6e8	75	Дрон DJI Mini 3 #177	Тестовий ресурс #12 для підрозділу north	171	SN-62752D0B	Склад Чернігів, стелаж B	FAIR	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
90025136-169b-4a1f-a23c-a329e55e34de	0997831c-654f-471b-934f-cedafbc54ea5	51a3c25e-42cf-4164-a5ab-91125b65a6e8	75	Олива моторна 5W-40 (4 л) #753	Тестовий ресурс #13 для підрозділу north	1	SN-CF464F34	Склад Суми, зона 1	DAMAGED	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
7eeeae4c-f776-4822-9065-134db5a107c2	0997831c-654f-471b-934f-cedafbc54ea5	51a3c25e-42cf-4164-a5ab-91125b65a6e8	75	Мультитул Leatherman Wave #930	Тестовий ресурс #14 для підрозділу north	129	SN-60D7FCE0	Склад Суми, зона 1	GOOD	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
9a074b6d-3c73-4d72-bea5-c9d3bdb6a395	0997831c-654f-471b-934f-cedafbc54ea5	d424ffe3-b576-47f0-b4bd-d671fb37e5e3	75	Рукавиці робочі (пара) #414	Тестовий ресурс #15 для підрозділу north	130	SN-D2FAFA2C	Склад Чернігів, стелаж A	DAMAGED	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
5140c5f8-e2d0-47c6-9fd2-adf617cc6ced	0997831c-654f-471b-934f-cedafbc54ea5	6c39e1a3-edbe-48bd-ae22-a08251726603	76	Захисна каска EN397 #311	Тестовий ресурс #1 для підрозділу north_chernihiv	55	SN-701AA9D8	Склад Чернігів, стелаж B	DAMAGED	5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
135b3f8d-4925-4349-a444-15c63869c2bc	0997831c-654f-471b-934f-cedafbc54ea5	51a3c25e-42cf-4164-a5ab-91125b65a6e8	76	Олива моторна 5W-40 (4 л) #882	Тестовий ресурс #2 для підрозділу north_chernihiv	67	SN-58051072	Склад Суми, зона 1	NEW	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
1fb6c925-ffa4-4dbb-a221-0e76c3ee47b2	0997831c-654f-471b-934f-cedafbc54ea5	51a3c25e-42cf-4164-a5ab-91125b65a6e8	76	Рюкзак тактичний 60 л #300	Тестовий ресурс #3 для підрозділу north_chernihiv	45	SN-FF5C316A	Склад Чернігів, стелаж B	DAMAGED	5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
367208fe-741d-4cd0-8fa5-383d07311708	0997831c-654f-471b-934f-cedafbc54ea5	51a3c25e-42cf-4164-a5ab-91125b65a6e8	76	Ліхтарик тактичний Fenix #798	Тестовий ресурс #4 для підрозділу north_chernihiv	187	SN-1FA7EF80	Склад Чернігів, стелаж B	DAMAGED	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
ae460877-95d6-426b-8daf-cec2a1ea93c7	0997831c-654f-471b-934f-cedafbc54ea5	6c39e1a3-edbe-48bd-ae22-a08251726603	76	Генератор бензиновий 5 кВт #504	Тестовий ресурс #5 для підрозділу north_chernihiv	155	SN-B2DBE06A	Склад Чернігів, стелаж A	NEW	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
c8528510-110e-4391-9aff-1dd6c53055a1	0997831c-654f-471b-934f-cedafbc54ea5	6c39e1a3-edbe-48bd-ae22-a08251726603	76	Тент армований 6×8 м #256	Тестовий ресурс #6 для підрозділу north_chernihiv	172	SN-4F5BCF67	Склад Суми, зона 1	GOOD	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
74570fa5-5a2f-4fc9-8c79-9eabcbab1af9	0997831c-654f-471b-934f-cedafbc54ea5	6c39e1a3-edbe-48bd-ae22-a08251726603	76	Вогнегасник ВП-5 #770	Тестовий ресурс #7 для підрозділу north_chernihiv	3	SN-5B6D2D2F	Склад Суми, зона 1	DAMAGED	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
51d5f692-f4e7-4ae3-93bf-b7e3e563ca75	0997831c-654f-471b-934f-cedafbc54ea5	d427d9e2-60a7-4c7a-8aa9-8abbb69f5159	76	Тент армований 6×8 м #465	Тестовий ресурс #8 для підрозділу north_chernihiv	133	SN-E1BFD45C	Склад Чернігів, стелаж B	GOOD	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
d95b0275-5b2d-47c0-9663-a01c096533c9	0997831c-654f-471b-934f-cedafbc54ea5	134c5fe0-51e6-43c4-ac46-670f00460450	76	Комплект аптечок IFAK #252	Тестовий ресурс #9 для підрозділу north_chernihiv	188	SN-820AE7CA	Склад Чернігів, стелаж B	FAIR	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
5b85ccad-a271-44d0-a03f-ff4d92e2b02f	0997831c-654f-471b-934f-cedafbc54ea5	d424ffe3-b576-47f0-b4bd-d671fb37e5e3	80	Радіостанція Motorola DP2400 #239	Тестовий ресурс #2 для підрозділу south	53	SN-DD1D6639	Склад Одеса, ангар 2	NEW	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
a0607946-91db-4244-bcf5-2e28ca6c157b	0997831c-654f-471b-934f-cedafbc54ea5	4701ed53-f70a-4a5d-85d4-dfdb8b63bb76	76	Фільтр повітряний Mann-Filter #732	Тестовий ресурс #10 для підрозділу north_chernihiv	85	SN-D71E0203	Склад Суми, зона 1	NEW	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
cac4f93b-3aea-493d-88c6-c44b922ba64d	0997831c-654f-471b-934f-cedafbc54ea5	51a3c25e-42cf-4164-a5ab-91125b65a6e8	76	Рукавиці робочі (пара) #980	Тестовий ресурс #11 для підрозділу north_chernihiv	66	SN-5CE24B68	Склад Чернігів, стелаж A	GOOD	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
49aaee56-7beb-4736-b71e-469dd6bd75fc	0997831c-654f-471b-934f-cedafbc54ea5	c6312649-caa3-4bb6-9a8d-66bf91c33b7d	76	Трос буксирувальний 5 т #251	Тестовий ресурс #12 для підрозділу north_chernihiv	167	SN-7DBA9F09	Склад Суми, зона 1	DAMAGED	5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
8033ea57-c691-4e8e-9d35-00945cf92dc5	0997831c-654f-471b-934f-cedafbc54ea5	3b6bf3aa-0eb5-43f5-bad7-ca0a1fd595ee	76	Аптечка автомобільна #174	Тестовий ресурс #13 для підрозділу north_chernihiv	47	SN-6EACDA17	Склад Суми, зона 1	GOOD	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
6a77c89f-0247-4e68-9185-2b5ecdcf746c	0997831c-654f-471b-934f-cedafbc54ea5	d424ffe3-b576-47f0-b4bd-d671fb37e5e3	76	Планшет Samsung Galaxy Tab #974	Тестовий ресурс #14 для підрозділу north_chernihiv	150	SN-CC6D371D	Склад Суми, зона 1	DAMAGED	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
3a132367-b565-4007-b4a5-00160d6e7fec	0997831c-654f-471b-934f-cedafbc54ea5	134c5fe0-51e6-43c4-ac46-670f00460450	76	Фільтр повітряний Mann-Filter #444	Тестовий ресурс #15 для підрозділу north_chernihiv	58	SN-0EEFA8BA	Склад Чернігів, стелаж B	FAIR	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
40a30a31-53e5-4f0f-a2d7-43820b9150d1	0997831c-654f-471b-934f-cedafbc54ea5	d424ffe3-b576-47f0-b4bd-d671fb37e5e3	76	Захисна каска EN397 #755	Тестовий ресурс #16 для підрозділу north_chernihiv	168	SN-0462971E	Склад Суми, зона 1	FAIR	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
bc87888f-4172-4560-8118-3a2e7c0b9775	0997831c-654f-471b-934f-cedafbc54ea5	134c5fe0-51e6-43c4-ac46-670f00460450	76	Рюкзак тактичний 60 л #832	Тестовий ресурс #17 для підрозділу north_chernihiv	144	SN-69F3C0ED	Склад Чернігів, стелаж B	DAMAGED	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
8b2a128f-aaa6-4df1-9ba0-213104ce0761	0997831c-654f-471b-934f-cedafbc54ea5	c6312649-caa3-4bb6-9a8d-66bf91c33b7d	76	Захисна каска EN397 #383	Тестовий ресурс #18 для підрозділу north_chernihiv	161	SN-1C0C38E5	Склад Чернігів, стелаж A	GOOD	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
6c176c9b-c3b3-49b8-86c8-5918d0672ef1	0997831c-654f-471b-934f-cedafbc54ea5	c6312649-caa3-4bb6-9a8d-66bf91c33b7d	76	Лопата складна армійська #644	Тестовий ресурс #19 для підрозділу north_chernihiv	142	SN-227AF16D	Склад Суми, зона 1	DAMAGED	5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
63f1f2f4-3e0c-42f8-b909-c5411f979886	0997831c-654f-471b-934f-cedafbc54ea5	4701ed53-f70a-4a5d-85d4-dfdb8b63bb76	76	Аптечка автомобільна #201	Тестовий ресурс #20 для підрозділу north_chernihiv	103	SN-D8E625CE	Склад Чернігів, стелаж A	GOOD	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
cd4f634a-2e3d-4fa8-b68f-459f70bae9f8	0997831c-654f-471b-934f-cedafbc54ea5	134c5fe0-51e6-43c4-ac46-670f00460450	79	Комплект аптечок IFAK #997	Тестовий ресурс #1 для підрозділу north_sumy	23	SN-89D4E5DB	Склад Чернігів, стелаж B	FAIR	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
a3599ebf-0fc4-4687-95c5-ce82ec067a0f	0997831c-654f-471b-934f-cedafbc54ea5	d427d9e2-60a7-4c7a-8aa9-8abbb69f5159	79	Фільтр для води LifeStraw #656	Тестовий ресурс #2 для підрозділу north_sumy	24	SN-F8F7628D	Склад Чернігів, стелаж B	NEW	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
029eb38a-e6b4-48d0-bcb8-b106497f841b	0997831c-654f-471b-934f-cedafbc54ea5	d424ffe3-b576-47f0-b4bd-d671fb37e5e3	79	Насос ручний для шин #475	Тестовий ресурс #3 для підрозділу north_sumy	34	SN-8F1125E1	Склад Чернігів, стелаж A	GOOD	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
a56f6cbe-69b2-4aa0-bd4e-8dbbf420867d	0997831c-654f-471b-934f-cedafbc54ea5	134c5fe0-51e6-43c4-ac46-670f00460450	79	Дрон DJI Mini 3 #363	Тестовий ресурс #4 для підрозділу north_sumy	32	SN-98A60EE7	Склад Суми, зона 1	NEW	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
d4f66d92-8902-428d-a763-bbe763e57338	0997831c-654f-471b-934f-cedafbc54ea5	d427d9e2-60a7-4c7a-8aa9-8abbb69f5159	79	Фільтр повітряний Mann-Filter #872	Тестовий ресурс #5 для підрозділу north_sumy	77	SN-8353A15E	Склад Чернігів, стелаж A	NEW	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
35e025fd-e7dc-4154-9ac2-a387aec2ba8e	0997831c-654f-471b-934f-cedafbc54ea5	134c5fe0-51e6-43c4-ac46-670f00460450	79	Вогнегасник ВП-5 #106	Тестовий ресурс #6 для підрозділу north_sumy	78	SN-07B0A243	Склад Чернігів, стелаж A	NEW	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
6bc0bce7-deeb-4900-a816-2f726271aac0	0997831c-654f-471b-934f-cedafbc54ea5	d427d9e2-60a7-4c7a-8aa9-8abbb69f5159	79	Трос буксирувальний 5 т #481	Тестовий ресурс #7 для підрозділу north_sumy	154	SN-A0427745	Склад Чернігів, стелаж A	DAMAGED	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
a8800ab1-d81a-413e-a681-8977babb85d9	0997831c-654f-471b-934f-cedafbc54ea5	6c39e1a3-edbe-48bd-ae22-a08251726603	79	Генератор бензиновий 5 кВт #256	Тестовий ресурс #8 для підрозділу north_sumy	68	SN-77D8FCEF	Склад Чернігів, стелаж B	DAMAGED	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
0dd9a7a0-1247-4f65-856d-15ccbd23c62d	0997831c-654f-471b-934f-cedafbc54ea5	51a3c25e-42cf-4164-a5ab-91125b65a6e8	79	Фільтр для води LifeStraw #581	Тестовий ресурс #9 для підрозділу north_sumy	52	SN-1DBEDEEB	Склад Чернігів, стелаж A	FAIR	5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
6e1a6814-c65a-42f1-8304-8483ae2efb14	0997831c-654f-471b-934f-cedafbc54ea5	134c5fe0-51e6-43c4-ac46-670f00460450	79	Рюкзак тактичний 60 л #844	Тестовий ресурс #10 для підрозділу north_sumy	38	SN-C8D08534	Склад Чернігів, стелаж B	DAMAGED	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
44a7e93b-1224-419f-b72b-f5874da8f20f	0997831c-654f-471b-934f-cedafbc54ea5	d424ffe3-b576-47f0-b4bd-d671fb37e5e3	79	Ліхтарик тактичний Fenix #344	Тестовий ресурс #11 для підрозділу north_sumy	108	SN-CF81283C	Склад Суми, зона 1	GOOD	5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
22bb6459-8545-4967-a71b-01666d7bb21a	0997831c-654f-471b-934f-cedafbc54ea5	6c39e1a3-edbe-48bd-ae22-a08251726603	79	Аптечка автомобільна #821	Тестовий ресурс #12 для підрозділу north_sumy	54	SN-9BCB79E5	Склад Чернігів, стелаж A	FAIR	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
482cf5a5-67a2-47ac-a9bd-d7beb68d1a4a	0997831c-654f-471b-934f-cedafbc54ea5	c6312649-caa3-4bb6-9a8d-66bf91c33b7d	80	Термобілизна комплект #992	Тестовий ресурс #1 для підрозділу south	75	SN-E8887C26	Склад Одеса, ангар 1	NEW	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
c8d9f854-781a-482b-9e68-23ad365aae0e	0997831c-654f-471b-934f-cedafbc54ea5	d424ffe3-b576-47f0-b4bd-d671fb37e5e3	80	Радіостанція Motorola DP2400 #941	Тестовий ресурс #3 для підрозділу south	33	SN-60516144	Склад Миколаїв, зона 3	NEW	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
9243a35c-b2f0-4c4f-853f-3dab9b8cae10	0997831c-654f-471b-934f-cedafbc54ea5	6c39e1a3-edbe-48bd-ae22-a08251726603	80	Ноутбук Lenovo ThinkPad E15 #242	Тестовий ресурс #4 для підрозділу south	40	SN-A5A96276	Склад Миколаїв, зона 3	FAIR	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
5b90e304-2073-45f8-8f64-34f7319e1264	0997831c-654f-471b-934f-cedafbc54ea5	d424ffe3-b576-47f0-b4bd-d671fb37e5e3	80	Генератор бензиновий 5 кВт #198	Тестовий ресурс #5 для підрозділу south	53	SN-6ECF951C	Склад Одеса, ангар 1	DAMAGED	5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
a28a5172-e9d3-4134-9168-70c71678d5a1	0997831c-654f-471b-934f-cedafbc54ea5	d427d9e2-60a7-4c7a-8aa9-8abbb69f5159	80	Комплект аптечок IFAK #894	Тестовий ресурс #6 для підрозділу south	53	SN-F71CD004	Склад Одеса, ангар 2	NEW	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
5c00c0d3-741a-4c55-9b4f-dc44b573369d	0997831c-654f-471b-934f-cedafbc54ea5	4701ed53-f70a-4a5d-85d4-dfdb8b63bb76	80	Олива моторна 5W-40 (4 л) #835	Тестовий ресурс #7 для підрозділу south	128	SN-ACB9FB8D	Склад Одеса, ангар 2	NEW	5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
fdf0395d-b716-4b64-a6ce-29aef92616dd	0997831c-654f-471b-934f-cedafbc54ea5	6c39e1a3-edbe-48bd-ae22-a08251726603	80	Мультитул Leatherman Wave #298	Тестовий ресурс #8 для підрозділу south	78	SN-AEA14552	Склад Одеса, ангар 2	GOOD	5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
297bf17f-035e-4f9f-bc1b-4de79317ebfa	0997831c-654f-471b-934f-cedafbc54ea5	c6312649-caa3-4bb6-9a8d-66bf91c33b7d	80	Аптечка автомобільна #260	Тестовий ресурс #9 для підрозділу south	95	SN-8A4CF8FB	Склад Миколаїв, зона 3	NEW	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
50f16782-6ed6-41f4-90b9-2c623623a3d9	0997831c-654f-471b-934f-cedafbc54ea5	51a3c25e-42cf-4164-a5ab-91125b65a6e8	80	Кабель-подовжувач 50 м #783	Тестовий ресурс #10 для підрозділу south	31	SN-EEF43DB2	Склад Одеса, ангар 2	NEW	5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
9803ce7c-ed61-4437-9074-928f3d4d9427	0997831c-654f-471b-934f-cedafbc54ea5	51a3c25e-42cf-4164-a5ab-91125b65a6e8	81	Радіостанція Motorola DP2400 #769	Тестовий ресурс #1 для підрозділу south_odesa	145	SN-5FA73C34	Склад Миколаїв, зона 3	GOOD	5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
5c7d8f7f-3bb5-490b-81fb-731622a4631d	0997831c-654f-471b-934f-cedafbc54ea5	3b6bf3aa-0eb5-43f5-bad7-ca0a1fd595ee	81	Рюкзак тактичний 60 л #577	Тестовий ресурс #2 для підрозділу south_odesa	102	SN-CCA2BCCE	Склад Одеса, ангар 1	GOOD	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
3617500f-9e88-4156-b126-29fa6d7af99b	0997831c-654f-471b-934f-cedafbc54ea5	51a3c25e-42cf-4164-a5ab-91125b65a6e8	81	Набір викруток 40 шт. #458	Тестовий ресурс #3 для підрозділу south_odesa	30	SN-0459E7E3	Склад Миколаїв, зона 3	NEW	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
fa2fe978-7ddf-428a-af70-102971fd262c	0997831c-654f-471b-934f-cedafbc54ea5	d424ffe3-b576-47f0-b4bd-d671fb37e5e3	81	Ноутбук Lenovo ThinkPad E15 #759	Тестовий ресурс #4 для підрозділу south_odesa	67	SN-B6FB9EF1	Склад Одеса, ангар 2	DAMAGED	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
9dcc826c-d6ae-4c8d-b6b4-76a5d53cd84d	0997831c-654f-471b-934f-cedafbc54ea5	134c5fe0-51e6-43c4-ac46-670f00460450	81	Вогнегасник ВП-5 #849	Тестовий ресурс #5 для підрозділу south_odesa	82	SN-A86B49F4	Склад Одеса, ангар 2	NEW	5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
5a5505d3-f063-4ce3-862b-3cad7a1743b9	0997831c-654f-471b-934f-cedafbc54ea5	6c39e1a3-edbe-48bd-ae22-a08251726603	81	Тент армований 6×8 м #172	Тестовий ресурс #6 для підрозділу south_odesa	154	SN-9292E3F7	Склад Миколаїв, зона 3	GOOD	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
e25452e0-151c-44e9-b25d-bde5a9d520a9	0997831c-654f-471b-934f-cedafbc54ea5	d427d9e2-60a7-4c7a-8aa9-8abbb69f5159	81	Генератор бензиновий 5 кВт #644	Тестовий ресурс #7 для підрозділу south_odesa	170	SN-3CC6DE20	Склад Одеса, ангар 2	FAIR	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
b4f92a5e-8049-4fe8-beac-4f9e427b9cae	0997831c-654f-471b-934f-cedafbc54ea5	3b6bf3aa-0eb5-43f5-bad7-ca0a1fd595ee	81	Захисна каска EN397 #733	Тестовий ресурс #8 для підрозділу south_odesa	32	SN-A90DE261	Склад Одеса, ангар 1	DAMAGED	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
90e79d59-e64a-4ae0-ac3a-402dde63863d	0997831c-654f-471b-934f-cedafbc54ea5	134c5fe0-51e6-43c4-ac46-670f00460450	81	Олива моторна 5W-40 (4 л) #152	Тестовий ресурс #9 для підрозділу south_odesa	165	SN-03742301	Склад Миколаїв, зона 3	NEW	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
8957a1e5-4977-46a3-b12c-64fb1dcd331f	0997831c-654f-471b-934f-cedafbc54ea5	51a3c25e-42cf-4164-a5ab-91125b65a6e8	81	Захисна каска EN397 #426	Тестовий ресурс #10 для підрозділу south_odesa	192	SN-D77A5EFC	Склад Одеса, ангар 2	NEW	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
f79a256c-f608-46ac-bb65-7fe57ba62daf	0997831c-654f-471b-934f-cedafbc54ea5	134c5fe0-51e6-43c4-ac46-670f00460450	81	Фільтр для води LifeStraw #223	Тестовий ресурс #11 для підрозділу south_odesa	196	SN-9BE081E1	Склад Одеса, ангар 1	NEW	5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
cc6933c2-f196-40d3-a294-bd9700c67de3	0997831c-654f-471b-934f-cedafbc54ea5	d424ffe3-b576-47f0-b4bd-d671fb37e5e3	81	Олива моторна 5W-40 (4 л) #753	Тестовий ресурс #12 для підрозділу south_odesa	112	SN-EFAA0399	Склад Одеса, ангар 1	DAMAGED	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
e2fa742b-181e-445a-912a-9325d79548f7	0997831c-654f-471b-934f-cedafbc54ea5	51a3c25e-42cf-4164-a5ab-91125b65a6e8	81	Аптечка автомобільна #419	Тестовий ресурс #13 для підрозділу south_odesa	30	SN-8A85D016	Склад Одеса, ангар 2	GOOD	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
03e662f2-9401-48ac-9644-6cbb190dc751	0997831c-654f-471b-934f-cedafbc54ea5	c6312649-caa3-4bb6-9a8d-66bf91c33b7d	81	Олива моторна 5W-40 (4 л) #541	Тестовий ресурс #14 для підрозділу south_odesa	129	SN-3B7EDC1B	Склад Одеса, ангар 2	GOOD	5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
c6752de9-49ca-4132-af29-661c853c7a4a	0997831c-654f-471b-934f-cedafbc54ea5	134c5fe0-51e6-43c4-ac46-670f00460450	81	Кабель-подовжувач 50 м #480	Тестовий ресурс #15 для підрозділу south_odesa	15	SN-087ACA8A	Склад Одеса, ангар 2	GOOD	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
fddc8701-46c1-4f5f-9f42-829e235ec1d3	0997831c-654f-471b-934f-cedafbc54ea5	c6312649-caa3-4bb6-9a8d-66bf91c33b7d	81	Рукавиці робочі (пара) #549	Тестовий ресурс #16 для підрозділу south_odesa	74	SN-0AD5A076	Склад Одеса, ангар 2	NEW	5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
7bf4fafb-8ab7-4323-ac96-ec8190fe3909	0997831c-654f-471b-934f-cedafbc54ea5	6c39e1a3-edbe-48bd-ae22-a08251726603	81	Рукавиці робочі (пара) #786	Тестовий ресурс #17 для підрозділу south_odesa	174	SN-87F5E526	Склад Одеса, ангар 1	GOOD	5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
4eb2ff5b-6f2b-4b83-86f8-19f0be1d0bb9	0997831c-654f-471b-934f-cedafbc54ea5	d424ffe3-b576-47f0-b4bd-d671fb37e5e3	81	Ліхтарик тактичний Fenix #363	Тестовий ресурс #18 для підрозділу south_odesa	165	SN-C6E4544C	Склад Одеса, ангар 2	GOOD	5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
9d85f1a5-ef55-455f-9436-909f00a6f330	0997831c-654f-471b-934f-cedafbc54ea5	134c5fe0-51e6-43c4-ac46-670f00460450	83	Генератор бензиновий 5 кВт #132	Тестовий ресурс #1 для підрозділу south_mykolaiv	26	SN-DF4AC90A	Склад Миколаїв, зона 3	GOOD	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
0edddd6b-37cc-4457-abf3-92e17877eaa4	0997831c-654f-471b-934f-cedafbc54ea5	c6312649-caa3-4bb6-9a8d-66bf91c33b7d	83	Термобілизна комплект #462	Тестовий ресурс #2 для підрозділу south_mykolaiv	29	SN-465C06E3	Склад Миколаїв, зона 3	FAIR	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
eff360a3-a9a2-4002-bd25-d1ff836b3528	0997831c-654f-471b-934f-cedafbc54ea5	d427d9e2-60a7-4c7a-8aa9-8abbb69f5159	83	Кабель-подовжувач 50 м #594	Тестовий ресурс #3 для підрозділу south_mykolaiv	109	SN-D70FEEB8	Склад Одеса, ангар 2	NEW	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
fb969c14-9953-42cc-b2ca-f00b13891477	0997831c-654f-471b-934f-cedafbc54ea5	3b6bf3aa-0eb5-43f5-bad7-ca0a1fd595ee	83	Вогнегасник ВП-5 #441	Тестовий ресурс #4 для підрозділу south_mykolaiv	173	SN-B8C93FC0	Склад Миколаїв, зона 3	GOOD	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
99219a48-3ee8-42fb-b313-b98dfafe4a89	0997831c-654f-471b-934f-cedafbc54ea5	51a3c25e-42cf-4164-a5ab-91125b65a6e8	83	Олива моторна 5W-40 (4 л) #468	Тестовий ресурс #5 для підрозділу south_mykolaiv	93	SN-68D8798D	Склад Одеса, ангар 2	NEW	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
b4ffa935-d012-48f8-983d-40cd15d5d0ea	0997831c-654f-471b-934f-cedafbc54ea5	4701ed53-f70a-4a5d-85d4-dfdb8b63bb76	83	Вогнегасник ВП-5 #108	Тестовий ресурс #6 для підрозділу south_mykolaiv	25	SN-1EE2DC63	Склад Одеса, ангар 1	FAIR	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
9dc3bbaa-2aa0-49df-9eed-d0bdaffc454a	0997831c-654f-471b-934f-cedafbc54ea5	4701ed53-f70a-4a5d-85d4-dfdb8b63bb76	83	Вогнегасник ВП-5 #148	Тестовий ресурс #7 для підрозділу south_mykolaiv	181	SN-A8336AD9	Склад Одеса, ангар 2	FAIR	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
8d705fd8-abd2-41eb-a9bf-d355a12a0210	0997831c-654f-471b-934f-cedafbc54ea5	c6312649-caa3-4bb6-9a8d-66bf91c33b7d	83	Аптечка автомобільна #546	Тестовий ресурс #8 для підрозділу south_mykolaiv	127	SN-C5EAD7CE	Склад Миколаїв, зона 3	NEW	5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
f08654c5-5f70-40e2-94b1-f6a5c5fc8188	0997831c-654f-471b-934f-cedafbc54ea5	d424ffe3-b576-47f0-b4bd-d671fb37e5e3	83	Ліхтарик тактичний Fenix #342	Тестовий ресурс #9 для підрозділу south_mykolaiv	14	SN-A4D12E44	Склад Миколаїв, зона 3	DAMAGED	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
e3bebdbd-6109-460d-9cc0-acc310b22a2e	0997831c-654f-471b-934f-cedafbc54ea5	3b6bf3aa-0eb5-43f5-bad7-ca0a1fd595ee	83	Олива моторна 5W-40 (4 л) #140	Тестовий ресурс #10 для підрозділу south_mykolaiv	59	SN-D7B66D3B	Склад Одеса, ангар 2	FAIR	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
db5a599d-c1b1-4a37-913f-560e6193c88f	0997831c-654f-471b-934f-cedafbc54ea5	c6312649-caa3-4bb6-9a8d-66bf91c33b7d	84	Термобілизна комплект #580	Тестовий ресурс #1 для підрозділу central	12	SN-6439A40D	Склад Київ, антресоль	GOOD	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
1e223058-290f-40ae-8086-a4b2234df1fe	0997831c-654f-471b-934f-cedafbc54ea5	4701ed53-f70a-4a5d-85d4-dfdb8b63bb76	84	Вогнегасник ВП-5 #238	Тестовий ресурс #2 для підрозділу central	64	SN-952D5D15	Склад Київ, поверх 2	GOOD	5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
ca27a10d-5e0c-4f4d-955e-bebc39a938e7	0997831c-654f-471b-934f-cedafbc54ea5	c6312649-caa3-4bb6-9a8d-66bf91c33b7d	84	Захисна каска EN397 #521	Тестовий ресурс #3 для підрозділу central	24	SN-D62F418C	Склад Київ, антресоль	DAMAGED	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
79b96d03-54ca-49a5-9f7b-bc5c28b77920	0997831c-654f-471b-934f-cedafbc54ea5	c6312649-caa3-4bb6-9a8d-66bf91c33b7d	84	Захисна каска EN397 #482	Тестовий ресурс #4 для підрозділу central	108	SN-B9D8FE09	Склад Київ, антресоль	FAIR	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
cc9d12bf-112b-4978-9653-67cfb241cbf9	0997831c-654f-471b-934f-cedafbc54ea5	6c39e1a3-edbe-48bd-ae22-a08251726603	84	Рукавиці робочі (пара) #877	Тестовий ресурс #5 для підрозділу central	98	SN-35BBDB51	Склад Київ, поверх 1	GOOD	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
d1a219f1-6ed2-4e5a-8676-c32e633aafc3	0997831c-654f-471b-934f-cedafbc54ea5	d424ffe3-b576-47f0-b4bd-d671fb37e5e3	84	Рюкзак тактичний 60 л #583	Тестовий ресурс #6 для підрозділу central	104	SN-886ADE1B	Склад Київ, антресоль	NEW	5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
a897369d-1b81-4daa-8bf9-bb68d9a63d6a	0997831c-654f-471b-934f-cedafbc54ea5	51a3c25e-42cf-4164-a5ab-91125b65a6e8	84	Лопата складна армійська #562	Тестовий ресурс #7 для підрозділу central	160	SN-9EF90157	Склад Київ, поверх 1	GOOD	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
7d73459d-9c7f-4014-bbef-16a2afb43a04	0997831c-654f-471b-934f-cedafbc54ea5	d424ffe3-b576-47f0-b4bd-d671fb37e5e3	84	Кабель-подовжувач 50 м #945	Тестовий ресурс #8 для підрозділу central	140	SN-8E7151B5	Склад Київ, поверх 2	NEW	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
8199727c-6c56-43c8-a1b6-35802f5a9599	0997831c-654f-471b-934f-cedafbc54ea5	51a3c25e-42cf-4164-a5ab-91125b65a6e8	85	Термобілизна комплект #147	Тестовий ресурс #1 для підрозділу central_kyiv	150	SN-A4D36A8F	Склад Київ, антресоль	DAMAGED	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
37cea2c2-b6c5-4f05-91c1-0dd673491045	0997831c-654f-471b-934f-cedafbc54ea5	134c5fe0-51e6-43c4-ac46-670f00460450	85	Захисна каска EN397 #650	Тестовий ресурс #2 для підрозділу central_kyiv	86	SN-E6D34367	Склад Київ, антресоль	GOOD	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
3913ffa9-12c7-460c-bf58-b6813bd67066	0997831c-654f-471b-934f-cedafbc54ea5	3b6bf3aa-0eb5-43f5-bad7-ca0a1fd595ee	85	Ліхтарик тактичний Fenix #636	Тестовий ресурс #3 для підрозділу central_kyiv	144	SN-7B99F75A	Склад Київ, антресоль	GOOD	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
e504a3d2-45d0-4572-ad41-b0aed16085ec	0997831c-654f-471b-934f-cedafbc54ea5	d427d9e2-60a7-4c7a-8aa9-8abbb69f5159	85	Вогнегасник ВП-5 #745	Тестовий ресурс #4 для підрозділу central_kyiv	152	SN-4E7A9160	Склад Київ, поверх 2	NEW	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
28febdc6-e9c0-411a-95fb-550d98d76ec0	0997831c-654f-471b-934f-cedafbc54ea5	134c5fe0-51e6-43c4-ac46-670f00460450	85	Кабель-подовжувач 50 м #545	Тестовий ресурс #5 для підрозділу central_kyiv	62	SN-675839E5	Склад Київ, антресоль	DAMAGED	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
bf1235cc-7e52-433d-86fd-7cbde828bf03	0997831c-654f-471b-934f-cedafbc54ea5	3b6bf3aa-0eb5-43f5-bad7-ca0a1fd595ee	85	Термобілизна комплект #574	Тестовий ресурс #6 для підрозділу central_kyiv	91	SN-8E319D2A	Склад Київ, поверх 2	GOOD	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
4f86352a-5b43-4072-addb-be32b5ca4673	0997831c-654f-471b-934f-cedafbc54ea5	51a3c25e-42cf-4164-a5ab-91125b65a6e8	85	Кабель-подовжувач 50 м #882	Тестовий ресурс #7 для підрозділу central_kyiv	92	SN-D289FE6C	Склад Київ, антресоль	NEW	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
dcfa5bd6-1265-4897-919a-54dc9f35ebae	0997831c-654f-471b-934f-cedafbc54ea5	6c39e1a3-edbe-48bd-ae22-a08251726603	85	Рюкзак тактичний 60 л #721	Тестовий ресурс #8 для підрозділу central_kyiv	144	SN-478A5D5C	Склад Київ, поверх 1	DAMAGED	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
1d11a833-d398-4f08-afc1-48a4a805bbce	0997831c-654f-471b-934f-cedafbc54ea5	d424ffe3-b576-47f0-b4bd-d671fb37e5e3	85	Тент армований 6×8 м #223	Тестовий ресурс #9 для підрозділу central_kyiv	30	SN-06AF5805	Склад Київ, поверх 1	NEW	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
70b5e56f-18a5-4774-a3bc-c92336ae03df	0997831c-654f-471b-934f-cedafbc54ea5	51a3c25e-42cf-4164-a5ab-91125b65a6e8	85	Ноутбук Lenovo ThinkPad E15 #668	Тестовий ресурс #10 для підрозділу central_kyiv	28	SN-216E56CF	Склад Київ, поверх 2	GOOD	5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
50ac682a-980f-405a-81e8-719c2c428df8	0997831c-654f-471b-934f-cedafbc54ea5	6c39e1a3-edbe-48bd-ae22-a08251726603	85	Кабель-подовжувач 50 м #541	Тестовий ресурс #11 для підрозділу central_kyiv	13	SN-F928B519	Склад Київ, антресоль	NEW	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
ce441595-1871-4afa-9567-21baaec2fde3	0997831c-654f-471b-934f-cedafbc54ea5	d427d9e2-60a7-4c7a-8aa9-8abbb69f5159	85	Фільтр для води LifeStraw #925	Тестовий ресурс #12 для підрозділу central_kyiv	87	SN-63CE3065	Склад Київ, поверх 2	NEW	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
b11c5667-e567-4eb1-9bba-f63dab77b42c	0997831c-654f-471b-934f-cedafbc54ea5	d424ffe3-b576-47f0-b4bd-d671fb37e5e3	85	Дрон DJI Mini 3 #922	Тестовий ресурс #13 для підрозділу central_kyiv	55	SN-815F49F5	Склад Київ, антресоль	GOOD	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
0d547486-28b3-4f49-8338-72902ce9ccab	0997831c-654f-471b-934f-cedafbc54ea5	d427d9e2-60a7-4c7a-8aa9-8abbb69f5159	85	Набір викруток 40 шт. #936	Тестовий ресурс #14 для підрозділу central_kyiv	34	SN-B9ED5BA7	Склад Київ, антресоль	GOOD	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
f1c6bd0d-b35c-41bd-bb5a-0ed4f039ef5b	0997831c-654f-471b-934f-cedafbc54ea5	134c5fe0-51e6-43c4-ac46-670f00460450	85	Спальний мішок -10°C #712	Тестовий ресурс #15 для підрозділу central_kyiv	96	SN-8EAAD714	Склад Київ, поверх 1	GOOD	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
96fa9d33-5741-4690-8d21-dbfd6cdbc986	0997831c-654f-471b-934f-cedafbc54ea5	4701ed53-f70a-4a5d-85d4-dfdb8b63bb76	85	Фільтр повітряний Mann-Filter #974	Тестовий ресурс #16 для підрозділу central_kyiv	70	SN-2744BDDA	Склад Київ, антресоль	NEW	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
117bb8dd-2568-410f-855f-22d788f82954	0997831c-654f-471b-934f-cedafbc54ea5	51a3c25e-42cf-4164-a5ab-91125b65a6e8	85	Планшет Samsung Galaxy Tab #203	Тестовий ресурс #17 для підрозділу central_kyiv	139	SN-AF6C750A	Склад Київ, антресоль	GOOD	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
e89beac9-7b87-4267-8185-06bc1dc6628a	0997831c-654f-471b-934f-cedafbc54ea5	d427d9e2-60a7-4c7a-8aa9-8abbb69f5159	85	Насос ручний для шин #892	Тестовий ресурс #18 для підрозділу central_kyiv	34	SN-9C232EF8	Склад Київ, антресоль	NEW	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
871120e1-973b-48b7-958e-111d2876f2a0	0997831c-654f-471b-934f-cedafbc54ea5	4701ed53-f70a-4a5d-85d4-dfdb8b63bb76	85	Трос буксирувальний 5 т #950	Тестовий ресурс #19 для підрозділу central_kyiv	5	SN-6C91F575	Склад Київ, поверх 1	FAIR	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
7430bdb4-2ea4-42d1-9de3-ecb88feced02	0997831c-654f-471b-934f-cedafbc54ea5	3b6bf3aa-0eb5-43f5-bad7-ca0a1fd595ee	85	Фільтр для води LifeStraw #665	Тестовий ресурс #20 для підрозділу central_kyiv	37	SN-BEBC9DF8	Склад Київ, поверх 1	GOOD	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
e4d3e550-7945-4a29-9dc3-445f52a081cb	0997831c-654f-471b-934f-cedafbc54ea5	134c5fe0-51e6-43c4-ac46-670f00460450	85	Радіостанція Motorola DP2400 #886	Тестовий ресурс #21 для підрозділу central_kyiv	97	SN-13F71CAB	Склад Київ, поверх 1	NEW	10	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
5722d1b5-f3af-4274-88c7-807b3965d4f3	0997831c-654f-471b-934f-cedafbc54ea5	51a3c25e-42cf-4164-a5ab-91125b65a6e8	85	Мультитул Leatherman Wave #348	Тестовий ресурс #22 для підрозділу central_kyiv	90	SN-A56909D4	Склад Київ, антресоль	NEW	20	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
c8174ea4-f78b-452a-84e9-d5b1df292f85	0997831c-654f-471b-934f-cedafbc54ea5	4701ed53-f70a-4a5d-85d4-dfdb8b63bb76	85	Фільтр для води LifeStraw #868	Тестовий ресурс #23 для підрозділу central_kyiv	125	SN-224B2A62	Склад Київ, поверх 1	NEW	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
52c5eabb-0145-4752-96a8-c484b0e03bcb	0997831c-654f-471b-934f-cedafbc54ea5	3b6bf3aa-0eb5-43f5-bad7-ca0a1fd595ee	85	Комплект аптечок IFAK #577	Тестовий ресурс #24 для підрозділу central_kyiv	196	SN-0115CA92	Склад Київ, поверх 2	NEW	5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
796f7a6a-6b11-4ca9-bec2-7b1d34c28a31	0997831c-654f-471b-934f-cedafbc54ea5	134c5fe0-51e6-43c4-ac46-670f00460450	85	Рюкзак тактичний 60 л #845	Тестовий ресурс #25 для підрозділу central_kyiv	171	SN-D2D5C7D3	Склад Київ, поверх 2	NEW	2	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	PCS	\N	1.00		0.00
\.


--
-- Data for Name: shipment_items; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.shipment_items (id, shipment_id, resource_id, quantity, request_id) FROM stdin;
88e7938e-3635-469c-ba85-0a88a864067d	34af91e3-95bb-4633-a8ef-3ce9a7336400	a19cdafb-a938-4bf7-83ac-58f9df641b4f	1	617e59a8-38a6-4cf7-8c5d-12c13a09ae2d
a6a82d8b-264b-4f81-8094-40193d44dd92	798d6084-4572-441c-9f52-3f05892de225	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	\N
abc08c7b-c415-45c3-b33b-73ffdc4c6bae	55a417c4-d5a8-4656-8a14-741f2c2c0273	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	\N
58bfae55-8612-4b1b-ba83-4d622ce2c9f7	8839bbd3-ea60-4e00-b25f-823449d5e004	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	\N
9db1850f-806a-436f-b5a1-1ab2b93815d6	5d03ed08-289e-4825-af7c-d9134ecd4107	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	\N
6f55b448-c48d-4d41-be71-bda58a119733	a5504524-8653-44aa-8c38-a86832a5963e	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	\N
732e5839-79f0-4906-9a99-a8441412883a	2f71fc41-c742-48b4-8d2f-9c7b6263c3e0	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	\N
84c2099d-6b80-4dce-89ca-0238e72e094a	7e631304-8a13-4a8b-9605-e83fbd90e834	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	\N
49adb59d-54d2-41f7-a11e-30a2aad4f876	20bb818e-2a8c-4058-b865-5f52388a30ea	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	\N
1a65dfd5-da01-4ea2-be5a-7fdb2482e688	9cd279eb-1626-4932-938c-63cd7f4829e2	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	\N
c7c45268-27da-4150-9c3e-bcccd383fb98	b05dd831-264e-4bed-a7a9-e3e3d06c1393	a19cdafb-a938-4bf7-83ac-58f9df641b4f	1	\N
9e5b125b-3b1c-4e77-b2f1-f2db020fc38a	8e1ca805-ed59-4b6f-bded-84ad66ae0138	a19cdafb-a938-4bf7-83ac-58f9df641b4f	1	\N
4bdafb66-1eb9-4e62-b1be-fd00f9229522	704ce237-6254-415d-aadb-7089f825e48b	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	44878184-a03c-4e93-9b39-bd9ba1b650e6
adfa935c-3cf2-40ad-a322-ad20940051b3	8ff169fe-d4d5-41d5-ab14-355008cc9a0e	a19cdafb-a938-4bf7-83ac-58f9df641b4f	1	\N
86fc2753-4453-43bc-8660-4e079bdec649	51af7bd7-5340-4970-ac2a-61bba55cf810	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	230457f2-442e-4eb3-80b0-7feb166ec89d
df23db9c-a7a6-43b2-9748-3a755305264c	570a39e6-bc00-4985-aa2c-bf24ed2d4375	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	8bcb06cf-5b5e-4e2c-b8b9-a766449e409a
fe92c3d4-697f-4fc7-832d-eeb9cc2857f1	59069426-28fa-45c7-8d5b-54503c66165c	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	f713b5c9-223b-45ad-94de-1122170b287f
a8b0d2bf-18f3-4d58-8c0b-3910618ad4ab	af0bd3a8-de7d-4f6f-ac96-db08cd1fb0c4	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	53db2ad3-fd81-4f75-974d-8a633b8a86f5
34af291b-3f5e-4df2-a364-38b9cf73a5bb	9e717db0-b052-4834-a3f9-f8ec4279e19f	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	a1d6a6a0-3749-4303-a799-30d09b346245
2b72270e-b16e-4af2-94aa-7874ed542ea4	9bbd6df5-403d-4b06-828f-42b3979ec189	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	6d4c35df-6eeb-4a42-a59c-fa023ba79096
e3b567d1-0768-4e95-b9f4-acb05d0c5c93	cfb9dcda-8034-49a7-82f1-8744022ee386	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	7a304c19-ae39-4673-802b-e7019c06ed81
012d49c6-b366-4ea2-865f-89b0ed6eda57	103a7304-ed77-4f89-9ba5-1746383c6cba	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	92d41b0a-acda-4500-90c4-cf6d697d0dfc
c6694892-e871-4f9c-af8f-eac71e9660cb	a1b425ea-0777-4d0d-965b-5082a819db50	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	87cd4448-e066-4c10-a0f8-1e1faa25b499
0cf96929-3af5-443b-b98d-3ca13f52d5f4	a190d65f-0532-4580-945a-a1efbedb5d4a	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	19031419-655e-44bf-9863-c68c70b42ad0
62da083d-4467-42a5-8028-291404cc426c	ad4a7da0-ca3b-4e94-a08a-160c9357931d	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	75f07327-df11-46d1-bbe3-08f089791c8b
c1ce9a8e-686c-4b42-af7f-6b9322b8a214	7aee9783-f7fc-41e5-b369-b295796f2366	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	c3de2651-f9f3-4b48-8ed7-e6474b5bf1d3
9a1ed387-664d-4e6e-b226-75a095bc6820	2925f911-f566-47fc-9ffe-a9488127b7fc	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	7cf6f81a-664b-47b6-949e-35385a3d3270
0b50a078-24b4-4c11-a65b-eb9bd0ffdeb5	ebd81f67-bf05-4329-8837-6a6c19dcd454	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	7a0d4a54-d7b7-4fde-8942-55a04d5577d8
8d50780b-f07c-4026-a5af-8a1489992d18	298b0158-fd71-4224-ba82-67f7d753e108	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	2fb378e7-d3be-41ba-bfae-fe84e7c985f2
861cd946-801e-4aec-a37d-2e1254c6808d	d7f8a04d-1d50-40cf-9be4-1cabc516f15b	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	c60622f1-fea1-4c60-b80c-4a83ed2be749
48b155b0-11b7-4ac9-9110-4cb995cc4e35	850d5366-3794-474d-b863-bf6beaeb2d36	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	a1484e78-3114-4773-9cd1-fd7dfdb3fe54
2d1087f6-6855-4381-9814-68bffe5f2eac	aedbe475-34da-4a74-801e-a3752873ef8f	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	9722f33c-6e4d-45ea-a1e8-346b318c1dda
d6a899ae-9f0a-46f4-9d25-a28179427516	afc7246a-4a83-44d2-81d7-278ebab69699	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	ac220e0a-de5e-45ed-8635-8a0eaae954aa
1a10d7c6-8aa6-44c8-852c-c03e525a9d51	d07d5ba7-5520-4782-bd60-a296058a22e4	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	0b3ab0c2-0ef3-42a3-a6ab-9bd4bcc9cd2f
8b2d7b15-0a6c-4d93-96f1-a31babf72c9b	d5c49bb6-339a-4a58-9ff0-c360087993b0	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	a55619ee-2d84-4572-8a92-f93f506fa5e7
955653fa-eb2c-494c-b23d-5a2ec8a41af2	a1093fd8-3282-473d-8713-bf8500ad2c5b	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	d83944c7-221b-4499-96ea-d22dd53b9e4e
46dfc90b-0362-4219-a34a-264417e379dc	59153e98-ae5d-45a0-8e0f-f3fd2fe0435d	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	e4bbb1a5-8e7b-46b3-be05-0c103fa57a30
4a94667f-42f8-4a6b-8160-af04985d6d63	b44f2f01-9882-48ab-bc7c-8ea84cc30f9c	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	960f3ecb-6ce2-4224-88b9-abf8241825f6
3fee6b67-ff95-4f41-927d-df968bd3d58e	b1f16ff5-0fb9-4380-b573-d5ffde39ba59	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	a58fc4df-9977-4289-ac15-b2defd44a741
983326ed-67fb-4103-84d6-ad70e870b84a	f7e1401b-8706-4858-9fa8-e4aa4cbba26e	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	627c4181-535a-466b-8af3-c896dbe98eb5
e87a9a4b-e85e-44be-9b0d-363f92f9c66f	7160f8f6-77a7-4864-8d43-1ef6fbdde1b5	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	4a06e118-4ae2-4ede-99fa-128dc83cd510
31d3ae62-f122-4fb0-bb10-79e03aab56ab	9329b2e7-5c15-447d-805b-b1f8d6d3bdf4	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	d40732b3-c42c-495b-8b33-fe4d1a7d80e1
02275d48-e3e3-43da-8c9e-60580157fcc7	a15772fc-002a-4e30-b90a-52c850191858	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	79e917e5-660f-4d4c-8ba8-32ef9d37c229
750d2c47-7816-46ea-9986-18b05f665bdb	4d1d8708-f792-45fe-9932-21de489833bf	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	4a8f8ca5-995a-4253-9536-23333720b7c4
2330dcd6-35e0-4ce5-a299-072c3ff52149	d3ed0394-d72e-4d4b-8023-72fd2cf72b91	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	0a5ba720-77d5-4fc6-8ee1-5ade154ef999
d2c03e79-5f05-42d3-96b5-07d0270f931d	84a00a91-1625-43db-83ac-78a9d92c0962	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	1c444910-7c78-409c-852a-11d186c31578
89c41314-f8e5-4458-abc4-46b30d903806	6c320299-71b2-4811-9d1a-0210fa33d864	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	67f8edf9-03ed-4387-b855-2f14a2a55035
36c8b7fb-7536-49ec-b7a4-268311ad7728	772a04e3-29ca-470c-b8f9-a54221d88de1	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	4ac04bc3-7209-4b80-bf41-de836563be7b
32a5d074-b0db-4da0-9649-5fbbce994184	d8e0dc0f-a992-4782-8fae-1b50b90185e1	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	fe351879-2b3a-49e0-80f7-e7b9ffb8203b
4dad1ee8-54cf-461b-a01b-1c3380ada0f0	c189fea9-e9ac-4722-8f9a-799037a4035c	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	2c5f4d28-7147-4d40-9638-78b4d6e4347b
9bb7b10b-8da4-42c6-8f55-9bd4d9a7fbf9	d0e4a50b-3888-4c08-85b7-bbc9b2bec7d9	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	79f1e0c2-1e2c-4418-9c3f-c61698e3b81e
a8d20f75-fa6f-416c-b4f6-8f1d5dec6d14	a1d3a9d8-3a6d-411d-a1a1-c8dfbe0e71ad	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	4fb347a1-a5b1-426a-9aee-80122f4b719e
11ad83dc-12a9-42a5-85dc-4c29ff492686	9bd63b6e-2560-48cf-bf88-54b89904271b	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	e7b7c7b0-1407-46f8-accf-c178408bbf17
dc7d5119-632b-4da7-8ecc-797a6fb97023	00eaa206-0780-4c97-89c7-5245593a747e	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	a7742526-a89e-4b0a-93fc-3925d487ca3f
822cfa97-77ba-42a6-aea0-ee8c3ff2ba5f	fb8f5382-6ba2-44ed-9170-c0ad9ead2d86	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	407cf304-0c12-4fb7-beaf-b1164247e366
ce26a795-bcd9-42d2-8ce7-adf920b6c195	84ac7e4e-76c9-4d76-a84e-87fa5adc65b8	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	46cb862f-6806-41be-ae82-8a003a2c53c7
1a9f5ef8-bc62-400c-8782-85875571895e	57663c7f-c1e6-4e4c-b12a-b7273eb2f2be	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	570b02c0-74b7-4c64-b541-33f15c12e0bc
7b1ca76f-1fcf-435c-8a14-19a5ce432430	c61746a4-50f6-48d6-907e-7bfffdd17776	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	a6853850-640e-464c-8f9e-a31544524197
9a2fc0c1-c528-47bb-a170-a8e1f5752c93	0dc954b2-401a-4e53-8b14-8d2e7a3eea16	a30253b6-4e38-43f9-9e93-2e8ab763d291	1	e26f128f-43d2-4bc0-a1d5-cf7748986488
\.


--
-- Data for Name: shipment_refuels; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.shipment_refuels (id, shipment_id, vehicle_id, liters, odometer_km, station_name, cost_uah, created_by, tenant_id, created_at) FROM stdin;
f80c7fc1-7e6d-4874-843a-231bbc55cf06	a1d3a9d8-3a6d-411d-a1a1-c8dfbe0e71ad	ea092cad-3147-4245-8ce6-ee7b8029806b	5.00	2	WOG	450.00	21833e68-101f-4378-a916-62c120a9f192	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-30 22:02:20.735099+00
\.


--
-- Data for Name: shipments; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.shipments (id, from_warehouse_id, to_warehouse_id, vehicle_id, priority, status, created_at, updated_at, tenant_id, started_at, delivered_at, direction, distance_km, actual_km) FROM stdin;
570a39e6-bc00-4985-aa2c-bf24ed2d4375	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-06 21:07:03.669682+00	2026-05-17 00:35:14.207682	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-06 21:13:06.198216+00	2026-05-17 00:35:14.207682+00	DOWNSTREAM	0	0
55a417c4-d5a8-4656-8a14-741f2c2c0273	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-03 14:13:30.063369+00	2026-05-03 14:31:25.470134	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	\N	\N	\N	0	0
8839bbd3-ea60-4e00-b25f-823449d5e004	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-03 15:40:42.319113+00	2026-05-03 15:49:26.029477	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	\N	\N	\N	0	0
5d03ed08-289e-4825-af7c-d9134ecd4107	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-03 16:05:32.132838+00	2026-05-03 16:09:39.334354	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	\N	\N	\N	0	0
a5504524-8653-44aa-8c38-a86832a5963e	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-03 16:20:44.673613+00	2026-05-03 16:24:14.641374	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	\N	\N	\N	0	0
2f71fc41-c742-48b4-8d2f-9c7b6263c3e0	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-03 16:25:00.82014+00	2026-05-03 16:38:21.830721	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	\N	\N	\N	0	0
20bb818e-2a8c-4058-b865-5f52388a30ea	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-03 16:38:53.56172+00	2026-05-03 16:44:15.773572	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	\N	\N	\N	0	0
0c54467b-ad81-4d0b-8554-2a2e44e46b22	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-03 16:44:40.546601+00	2026-05-03 16:58:00.958671	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	\N	\N	\N	0	0
798d6084-4572-441c-9f52-3f05892de225	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-03 17:19:42.471218+00	2026-05-03 17:45:33.292524	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-03 17:39:58.799884+00	2026-05-03 17:45:33.292524+00	\N	0	0
59069426-28fa-45c7-8d5b-54503c66165c	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-17 00:42:45.685204+00	2026-05-27 21:14:32.595674	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-17 00:43:50.162876+00	2026-05-27 21:14:32.595674+00	DOWNSTREAM	0	0
7e631304-8a13-4a8b-9605-e83fbd90e834	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-03 17:49:31.41575+00	2026-05-03 17:50:16.311385	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-03 17:50:06.02148+00	2026-05-03 17:50:16.311385+00	\N	0	0
9cd279eb-1626-4932-938c-63cd7f4829e2	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-03 17:58:11.756454+00	2026-05-03 18:12:18.672828	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-03 18:12:15.461233+00	2026-05-03 18:12:18.672828+00	\N	0	0
a190d65f-0532-4580-945a-a1efbedb5d4a	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	ea092cad-3147-4245-8ce6-ee7b8029806b	NORMAL	DELIVERED	2026-05-28 12:23:37.856034+00	2026-05-28 12:23:46.884941	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-28 12:23:46.017705+00	2026-05-28 12:23:46.884941+00	DOWNSTREAM	0	0
b05dd831-264e-4bed-a7a9-e3e3d06c1393	f64e8882-623d-41cd-8f72-2be30b783f8d	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-03 19:49:15.058128+00	2026-05-03 19:52:19.448182	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-03 19:52:17.854568+00	2026-05-03 19:52:19.448182+00	UPSTREAM	0	0
af0bd3a8-de7d-4f6f-ac96-db08cd1fb0c4	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-27 21:21:41.214277+00	2026-05-27 21:21:55.867446	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-27 21:21:54.286562+00	2026-05-27 21:21:55.867446+00	DOWNSTREAM	0	0
8e1ca805-ed59-4b6f-bded-84ad66ae0138	f64e8882-623d-41cd-8f72-2be30b783f8d	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-03 20:01:53.807223+00	2026-05-03 20:28:14.895037	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-03 20:02:12.395434+00	2026-05-03 20:28:14.895037+00	UPSTREAM	0	0
34af91e3-95bb-4633-a8ef-3ce9a7336400	f64e8882-623d-41cd-8f72-2be30b783f8d	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-03 21:38:16.967059+00	2026-05-05 12:33:33.077516	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-05 12:33:32.185239+00	2026-05-05 12:33:33.077516+00	UPSTREAM	0	0
704ce237-6254-415d-aadb-7089f825e48b	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-05 12:33:45.632358+00	2026-05-05 12:34:04.979438	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-05 12:33:57.542961+00	2026-05-05 12:34:04.979438+00	DOWNSTREAM	0	0
9e717db0-b052-4834-a3f9-f8ec4279e19f	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-27 21:49:18.529334+00	2026-05-27 21:49:34.558178	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-27 21:49:24.740865+00	2026-05-27 21:49:34.558178+00	DOWNSTREAM	0	0
8ff169fe-d4d5-41d5-ab14-355008cc9a0e	f64e8882-623d-41cd-8f72-2be30b783f8d	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-05 18:32:10.657133+00	2026-05-06 21:04:26.387872	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-06 17:46:49.346229+00	2026-05-06 21:04:26.387872+00	UPSTREAM	0	0
51af7bd7-5340-4970-ac2a-61bba55cf810	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-06 21:04:40.169544+00	2026-05-06 21:05:19.823429	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-06 21:05:19.099879+00	2026-05-06 21:05:19.823429+00	DOWNSTREAM	0	0
298b0158-fd71-4224-ba82-67f7d753e108	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	ea092cad-3147-4245-8ce6-ee7b8029806b	NORMAL	DELIVERED	2026-05-28 14:31:54.541627+00	2026-05-28 14:40:29.472113	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-28 14:36:52.414023+00	2026-05-28 14:40:29.472113+00	DOWNSTREAM	0	0
ad4a7da0-ca3b-4e94-a08a-160c9357931d	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	dac43200-3bfe-47c8-8e60-17bf1ce6cf40	NORMAL	DELIVERED	2026-05-28 12:24:08.989997+00	2026-05-28 12:24:15.621419	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-28 12:24:14.992227+00	2026-05-28 12:24:15.621419+00	DOWNSTREAM	0	0
9bbd6df5-403d-4b06-828f-42b3979ec189	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-27 21:52:09.657879+00	2026-05-28 07:45:38.395514	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-28 07:45:35.963456+00	2026-05-28 07:45:38.395514+00	DOWNSTREAM	0	0
cfb9dcda-8034-49a7-82f1-8744022ee386	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	8355c89e-b072-4113-a6d2-adecc94a5a31	NORMAL	DELIVERED	2026-05-28 08:12:02.118154+00	2026-05-28 08:27:33.748915	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-28 08:14:43.613106+00	2026-05-28 08:27:33.748915+00	DOWNSTREAM	0	0
a1b425ea-0777-4d0d-965b-5082a819db50	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	ea092cad-3147-4245-8ce6-ee7b8029806b	NORMAL	DELIVERED	2026-05-28 12:11:43.982237+00	2026-05-28 12:12:11.336955	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-28 12:12:08.0196+00	2026-05-28 12:12:11.336955+00	DOWNSTREAM	0	0
103a7304-ed77-4f89-9ba5-1746383c6cba	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	8355c89e-b072-4113-a6d2-adecc94a5a31	NORMAL	DELIVERED	2026-05-28 08:29:12.905425+00	2026-05-28 12:12:12.103867	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-28 12:12:10.543321+00	2026-05-28 12:12:12.103867+00	DOWNSTREAM	0	0
d7f8a04d-1d50-40cf-9be4-1cabc516f15b	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	1565aa02-d4ff-464c-bc5b-c96f63fd9be1	NORMAL	DELIVERED	2026-05-28 14:32:28.701979+00	2026-05-28 14:40:32.00968	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-28 14:40:31.310923+00	2026-05-28 14:40:32.00968+00	DOWNSTREAM	0	0
7aee9783-f7fc-41e5-b369-b295796f2366	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	ea092cad-3147-4245-8ce6-ee7b8029806b	NORMAL	DELIVERED	2026-05-28 13:07:33.995494+00	2026-05-28 14:31:06.923009	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-28 13:11:26.94303+00	2026-05-28 14:31:06.923009+00	DOWNSTREAM	0	0
2925f911-f566-47fc-9ffe-a9488127b7fc	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	1565aa02-d4ff-464c-bc5b-c96f63fd9be1	NORMAL	DELIVERED	2026-05-28 13:12:06.036328+00	2026-05-28 14:31:07.276485	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-28 14:31:06.393622+00	2026-05-28 14:31:07.276485+00	DOWNSTREAM	0	0
ebd81f67-bf05-4329-8837-6a6c19dcd454	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-28 13:12:40.452481+00	2026-05-28 14:31:08.163974	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-28 14:31:05.35296+00	2026-05-28 14:31:08.163974+00	DOWNSTREAM	0	0
aedbe475-34da-4a74-801e-a3752873ef8f	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-28 15:46:22.751759+00	2026-05-28 16:15:35.31537	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-28 16:15:30.510712+00	2026-05-28 16:15:35.31537+00	DOWNSTREAM	0	0
850d5366-3794-474d-b863-bf6beaeb2d36	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	ea092cad-3147-4245-8ce6-ee7b8029806b	NORMAL	DELIVERED	2026-05-28 15:44:49.568196+00	2026-05-28 15:47:43.65876	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-28 15:47:27.852755+00	2026-05-28 15:47:43.65876+00	DOWNSTREAM	0	0
d5c49bb6-339a-4a58-9ff0-c360087993b0	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	dac43200-3bfe-47c8-8e60-17bf1ce6cf40	NORMAL	DELIVERED	2026-05-28 16:14:32.876511+00	2026-05-28 16:15:33.480339	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-28 16:15:32.958544+00	2026-05-28 16:15:33.480339+00	DOWNSTREAM	0	0
d07d5ba7-5520-4782-bd60-a296058a22e4	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	ea092cad-3147-4245-8ce6-ee7b8029806b	NORMAL	DELIVERED	2026-05-28 15:48:26.001626+00	2026-05-28 16:15:34.135245	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-28 16:15:31.905808+00	2026-05-28 16:15:34.135245+00	DOWNSTREAM	0	0
afc7246a-4a83-44d2-81d7-278ebab69699	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	1565aa02-d4ff-464c-bc5b-c96f63fd9be1	NORMAL	DELIVERED	2026-05-28 15:46:27.8403+00	2026-05-28 16:15:34.745336	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-28 16:15:31.044457+00	2026-05-28 16:15:34.745336+00	DOWNSTREAM	0	0
b44f2f01-9882-48ab-bc7c-8ea84cc30f9c	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-28 20:58:42.22901+00	2026-05-29 07:31:30.539102	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-28 21:51:56.649445+00	2026-05-29 07:31:30.539102+00	DOWNSTREAM	0	1
f7e1401b-8706-4858-9fa8-e4aa4cbba26e	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	8355c89e-b072-4113-a6d2-adecc94a5a31	NORMAL	DELIVERED	2026-05-28 21:30:08.941903+00	2026-05-29 07:55:28.763902	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-29 07:55:26.549999+00	2026-05-29 07:55:28.763902+00	DOWNSTREAM	0	1
b1f16ff5-0fb9-4380-b573-d5ffde39ba59	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	dac43200-3bfe-47c8-8e60-17bf1ce6cf40	NORMAL	DELIVERED	2026-05-28 21:29:50.232007+00	2026-05-29 07:55:32.737308	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-29 07:55:29.581854+00	2026-05-29 07:55:32.737308+00	DOWNSTREAM	0	1
59153e98-ae5d-45a0-8e0f-f3fd2fe0435d	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	1565aa02-d4ff-464c-bc5b-c96f63fd9be1	NORMAL	DELIVERED	2026-05-28 20:58:32.567807+00	2026-05-29 07:55:36.162222	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-29 07:55:33.929297+00	2026-05-29 07:55:36.162222+00	DOWNSTREAM	0	1
a1093fd8-3282-473d-8713-bf8500ad2c5b	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	ea092cad-3147-4245-8ce6-ee7b8029806b	NORMAL	DELIVERED	2026-05-28 20:58:22.145394+00	2026-05-29 07:55:39.657918	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-29 07:55:37.353945+00	2026-05-29 07:55:39.657918+00	DOWNSTREAM	0	1
7160f8f6-77a7-4864-8d43-1ef6fbdde1b5	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-29 07:33:02.530928+00	2026-05-29 07:55:25.962499	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-29 07:55:23.985182+00	2026-05-29 07:55:25.962499+00	DOWNSTREAM	0	1
0dc954b2-401a-4e53-8b14-8d2e7a3eea16	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	dac43200-3bfe-47c8-8e60-17bf1ce6cf40	NORMAL	PENDING	2026-05-30 21:30:31.999683+00	2026-05-30 21:30:31.999683	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	\N	\N	DOWNSTREAM	8.1408	0
d0e4a50b-3888-4c08-85b7-bbc9b2bec7d9	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	ea092cad-3147-4245-8ce6-ee7b8029806b	NORMAL	DELIVERED	2026-05-29 22:48:54.567791+00	2026-05-29 23:07:38.798591	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-29 22:49:17.00594+00	2026-05-29 23:07:38.798591+00	DOWNSTREAM	0	1
84a00a91-1625-43db-83ac-78a9d92c0962	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	8355c89e-b072-4113-a6d2-adecc94a5a31	NORMAL	DELIVERED	2026-05-29 09:18:19.015723+00	2026-05-29 09:36:48.874212	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-29 09:36:47.335212+00	2026-05-29 09:36:48.874212+00	DOWNSTREAM	0	1
d3ed0394-d72e-4d4b-8023-72fd2cf72b91	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	dac43200-3bfe-47c8-8e60-17bf1ce6cf40	NORMAL	DELIVERED	2026-05-29 08:35:34.704693+00	2026-05-29 09:36:51.861145	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-29 09:36:50.018278+00	2026-05-29 09:36:51.861145+00	DOWNSTREAM	0	1
4d1d8708-f792-45fe-9932-21de489833bf	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	1565aa02-d4ff-464c-bc5b-c96f63fd9be1	NORMAL	DELIVERED	2026-05-29 07:57:30.495852+00	2026-05-29 09:36:54.577963	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-29 09:36:52.635856+00	2026-05-29 09:36:54.577963+00	DOWNSTREAM	0	1
a1d3a9d8-3a6d-411d-a1a1-c8dfbe0e71ad	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	ea092cad-3147-4245-8ce6-ee7b8029806b	NORMAL	DELIVERED	2026-05-29 23:07:58.024204+00	2026-05-31 22:53:35.193414	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-29 23:08:14.815047+00	2026-05-31 22:53:35.193414+00	DOWNSTREAM	0	16
9bd63b6e-2560-48cf-bf88-54b89904271b	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	1565aa02-d4ff-464c-bc5b-c96f63fd9be1	NORMAL	DELIVERED	2026-05-29 23:24:20.759512+00	2026-05-30 10:17:25.310516	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-29 23:24:25.868492+00	2026-05-30 10:17:25.310516+00	DOWNSTREAM	0	8
a15772fc-002a-4e30-b90a-52c850191858	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	48471a6b-ac6a-4327-b13c-dce6768ac4a0	NORMAL	DELIVERED	2026-05-29 07:57:08.396853+00	2026-05-29 09:36:57.56492	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-29 09:36:55.552893+00	2026-05-29 09:36:57.56492+00	DOWNSTREAM	0	1
9329b2e7-5c15-447d-805b-b1f8d6d3bdf4	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	ea092cad-3147-4245-8ce6-ee7b8029806b	NORMAL	DELIVERED	2026-05-29 07:56:59.674937+00	2026-05-29 09:37:00.091532	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-29 09:36:58.337805+00	2026-05-29 09:37:00.091532+00	DOWNSTREAM	0	1
00eaa206-0780-4c97-89c7-5245593a747e	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	1565aa02-d4ff-464c-bc5b-c96f63fd9be1	NORMAL	DELIVERED	2026-05-30 10:17:43.748931+00	2026-05-30 12:48:02.825064	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-30 10:17:50.894727+00	2026-05-30 12:48:02.825064+00	DOWNSTREAM	0	8
6c320299-71b2-4811-9d1a-0210fa33d864	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	ea092cad-3147-4245-8ce6-ee7b8029806b	NORMAL	DELIVERED	2026-05-29 14:59:26.595929+00	2026-05-29 21:31:58.173825	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-29 14:59:39.583147+00	2026-05-29 21:31:58.173825+00	DOWNSTREAM	0	1
d8e0dc0f-a992-4782-8fae-1b50b90185e1	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	1565aa02-d4ff-464c-bc5b-c96f63fd9be1	NORMAL	DELIVERED	2026-05-29 22:05:13.440684+00	2026-05-29 22:33:47.197123	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-29 22:05:45.390935+00	2026-05-29 22:33:47.197123+00	DOWNSTREAM	0	1
fb8f5382-6ba2-44ed-9170-c0ad9ead2d86	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	1565aa02-d4ff-464c-bc5b-c96f63fd9be1	NORMAL	DELIVERED	2026-05-30 12:48:21.535056+00	2026-05-30 12:52:08.804323	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-30 12:48:30.997518+00	2026-05-30 12:52:08.804323+00	DOWNSTREAM	0	8
772a04e3-29ca-470c-b8f9-a54221d88de1	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	ea092cad-3147-4245-8ce6-ee7b8029806b	NORMAL	DELIVERED	2026-05-29 21:50:55.289762+00	2026-05-29 22:33:50.827435	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-29 21:51:14.011361+00	2026-05-29 22:33:50.827435+00	DOWNSTREAM	0	1
c189fea9-e9ac-4722-8f9a-799037a4035c	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	ea092cad-3147-4245-8ce6-ee7b8029806b	NORMAL	DELIVERED	2026-05-29 22:34:39.153235+00	2026-05-29 22:48:29.468411	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-29 22:40:44.706081+00	2026-05-29 22:48:29.468411+00	DOWNSTREAM	0	1
84ac7e4e-76c9-4d76-a84e-87fa5adc65b8	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	1565aa02-d4ff-464c-bc5b-c96f63fd9be1	NORMAL	DELIVERED	2026-05-30 12:54:09.853363+00	2026-05-30 13:12:25.804356	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-30 12:54:41.371773+00	2026-05-30 13:12:25.804356+00	DOWNSTREAM	0	16
57663c7f-c1e6-4e4c-b12a-b7273eb2f2be	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	1565aa02-d4ff-464c-bc5b-c96f63fd9be1	NORMAL	IN_TRANSIT	2026-05-30 13:13:26.000459+00	2026-05-30 13:15:06.599269	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-30 13:15:06.599269+00	\N	DOWNSTREAM	8.1408	0
c61746a4-50f6-48d6-907e-7bfffdd17776	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	f64e8882-623d-41cd-8f72-2be30b783f8d	dac43200-3bfe-47c8-8e60-17bf1ce6cf40	NORMAL	DELIVERED	2026-05-30 15:15:16.429311+00	2026-05-30 15:16:51.525543	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2026-05-30 15:16:26.741637+00	2026-05-30 15:16:51.525543+00	DOWNSTREAM	8.1408	16
\.


--
-- Data for Name: supply_requests; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.supply_requests (id, tenant_id, created_by, resource_id, quantity, status, approved_by, approved_at, comment, created_at, updated_at, target_warehouse_id, resource_name, resource_category_id) FROM stdin;
e948e8a6-1151-4086-8a32-41f1efd5350e	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	100	APPROVED	21833e68-101f-4378-a916-62c120a9f192	2026-05-03 21:44:25.403585+00		2026-05-03 21:44:25.400083+00	2026-05-03 21:44:25.403585+00	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	Гідравлічний візок (Рохля) 2.5т	6a273219-ea10-4ea7-aff4-9647696668f4
4fb347a1-a5b1-426a-9aee-80122f4b719e	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 23:07:49.667844+00		2026-05-29 23:07:49.654453+00	2026-05-31 22:53:35.193414+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
617e59a8-38a6-4cf7-8c5d-12c13a09ae2d	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-03 21:21:13.84631+00	Автоматична ескалація SLA	2026-05-03 21:21:13.831048+00	2026-05-05 12:33:33.077516+00	c3d3f8df-0be9-4931-91f4-a76554bfb5f9	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
59f2d6d8-3472-46e7-8879-502ade5630c4	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	DISPATCHED	21833e68-101f-4378-a916-62c120a9f192	2026-05-28 21:30:28.057574+00		2026-05-28 21:30:28.044693+00	2026-05-28 21:51:56.649445+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
44878184-a03c-4e93-9b39-bd9ba1b650e6	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-05 12:33:10.098974+00		2026-05-05 12:33:10.092543+00	2026-05-05 12:34:04.979438+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
230457f2-442e-4eb3-80b0-7feb166ec89d	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-06 17:45:40.178951+00		2026-05-06 17:45:40.163307+00	2026-05-06 21:05:19.823429+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
3feb3240-4a1b-4f0d-834c-037b62a8efa9	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	DISPATCHED	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 07:51:51.003271+00		2026-05-29 07:51:50.988392+00	2026-05-29 07:55:23.985182+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
4a06e118-4ae2-4ede-99fa-128dc83cd510	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 07:31:52.889777+00		2026-05-29 07:31:52.875979+00	2026-05-29 07:55:25.962499+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
627c4181-535a-466b-8af3-c896dbe98eb5	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-28 21:30:03.834568+00		2026-05-28 21:30:03.820003+00	2026-05-29 07:55:28.763902+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
1c444910-7c78-409c-852a-11d186c31578	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 08:14:48.363844+00		2026-05-29 08:14:48.354156+00	2026-05-29 09:36:48.874212+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
0a5ba720-77d5-4fc6-8ee1-5ade154ef999	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 07:59:12.19502+00		2026-05-29 07:59:12.185292+00	2026-05-29 09:36:51.861145+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
ecc2d04a-15bb-439f-a7fc-9f070cff31ff	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	DISPATCHED	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 09:33:12.194574+00		2026-05-29 09:33:12.188947+00	2026-05-29 09:36:47.335212+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
79e917e5-660f-4d4c-8ba8-32ef9d37c229	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 07:56:32.922681+00		2026-05-29 07:56:32.906387+00	2026-05-29 09:36:57.56492+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
d40732b3-c42c-495b-8b33-fe4d1a7d80e1	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 07:56:07.275609+00		2026-05-29 07:56:07.260304+00	2026-05-29 09:37:00.091532+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
b3ef8912-7968-4d8f-8057-4a7b5109ee4d	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	APPROVED	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 12:13:24.250339+00		2026-05-29 12:13:24.238463+00	2026-05-29 12:13:24.250339+00	62ba10c5-14ee-4628-bfcc-cfb2f99d4e46	GPS-трекер Teltonika FMB120	08a0ee32-49dc-4cec-aa93-1f95909488af
3482bf32-472b-48af-9db0-c40e10294bb8	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	APPROVED	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 12:15:03.654622+00		2026-05-29 12:15:03.641455+00	2026-05-29 12:15:03.654622+00	29918579-9237-4f65-9ce3-db58b825b86c	GPS-трекер Teltonika FMB120	08a0ee32-49dc-4cec-aa93-1f95909488af
67f8edf9-03ed-4387-b855-2f14a2a55035	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 09:37:30.413969+00		2026-05-29 09:37:30.402202+00	2026-05-29 21:31:58.173825+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
217ea66a-0efb-4ad6-9954-4ee380edd4f0	0997831c-654f-471b-934f-cedafbc54ea5	3e59277f-0c1a-4aa3-973c-cf3db1e19497	615d5262-b41b-488f-b355-c71185ebdcb0	5	REJECTED	3e59277f-0c1a-4aa3-973c-cf3db1e19497	2026-04-27 10:33:21.261126+00	Терміново	2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
25a0adff-f4dc-491d-beed-72302c6d803c	0997831c-654f-471b-934f-cedafbc54ea5	3e59277f-0c1a-4aa3-973c-cf3db1e19497	615d5262-b41b-488f-b355-c71185ebdcb0	21	REJECTED	3e59277f-0c1a-4aa3-973c-cf3db1e19497	2026-04-18 10:33:21.261951+00	Терміново	2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
3e9f6966-09ac-4818-86a9-c9cd9b3df69c	0997831c-654f-471b-934f-cedafbc54ea5	3e59277f-0c1a-4aa3-973c-cf3db1e19497	00737572-6cf5-4e7c-81be-aa295129cc72	11	COMPLETED	3e59277f-0c1a-4aa3-973c-cf3db1e19497	2026-04-22 10:33:21.262676+00		2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
a80b339b-82a9-4a97-a5ce-220d3262f955	0997831c-654f-471b-934f-cedafbc54ea5	3e59277f-0c1a-4aa3-973c-cf3db1e19497	9a074b6d-3c73-4d72-bea5-c9d3bdb6a395	3	REJECTED	3e59277f-0c1a-4aa3-973c-cf3db1e19497	2026-04-23 10:33:21.263277+00	Планова поставка	2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
fe351879-2b3a-49e0-80f7-e7b9ffb8203b	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 22:05:03.82078+00		2026-05-29 22:05:03.810039+00	2026-05-29 22:33:47.197123+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
b6257ecb-250c-49cf-ad81-287d18a3e20a	0997831c-654f-471b-934f-cedafbc54ea5	601e6d16-947b-4eb0-a67c-5304edafedd8	38df8f83-898c-4bf2-8b4b-08721c9a3bc7	23	REJECTED	3e59277f-0c1a-4aa3-973c-cf3db1e19497	2026-04-23 10:33:21.264364+00	Терміново	2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
955f81dd-ac22-4f18-bc19-783c827ddebe	0997831c-654f-471b-934f-cedafbc54ea5	601e6d16-947b-4eb0-a67c-5304edafedd8	00737572-6cf5-4e7c-81be-aa295129cc72	15	COMPLETED	3e59277f-0c1a-4aa3-973c-cf3db1e19497	2026-04-23 10:33:21.264825+00		2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
4ac04bc3-7209-4b80-bf41-de836563be7b	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 21:33:29.581486+00		2026-05-29 21:33:29.567485+00	2026-05-29 22:33:50.827435+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
d2561c47-f578-48da-a1bb-259361ebf322	0997831c-654f-471b-934f-cedafbc54ea5	601e6d16-947b-4eb0-a67c-5304edafedd8	d902e3a4-c9ed-4dca-a8f9-ea0bd6b82b55	13	REJECTED	3e59277f-0c1a-4aa3-973c-cf3db1e19497	2026-05-03 10:33:21.266385+00	Терміново	2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
8cca3a73-5564-4f47-a697-4539ba62045d	0997831c-654f-471b-934f-cedafbc54ea5	3e59277f-0c1a-4aa3-973c-cf3db1e19497	00737572-6cf5-4e7c-81be-aa295129cc72	30	ESCALATED	\N	\N	Автоматична ескалація SLA	2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
592d1bb0-aca0-4696-ba3a-afaf0a91ee38	0997831c-654f-471b-934f-cedafbc54ea5	601e6d16-947b-4eb0-a67c-5304edafedd8	ab810d5d-ca53-4668-a44b-91f438e3db19	30	ESCALATED	\N	\N	Автоматична ескалація SLA	2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
3b670dbb-11af-4eeb-9e09-fd743b30f14b	0997831c-654f-471b-934f-cedafbc54ea5	601e6d16-947b-4eb0-a67c-5304edafedd8	54b243f1-bd9d-489e-afbf-9a44b3bf0a28	11	ESCALATED	\N	\N	Автоматична ескалація SLA	2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
b967d3a4-1450-41d5-8d3b-144e75530da8	0997831c-654f-471b-934f-cedafbc54ea5	601e6d16-947b-4eb0-a67c-5304edafedd8	d902e3a4-c9ed-4dca-a8f9-ea0bd6b82b55	19	ESCALATED	\N	\N	Автоматична ескалація SLA	2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
33fa3eb0-15e1-4bf4-aec8-276d818742bd	0997831c-654f-471b-934f-cedafbc54ea5	3e59277f-0c1a-4aa3-973c-cf3db1e19497	38df8f83-898c-4bf2-8b4b-08721c9a3bc7	1	ESCALATED	\N	\N	Автоматична ескалація SLA	2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
8bcb06cf-5b5e-4e2c-b8b9-a766449e409a	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-06 21:05:36.439837+00		2026-05-06 21:05:36.437634+00	2026-05-17 00:35:14.207682+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
2c5f4d28-7147-4d40-9638-78b4d6e4347b	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 22:34:14.936846+00		2026-05-29 22:34:14.923227+00	2026-05-29 22:48:29.468411+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
25b6bbb9-853e-43df-a565-6e708566f756	0997831c-654f-471b-934f-cedafbc54ea5	135fe7a0-2fc0-4b17-9e22-4d1c32d9658b	49aaee56-7beb-4736-b71e-469dd6bd75fc	8	REJECTED	3e59277f-0c1a-4aa3-973c-cf3db1e19497	2026-04-28 10:33:21.267789+00		2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
c6366438-d224-40fd-a359-735e829aa4dc	0997831c-654f-471b-934f-cedafbc54ea5	135fe7a0-2fc0-4b17-9e22-4d1c32d9658b	d95b0275-5b2d-47c0-9663-a01c096533c9	13	APPROVED	3e59277f-0c1a-4aa3-973c-cf3db1e19497	2026-04-28 10:33:21.26815+00	Поповнення резерву	2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
b9a3679c-f30a-4cc4-9841-9d5436f07351	0997831c-654f-471b-934f-cedafbc54ea5	135fe7a0-2fc0-4b17-9e22-4d1c32d9658b	1fb6c925-ffa4-4dbb-a221-0e76c3ee47b2	4	APPROVED	3e59277f-0c1a-4aa3-973c-cf3db1e19497	2026-05-04 10:33:21.268949+00		2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
d22e9be0-b1cf-4172-b06d-f398201478db	0997831c-654f-471b-934f-cedafbc54ea5	f03a133e-617b-4c6f-9bab-3560b041c358	a28a5172-e9d3-4134-9168-70c71678d5a1	21	COMPLETED	f03a133e-617b-4c6f-9bab-3560b041c358	2026-04-18 10:33:21.269453+00		2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
e7d8ae49-e44c-4a55-b6a0-9211fd946c4d	0997831c-654f-471b-934f-cedafbc54ea5	f03a133e-617b-4c6f-9bab-3560b041c358	5c00c0d3-741a-4c55-9b4f-dc44b573369d	16	REJECTED	f03a133e-617b-4c6f-9bab-3560b041c358	2026-05-04 10:33:21.270196+00	Терміново	2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
168365ea-cf7a-48cd-9b0c-eb6c1044100a	0997831c-654f-471b-934f-cedafbc54ea5	f03a133e-617b-4c6f-9bab-3560b041c358	482cf5a5-67a2-47ac-a9bd-d7beb68d1a4a	8	COMPLETED	f03a133e-617b-4c6f-9bab-3560b041c358	2026-04-29 10:33:21.270704+00	Планова поставка	2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
708e528d-9c25-4352-b0d9-a9987e5afe87	0997831c-654f-471b-934f-cedafbc54ea5	612e483e-e25c-4d11-89b7-8327b8257c4f	7bf4fafb-8ab7-4323-ac96-ec8190fe3909	19	COMPLETED	f03a133e-617b-4c6f-9bab-3560b041c358	2026-05-04 10:33:21.271574+00		2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
8bcf1b60-1e37-443a-a851-1b9b854031e3	0997831c-654f-471b-934f-cedafbc54ea5	e169450b-acde-48e4-b894-fc4538fa26a9	1d11a833-d398-4f08-afc1-48a4a805bbce	29	REJECTED	81810dfd-0d94-402c-a398-fa6d5559686e	2026-05-05 10:33:21.272449+00	Планова поставка	2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
5917ffef-ef71-4ffd-943f-1fa17bd07f96	0997831c-654f-471b-934f-cedafbc54ea5	e169450b-acde-48e4-b894-fc4538fa26a9	117bb8dd-2568-410f-855f-22d788f82954	1	APPROVED	81810dfd-0d94-402c-a398-fa6d5559686e	2026-04-29 10:33:21.27291+00	Планова поставка	2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
2ec176b4-f455-41e1-ba4f-c283859c2c30	0997831c-654f-471b-934f-cedafbc54ea5	e169450b-acde-48e4-b894-fc4538fa26a9	e4d3e550-7945-4a29-9dc3-445f52a081cb	13	REJECTED	81810dfd-0d94-402c-a398-fa6d5559686e	2026-04-27 10:33:21.273394+00		2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
154ed374-3d5b-405a-9404-41b867503f2f	0997831c-654f-471b-934f-cedafbc54ea5	612e483e-e25c-4d11-89b7-8327b8257c4f	9803ce7c-ed61-4437-9074-928f3d4d9427	22	ESCALATED	\N	\N	Автоматична ескалація SLA	2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
f7d3c5d0-45bd-4d23-88e5-e7bfa3765b65	0997831c-654f-471b-934f-cedafbc54ea5	e169450b-acde-48e4-b894-fc4538fa26a9	8199727c-6c56-43c8-a1b6-35802f5a9599	8	ESCALATED	\N	\N	Автоматична ескалація SLA	2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
aa978e04-66ca-417f-aff4-5b1a92334e11	0997831c-654f-471b-934f-cedafbc54ea5	135fe7a0-2fc0-4b17-9e22-4d1c32d9658b	74570fa5-5a2f-4fc9-8c79-9eabcbab1af9	27	ESCALATED	\N	\N	Автоматична ескалація SLA	2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
2bf1df38-c5d2-4a55-84de-d1f1a8b064c7	0997831c-654f-471b-934f-cedafbc54ea5	612e483e-e25c-4d11-89b7-8327b8257c4f	cc6933c2-f196-40d3-a294-bd9700c67de3	5	ESCALATED	\N	\N	Автоматична ескалація SLA	2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
f92a5020-2bf8-4a34-bfeb-a82d047bd09d	0997831c-654f-471b-934f-cedafbc54ea5	135fe7a0-2fc0-4b17-9e22-4d1c32d9658b	a0607946-91db-4244-bcf5-2e28ca6c157b	6	ESCALATED	\N	\N	Автоматична ескалація SLA	2026-04-28 10:33:21.041491+00	2026-04-27 10:33:21.041491+00	\N	\N	\N
96929ab1-f11f-431a-91b9-e0ead5785672	0997831c-654f-471b-934f-cedafbc54ea5	8e4524be-8fb5-46a9-abb3-9d132f4fcda9	\N	1	APPROVED	8e4524be-8fb5-46a9-abb3-9d132f4fcda9	2026-05-07 10:39:03.648945+00		2026-05-07 10:39:03.635365+00	2026-05-07 10:39:03.648945+00	7eaffe2b-144b-4537-9829-abe45241fc5a	Ноутбук Lenovo ThinkPad E15 #668	51a3c25e-42cf-4164-a5ab-91125b65a6e8
c3b96779-2081-4494-be9b-ef7cccc84bf3	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	9dd794cb-2a09-48f9-91e9-90822118a266	30	APPROVED	\N	\N	Автоматичне замовлення через Smart-модуль	2026-05-07 12:49:37.965457+00	2026-05-07 12:49:37.965457+00	\N	Ноут	\N
f713b5c9-223b-45ad-94de-1122170b287f	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-17 00:34:53.854204+00		2026-05-17 00:34:53.848844+00	2026-05-27 21:14:32.595674+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
4b6531da-0415-46b1-a90c-ab6bf6a48c9e	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	DISPATCHED	21833e68-101f-4378-a916-62c120a9f192	2026-05-17 11:44:42.359844+00		2026-05-17 11:44:42.347318+00	2026-05-27 21:21:54.286562+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
53db2ad3-fd81-4f75-974d-8a633b8a86f5	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-22 08:42:51.299402+00		2026-05-22 08:42:51.283583+00	2026-05-27 21:21:55.867446+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
a1d6a6a0-3749-4303-a799-30d09b346245	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-27 21:48:57.395839+00		2026-05-27 21:48:57.38736+00	2026-05-27 21:49:34.558178+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
6d4c35df-6eeb-4a42-a59c-fa023ba79096	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-27 21:51:24.487897+00		2026-05-27 21:51:24.481226+00	2026-05-28 07:45:38.395514+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
7a304c19-ae39-4673-802b-e7019c06ed81	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-28 07:47:02.74887+00		2026-05-28 07:47:02.737927+00	2026-05-28 08:27:33.748915+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
19031419-655e-44bf-9863-c68c70b42ad0	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-28 12:12:41.201111+00		2026-05-28 12:12:41.187893+00	2026-05-28 12:23:46.884941+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
2c639beb-e3c1-4d18-a672-651deca05a3b	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	DISPATCHED	21833e68-101f-4378-a916-62c120a9f192	2026-05-28 08:42:25.260555+00		2026-05-28 08:42:25.250172+00	2026-05-28 12:12:08.0196+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
87cd4448-e066-4c10-a0f8-1e1faa25b499	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-28 12:11:29.746218+00		2026-05-28 12:11:29.739744+00	2026-05-28 12:12:11.336955+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
92d41b0a-acda-4500-90c4-cf6d697d0dfc	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-28 08:28:57.590871+00		2026-05-28 08:28:57.577322+00	2026-05-28 12:12:12.103867+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
75f07327-df11-46d1-bbe3-08f089791c8b	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-28 12:24:02.648264+00		2026-05-28 12:24:02.631501+00	2026-05-28 12:24:15.621419+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
c3de2651-f9f3-4b48-8ed7-e6474b5bf1d3	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-28 13:07:22.587482+00		2026-05-28 13:07:22.579203+00	2026-05-28 14:31:06.923009+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
2f064ad8-9087-4fc4-8b64-73afe2c25925	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	DISPATCHED	21833e68-101f-4378-a916-62c120a9f192	2026-05-28 13:16:13.145671+00		2026-05-28 13:16:13.135457+00	2026-05-28 14:31:05.35296+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
7cf6f81a-664b-47b6-949e-35385a3d3270	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-28 13:12:00.663812+00		2026-05-28 13:12:00.651787+00	2026-05-28 14:31:07.276485+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
7a0d4a54-d7b7-4fde-8942-55a04d5577d8	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-28 13:12:35.274073+00		2026-05-28 13:12:35.262535+00	2026-05-28 14:31:08.163974+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
2fb378e7-d3be-41ba-bfae-fe84e7c985f2	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-28 14:31:36.56125+00		2026-05-28 14:31:36.550399+00	2026-05-28 14:40:29.472113+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
c60622f1-fea1-4c60-b80c-4a83ed2be749	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-28 14:32:21.078958+00		2026-05-28 14:32:21.066027+00	2026-05-28 14:40:32.00968+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
570b02c0-74b7-4c64-b541-33f15c12e0bc	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	DISPATCHED	21833e68-101f-4378-a916-62c120a9f192	2026-05-30 13:12:55.66425+00		2026-05-30 13:12:55.651105+00	2026-05-30 13:15:06.599269+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
a6853850-640e-464c-8f9e-a31544524197	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-30 15:12:33.832629+00		2026-05-30 15:12:33.820784+00	2026-05-30 15:16:51.525543+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
e26f128f-43d2-4bc0-a1d5-cf7748986488	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	LOADING	21833e68-101f-4378-a916-62c120a9f192	2026-05-30 21:30:23.733884+00		2026-05-30 21:30:23.715604+00	2026-05-30 21:30:31.999683+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
960f3ecb-6ce2-4224-88b9-abf8241825f6	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-28 20:30:38.637+00		2026-05-28 20:30:38.626584+00	2026-05-29 07:31:30.539102+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
1d2b8d5a-2259-4641-a0d3-275a1dcd3647	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	DISPATCHED	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 07:33:38.933061+00		2026-05-29 07:33:38.92109+00	2026-05-29 07:55:23.985182+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
a1484e78-3114-4773-9cd1-fd7dfdb3fe54	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-28 15:30:57.281872+00		2026-05-28 15:30:57.274879+00	2026-05-28 15:47:43.65876+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
a58fc4df-9977-4289-ac15-b2defd44a741	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-28 20:59:11.152522+00		2026-05-28 20:59:11.149428+00	2026-05-29 07:55:32.737308+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
e4bbb1a5-8e7b-46b3-be05-0c103fa57a30	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-28 20:30:28.969729+00		2026-05-28 20:30:28.958023+00	2026-05-29 07:55:36.162222+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
d83944c7-221b-4499-96ea-d22dd53b9e4e	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-28 16:16:45.105403+00		2026-05-28 16:16:45.095495+00	2026-05-29 07:55:39.657918+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
e7b7c7b0-1407-46f8-accf-c178408bbf17	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 23:24:05.934042+00		2026-05-29 23:24:05.920043+00	2026-05-30 10:17:25.310516+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
a55619ee-2d84-4572-8a92-f93f506fa5e7	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-28 16:14:24.156659+00		2026-05-28 16:14:24.141518+00	2026-05-28 16:15:33.480339+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
0b3ab0c2-0ef3-42a3-a6ab-9bd4bcc9cd2f	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-28 15:48:18.613673+00		2026-05-28 15:48:18.610694+00	2026-05-28 16:15:34.135245+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
ac220e0a-de5e-45ed-8635-8a0eaae954aa	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-28 15:46:14.619045+00		2026-05-28 15:46:14.606046+00	2026-05-28 16:15:34.745336+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
9722f33c-6e4d-45ea-a1e8-346b318c1dda	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-28 15:46:03.10305+00		2026-05-28 15:46:03.087553+00	2026-05-28 16:15:35.31537+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
4a8f8ca5-995a-4253-9536-23333720b7c4	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 07:56:50.210281+00		2026-05-29 07:56:50.200098+00	2026-05-29 09:36:54.577963+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
79f1e0c2-1e2c-4418-9c3f-c61698e3b81e	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-29 22:48:44.110713+00		2026-05-29 22:48:44.098549+00	2026-05-29 23:07:38.798591+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
a7742526-a89e-4b0a-93fc-3925d487ca3f	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-30 10:16:42.904596+00		2026-05-30 10:16:42.893673+00	2026-05-30 12:48:02.825064+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
407cf304-0c12-4fb7-beaf-b1164247e366	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-30 12:47:19.875604+00		2026-05-30 12:47:19.864617+00	2026-05-30 12:52:08.804323+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
46cb862f-6806-41be-ae82-8a003a2c53c7	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	21833e68-101f-4378-a916-62c120a9f192	\N	1	COMPLETED	21833e68-101f-4378-a916-62c120a9f192	2026-05-30 12:53:52.988795+00		2026-05-30 12:53:52.976414+00	2026-05-30 13:12:25.804356+00	f64e8882-623d-41cd-8f72-2be30b783f8d	Стретч-плівка пакувальна (рулон 500мм, 20мкм)	733c1f11-0b6b-4f34-ab03-0c16946b9156
\.


--
-- Data for Name: tenants; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.tenants (id, name, slug, subscription_tier, subscription_expires_at, owner_email, is_active, created_at, updated_at) FROM stdin;
d8729234-8e3b-41d9-83bc-9c725fe65838	Default Organization	default	ENTERPRISE	\N	\N	t	2026-04-22 23:39:00.29459+00	2026-04-22 23:39:00.29459+00
8cf55fa1-69a0-4a2c-b3a1-50143ec86428	Пійлівський Логістичний Центр	piylo	ENTERPRISE	\N	markostrutinsky@gmail.com	t	2026-04-23 00:11:12.145828+00	2026-05-06 17:01:03.46493+00
0997831c-654f-471b-934f-cedafbc54ea5	ТОВ «Карго-Логістика»	cargo-logistics	PRO	\N	owner@cargo-logistics.local	t	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
\.


--
-- Data for Name: units; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.units (id, tenant_id, parent_id, name, unit_type, subscription_tier) FROM stdin;
2	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	\N	Центральний регіон управління постачанням	REGION	BASIC
3	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	2	Київський головний розподільчий центр	BRANCH	BASIC
4	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	3	Відділ складської логістики	DEPARTMENT	BASIC
5	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	4	Зміна денних комплектувальників	TEAM	BASIC
6	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	4	Група операторів висотних штабелерів	TEAM	BASIC
7	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	3	Диспетчерський департамент	DEPARTMENT	BASIC
8	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	7	Команда моніторингу GPS	TEAM	BASIC
9	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	7	Бригада водіїв великогабаритного транспорту	TEAM	BASIC
10	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	3	Відділ зворотної логістики	DEPARTMENT	BASIC
11	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	10	Група інспекції повернутих товарів	TEAM	BASIC
12	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	\N	Західний логістичний хаб	REGION	BASIC
13	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	12	Івано-Франківський транзитний термінал	BRANCH	BASIC
14	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	13	Відділ сортування та крос-докінгу	DEPARTMENT	BASIC
15	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	14	Команда швидкого розвантаження	TEAM	BASIC
16	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	14	Бригада пакувальників	TEAM	BASIC
17	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	13	Відділ технічного забезпечення	DEPARTMENT	BASIC
18	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	17	Мобільна ремонтна бригада автопарку	TEAM	BASIC
19	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	12	Львівський митно-ліцензійний комплекс	BRANCH	BASIC
20	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	19	Відділ міжнародного транзиту	DEPARTMENT	BASIC
21	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	20	Група оформлення супровідної документації	TEAM	BASIC
22	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	20	Команда приймання імпорту	TEAM	BASIC
23	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	\N	Південний регіон	REGION	BASIC
24	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	23	Одеський портовий термінал	BRANCH	BASIC
25	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	24	Відділ контейнерних перевезень	DEPARTMENT	BASIC
26	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	25	Команда крановщиків	TEAM	BASIC
27	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	25	Група кріплення вантажів	TEAM	BASIC
28	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	24	Департамент мультимодальних перевезень	DEPARTMENT	BASIC
29	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	24	Команда координації "море-залізниця"	DEPARTMENT	BASIC
30	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	23	Миколаївський центр зберігання	BRANCH	BASIC
31	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	30	Відділ управління запасами	DEPARTMENT	BASIC
32	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	31	Команда інвентаризації	TEAM	BASIC
75	0997831c-654f-471b-934f-cedafbc54ea5	\N	Регіон «Північ»	REGION	BASIC
76	0997831c-654f-471b-934f-cedafbc54ea5	75	Філія Чернігів	BRANCH	BASIC
77	0997831c-654f-471b-934f-cedafbc54ea5	76	Відділ логістики Чернігів	DEPARTMENT	BASIC
78	0997831c-654f-471b-934f-cedafbc54ea5	77	Команда доставки	TEAM	BASIC
79	0997831c-654f-471b-934f-cedafbc54ea5	75	Філія Суми	BRANCH	BASIC
80	0997831c-654f-471b-934f-cedafbc54ea5	\N	Регіон «Південь»	REGION	BASIC
81	0997831c-654f-471b-934f-cedafbc54ea5	80	Філія Одеса	BRANCH	BASIC
82	0997831c-654f-471b-934f-cedafbc54ea5	81	Відділ митної логістики	DEPARTMENT	BASIC
83	0997831c-654f-471b-934f-cedafbc54ea5	80	Філія Миколаїв	BRANCH	BASIC
84	0997831c-654f-471b-934f-cedafbc54ea5	\N	Регіон «Центр»	REGION	BASIC
85	0997831c-654f-471b-934f-cedafbc54ea5	84	Філія Київ	BRANCH	BASIC
86	0997831c-654f-471b-934f-cedafbc54ea5	85	Відділ складської логістики	DEPARTMENT	BASIC
87	0997831c-654f-471b-934f-cedafbc54ea5	86	Команда А (нічна зміна)	TEAM	BASIC
88	0997831c-654f-471b-934f-cedafbc54ea5	86	Команда Б (денна зміна)	TEAM	BASIC
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.users (id, tenant_id, username, email, full_name, phone, password_hash, role, status, unit_id, created_at, updated_at) FROM stdin;
21833e68-101f-4378-a916-62c120a9f192	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	markostrutinsky	markostrutinsky@gmail.com	Струтинський Марко Валерійович	\N	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	TENANT_ADMIN	ACTIVE	\N	2026-04-23 00:11:12.200434+00	2026-04-23 00:11:12.200434+00
33a91fae-503d-45d2-aa42-4fdf78fcc983	d8729234-8e3b-41d9-83bc-9c725fe65838	system_admin	platform@omnilog.system	Platform Administrator	\N	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	SYSTEM_ADMIN	ACTIVE	\N	2026-04-23 00:30:15.706342+00	2026-04-23 00:30:15.706342+00
acc38b6c-89a4-4403-8cb2-916107ff017a	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	o.kovalenko	o.kovalenko@logistics.ua	Коваленко Олена Вікторівна	+380631112233	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	REGION_DIRECTOR	ACTIVE	2	2026-05-02 14:40:04.767245+00	2026-05-02 14:40:04.767245+00
78205adb-41af-4988-9cd4-6bee3337f215	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	n.boiko	n.boiko@logistics.ua	Бойко Наталія Сергіївна	+380502223344	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	BRANCH_MANAGER	ACTIVE	3	2026-05-02 14:40:59.210161+00	2026-05-02 14:40:59.210161+00
bb96be7a-c867-48c2-a128-fb143bf6409a	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	y.kravchenko	y.kravchenko@logistics.ua	Кравченко Юлія Анатоліївна	+380503334455	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	DEPT_MANAGER	ACTIVE	4	2026-05-02 14:42:35.4651+00	2026-05-02 14:42:35.4651+00
744c25e6-ad82-4471-8cd6-8da1cb476d26	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	d.tkachenko	d.tkachenko@logistics.ua	Ткаченко Дмитро Олегович	+380631112233	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	REGION_LOGISTICIAN	ACTIVE	2	2026-05-02 14:44:04.575577+00	2026-05-02 14:44:04.575577+00
44382f55-1ab8-496b-a63e-631171e55486	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	i.melnyk	i.melnyk@logistics.ua	Мельник Ігор Володимирович	+380931112233	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	REGION_STOREKEEPER	ACTIVE	2	2026-05-02 14:44:53.771363+00	2026-05-02 14:44:53.771363+00
1baed435-4416-444a-9939-724c8858ab78	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	m.hrytsenko	m.hrytsenko@logistics.ua	Гриценко Максим Іванович	+380672223344	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	BRANCH_LOGISTICIAN	ACTIVE	3	2026-05-02 14:45:58.00833+00	2026-05-02 14:45:58.00833+00
8aa95ca9-5e34-4ab6-a033-be734bf87629	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	o.lysenko	o.lysenko@logistics.ua	Лисенко Олександр Петрович	+380632223344	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	BRANCH_STOREKEEPER	ACTIVE	3	2026-05-02 14:46:42.081748+00	2026-05-02 14:46:42.081748+00
6b43bc82-4a5b-46cc-af09-6ef4adecd95b	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	s.marchenko	s.marchenko@logistics.ua	Марченко Сергій Олександрович	+380673334455	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	DEPT_SUPERVISOR	ACTIVE	4	2026-05-02 14:47:46.761666+00	2026-05-02 14:47:46.761666+00
b6635a10-c63a-4151-b9a4-ec34d35c99ee	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	v.rudenko	v.rudenko@logistics.ua	Руденко Віталій Миколайович	+380933334455	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	DEPT_MANAGER	ACTIVE	7	2026-05-02 14:48:32.987632+00	2026-05-02 14:48:32.987632+00
05eb7a5c-0413-4a6b-ad6a-5ad545e70a33	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	i.savchenko	i.savchenko@logistics.ua	Савченко Ірина Борисівна	+380504445566	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	DEPT_SUPERVISOR	ACTIVE	7	2026-05-02 14:49:14.931234+00	2026-05-02 14:49:14.931234+00
c22f350c-6f18-4839-8d2d-79fd623e6a85	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	v.pavlenko	v.pavlenko@logistics.ua	Павленко Віктор Андрійович	+380674445566	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	DEPT_MANAGER	ACTIVE	10	2026-05-02 14:50:03.977437+00	2026-05-02 14:50:03.977437+00
39f7d374-63c3-4b75-b412-1a3e03a618b2	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	a.havryliuk	a.havryliuk@logistics.ua	Гаврилюк Анна Тарасівна	+380634445566	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	DEPT_SUPERVISOR	ACTIVE	10	2026-05-02 14:50:48.925629+00	2026-05-02 14:50:48.925629+00
200577bb-7a24-4a4e-a4a8-7146fb112d37	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	t.romanenko	t.romanenko@logistics.ua	Романенко Тарас Юрійович	+380505556677	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	TEAM_LEAD	ACTIVE	5	2026-05-02 14:59:29.456797+00	2026-05-02 14:59:29.456797+00
94579c72-8ab3-4faf-8d91-8d248206ebe6	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	o.sydorenko	o.sydorenko@logistics.ua	Сидоренко Олег Віталійович	+380675556677	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	EMPLOYEE	ACTIVE	5	2026-05-02 15:05:33.137408+00	2026-05-02 15:05:33.137408+00
26de8e1e-4ea0-46f8-907a-4c7a3390904b	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	v.polishchuk	v.polishchuk@logistics.ua	Поліщук Василь Степанович	+380935556677	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	TEAM_LEAD	ACTIVE	6	2026-05-02 15:06:28.233677+00	2026-05-02 15:06:28.233677+00
0a4bc519-6dc3-4f19-9ec1-ecff82c91d03	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	i.kozak	i.kozak@logistics.ua	Козак Іван Павлович	+380506667788	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	EMPLOYEE	ACTIVE	6	2026-05-02 15:07:15.985483+00	2026-05-02 15:07:15.985483+00
b4bb3fff-a72f-437e-8d2f-48ed65d12eda	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	m.honcharenko	m.honcharenko@logistics.ua	Гончаренко Михайло Сергійович	+380676667788	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	TEAM_LEAD	ACTIVE	8	2026-05-02 15:08:23.019411+00	2026-05-02 15:08:23.019411+00
dab88a27-f886-47b0-add8-8569e62942fa	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	a.lytvyn	a.lytvyn@logistics.ua	Литвин Андрій Анатолійович	+380636667788	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	EMPLOYEE	ACTIVE	8	2026-05-02 15:09:01.668555+00	2026-05-02 15:09:01.668555+00
4817a99c-f2c5-4157-8e3c-eedc2241fbc2	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	v.moroz	v.moroz@logistics.ua	Мороз Володимир Ігорович	+380507778899	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	TEAM_LEAD	ACTIVE	9	2026-05-02 15:09:43.370225+00	2026-05-02 15:09:43.370225+00
2b080e9c-a604-4935-98c3-3bed5932986a	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	m.petrenko	m.petrenko@logistics.ua	Петренко Микола Васильович	+380677778899	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	EMPLOYEE	ACTIVE	9	2026-05-02 15:10:29.851338+00	2026-05-02 15:10:29.851338+00
3cd8d6bf-fad2-45f6-b6fb-cfce134df557	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	o.zakharchenko	o.zakharchenko@logistics.ua	Захарченко Ольга Павлівна	+380937778899	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	TEAM_LEAD	ACTIVE	11	2026-05-02 15:11:27.366143+00	2026-05-02 15:11:27.366143+00
471f0149-166b-4f29-803b-3c910fa974c1	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	m.vlasenko	m.vlasenko@logistics.ua	Власенко Марина Юріївна	+380508889900	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	EMPLOYEE	ACTIVE	11	2026-05-02 15:12:09.113938+00	2026-05-02 15:12:09.113938+00
389997eb-96b7-4e37-a6ab-2db5692a6255	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	ipetrenko	i.petrenko@example.ua	Петренко Іван Миколайович	+380501112233	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	EMPLOYEE	ACTIVE	2	2026-05-03 13:30:14.308251+00	2026-05-03 13:30:14.308251+00
b1b097e6-452e-4b2a-b62f-ff8e543a59a0	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	otkachenko	o.tkachenko@example.ua	Ткаченко Олена Василівна	+380674445566	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	EMPLOYEE	ACTIVE	3	2026-05-03 13:31:18.104658+00	2026-05-03 13:31:18.104658+00
dd8a8b66-805e-4b3f-92e5-ea158a88f421	0997831c-654f-471b-934f-cedafbc54ea5	admin	admin@cargo-logistics.local	Адміністратор Тенанта	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	TENANT_ADMIN	ACTIVE	\N	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
8e4524be-8fb5-46a9-abb3-9d132f4fcda9	0997831c-654f-471b-934f-cedafbc54ea5	owner	owner@cargo-logistics.local	Власник Карго-Логістика	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	TENANT_ADMIN	ACTIVE	\N	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
3e59277f-0c1a-4aa3-973c-cf3db1e19497	0997831c-654f-471b-934f-cedafbc54ea5	dir.north	dir.north@cargo-logistics.local	Олексій Дирєкторенко	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	REGION_DIRECTOR	ACTIVE	75	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
601e6d16-947b-4eb0-a67c-5304edafedd8	0997831c-654f-471b-934f-cedafbc54ea5	log.north	log.north@cargo-logistics.local	Ірина Логістенко	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	REGION_LOGISTICIAN	ACTIVE	75	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
6725c6a2-2570-4368-b33c-fe98f28ef4b8	0997831c-654f-471b-934f-cedafbc54ea5	store.north	store.north@cargo-logistics.local	Дмитро Комірниченко	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	REGION_STOREKEEPER	ACTIVE	75	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
135fe7a0-2fc0-4b17-9e22-4d1c32d9658b	0997831c-654f-471b-934f-cedafbc54ea5	mgr.chernihiv	mgr.chernihiv@cargo-logistics.local	Василь Менеджеренко	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	BRANCH_MANAGER	ACTIVE	76	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
58b58440-83aa-463d-85c7-016c6c912d90	0997831c-654f-471b-934f-cedafbc54ea5	log.chernihiv	log.chernihiv@cargo-logistics.local	Наталія Логіст Чернігів	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	BRANCH_LOGISTICIAN	ACTIVE	76	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
e80645b9-d2bc-4821-b137-78af6069db5c	0997831c-654f-471b-934f-cedafbc54ea5	store.ch	store.ch@cargo-logistics.local	Роман Комірник Чернігів	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	BRANCH_STOREKEEPER	ACTIVE	76	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
2b6f42a8-76d9-421a-b160-7f77cd43e82c	0997831c-654f-471b-934f-cedafbc54ea5	dept.ch	dept.ch@cargo-logistics.local	Андрій Начвідділу Чернігів	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	DEPT_MANAGER	ACTIVE	77	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
dd827d8d-ca04-4407-8fac-0c507185f8d0	0997831c-654f-471b-934f-cedafbc54ea5	sup.ch	sup.ch@cargo-logistics.local	Юлія Супервайзер Чернігів	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	DEPT_SUPERVISOR	ACTIVE	77	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
cc4504c7-ebfb-4ff0-b1b5-8967f7bd3df5	0997831c-654f-471b-934f-cedafbc54ea5	lead.ch	lead.ch@cargo-logistics.local	Сергій Тімлід Чернігів	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	TEAM_LEAD	ACTIVE	78	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
a9102433-2f88-4476-8cd6-23ff468ba599	0997831c-654f-471b-934f-cedafbc54ea5	emp1.ch	emp1.ch@cargo-logistics.local	Микола Працівник 1	+380671234567	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	EMPLOYEE	ACTIVE	78	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
454f6152-f039-4a95-bde1-a69fade0a92a	0997831c-654f-471b-934f-cedafbc54ea5	emp2.ch	emp2.ch@cargo-logistics.local	Оксана Працівниця 2	+380671234568	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	EMPLOYEE	ACTIVE	78	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
087e4e1d-c144-421d-91c8-031159dff1e5	0997831c-654f-471b-934f-cedafbc54ea5	mgr.sumy	mgr.sumy@cargo-logistics.local	Катерина Менеджер Суми	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	BRANCH_MANAGER	ACTIVE	79	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
70097b9c-a330-4283-a548-321e66624052	0997831c-654f-471b-934f-cedafbc54ea5	log.sumy	log.sumy@cargo-logistics.local	Євген Логіст Суми	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	BRANCH_LOGISTICIAN	ACTIVE	79	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
f03a133e-617b-4c6f-9bab-3560b041c358	0997831c-654f-471b-934f-cedafbc54ea5	dir.south	dir.south@cargo-logistics.local	Тетяна Директор Південь	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	REGION_DIRECTOR	ACTIVE	80	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
891becf7-f500-4375-939b-f29144c7ff65	0997831c-654f-471b-934f-cedafbc54ea5	log.south	log.south@cargo-logistics.local	Богдан Логіст Південь	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	REGION_LOGISTICIAN	ACTIVE	80	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
3c60c966-abc0-4f9e-84b0-996d86aedc07	0997831c-654f-471b-934f-cedafbc54ea5	store.south	store.south@cargo-logistics.local	Людмила Комірник Південь	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	REGION_STOREKEEPER	ACTIVE	80	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
223e9c26-98f3-4f00-8af8-dc066e2de9f4	0997831c-654f-471b-934f-cedafbc54ea5	mgr.odesa	mgr.odesa@cargo-logistics.local	Вікторія Менеджер Одеса	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	BRANCH_MANAGER	ACTIVE	81	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
612e483e-e25c-4d11-89b7-8327b8257c4f	0997831c-654f-471b-934f-cedafbc54ea5	log.odesa	log.odesa@cargo-logistics.local	Ярослав Логіст Одеса	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	BRANCH_LOGISTICIAN	ACTIVE	81	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
50ab24a9-ec74-4786-8e87-0eab3857560d	0997831c-654f-471b-934f-cedafbc54ea5	dept.od	dept.od@cargo-logistics.local	Ганна Начвідділу Митниця	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	DEPT_MANAGER	ACTIVE	82	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
60edac2d-4962-4ae1-9757-feb98bb6ed5e	0997831c-654f-471b-934f-cedafbc54ea5	mgr.mykolaiv	mgr.mykolaiv@cargo-logistics.local	Петро Менеджер Миколаїв	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	BRANCH_MANAGER	ACTIVE	83	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
81810dfd-0d94-402c-a398-fa6d5559686e	0997831c-654f-471b-934f-cedafbc54ea5	dir.central	dir.central@cargo-logistics.local	Ростислав Директор Центр	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	REGION_DIRECTOR	ACTIVE	84	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
2ba13f13-bda2-4ae4-a6aa-2379e273419f	0997831c-654f-471b-934f-cedafbc54ea5	log.central	log.central@cargo-logistics.local	Марина Логіст Центр	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	REGION_LOGISTICIAN	ACTIVE	84	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
0ad13c35-c648-4bff-ad65-89e9a87c7598	0997831c-654f-471b-934f-cedafbc54ea5	store.central	store.central@cargo-logistics.local	Іван Комірник Центр	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	REGION_STOREKEEPER	ACTIVE	84	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
e169450b-acde-48e4-b894-fc4538fa26a9	0997831c-654f-471b-934f-cedafbc54ea5	mgr.kyiv	mgr.kyiv@cargo-logistics.local	Олена Менеджер Київ	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	BRANCH_MANAGER	ACTIVE	85	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
6505e3f4-cd82-4eca-8544-9fc103a4bf50	0997831c-654f-471b-934f-cedafbc54ea5	log.kyiv	log.kyiv@cargo-logistics.local	Павло Логіст Київ	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	BRANCH_LOGISTICIAN	ACTIVE	85	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
76eff538-8750-4abb-95bc-7164f9ca90a4	0997831c-654f-471b-934f-cedafbc54ea5	store.kyiv	store.kyiv@cargo-logistics.local	Аліна Комірник Київ	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	BRANCH_STOREKEEPER	ACTIVE	85	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
03ced0df-4ca8-4e03-9115-a390abdf9e15	0997831c-654f-471b-934f-cedafbc54ea5	dept.kyiv	dept.kyiv@cargo-logistics.local	Максим Начвідділу Київ	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	DEPT_MANAGER	ACTIVE	86	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
09762e7f-3b5c-4e10-bacc-05044d2f4a20	0997831c-654f-471b-934f-cedafbc54ea5	sup.kyiv	sup.kyiv@cargo-logistics.local	Світлана Супервайзер Київ	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	DEPT_SUPERVISOR	ACTIVE	86	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
9d623dfc-ed11-4d97-8743-ba56fe888e51	0997831c-654f-471b-934f-cedafbc54ea5	lead.kyiv.a	lead.kyiv.a@cargo-logistics.local	Тарас Тімлід A	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	TEAM_LEAD	ACTIVE	87	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
cbaaea35-191f-4098-92f3-66e2e15c240f	0997831c-654f-471b-934f-cedafbc54ea5	lead.kyiv.b	lead.kyiv.b@cargo-logistics.local	Лариса Тімлід Б	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	TEAM_LEAD	ACTIVE	88	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
d2d32b8c-46ae-432e-80eb-281d903e6ca4	0997831c-654f-471b-934f-cedafbc54ea5	emp3.kyiv	emp3.kyiv@cargo-logistics.local	Денис Працівник 3	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	EMPLOYEE	ACTIVE	87	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
a00d600c-8a93-477b-87d3-e29170ae267e	0997831c-654f-471b-934f-cedafbc54ea5	emp4.kyiv	emp4.kyiv@cargo-logistics.local	Вероніка Працівниця 4	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	EMPLOYEE	ACTIVE	88	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
33654277-e2a8-43df-b04a-1c4cc384c000	0997831c-654f-471b-934f-cedafbc54ea5	blocked	blocked@cargo-logistics.local	Заблокований Тест	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	EMPLOYEE	BLOCKED	78	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
cd2e4bc4-c122-4c30-9f7b-2cdd0d7fd4b3	0997831c-654f-471b-934f-cedafbc54ea5	pending	pending@cargo-logistics.local	Новий Тест (PENDING)	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	EMPLOYEE	PENDING	86	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
a2fbe6df-f0aa-432e-a5df-6759d8b5dff6	\N	contractor1	contractor1@cargo-logistics.local	Волонтер Антон	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	CONTRACTOR	ACTIVE	\N	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
e4e50963-10f3-46c5-86cd-988e5799c3fe	\N	contractor2	contractor2@cargo-logistics.local	Волонтер Марія	\N	$2b$12$NP9gW59DE.aHpdmD7l3n/uIp/c8T36VnabGcGQQD6quTZLyNBgIkW	CONTRACTOR	ACTIVE	\N	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00
c301c093-ebc6-491a-9c40-eb425e472857	\N	markostrutinsky1	markostrutinsky1@gmail.com	Петренко Іван Дмитрович	\N	$2a$10$47I7pZajwNFnU43tKru6E.TTIqJpe1DyJXj6GKcmVV/gpfIVMaVDO	CONTRACTOR	ACTIVE	\N	2026-06-08 10:18:10.987536+00	2026-06-08 10:18:10.987536+00
\.


--
-- Data for Name: vehicle_driver_history; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.vehicle_driver_history (id, vehicle_id, driver_id, assigned_at) FROM stdin;
2f4e7b7d-755e-46ae-ba03-23dc3c92bb6c	ea092cad-3147-4245-8ce6-ee7b8029806b	0a4bc519-6dc3-4f19-9ec1-ecff82c91d03	2026-05-03 13:23:41.067183+00
f7bcacc2-2afe-41e7-b9a9-987954584f3b	48471a6b-ac6a-4327-b13c-dce6768ac4a0	389997eb-96b7-4e37-a6ab-2db5692a6255	2026-05-03 13:41:00.649047+00
9e1203f0-70fb-4660-a66e-7b80b770a8dd	48471a6b-ac6a-4327-b13c-dce6768ac4a0	\N	2026-05-03 15:06:30.741244+00
a5ddc117-3073-4a9f-ba6e-eb04dd95158d	48471a6b-ac6a-4327-b13c-dce6768ac4a0	389997eb-96b7-4e37-a6ab-2db5692a6255	2026-05-03 15:07:31.055869+00
7d5ab450-ae67-448b-864f-e077ebf96994	dac43200-3bfe-47c8-8e60-17bf1ce6cf40	b1b097e6-452e-4b2a-b62f-ff8e543a59a0	2026-05-28 12:21:28.423531+00
cf8fcb73-058e-4667-abf2-1f99ad00f20c	48471a6b-ac6a-4327-b13c-dce6768ac4a0	\N	2026-05-29 12:47:31.395474+00
\.


--
-- Data for Name: vehicles; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.vehicles (id, tenant_id, brand, model, plate_number, status, driver_id, created_at, updated_at, tank_capacity, fuel_norm, maintenance_interval_km, last_maintenance_odometer, status_reason, type, capacity_kg, current_warehouse_id, home_warehouse_id) FROM stdin;
5482a426-5bcf-48a2-82cd-d61441e76810	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	Isuzu	NQR 90	AI4455EE	ACTIVE	2b080e9c-a604-4935-98c3-3bed5932986a	2026-05-02 15:33:37.935454+00	2026-05-02 15:33:37.935454+00	140.00	18.00	10000	0	\N	TRUCK	5000.00	\N	\N
1565aa02-d4ff-464c-bc5b-c96f63fd9be1	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	Ford	Mondeo	AT1234BH	ON_MISSION	94579c72-8ab3-4faf-8d91-8d248206ebe6	2026-05-28 07:44:53.438673+00	2026-05-30 13:13:26.000459+00	70.00	8.00	10000	0	\N	VAN	1500.00	f64e8882-623d-41cd-8f72-2be30b783f8d	f64e8882-623d-41cd-8f72-2be30b783f8d
8355c89e-b072-4113-a6d2-adecc94a5a31	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	MAN	TGX 18.440	KA1234AA	ACTIVE	4817a99c-f2c5-4157-8e3c-eedc2241fbc2	2026-05-02 15:32:39.845265+00	2026-05-29 09:36:48.874212+00	800.00	32.50	10000	0	\N	TRUCK	22000.00	f64e8882-623d-41cd-8f72-2be30b783f8d	f64e8882-623d-41cd-8f72-2be30b783f8d
dac43200-3bfe-47c8-8e60-17bf1ce6cf40	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	Mercedes-Benz	Sprinter 316 CDI	BC5678CX	ON_MISSION	b1b097e6-452e-4b2a-b62f-ff8e543a59a0	2026-05-02 15:34:23.067636+00	2026-05-30 21:30:31.999683+00	75.00	11.20	10000	0	\N	VAN	2500.00	f64e8882-623d-41cd-8f72-2be30b783f8d	f64e8882-623d-41cd-8f72-2be30b783f8d
48471a6b-ac6a-4327-b13c-dce6768ac4a0	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	Renault	Master	AB3322OO	ACTIVE	\N	2026-05-02 15:34:56.400134+00	2026-05-29 09:36:57.56492+00	80.00	10.50	10000	0	\N	VAN	1500.00	f64e8882-623d-41cd-8f72-2be30b783f8d	f64e8882-623d-41cd-8f72-2be30b783f8d
ea092cad-3147-4245-8ce6-ee7b8029806b	8cf55fa1-69a0-4a2c-b3a1-50143ec86428	Mitsubishi	L200	CE9012MP	ACTIVE	0a4bc519-6dc3-4f19-9ec1-ecff82c91d03	2026-05-02 15:35:32.632995+00	2026-05-31 22:53:35.193414+00	75.00	9.50	10000	0	\N	PICKUP	1000.00	f64e8882-623d-41cd-8f72-2be30b783f8d	f64e8882-623d-41cd-8f72-2be30b783f8d
1fe8f920-1f68-4f98-ac28-06256ac21c55	0997831c-654f-471b-934f-cedafbc54ea5	Renault	Master 2.3 dCi (2023)	AP8906AK	ACTIVE	454f6152-f039-4a95-bde1-a69fade0a92a	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0.00	0.00	10000	0	\N	VAN	1500.00	\N	\N
1d9b95fd-005c-4ad4-ba63-4d468ce9e913	0997831c-654f-471b-934f-cedafbc54ea5	Ford	Transit 2.0 TDCi (2024)	OP5112AM	ACTIVE	a9102433-2f88-4476-8cd6-23ff468ba599	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0.00	0.00	10000	0	\N	VAN	1500.00	\N	\N
f678ceb7-8984-4534-abac-41d262aaea81	0997831c-654f-471b-934f-cedafbc54ea5	Mercedes	Sprinter 319 CDI (2022)	EO9000TC	ACTIVE	cc4504c7-ebfb-4ff0-b1b5-8967f7bd3df5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0.00	0.00	10000	0	\N	VAN	1500.00	\N	\N
6a4d89a0-f357-4ccf-bfaf-8f0791f22192	0997831c-654f-471b-934f-cedafbc54ea5	Toyota	Hilux 2.8 D-4D (2023)	EE7382CT	ACTIVE	a9102433-2f88-4476-8cd6-23ff468ba599	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0.00	0.00	10000	0	\N	VAN	1500.00	\N	\N
9e479722-a3d7-4cc7-ad11-0f26f67be07b	0997831c-654f-471b-934f-cedafbc54ea5	Volkswagen	Crafter 35 (2020)	KT4438CM	ACTIVE	cc4504c7-ebfb-4ff0-b1b5-8967f7bd3df5	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0.00	0.00	10000	0	\N	VAN	1500.00	\N	\N
4775a9bb-25f8-49ef-a65b-1002543d468b	0997831c-654f-471b-934f-cedafbc54ea5	MAN	TGE 5.180 (2023)	MA7721PP	ACTIVE	454f6152-f039-4a95-bde1-a69fade0a92a	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0.00	0.00	10000	0	\N	VAN	1500.00	\N	\N
75da1df1-ee2f-41b8-90da-385bfaa4d080	0997831c-654f-471b-934f-cedafbc54ea5	Renault	Master 2.3 dCi (2021)	EM9676PB	ACTIVE	70097b9c-a330-4283-a548-321e66624052	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0.00	0.00	10000	0	\N	VAN	1500.00	\N	\N
3863dff1-509f-4519-b6fd-eb60da807648	0997831c-654f-471b-934f-cedafbc54ea5	Ford	Transit 2.0 TDCi (2019)	MB1989CH	ACTIVE	612e483e-e25c-4d11-89b7-8327b8257c4f	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0.00	0.00	10000	0	\N	VAN	1500.00	\N	\N
9e508594-5002-4e3e-bb16-32d710685d17	0997831c-654f-471b-934f-cedafbc54ea5	Mercedes	Sprinter 319 CDI (2023)	EX5419ET	ACTIVE	612e483e-e25c-4d11-89b7-8327b8257c4f	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0.00	0.00	10000	0	\N	VAN	1500.00	\N	\N
23fbee13-4b35-48cb-81fe-d0c496119765	0997831c-654f-471b-934f-cedafbc54ea5	Toyota	Hilux 2.8 D-4D (2024)	PM2890KT	ACTIVE	70097b9c-a330-4283-a548-321e66624052	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0.00	0.00	10000	0	\N	VAN	1500.00	\N	\N
7c691151-27c5-420c-80e8-cf4bedb56aa9	0997831c-654f-471b-934f-cedafbc54ea5	Volkswagen	Crafter 35 (2020)	OE0677TP	ACTIVE	612e483e-e25c-4d11-89b7-8327b8257c4f	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0.00	0.00	10000	0	\N	VAN	1500.00	\N	\N
131f4e80-eafd-458f-a070-9c51bb8605b7	0997831c-654f-471b-934f-cedafbc54ea5	Renault	Master 2.3 dCi (2020)	XM5108EP	ACTIVE	a00d600c-8a93-477b-87d3-e29170ae267e	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0.00	0.00	10000	0	\N	VAN	1500.00	\N	\N
ad72c4ea-32b5-41a4-bffa-1d9b31b545d8	0997831c-654f-471b-934f-cedafbc54ea5	Ford	Transit 2.0 TDCi (2020)	CE5914EK	ACTIVE	d2d32b8c-46ae-432e-80eb-281d903e6ca4	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0.00	0.00	10000	0	\N	VAN	1500.00	\N	\N
d5605709-91a8-4772-92c4-d4bd92b65717	0997831c-654f-471b-934f-cedafbc54ea5	Mercedes	Sprinter 319 CDI (2021)	PK6393CO	ACTIVE	d2d32b8c-46ae-432e-80eb-281d903e6ca4	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0.00	0.00	10000	0	\N	VAN	1500.00	\N	\N
eff2cb75-547b-42eb-89f5-233a879ac549	0997831c-654f-471b-934f-cedafbc54ea5	Toyota	Hilux 2.8 D-4D (2021)	TX3135EE	ACTIVE	a00d600c-8a93-477b-87d3-e29170ae267e	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0.00	0.00	10000	0	\N	VAN	1500.00	\N	\N
d3c1b403-57b5-4f49-8828-7b31b580ddd0	0997831c-654f-471b-934f-cedafbc54ea5	Volkswagen	Crafter 35 (2019)	HE4320ME	ACTIVE	d2d32b8c-46ae-432e-80eb-281d903e6ca4	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0.00	0.00	10000	0	\N	VAN	1500.00	\N	\N
3a1de34e-544b-4b67-aa92-20e6cfb09e72	0997831c-654f-471b-934f-cedafbc54ea5	MAN	TGE 5.180 (2019)	KO5837OE	ACTIVE	9d623dfc-ed11-4d97-8743-ba56fe888e51	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0.00	0.00	10000	0	\N	VAN	1500.00	\N	\N
e51f66b4-9055-42e2-8790-5063fe231566	0997831c-654f-471b-934f-cedafbc54ea5	Iveco	Daily 70C18 (2024)	KP6230BX	ACTIVE	9d623dfc-ed11-4d97-8743-ba56fe888e51	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0.00	0.00	10000	0	\N	VAN	1500.00	\N	\N
\.


--
-- Data for Name: warehouses; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.warehouses (id, unit_id, name, location_type, latitude, longitude, created_at, updated_at, tenant_id) FROM stdin;
c3d3f8df-0be9-4931-91f4-a76554bfb5f9	2	Головний транзитний хаб "Центр"	STATIONARY	50.4501	30.5234	2026-05-02 15:20:16.909122+00	2026-05-02 15:20:16.909122+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
f64e8882-623d-41cd-8f72-2be30b783f8d	3	Розподільчий центр "Київ-Північ"	STATIONARY	50.48765	30.48552	2026-05-02 15:20:46.471073+00	2026-05-02 15:20:46.471073+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
28f6a833-60fc-48bd-bb5c-30dfca2e3ace	4	Склад палетного зберігання №1	STATIONARY	50.48801	30.48605	2026-05-02 15:21:11.836824+00	2026-05-02 15:21:11.836824+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
b0981b66-2238-4044-ad6a-eef5917728f0	7	Майданчик крос-докінгу (Зона А)	STATIONARY	50.48912	30.487	2026-05-02 15:21:43.46294+00	2026-05-02 15:21:43.46294+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
47bbf66e-ebd4-4c59-ad77-87314f7a302d	10	Склад браку та повернень	STATIONARY	50.485	30.481	2026-05-02 15:22:18.5305+00	2026-05-02 15:22:18.5305+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
30efffba-f59c-4990-b1f6-e371a2dffca7	5	Зона дрібноштучної комплектації	STATIONARY	50.48825	30.4862	2026-05-02 15:23:07.880178+00	2026-05-02 15:23:07.880178+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
29918579-9237-4f65-9ce3-db58b825b86c	6	Зона висотного зберігання (Ряди 10-25)	STATIONARY	50.4885	30.48655	2026-05-02 15:23:35.525554+00	2026-05-02 15:23:35.525554+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
7a1639fb-d584-409a-89d9-e2b09fa379cd	8	Склад навігаційного обладнання та трекерів	STATIONARY	50.45105	30.525	2026-05-02 15:24:07.658419+00	2026-05-02 15:24:07.658419+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
b7de7c00-1ff9-4234-a909-40c7106ace7a	9	Автопарк ВГТ (Мобільний хаб)	MOBILE	50.491	30.49	2026-05-02 15:24:41.934823+00	2026-05-02 15:24:41.934823+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
62ba10c5-14ee-4628-bfcc-cfb2f99d4e46	11	Ізолятор карантинної продукції	STATIONARY	50.48525	30.4815	2026-05-02 15:25:05.612261+00	2026-05-02 15:25:05.612261+00	8cf55fa1-69a0-4a2c-b3a1-50143ec86428
952b60eb-3d73-478b-9146-c3e74a718ded	76	Головний склад Чернігів	STATIONARY	51.4982	31.2893	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
f9e374f6-00aa-4ee1-9f4e-f9e954991f3e	79	Склад Суми	STATIONARY	50.9077	34.7981	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
529304e6-6869-44a8-ba72-4fb470ba14d2	81	Склад Одеса-Порт	STATIONARY	46.4774	30.7326	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
92b1cec4-0a85-4e51-a087-3fd7543f8bbb	81	Склад Одеса-Центр	STATIONARY	46.4825	30.7233	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
fbde5ef3-81ea-4084-b588-c5940358d9b5	83	Склад Миколаїв	STATIONARY	46.975	31.9946	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
7eaffe2b-144b-4537-9829-abe45241fc5a	85	Склад Київ-Північ	STATIONARY	50.526	30.51	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
c20b9d17-af0d-4f5c-8d14-a6cf9e782115	85	Склад Київ-Південь	STATIONARY	50.393	30.545	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
0c0df47a-59d3-422a-bb25-458d50648b6b	85	Мобільний хаб Київ	MOBILE	50.4501	30.5234	2026-05-07 10:33:21.041491+00	2026-05-07 10:33:21.041491+00	0997831c-654f-471b-934f-cedafbc54ea5
\.


--
-- Name: audit_logs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.audit_logs_id_seq', 509, true);


--
-- Name: geofence_alerts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.geofence_alerts_id_seq', 1, false);


--
-- Name: geofences_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.geofences_id_seq', 16, true);


--
-- Name: gps_locations_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.gps_locations_id_seq', 1099, true);


--
-- Name: units_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.units_id_seq', 88, true);


--
-- Name: audit_logs audit_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.audit_logs
    ADD CONSTRAINT audit_logs_pkey PRIMARY KEY (id);


--
-- Name: contractor_memberships contractor_memberships_contractor_id_tenant_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.contractor_memberships
    ADD CONSTRAINT contractor_memberships_contractor_id_tenant_id_key UNIQUE (contractor_id, tenant_id);


--
-- Name: contractor_memberships contractor_memberships_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.contractor_memberships
    ADD CONSTRAINT contractor_memberships_pkey PRIMARY KEY (id);


--
-- Name: contractor_requests contractor_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.contractor_requests
    ADD CONSTRAINT contractor_requests_pkey PRIMARY KEY (id);


--
-- Name: fuel_records fuel_records_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.fuel_records
    ADD CONSTRAINT fuel_records_pkey PRIMARY KEY (id);


--
-- Name: geofence_alerts geofence_alerts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.geofence_alerts
    ADD CONSTRAINT geofence_alerts_pkey PRIMARY KEY (id);


--
-- Name: geofences geofences_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.geofences
    ADD CONSTRAINT geofences_pkey PRIMARY KEY (id);


--
-- Name: gps_locations gps_locations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.gps_locations
    ADD CONSTRAINT gps_locations_pkey PRIMARY KEY (id);


--
-- Name: inventory_check_items inventory_check_items_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inventory_check_items
    ADD CONSTRAINT inventory_check_items_pkey PRIMARY KEY (id);


--
-- Name: inventory_checks inventory_checks_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inventory_checks
    ADD CONSTRAINT inventory_checks_pkey PRIMARY KEY (id);


--
-- Name: invite_tokens invite_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.invite_tokens
    ADD CONSTRAINT invite_tokens_pkey PRIMARY KEY (id);


--
-- Name: invite_tokens invite_tokens_user_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.invite_tokens
    ADD CONSTRAINT invite_tokens_user_id_key UNIQUE (user_id);


--
-- Name: maintenance_records maintenance_records_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.maintenance_records
    ADD CONSTRAINT maintenance_records_pkey PRIMARY KEY (id);


--
-- Name: notifications notifications_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_pkey PRIMARY KEY (id);


--
-- Name: refresh_tokens refresh_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id);


--
-- Name: resource_assignments resource_assignments_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_assignments
    ADD CONSTRAINT resource_assignments_pkey PRIMARY KEY (id);


--
-- Name: resource_categories resource_categories_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_categories
    ADD CONSTRAINT resource_categories_pkey PRIMARY KEY (id);


--
-- Name: resource_categories resource_categories_tenant_id_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_categories
    ADD CONSTRAINT resource_categories_tenant_id_name_key UNIQUE (tenant_id, name);


--
-- Name: resource_categories resource_categories_tenant_name_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_categories
    ADD CONSTRAINT resource_categories_tenant_name_unique UNIQUE (tenant_id, name);


--
-- Name: resources resources_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resources
    ADD CONSTRAINT resources_pkey PRIMARY KEY (id);


--
-- Name: shipment_items shipment_items_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipment_items
    ADD CONSTRAINT shipment_items_pkey PRIMARY KEY (id);


--
-- Name: shipment_refuels shipment_refuels_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipment_refuels
    ADD CONSTRAINT shipment_refuels_pkey PRIMARY KEY (id);


--
-- Name: shipments shipments_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipments
    ADD CONSTRAINT shipments_pkey PRIMARY KEY (id);


--
-- Name: supply_requests supply_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supply_requests
    ADD CONSTRAINT supply_requests_pkey PRIMARY KEY (id);


--
-- Name: tenants tenants_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tenants
    ADD CONSTRAINT tenants_name_key UNIQUE (name);


--
-- Name: tenants tenants_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tenants
    ADD CONSTRAINT tenants_pkey PRIMARY KEY (id);


--
-- Name: tenants tenants_slug_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tenants
    ADD CONSTRAINT tenants_slug_key UNIQUE (slug);


--
-- Name: units units_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.units
    ADD CONSTRAINT units_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_tenant_email_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_tenant_email_unique UNIQUE (tenant_id, email);


--
-- Name: users users_tenant_id_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_tenant_id_email_key UNIQUE (tenant_id, email);


--
-- Name: users users_tenant_id_username_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_tenant_id_username_key UNIQUE (tenant_id, username);


--
-- Name: users users_tenant_username_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_tenant_username_unique UNIQUE (tenant_id, username);


--
-- Name: vehicle_driver_history vehicle_driver_history_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vehicle_driver_history
    ADD CONSTRAINT vehicle_driver_history_pkey PRIMARY KEY (id);


--
-- Name: vehicles vehicles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vehicles
    ADD CONSTRAINT vehicles_pkey PRIMARY KEY (id);


--
-- Name: vehicles vehicles_tenant_id_plate_number_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vehicles
    ADD CONSTRAINT vehicles_tenant_id_plate_number_key UNIQUE (tenant_id, plate_number);


--
-- Name: vehicles vehicles_tenant_plate_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vehicles
    ADD CONSTRAINT vehicles_tenant_plate_unique UNIQUE (tenant_id, plate_number);


--
-- Name: warehouses warehouses_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouses
    ADD CONSTRAINT warehouses_pkey PRIMARY KEY (id);


--
-- Name: idx_audit_logs_tenant; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_audit_logs_tenant ON public.audit_logs USING btree (tenant_id);


--
-- Name: idx_contractor_memberships_contractor; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_contractor_memberships_contractor ON public.contractor_memberships USING btree (contractor_id);


--
-- Name: idx_contractor_memberships_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_contractor_memberships_status ON public.contractor_memberships USING btree (status);


--
-- Name: idx_contractor_memberships_tenant; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_contractor_memberships_tenant ON public.contractor_memberships USING btree (tenant_id);


--
-- Name: idx_contractor_requests_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_contractor_requests_status ON public.contractor_requests USING btree (status);


--
-- Name: idx_contractor_requests_tenant; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_contractor_requests_tenant ON public.contractor_requests USING btree (tenant_id);


--
-- Name: idx_contractor_requests_warehouse; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_contractor_requests_warehouse ON public.contractor_requests USING btree (target_warehouse_id);


--
-- Name: idx_fuel_records_tenant; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_fuel_records_tenant ON public.fuel_records USING btree (tenant_id);


--
-- Name: idx_fuel_records_vehicle; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_fuel_records_vehicle ON public.fuel_records USING btree (vehicle_id);


--
-- Name: idx_geofence_alerts_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_geofence_alerts_created ON public.geofence_alerts USING btree (created_at DESC);


--
-- Name: idx_geofence_alerts_geofence; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_geofence_alerts_geofence ON public.geofence_alerts USING btree (geofence_id);


--
-- Name: idx_geofence_alerts_vehicle; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_geofence_alerts_vehicle ON public.geofence_alerts USING btree (vehicle_id);


--
-- Name: idx_geofences_tenant; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_geofences_tenant ON public.geofences USING btree (tenant_id);


--
-- Name: idx_geofences_unit; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_geofences_unit ON public.geofences USING btree (unit_id);


--
-- Name: idx_gps_locations_tenant; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_gps_locations_tenant ON public.gps_locations USING btree (tenant_id);


--
-- Name: idx_gps_locations_timestamp; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_gps_locations_timestamp ON public.gps_locations USING btree ("timestamp" DESC);


--
-- Name: idx_gps_locations_unit; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_gps_locations_unit ON public.gps_locations USING btree (unit_id);


--
-- Name: idx_gps_locations_vehicle; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_gps_locations_vehicle ON public.gps_locations USING btree (vehicle_id);


--
-- Name: idx_gps_locations_vehicle_time; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_gps_locations_vehicle_time ON public.gps_locations USING btree (vehicle_id, "timestamp" DESC);


--
-- Name: idx_invite_tokens_hash; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_invite_tokens_hash ON public.invite_tokens USING btree (token_hash);


--
-- Name: idx_maintenance_records_vehicle_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_maintenance_records_vehicle_id ON public.maintenance_records USING btree (vehicle_id);


--
-- Name: idx_notifications_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_notifications_created_at ON public.notifications USING btree (created_at DESC);


--
-- Name: idx_notifications_is_read; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_notifications_is_read ON public.notifications USING btree (is_read);


--
-- Name: idx_notifications_tenant_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_notifications_tenant_id ON public.notifications USING btree (tenant_id);


--
-- Name: idx_notifications_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_notifications_user_id ON public.notifications USING btree (user_id);


--
-- Name: idx_refresh_tokens_hash; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_refresh_tokens_hash ON public.refresh_tokens USING btree (token_hash);


--
-- Name: idx_refresh_tokens_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_refresh_tokens_user ON public.refresh_tokens USING btree (user_id);


--
-- Name: idx_resource_assignments_resource; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_resource_assignments_resource ON public.resource_assignments USING btree (resource_id);


--
-- Name: idx_resource_assignments_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_resource_assignments_user ON public.resource_assignments USING btree (user_id);


--
-- Name: idx_resource_categories_tenant; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_resource_categories_tenant ON public.resource_categories USING btree (tenant_id);


--
-- Name: idx_resources_category; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_resources_category ON public.resources USING btree (category_id);


--
-- Name: idx_resources_tenant; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_resources_tenant ON public.resources USING btree (tenant_id);


--
-- Name: idx_resources_unit; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_resources_unit ON public.resources USING btree (unit_id);


--
-- Name: idx_resources_warehouse_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_resources_warehouse_id ON public.resources USING btree (warehouse_id);


--
-- Name: idx_shipment_items_shipment; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_shipment_items_shipment ON public.shipment_items USING btree (shipment_id);


--
-- Name: idx_shipment_refuels_shipment; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_shipment_refuels_shipment ON public.shipment_refuels USING btree (shipment_id);


--
-- Name: idx_shipment_refuels_vehicle; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_shipment_refuels_vehicle ON public.shipment_refuels USING btree (vehicle_id);


--
-- Name: idx_shipments_direction; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_shipments_direction ON public.shipments USING btree (direction);


--
-- Name: idx_shipments_from; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_shipments_from ON public.shipments USING btree (from_warehouse_id);


--
-- Name: idx_shipments_tenant; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_shipments_tenant ON public.shipments USING btree (tenant_id);


--
-- Name: idx_shipments_to; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_shipments_to ON public.shipments USING btree (to_warehouse_id);


--
-- Name: idx_supply_requests_created_by; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_supply_requests_created_by ON public.supply_requests USING btree (created_by);


--
-- Name: idx_supply_requests_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_supply_requests_status ON public.supply_requests USING btree (status);


--
-- Name: idx_supply_requests_tenant; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_supply_requests_tenant ON public.supply_requests USING btree (tenant_id);


--
-- Name: idx_tenants_slug; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_tenants_slug ON public.tenants USING btree (slug);


--
-- Name: idx_units_tenant; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_units_tenant ON public.units USING btree (tenant_id);


--
-- Name: idx_users_role; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_role ON public.users USING btree (role);


--
-- Name: idx_users_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_status ON public.users USING btree (status);


--
-- Name: idx_users_tenant; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_tenant ON public.users USING btree (tenant_id);


--
-- Name: idx_vehicle_driver_history_assigned_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_vehicle_driver_history_assigned_at ON public.vehicle_driver_history USING btree (assigned_at);


--
-- Name: idx_vehicle_driver_history_driver_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_vehicle_driver_history_driver_id ON public.vehicle_driver_history USING btree (driver_id);


--
-- Name: idx_vehicle_driver_history_vehicle_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_vehicle_driver_history_vehicle_id ON public.vehicle_driver_history USING btree (vehicle_id);


--
-- Name: idx_vehicles_home_warehouse; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_vehicles_home_warehouse ON public.vehicles USING btree (home_warehouse_id);


--
-- Name: idx_vehicles_tenant; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_vehicles_tenant ON public.vehicles USING btree (tenant_id);


--
-- Name: idx_vehicles_warehouse; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_vehicles_warehouse ON public.vehicles USING btree (current_warehouse_id);


--
-- Name: idx_warehouses_tenant; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_warehouses_tenant ON public.warehouses USING btree (tenant_id);


--
-- Name: audit_logs audit_logs_tenant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.audit_logs
    ADD CONSTRAINT audit_logs_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES public.tenants(id) ON DELETE SET NULL;


--
-- Name: audit_logs audit_logs_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.audit_logs
    ADD CONSTRAINT audit_logs_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: contractor_memberships contractor_memberships_contractor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.contractor_memberships
    ADD CONSTRAINT contractor_memberships_contractor_id_fkey FOREIGN KEY (contractor_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: contractor_memberships contractor_memberships_decided_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.contractor_memberships
    ADD CONSTRAINT contractor_memberships_decided_by_fkey FOREIGN KEY (decided_by) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: contractor_memberships contractor_memberships_tenant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.contractor_memberships
    ADD CONSTRAINT contractor_memberships_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES public.tenants(id) ON DELETE CASCADE;


--
-- Name: contractor_requests contractor_requests_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.contractor_requests
    ADD CONSTRAINT contractor_requests_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: contractor_requests contractor_requests_taken_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.contractor_requests
    ADD CONSTRAINT contractor_requests_taken_by_fkey FOREIGN KEY (taken_by) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: contractor_requests contractor_requests_target_warehouse_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.contractor_requests
    ADD CONSTRAINT contractor_requests_target_warehouse_id_fkey FOREIGN KEY (target_warehouse_id) REFERENCES public.warehouses(id) ON DELETE SET NULL;


--
-- Name: contractor_requests contractor_requests_tenant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.contractor_requests
    ADD CONSTRAINT contractor_requests_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES public.tenants(id) ON DELETE CASCADE;


--
-- Name: contractor_requests contractor_requests_unit_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.contractor_requests
    ADD CONSTRAINT contractor_requests_unit_id_fkey FOREIGN KEY (unit_id) REFERENCES public.units(id) ON DELETE SET NULL;


--
-- Name: fuel_records fuel_records_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.fuel_records
    ADD CONSTRAINT fuel_records_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: fuel_records fuel_records_tenant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.fuel_records
    ADD CONSTRAINT fuel_records_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES public.tenants(id) ON DELETE CASCADE;


--
-- Name: fuel_records fuel_records_vehicle_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.fuel_records
    ADD CONSTRAINT fuel_records_vehicle_id_fkey FOREIGN KEY (vehicle_id) REFERENCES public.vehicles(id) ON DELETE CASCADE;


--
-- Name: geofence_alerts geofence_alerts_geofence_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.geofence_alerts
    ADD CONSTRAINT geofence_alerts_geofence_id_fkey FOREIGN KEY (geofence_id) REFERENCES public.geofences(id) ON DELETE CASCADE;


--
-- Name: geofence_alerts geofence_alerts_vehicle_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.geofence_alerts
    ADD CONSTRAINT geofence_alerts_vehicle_id_fkey FOREIGN KEY (vehicle_id) REFERENCES public.vehicles(id) ON DELETE CASCADE;


--
-- Name: geofences geofences_tenant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.geofences
    ADD CONSTRAINT geofences_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES public.tenants(id) ON DELETE CASCADE;


--
-- Name: geofences geofences_unit_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.geofences
    ADD CONSTRAINT geofences_unit_id_fkey FOREIGN KEY (unit_id) REFERENCES public.units(id) ON DELETE CASCADE;


--
-- Name: gps_locations gps_locations_tenant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.gps_locations
    ADD CONSTRAINT gps_locations_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES public.tenants(id) ON DELETE CASCADE;


--
-- Name: gps_locations gps_locations_unit_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.gps_locations
    ADD CONSTRAINT gps_locations_unit_id_fkey FOREIGN KEY (unit_id) REFERENCES public.units(id) ON DELETE CASCADE;


--
-- Name: gps_locations gps_locations_vehicle_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.gps_locations
    ADD CONSTRAINT gps_locations_vehicle_id_fkey FOREIGN KEY (vehicle_id) REFERENCES public.vehicles(id) ON DELETE CASCADE;


--
-- Name: inventory_check_items inventory_check_items_check_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inventory_check_items
    ADD CONSTRAINT inventory_check_items_check_id_fkey FOREIGN KEY (check_id) REFERENCES public.inventory_checks(id) ON DELETE CASCADE;


--
-- Name: inventory_check_items inventory_check_items_resource_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inventory_check_items
    ADD CONSTRAINT inventory_check_items_resource_id_fkey FOREIGN KEY (resource_id) REFERENCES public.resources(id) ON DELETE CASCADE;


--
-- Name: inventory_checks inventory_checks_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inventory_checks
    ADD CONSTRAINT inventory_checks_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: inventory_checks inventory_checks_warehouse_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.inventory_checks
    ADD CONSTRAINT inventory_checks_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES public.warehouses(id) ON DELETE CASCADE;


--
-- Name: invite_tokens invite_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.invite_tokens
    ADD CONSTRAINT invite_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: maintenance_records maintenance_records_driver_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.maintenance_records
    ADD CONSTRAINT maintenance_records_driver_id_fkey FOREIGN KEY (driver_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: maintenance_records maintenance_records_vehicle_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.maintenance_records
    ADD CONSTRAINT maintenance_records_vehicle_id_fkey FOREIGN KEY (vehicle_id) REFERENCES public.vehicles(id) ON DELETE CASCADE;


--
-- Name: notifications notifications_tenant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES public.tenants(id) ON DELETE CASCADE;


--
-- Name: notifications notifications_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: refresh_tokens refresh_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: resource_assignments resource_assignments_resource_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_assignments
    ADD CONSTRAINT resource_assignments_resource_id_fkey FOREIGN KEY (resource_id) REFERENCES public.resources(id) ON DELETE CASCADE;


--
-- Name: resource_assignments resource_assignments_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_assignments
    ADD CONSTRAINT resource_assignments_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: resource_categories resource_categories_tenant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_categories
    ADD CONSTRAINT resource_categories_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES public.tenants(id) ON DELETE CASCADE;


--
-- Name: resources resources_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resources
    ADD CONSTRAINT resources_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.resource_categories(id) ON DELETE RESTRICT;


--
-- Name: resources resources_tenant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resources
    ADD CONSTRAINT resources_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES public.tenants(id) ON DELETE CASCADE;


--
-- Name: resources resources_unit_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resources
    ADD CONSTRAINT resources_unit_id_fkey FOREIGN KEY (unit_id) REFERENCES public.units(id) ON DELETE SET NULL;


--
-- Name: resources resources_warehouse_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resources
    ADD CONSTRAINT resources_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES public.warehouses(id) ON DELETE SET NULL;


--
-- Name: shipment_items shipment_items_request_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipment_items
    ADD CONSTRAINT shipment_items_request_id_fkey FOREIGN KEY (request_id) REFERENCES public.supply_requests(id) ON DELETE SET NULL;


--
-- Name: shipment_items shipment_items_resource_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipment_items
    ADD CONSTRAINT shipment_items_resource_id_fkey FOREIGN KEY (resource_id) REFERENCES public.resources(id) ON DELETE RESTRICT;


--
-- Name: shipment_items shipment_items_shipment_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipment_items
    ADD CONSTRAINT shipment_items_shipment_id_fkey FOREIGN KEY (shipment_id) REFERENCES public.shipments(id) ON DELETE CASCADE;


--
-- Name: shipment_refuels shipment_refuels_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipment_refuels
    ADD CONSTRAINT shipment_refuels_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: shipment_refuels shipment_refuels_shipment_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipment_refuels
    ADD CONSTRAINT shipment_refuels_shipment_id_fkey FOREIGN KEY (shipment_id) REFERENCES public.shipments(id) ON DELETE CASCADE;


--
-- Name: shipment_refuels shipment_refuels_tenant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipment_refuels
    ADD CONSTRAINT shipment_refuels_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES public.tenants(id) ON DELETE CASCADE;


--
-- Name: shipment_refuels shipment_refuels_vehicle_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipment_refuels
    ADD CONSTRAINT shipment_refuels_vehicle_id_fkey FOREIGN KEY (vehicle_id) REFERENCES public.vehicles(id) ON DELETE CASCADE;


--
-- Name: shipments shipments_from_warehouse_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipments
    ADD CONSTRAINT shipments_from_warehouse_id_fkey FOREIGN KEY (from_warehouse_id) REFERENCES public.warehouses(id) ON DELETE CASCADE;


--
-- Name: shipments shipments_tenant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipments
    ADD CONSTRAINT shipments_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES public.tenants(id) ON DELETE CASCADE;


--
-- Name: shipments shipments_to_warehouse_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipments
    ADD CONSTRAINT shipments_to_warehouse_id_fkey FOREIGN KEY (to_warehouse_id) REFERENCES public.warehouses(id) ON DELETE CASCADE;


--
-- Name: shipments shipments_vehicle_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipments
    ADD CONSTRAINT shipments_vehicle_id_fkey FOREIGN KEY (vehicle_id) REFERENCES public.vehicles(id) ON DELETE RESTRICT;


--
-- Name: supply_requests supply_requests_approved_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supply_requests
    ADD CONSTRAINT supply_requests_approved_by_fkey FOREIGN KEY (approved_by) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: supply_requests supply_requests_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supply_requests
    ADD CONSTRAINT supply_requests_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: supply_requests supply_requests_resource_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supply_requests
    ADD CONSTRAINT supply_requests_resource_category_id_fkey FOREIGN KEY (resource_category_id) REFERENCES public.resource_categories(id) ON DELETE SET NULL;


--
-- Name: supply_requests supply_requests_resource_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supply_requests
    ADD CONSTRAINT supply_requests_resource_id_fkey FOREIGN KEY (resource_id) REFERENCES public.resources(id) ON DELETE CASCADE;


--
-- Name: supply_requests supply_requests_target_warehouse_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supply_requests
    ADD CONSTRAINT supply_requests_target_warehouse_id_fkey FOREIGN KEY (target_warehouse_id) REFERENCES public.warehouses(id) ON DELETE SET NULL;


--
-- Name: supply_requests supply_requests_tenant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supply_requests
    ADD CONSTRAINT supply_requests_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES public.tenants(id) ON DELETE CASCADE;


--
-- Name: units units_parent_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.units
    ADD CONSTRAINT units_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES public.units(id) ON DELETE SET NULL;


--
-- Name: units units_tenant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.units
    ADD CONSTRAINT units_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES public.tenants(id) ON DELETE CASCADE;


--
-- Name: users users_tenant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES public.tenants(id) ON DELETE CASCADE;


--
-- Name: vehicle_driver_history vehicle_driver_history_driver_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vehicle_driver_history
    ADD CONSTRAINT vehicle_driver_history_driver_id_fkey FOREIGN KEY (driver_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: vehicle_driver_history vehicle_driver_history_vehicle_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vehicle_driver_history
    ADD CONSTRAINT vehicle_driver_history_vehicle_id_fkey FOREIGN KEY (vehicle_id) REFERENCES public.vehicles(id) ON DELETE CASCADE;


--
-- Name: vehicles vehicles_current_warehouse_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vehicles
    ADD CONSTRAINT vehicles_current_warehouse_id_fkey FOREIGN KEY (current_warehouse_id) REFERENCES public.warehouses(id) ON DELETE SET NULL;


--
-- Name: vehicles vehicles_driver_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vehicles
    ADD CONSTRAINT vehicles_driver_id_fkey FOREIGN KEY (driver_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: vehicles vehicles_home_warehouse_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vehicles
    ADD CONSTRAINT vehicles_home_warehouse_id_fkey FOREIGN KEY (home_warehouse_id) REFERENCES public.warehouses(id) ON DELETE SET NULL;


--
-- Name: vehicles vehicles_tenant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vehicles
    ADD CONSTRAINT vehicles_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES public.tenants(id) ON DELETE CASCADE;


--
-- Name: warehouses warehouses_tenant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouses
    ADD CONSTRAINT warehouses_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES public.tenants(id) ON DELETE CASCADE;


--
-- Name: warehouses warehouses_unit_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.warehouses
    ADD CONSTRAINT warehouses_unit_id_fkey FOREIGN KEY (unit_id) REFERENCES public.units(id) ON DELETE CASCADE;


--
-- Name: audit_logs; Type: ROW SECURITY; Schema: public; Owner: -
--

ALTER TABLE public.audit_logs ENABLE ROW LEVEL SECURITY;

--
-- Name: contractor_requests; Type: ROW SECURITY; Schema: public; Owner: -
--

ALTER TABLE public.contractor_requests ENABLE ROW LEVEL SECURITY;

--
-- Name: fuel_records; Type: ROW SECURITY; Schema: public; Owner: -
--

ALTER TABLE public.fuel_records ENABLE ROW LEVEL SECURITY;

--
-- Name: geofences; Type: ROW SECURITY; Schema: public; Owner: -
--

ALTER TABLE public.geofences ENABLE ROW LEVEL SECURITY;

--
-- Name: gps_locations; Type: ROW SECURITY; Schema: public; Owner: -
--

ALTER TABLE public.gps_locations ENABLE ROW LEVEL SECURITY;

--
-- Name: notifications; Type: ROW SECURITY; Schema: public; Owner: -
--

ALTER TABLE public.notifications ENABLE ROW LEVEL SECURITY;

--
-- Name: resources; Type: ROW SECURITY; Schema: public; Owner: -
--

ALTER TABLE public.resources ENABLE ROW LEVEL SECURITY;

--
-- Name: shipment_refuels; Type: ROW SECURITY; Schema: public; Owner: -
--

ALTER TABLE public.shipment_refuels ENABLE ROW LEVEL SECURITY;

--
-- Name: shipments; Type: ROW SECURITY; Schema: public; Owner: -
--

ALTER TABLE public.shipments ENABLE ROW LEVEL SECURITY;

--
-- Name: supply_requests; Type: ROW SECURITY; Schema: public; Owner: -
--

ALTER TABLE public.supply_requests ENABLE ROW LEVEL SECURITY;

--
-- Name: audit_logs tenant_isolation; Type: POLICY; Schema: public; Owner: -
--

CREATE POLICY tenant_isolation ON public.audit_logs USING (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true)))) WITH CHECK (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true))));


--
-- Name: contractor_requests tenant_isolation; Type: POLICY; Schema: public; Owner: -
--

CREATE POLICY tenant_isolation ON public.contractor_requests USING (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true)))) WITH CHECK (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true))));


--
-- Name: fuel_records tenant_isolation; Type: POLICY; Schema: public; Owner: -
--

CREATE POLICY tenant_isolation ON public.fuel_records USING (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true)))) WITH CHECK (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true))));


--
-- Name: geofences tenant_isolation; Type: POLICY; Schema: public; Owner: -
--

CREATE POLICY tenant_isolation ON public.geofences USING (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true)))) WITH CHECK (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true))));


--
-- Name: gps_locations tenant_isolation; Type: POLICY; Schema: public; Owner: -
--

CREATE POLICY tenant_isolation ON public.gps_locations USING (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true)))) WITH CHECK (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true))));


--
-- Name: notifications tenant_isolation; Type: POLICY; Schema: public; Owner: -
--

CREATE POLICY tenant_isolation ON public.notifications USING (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true)))) WITH CHECK (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true))));


--
-- Name: resources tenant_isolation; Type: POLICY; Schema: public; Owner: -
--

CREATE POLICY tenant_isolation ON public.resources USING (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true)))) WITH CHECK (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true))));


--
-- Name: shipment_refuels tenant_isolation; Type: POLICY; Schema: public; Owner: -
--

CREATE POLICY tenant_isolation ON public.shipment_refuels USING (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true)))) WITH CHECK (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true))));


--
-- Name: shipments tenant_isolation; Type: POLICY; Schema: public; Owner: -
--

CREATE POLICY tenant_isolation ON public.shipments USING (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true)))) WITH CHECK (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true))));


--
-- Name: supply_requests tenant_isolation; Type: POLICY; Schema: public; Owner: -
--

CREATE POLICY tenant_isolation ON public.supply_requests USING (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true)))) WITH CHECK (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true))));


--
-- Name: units tenant_isolation; Type: POLICY; Schema: public; Owner: -
--

CREATE POLICY tenant_isolation ON public.units USING (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true)))) WITH CHECK (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true))));


--
-- Name: users tenant_isolation; Type: POLICY; Schema: public; Owner: -
--

CREATE POLICY tenant_isolation ON public.users USING (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true)))) WITH CHECK (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true))));


--
-- Name: vehicles tenant_isolation; Type: POLICY; Schema: public; Owner: -
--

CREATE POLICY tenant_isolation ON public.vehicles USING (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true)))) WITH CHECK (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true))));


--
-- Name: warehouses tenant_isolation; Type: POLICY; Schema: public; Owner: -
--

CREATE POLICY tenant_isolation ON public.warehouses USING (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true)))) WITH CHECK (((current_setting('app.tenant_id'::text, true) IS NULL) OR (current_setting('app.tenant_id'::text, true) = ''::text) OR ((tenant_id)::text = current_setting('app.tenant_id'::text, true))));


--
-- Name: units; Type: ROW SECURITY; Schema: public; Owner: -
--

ALTER TABLE public.units ENABLE ROW LEVEL SECURITY;

--
-- Name: users; Type: ROW SECURITY; Schema: public; Owner: -
--

ALTER TABLE public.users ENABLE ROW LEVEL SECURITY;

--
-- Name: vehicles; Type: ROW SECURITY; Schema: public; Owner: -
--

ALTER TABLE public.vehicles ENABLE ROW LEVEL SECURITY;

--
-- Name: warehouses; Type: ROW SECURITY; Schema: public; Owner: -
--

ALTER TABLE public.warehouses ENABLE ROW LEVEL SECURITY;

--
-- PostgreSQL database dump complete
--

\unrestrict 4IA6MIzqh232qwcJZIukbtnlb0enZHiG1RSuwZRJwfn8tQlvm8EL84pjAZWDmjt

