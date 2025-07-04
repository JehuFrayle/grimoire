--
-- PostgreSQL database dump
--

-- Dumped from database version 15.13 (Debian 15.13-1.pgdg120+1)
-- Dumped by pg_dump version 15.13 (Debian 15.13-1.pgdg120+1)

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

--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: notes; Type: TABLE; Schema: public; Owner: grimoire_user
--

CREATE TABLE public.notes (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    title character varying(200) NOT NULL,
    content text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    user_id uuid NOT NULL,
    is_public boolean DEFAULT false,
    deleted_at timestamp with time zone
);


ALTER TABLE public.notes OWNER TO grimoire_user;

--
-- Name: active_notes; Type: VIEW; Schema: public; Owner: grimoire_user
--

CREATE VIEW public.active_notes AS
 SELECT notes.id,
    notes.title,
    notes.content,
    notes.created_at,
    notes.updated_at,
    notes.user_id,
    notes.is_public,
    notes.deleted_at
   FROM public.notes
  WHERE (notes.deleted_at IS NULL);


ALTER TABLE public.active_notes OWNER TO grimoire_user;

--
-- Name: users; Type: TABLE; Schema: public; Owner: grimoire_user
--

CREATE TABLE public.users (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    username character varying(50) NOT NULL,
    email character varying(100) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    role character varying(20) NOT NULL,
    active boolean DEFAULT true NOT NULL,
    password_hash text DEFAULT ''::text NOT NULL,
    deleted_at timestamp with time zone,
    last_login timestamp with time zone,
    CONSTRAINT users_role_check CHECK (((role)::text = ANY ((ARRAY['admin'::character varying, 'user'::character varying, 'guest'::character varying])::text[])))
);


ALTER TABLE public.users OWNER TO grimoire_user;

--
-- Name: active_users; Type: VIEW; Schema: public; Owner: grimoire_user
--

CREATE VIEW public.active_users AS
 SELECT users.id,
    users.username,
    users.email,
    users.created_at,
    users.updated_at,
    users.role,
    users.active,
    users.password_hash,
    users.deleted_at,
    users.last_login
   FROM public.users
  WHERE (users.deleted_at IS NULL);


ALTER TABLE public.active_users OWNER TO grimoire_user;

--
-- Name: links; Type: TABLE; Schema: public; Owner: grimoire_user
--

CREATE TABLE public.links (
    id integer NOT NULL,
    user_id uuid NOT NULL,
    url text NOT NULL,
    title character varying(100),
    icon text,
    active boolean DEFAULT true NOT NULL
);


ALTER TABLE public.links OWNER TO grimoire_user;

--
-- Name: links_id_seq; Type: SEQUENCE; Schema: public; Owner: grimoire_user
--

CREATE SEQUENCE public.links_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.links_id_seq OWNER TO grimoire_user;

--
-- Name: links_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: grimoire_user
--

ALTER SEQUENCE public.links_id_seq OWNED BY public.links.id;


--
-- Name: note_tags; Type: TABLE; Schema: public; Owner: grimoire_user
--

CREATE TABLE public.note_tags (
    note_id uuid NOT NULL,
    tag_id uuid NOT NULL
);


ALTER TABLE public.note_tags OWNER TO grimoire_user;

--
-- Name: profiles; Type: TABLE; Schema: public; Owner: grimoire_user
--

CREATE TABLE public.profiles (
    user_id uuid NOT NULL,
    first_name character varying(50),
    last_name character varying(50),
    bio text,
    avatar_url text
);


ALTER TABLE public.profiles OWNER TO grimoire_user;

--
-- Name: shared_notes; Type: TABLE; Schema: public; Owner: grimoire_user
--

CREATE TABLE public.shared_notes (
    note_id uuid NOT NULL,
    shared_with_user_id uuid NOT NULL,
    shared_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    can_edit boolean DEFAULT false
);


ALTER TABLE public.shared_notes OWNER TO grimoire_user;

