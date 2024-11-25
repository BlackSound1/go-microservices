--
-- Name: user_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE 
    public.user_id_seq
START WITH 
    1
INCREMENT BY 
    1
NO MINVALUE
NO MAXVALUE
CACHE 
    1;


ALTER TABLE public.user_id_seq OWNER TO postgres;


SET default_tablespace = '';


SET default_table_access_method = heap;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE 
    public.users 
        (
            id INTEGER DEFAULT nextval('public.user_id_seq'::regclass) NOT NULL,
            email CHARACTER varying(255),
            first_name CHARACTER varying(255),
            last_name CHARACTER varying(255),
            password CHARACTER varying(60),
            user_active INTEGER DEFAULT 0,
            created_at TIMESTAMP WITHOUT TIME ZONE,
            updated_at TIMESTAMP WITHOUT TIME ZONE
        );


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_id_seq', 1, true);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY 
    public.users
ADD CONSTRAINT 
    users_pkey PRIMARY KEY (id);


INSERT INTO 
    "public"."users"
        ("email","first_name","last_name","password","user_active","created_at","updated_at")
VALUES
    (E'admin@example.com',E'Admin',E'User',E'$2a$12$1zGLuYDDNvATh4RA4avbKuheAMpb1svexSzrQm7up.bnpwQHs0jNe',1,E'2022-03-14 00:00:00',E'2022-03-14 00:00:00');