--
-- Name: tags; Type: TABLE; Schema: public; Owner: grimoire_user
--

CREATE TABLE public.tags (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(50) NOT NULL
);


ALTER TABLE public.tags OWNER TO grimoire_user;

--
-- Name: links id; Type: DEFAULT; Schema: public; Owner: grimoire_user
--

ALTER TABLE ONLY public.links ALTER COLUMN id SET DEFAULT nextval('public.links_id_seq'::regclass);


--
-- Name: links links_pkey; Type: CONSTRAINT; Schema: public; Owner: grimoire_user
--

ALTER TABLE ONLY public.links
    ADD CONSTRAINT links_pkey PRIMARY KEY (id);


--
-- Name: note_tags note_tags_pkey; Type: CONSTRAINT; Schema: public; Owner: grimoire_user
--

ALTER TABLE ONLY public.note_tags
    ADD CONSTRAINT note_tags_pkey PRIMARY KEY (note_id, tag_id);


--
-- Name: notes notes_pkey; Type: CONSTRAINT; Schema: public; Owner: grimoire_user
--

ALTER TABLE ONLY public.notes
    ADD CONSTRAINT notes_pkey PRIMARY KEY (id);


--
-- Name: profiles profiles_pkey; Type: CONSTRAINT; Schema: public; Owner: grimoire_user
--

ALTER TABLE ONLY public.profiles
    ADD CONSTRAINT profiles_pkey PRIMARY KEY (user_id);


--
-- Name: shared_notes shared_notes_pkey; Type: CONSTRAINT; Schema: public; Owner: grimoire_user
--

ALTER TABLE ONLY public.shared_notes
    ADD CONSTRAINT shared_notes_pkey PRIMARY KEY (note_id, shared_with_user_id);


--
-- Name: tags tags_name_key; Type: CONSTRAINT; Schema: public; Owner: grimoire_user
--

ALTER TABLE ONLY public.tags
    ADD CONSTRAINT tags_name_key UNIQUE (name);


--
-- Name: tags tags_pkey; Type: CONSTRAINT; Schema: public; Owner: grimoire_user
--

ALTER TABLE ONLY public.tags
    ADD CONSTRAINT tags_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: grimoire_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: grimoire_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: grimoire_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: links fk_links_user; Type: FK CONSTRAINT; Schema: public; Owner: grimoire_user
--

ALTER TABLE ONLY public.links
    ADD CONSTRAINT fk_links_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: note_tags fk_notetags_note; Type: FK CONSTRAINT; Schema: public; Owner: grimoire_user
--

ALTER TABLE ONLY public.note_tags
    ADD CONSTRAINT fk_notetags_note FOREIGN KEY (note_id) REFERENCES public.notes(id) ON DELETE CASCADE;


--
-- Name: note_tags fk_notetags_tag; Type: FK CONSTRAINT; Schema: public; Owner: grimoire_user
--

ALTER TABLE ONLY public.note_tags
    ADD CONSTRAINT fk_notetags_tag FOREIGN KEY (tag_id) REFERENCES public.tags(id) ON DELETE CASCADE;


--
-- Name: profiles fk_profiles_user; Type: FK CONSTRAINT; Schema: public; Owner: grimoire_user
--

ALTER TABLE ONLY public.profiles
    ADD CONSTRAINT fk_profiles_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: notes notes_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: grimoire_user
--

ALTER TABLE ONLY public.notes
    ADD CONSTRAINT notes_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: shared_notes shared_notes_note_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: grimoire_user
--

ALTER TABLE ONLY public.shared_notes
    ADD CONSTRAINT shared_notes_note_id_fkey FOREIGN KEY (note_id) REFERENCES public.notes(id) ON DELETE CASCADE;


--
-- Name: shared_notes shared_notes_shared_with_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: grimoire_user
--

ALTER TABLE ONLY public.shared_notes
    ADD CONSTRAINT shared_notes_shared_with_user_id_fkey FOREIGN KEY (shared_with_user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

