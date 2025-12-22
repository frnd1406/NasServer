--
-- PostgreSQL database dump
--

\restrict X1vg6fAut3rTk8KhZt0ByhpGUo9w9V89EyWOOpEaHhcGcpyCUR9kRZvW5mok4lR

-- Dumped from database version 16.11
-- Dumped by pg_dump version 16.11

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
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: public; Owner: nas_user
--

CREATE FUNCTION public.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_updated_at_column() OWNER TO nas_user;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: monitoring_samples; Type: TABLE; Schema: public; Owner: nas_user
--

CREATE TABLE public.monitoring_samples (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    source text NOT NULL,
    cpu_percent numeric(5,2) NOT NULL,
    ram_percent numeric(5,2) NOT NULL,
    created_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.monitoring_samples OWNER TO nas_user;

--
-- Name: refresh_tokens; Type: TABLE; Schema: public; Owner: nas_user
--

CREATE TABLE public.refresh_tokens (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    token_hash character varying(255) NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    revoked boolean DEFAULT false
);


ALTER TABLE public.refresh_tokens OWNER TO nas_user;

--
-- Name: system_alerts; Type: TABLE; Schema: public; Owner: nas_user
--

CREATE TABLE public.system_alerts (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    severity character varying(20) NOT NULL,
    message text NOT NULL,
    is_resolved boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    ai_analysis text
);


ALTER TABLE public.system_alerts OWNER TO nas_user;

--
-- Name: system_metrics; Type: TABLE; Schema: public; Owner: nas_user
--

CREATE TABLE public.system_metrics (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    agent_id text NOT NULL,
    cpu_usage numeric(5,2) NOT NULL,
    ram_usage numeric(5,2) NOT NULL,
    disk_usage numeric(5,2) NOT NULL,
    created_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.system_metrics OWNER TO nas_user;

--
-- Name: system_settings; Type: TABLE; Schema: public; Owner: nas_user
--

CREATE TABLE public.system_settings (
    key text NOT NULL,
    value text NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.system_settings OWNER TO nas_user;

--
-- Name: users; Type: TABLE; Schema: public; Owner: nas_user
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    username character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    password_hash character varying(255) NOT NULL,
    email_verified boolean DEFAULT false,
    verified_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    role character varying(20) DEFAULT 'user'::character varying NOT NULL,
    CONSTRAINT email_format CHECK (((email)::text ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'::text)),
    CONSTRAINT username_min_length CHECK ((length((username)::text) >= 3)),
    CONSTRAINT users_role_check CHECK (((role)::text = ANY ((ARRAY['user'::character varying, 'admin'::character varying])::text[])))
);


ALTER TABLE public.users OWNER TO nas_user;

--
-- Data for Name: monitoring_samples; Type: TABLE DATA; Schema: public; Owner: nas_user
--

COPY public.monitoring_samples (id, source, cpu_percent, ram_percent, created_at) FROM stdin;
\.


--
-- Data for Name: refresh_tokens; Type: TABLE DATA; Schema: public; Owner: nas_user
--

COPY public.refresh_tokens (id, user_id, token_hash, expires_at, created_at, revoked) FROM stdin;
\.


--
-- Data for Name: system_alerts; Type: TABLE DATA; Schema: public; Owner: nas_user
--

COPY public.system_alerts (id, severity, message, is_resolved, created_at, ai_analysis) FROM stdin;
af2fd7b4-d281-4cb3-a38f-b653b8124760	INFO	DB Port accessible internally	f	2025-11-27 20:12:23.816725+00	\N
cf86f229-b0cf-48dc-9586-970c518e8f6c	INFO	DB Port accessible internally	f	2025-11-27 20:13:08.341707+00	\N
b770200a-02a1-4e96-a12d-e3f8817c745c	INFO	DB Port accessible internally	f	2025-11-27 20:14:08.389881+00	\N
661c3b49-bfdc-4271-b6e4-fcccd8d50f40	INFO	DB Port accessible internally	f	2025-11-27 20:15:08.367768+00	\N
78880c33-71d4-4b47-9099-879f02108eb3	INFO	DB Port accessible internally	f	2025-11-27 20:16:08.358981+00	\N
42661c77-e198-4f5a-9e42-80addae3ef24	INFO	DB Port accessible internally	f	2025-11-27 20:17:08.387675+00	\N
02f5d8e7-62c7-4fe2-912c-24904a267a2e	INFO	DB Port accessible internally	f	2025-11-27 20:18:08.374698+00	\N
259a2177-6ea9-47f1-85d2-aa346c99ce6d	INFO	DB Port accessible internally	f	2025-11-27 20:19:08.380727+00	\N
b4a36cc6-62c7-4805-951f-e83db4d0f91f	INFO	DB Port accessible internally	f	2025-11-27 20:20:08.371822+00	\N
e8286b63-f0fa-4a36-997e-83346e2108dd	INFO	DB Port accessible internally	f	2025-11-28 10:54:35.245027+00	\N
b2635c89-c80d-44a9-a839-74c9fdb9203d	INFO	DB Port accessible internally	f	2025-11-28 10:55:35.225039+00	\N
5932bbaf-6e3b-4075-87d6-3c35d1d32fe6	INFO	DB Port accessible internally	f	2025-11-28 10:56:35.283414+00	\N
f0e134a6-4f59-4b54-b553-9aad892b4640	INFO	DB Port accessible internally	f	2025-11-28 10:57:35.249552+00	\N
dd3c61ca-dec5-47c8-9a95-8448255945b9	INFO	DB Port accessible internally	f	2025-11-28 10:58:35.281083+00	\N
1662de86-4ef4-4006-9079-0a1cf4b20399	INFO	DB Port accessible internally	f	2025-11-28 10:59:35.296674+00	\N
ccb47b4d-7c30-4223-bcfc-147d44f0768f	INFO	DB Port accessible internally	f	2025-11-28 11:00:35.291928+00	\N
187dc14b-321d-4cc6-93d2-639e220ad167	INFO	DB Port accessible internally	f	2025-11-28 11:01:35.257178+00	\N
3de7aa76-5e86-498b-b4f3-83de25025428	INFO	DB Port accessible internally	f	2025-11-28 11:02:35.269914+00	\N
8b86959c-df5f-4e43-8c5e-3b617e787167	INFO	DB Port accessible internally	f	2025-11-28 11:03:35.278086+00	\N
ec881822-c584-4012-9a1b-e6a9516ecbbf	INFO	DB Port accessible internally	f	2025-11-28 11:04:35.283044+00	\N
12e469b8-68c6-4d3e-b460-983c34eeef61	INFO	DB Port accessible internally	f	2025-11-28 11:05:35.278207+00	\N
3805d609-f039-4c5a-a820-71c339c180ee	INFO	DB Port accessible internally	f	2025-11-28 11:06:35.29752+00	\N
63b87a33-9059-4722-94e3-eac179f85a8c	INFO	DB Port accessible internally	f	2025-11-28 11:07:35.286116+00	\N
69604b49-4786-4d10-903c-1a9280f30fc4	INFO	DB Port accessible internally	f	2025-11-28 11:08:35.280909+00	\N
398f011d-dcbd-43c5-a3a7-27f986be6acc	INFO	DB Port accessible internally	f	2025-11-28 11:09:35.281273+00	\N
23cbea90-a720-44c8-82b2-bab28c861f1a	INFO	DB Port accessible internally	f	2025-11-28 11:10:35.284611+00	\N
dbc316a5-8a6b-47d3-bc6f-262175a33443	INFO	DB Port accessible internally	f	2025-11-28 11:11:35.297646+00	\N
1a57ef3a-410b-4567-9e60-77889572b630	INFO	DB Port accessible internally	f	2025-11-28 11:12:35.282927+00	\N
95a0f564-bb61-4a54-9b8e-76d2bc5620b5	INFO	DB Port accessible internally	f	2025-11-28 11:13:35.278763+00	\N
899c9b5a-2f78-4642-bb8e-dbd51930064d	INFO	DB Port accessible internally	f	2025-11-28 11:14:35.232845+00	\N
344f9e34-c90b-451c-8cff-387112c84e25	INFO	DB Port accessible internally	f	2025-11-28 11:15:35.264313+00	\N
56211ca1-afe3-4f76-8e5c-a0f178abd2dd	INFO	DB Port accessible internally	f	2025-11-28 11:16:35.257144+00	\N
c6ce1a73-0458-4a3b-bd7c-e54bc94bda54	INFO	DB Port accessible internally	f	2025-11-28 11:17:35.283837+00	\N
ae4723cb-394d-4d54-8c2b-3b627a0db4ea	INFO	DB Port accessible internally	f	2025-11-28 11:18:35.280053+00	\N
b1862843-6d79-4260-b207-049a3b046698	INFO	DB Port accessible internally	f	2025-11-28 11:19:35.284968+00	\N
db669ecb-bece-49d5-82da-d6d0b6b33899	INFO	DB Port accessible internally	f	2025-11-28 11:20:35.321978+00	\N
83165b56-3e44-4ecd-8a4f-ebc98e96f6e4	INFO	DB Port accessible internally	f	2025-11-28 11:21:35.28966+00	\N
74de370f-d6d2-420e-a710-6f413065cf8c	INFO	DB Port accessible internally	f	2025-11-28 11:22:35.283529+00	\N
615bd3bb-df65-43a9-a309-f83241e1737c	INFO	DB Port accessible internally	f	2025-11-28 11:23:35.245002+00	\N
8fe3d8b0-ecc2-4cb2-b5fd-04afb67178f2	INFO	DB Port accessible internally	f	2025-11-28 11:24:35.23283+00	\N
d6972aef-cc6d-43a4-958b-17a1e9e5d1ea	INFO	DB Port accessible internally	f	2025-11-28 11:25:35.282285+00	\N
bb0f985a-ba62-491f-a7c5-636abb68a408	INFO	DB Port accessible internally	f	2025-11-28 11:26:35.249317+00	\N
6517599a-c512-47ce-9440-cdee74926f14	INFO	DB Port accessible internally	f	2025-11-28 11:27:35.283177+00	\N
6057114d-2c59-483f-b154-b63ea196e4cb	INFO	DB Port accessible internally	f	2025-11-28 11:28:35.281779+00	\N
2ab45b18-8d00-40b2-a872-1426f589576b	INFO	DB Port accessible internally	f	2025-11-28 11:29:35.283358+00	\N
b57deabc-161b-492b-9a89-e6dd91072458	INFO	DB Port accessible internally	f	2025-11-28 11:30:35.28292+00	\N
d7f27b46-6899-4743-973d-f1eee2ceedb8	INFO	DB Port accessible internally	f	2025-11-28 11:31:35.279973+00	\N
25fdb18b-c256-4515-8c9e-71b4adfe62c3	INFO	DB Port accessible internally	f	2025-11-28 11:43:07.828253+00	\N
cae9bb80-7a8c-45c0-adab-849997280870	INFO	DB Port accessible internally	f	2025-11-28 11:44:07.869973+00	\N
d4ca4d6c-e266-4a80-b34c-56d8379f220a	INFO	DB Port accessible internally	f	2025-11-28 11:45:07.870356+00	\N
f7f67965-9802-483e-a7b2-2bb2f16c2163	INFO	DB Port accessible internally	f	2025-11-28 11:46:07.813324+00	\N
93c9a211-78ac-4fa6-9117-35266863f268	INFO	DB Port accessible internally	f	2025-11-28 11:47:07.867692+00	\N
50a57b9e-fcc9-4a89-bb64-d1e064046fd2	INFO	DB Port accessible internally	f	2025-11-28 11:48:07.865527+00	\N
54680e85-c533-489c-8296-73687be618dc	INFO	DB Port accessible internally	f	2025-11-28 11:49:07.829571+00	\N
b530099d-c7c9-4374-90f4-2dc49c19fc4c	INFO	DB Port accessible internally	f	2025-11-28 11:50:07.837044+00	\N
09cdcccf-7b79-4478-8b9c-4f6ae898302b	INFO	DB Port accessible internally	f	2025-11-28 11:51:07.826598+00	\N
f8264a2e-49dd-4f6d-8f42-540dbd3a7bc5	INFO	DB Port accessible internally	f	2025-11-28 11:52:07.893929+00	\N
4bd2ceb4-89ef-4ae0-9c8c-c743e2504422	CRITICAL	High CPU usage: avg 85.84% over last 60 seconds	f	2025-11-28 11:52:52.367763+00	\N
13233aaa-fb02-4d9e-8ae5-c6d1aa1a0111	INFO	DB Port accessible internally	f	2025-11-28 11:53:07.872708+00	\N
53b0fe38-54ee-4a6b-9d3c-e73cbb3e2322	INFO	DB Port accessible internally	f	2025-11-28 11:54:07.884065+00	\N
0bd8baf8-2c54-4c4a-987c-15cedd591d4c	INFO	DB Port accessible internally	f	2025-11-28 11:55:07.898844+00	\N
1856b70c-f6c0-4edc-969e-a358d0e1462b	INFO	DB Port accessible internally	f	2025-11-28 11:56:07.868367+00	\N
4cb91aa5-ab60-428e-9227-21d65ca93b90	INFO	DB Port accessible internally	f	2025-11-28 11:57:07.856675+00	\N
66bacc5e-c7ed-442f-a81f-50154f36196f	INFO	DB Port accessible internally	f	2025-11-28 11:59:17.430126+00	\N
d2ba5fd2-0bca-4e54-8d69-1d76ee9c1e28	INFO	DB Port accessible internally	f	2025-11-28 12:00:17.470772+00	\N
629fa9d4-284d-456a-9b48-fad73518ad8d	INFO	DB Port accessible internally	f	2025-11-28 12:01:17.463967+00	\N
7bb0a27d-5df9-422d-9f70-f23766c5309b	INFO	DB Port accessible internally	f	2025-11-28 12:02:17.483555+00	\N
76fd1e86-fbd7-4261-8555-a932bcfd2c93	INFO	DB Port accessible internally	f	2025-11-28 12:03:17.460357+00	\N
ee5a3db1-4281-4ba6-b053-d47100b59d79	INFO	DB Port accessible internally	f	2025-11-28 12:04:17.472248+00	\N
c365ef63-c797-495b-81a9-d637c50ab451	INFO	DB Port accessible internally	f	2025-11-28 12:05:17.487827+00	\N
147acf24-6306-400d-8b42-e199c0a86d1c	INFO	DB Port accessible internally	f	2025-11-28 12:06:17.483397+00	\N
840de0d3-3a30-408e-9d42-ca57d0220697	INFO	DB Port accessible internally	f	2025-11-28 12:07:17.487173+00	\N
65bc8a87-4f5c-45ff-932b-61510b7465f8	INFO	DB Port accessible internally	f	2025-11-28 12:08:17.487713+00	\N
07850130-f00c-423b-b1a0-f40d6c1a2623	INFO	DB Port accessible internally	f	2025-11-28 12:09:17.485017+00	\N
89b90d17-76dc-48c3-8823-94c29a8188b7	INFO	DB Port accessible internally	f	2025-11-28 12:10:17.502517+00	\N
d877efe7-befa-4684-953f-ced76cea5c96	INFO	DB Port accessible internally	f	2025-11-28 12:11:17.468498+00	\N
5955a6e1-6c89-4bff-b3aa-235996f6b596	INFO	DB Port accessible internally	f	2025-11-28 12:12:17.477446+00	\N
4e383049-5a8d-4daf-92f1-121f0dbb4748	INFO	DB Port accessible internally	f	2025-11-28 12:13:17.453143+00	\N
107cb428-efe3-4a48-ab92-12921981a650	INFO	DB Port accessible internally	f	2025-11-28 12:14:17.551802+00	\N
e343941e-0a2d-4572-bb4d-4bed2a199506	INFO	DB Port accessible internally	f	2025-11-28 12:15:17.504574+00	\N
3746b4ba-98ef-43f6-bb7e-7c5d4ef40384	INFO	DB Port accessible internally	f	2025-11-28 12:16:17.467339+00	\N
8238ba5d-e630-44f7-9994-1e986d85df59	INFO	DB Port accessible internally	f	2025-11-28 12:17:17.488307+00	\N
1a7a112e-dc90-4078-a33e-130a5971d8f4	INFO	DB Port accessible internally	f	2025-11-28 12:18:17.484954+00	\N
ee468a1f-1b03-4761-9f0d-6d4d3a491966	INFO	DB Port accessible internally	f	2025-11-28 12:19:17.486974+00	\N
beb25b4b-e45b-4e23-ae37-e86f36267126	INFO	DB Port accessible internally	f	2025-11-28 12:20:17.482216+00	\N
caf75646-cfa8-4837-b28b-4667b76e94c0	INFO	DB Port accessible internally	f	2025-11-28 12:21:17.489614+00	\N
6cdf9240-eed7-4693-902f-d301be1b3ece	INFO	DB Port accessible internally	f	2025-11-28 12:22:17.484356+00	\N
98acc65c-0be7-4e9a-9f96-f97a66bc8d4d	INFO	DB Port accessible internally	f	2025-11-28 12:23:17.478157+00	\N
5d6dc71b-5e4f-4101-9994-ce4ba6b22693	INFO	DB Port accessible internally	f	2025-11-28 12:24:17.455936+00	\N
053d5e96-96ac-4d8c-a315-681862d6e95f	INFO	DB Port accessible internally	f	2025-11-28 12:25:17.474515+00	\N
1ccec635-a2d1-4036-bec6-19dd2bea49a7	INFO	DB Port accessible internally	f	2025-11-28 12:26:17.472963+00	\N
fb34b96e-f9bb-4ea0-b384-fbf61e3ce811	INFO	DB Port accessible internally	f	2025-11-28 12:27:17.456092+00	\N
d0ade184-09cc-4667-8410-df1493f97929	INFO	DB Port accessible internally	f	2025-11-28 12:28:17.484855+00	\N
2e2afb5d-1d1e-4f42-a91f-c6e406aab571	INFO	DB Port accessible internally	f	2025-11-28 12:29:17.459323+00	\N
f4212c8d-85e4-4f18-9c0b-f9cb8ad7661c	INFO	DB Port accessible internally	f	2025-11-28 12:30:17.472796+00	\N
3f5b70e3-dde3-4e7b-b3ab-0c1b37776ea7	INFO	DB Port accessible internally	f	2025-11-28 12:31:17.455939+00	\N
e464618f-d48d-45be-9eec-d70c1cd27b39	INFO	DB Port accessible internally	f	2025-11-28 12:32:17.444164+00	\N
ec9f9fbc-cec9-4887-b867-a6d4522a2220	INFO	DB Port accessible internally	f	2025-11-28 12:33:17.483855+00	\N
7fc6f752-fa57-42bb-8f0b-67bc1b154b2d	INFO	DB Port accessible internally	f	2025-11-28 12:34:17.448782+00	\N
957da5cc-accf-4524-9705-d5fb674be500	INFO	DB Port accessible internally	f	2025-11-28 12:35:17.48629+00	\N
21685ef7-de0a-4506-95e4-5a26c06cedac	INFO	DB Port accessible internally	f	2025-11-28 12:36:17.452349+00	\N
ecabd8c7-8e38-4d78-8d7e-3e44a7664eff	INFO	DB Port accessible internally	f	2025-11-28 12:37:17.482493+00	\N
a7f8024a-52bb-4a6c-8adb-266a0ad3d7a1	INFO	DB Port accessible internally	f	2025-11-28 12:38:17.45729+00	\N
5f35eec0-f907-4235-9e3b-39bb27213d48	INFO	DB Port accessible internally	f	2025-11-28 12:39:17.481255+00	\N
3980d878-fee1-427a-95da-143e493c7f09	INFO	DB Port accessible internally	f	2025-11-28 12:40:17.49938+00	\N
a37072e1-7c09-4f82-80bd-ce55452f7f86	INFO	DB Port accessible internally	f	2025-11-28 12:41:17.47609+00	\N
f1301e7e-36d9-4209-8754-f42feee73c59	INFO	DB Port accessible internally	f	2025-11-28 12:42:17.466049+00	\N
7c364be2-e316-441b-8034-05716ff0184d	INFO	DB Port accessible internally	f	2025-11-28 12:43:17.461028+00	\N
b26ad34d-362c-4817-a5ce-2f856b5292d7	INFO	DB Port accessible internally	f	2025-11-28 12:44:17.486823+00	\N
3cbb9590-e9c8-4ca2-97e2-a292634b15b3	INFO	DB Port accessible internally	f	2025-11-28 12:45:17.486998+00	\N
28860570-33c1-4bfe-bd0c-788a7ab423b2	INFO	DB Port accessible internally	f	2025-11-28 12:46:17.457925+00	\N
4f9c6db9-a0f8-46c6-bb26-d11925ff2dc5	INFO	DB Port accessible internally	f	2025-11-28 12:47:17.468016+00	\N
d69da591-8efb-4647-9086-568553aadb48	INFO	DB Port accessible internally	f	2025-11-28 12:48:17.486807+00	\N
e3a2c835-d945-4ee6-bbff-557d79271a29	INFO	DB Port accessible internally	f	2025-11-28 12:49:17.454975+00	\N
564207e7-0ef3-4b8d-bd09-f1ade7437635	INFO	DB Port accessible internally	f	2025-11-28 12:50:17.551157+00	\N
58d87182-9074-44c9-9253-f5bd20881fe2	INFO	DB Port accessible internally	f	2025-11-28 12:51:17.493459+00	\N
cd93d7dc-7ec9-4929-989c-eb1ac45e9e06	INFO	DB Port accessible internally	f	2025-11-28 12:52:17.453661+00	\N
1758987d-c79b-45eb-9536-bce795b3a7e7	INFO	DB Port accessible internally	f	2025-11-28 12:53:17.438577+00	\N
fcc93084-3daa-4526-93ce-22fb48b7d855	INFO	DB Port accessible internally	f	2025-11-28 12:54:17.492998+00	\N
582d3e9f-3ced-4fa3-9fba-db62619fa005	INFO	DB Port accessible internally	f	2025-11-28 12:55:17.455382+00	\N
825c97b5-24d6-4cea-a351-428e2fb94033	INFO	DB Port accessible internally	f	2025-11-28 12:56:17.478218+00	\N
e60ce489-eaf1-4256-bb71-484af4d3d6bb	INFO	DB Port accessible internally	f	2025-11-28 12:57:17.463703+00	\N
1c3fd716-3d4b-425b-8067-ac3edfdb3e5e	INFO	DB Port accessible internally	f	2025-11-28 12:58:17.483557+00	\N
c56dec37-4522-4fa6-813f-9b8560c1baac	INFO	DB Port accessible internally	f	2025-11-28 12:59:17.474544+00	\N
f759cc0f-8648-4d0c-a1bd-7b0ce24a47b3	INFO	DB Port accessible internally	f	2025-11-28 13:00:17.469874+00	\N
323fbc6d-7706-4d2b-820c-4253d53c9b81	INFO	DB Port accessible internally	f	2025-11-28 13:01:17.467041+00	\N
92ce9ed3-209f-46dd-b90e-04ce62855b8f	INFO	DB Port accessible internally	f	2025-11-28 13:02:17.489904+00	\N
e32cbf25-6f1d-45a3-8079-8dc24562a9a9	INFO	DB Port accessible internally	f	2025-11-28 13:03:17.486842+00	\N
8b9acc2f-2731-4de9-b51a-744b35e43622	INFO	DB Port accessible internally	f	2025-11-28 13:04:17.479672+00	\N
7bad25f7-730d-4903-b06d-fd14c7c15e66	INFO	DB Port accessible internally	f	2025-11-28 13:05:17.481818+00	\N
5ae0a2b3-b88e-4d4e-a1d8-d2b18a63c42c	INFO	DB Port accessible internally	f	2025-11-28 13:06:17.507912+00	\N
0585c8df-7020-4424-9300-683269a7b5bb	INFO	DB Port accessible internally	f	2025-11-28 13:07:17.486076+00	\N
ae3c00d5-dcdd-429b-bd1f-e8d60b1dd8c3	INFO	DB Port accessible internally	f	2025-11-28 13:08:17.455372+00	\N
952fc35c-3e68-4b7e-a3b6-cf7cfbabbcc6	INFO	DB Port accessible internally	f	2025-11-28 13:09:17.451459+00	\N
da8bf1e0-b4d0-446c-aaf4-9bf5cb93de82	INFO	DB Port accessible internally	f	2025-11-28 13:10:17.450183+00	\N
f90e4eb4-2dbd-4f38-a87b-db3e049c8430	INFO	DB Port accessible internally	f	2025-11-28 13:11:17.494381+00	\N
e7699de7-f2bf-4a47-9d9b-1ebe864ad423	INFO	DB Port accessible internally	f	2025-11-28 13:12:17.443372+00	\N
f295e667-2da2-4b0c-9ff9-c33356d8c1ad	INFO	DB Port accessible internally	f	2025-11-28 13:13:17.475219+00	\N
1bbc0d36-7ebf-423c-8587-4ebfffa399ce	INFO	DB Port accessible internally	f	2025-11-28 13:14:17.446576+00	\N
40d2a2a4-1f55-462a-8953-f055c16b9fa6	INFO	DB Port accessible internally	f	2025-11-28 13:15:17.475001+00	\N
8ff222a9-e7e1-42da-bae9-12fbdb351625	INFO	DB Port accessible internally	f	2025-11-28 13:16:17.452783+00	\N
d4348e86-8dc9-4e77-a969-b274b7591b9d	INFO	DB Port accessible internally	f	2025-11-28 13:17:17.479111+00	\N
a8f40378-9e0f-41ff-aa41-4d603fb5e481	INFO	DB Port accessible internally	f	2025-11-28 13:18:17.466854+00	\N
30baa676-2b08-462a-ac61-dc18b6da319f	INFO	DB Port accessible internally	f	2025-11-28 13:19:17.465171+00	\N
b98fd590-f097-465f-ad96-cd87933ff424	INFO	DB Port accessible internally	f	2025-11-28 13:20:17.501951+00	\N
7a79afc5-24e3-42f9-959f-792a68d7e702	INFO	DB Port accessible internally	f	2025-11-28 13:21:17.476179+00	\N
8d10a690-67df-4a8c-bcc5-861d5dc9876d	INFO	DB Port accessible internally	f	2025-11-28 13:22:17.488854+00	\N
83028143-07f9-4c70-852a-2ad38e208dcd	INFO	DB Port accessible internally	f	2025-11-28 13:23:17.467148+00	\N
66009e20-d7db-48e7-b2d0-616ec643cdec	INFO	DB Port accessible internally	f	2025-11-28 13:24:17.470485+00	\N
07b66bba-623b-44a4-a7d9-e5bfce6f2b36	INFO	DB Port accessible internally	f	2025-11-28 13:25:17.467582+00	\N
2bfaca17-b7dd-4bf2-a6b8-53e8d0f86c90	INFO	DB Port accessible internally	f	2025-11-28 13:26:17.515033+00	\N
9543ff06-1249-46f0-8cef-2b554e76c440	INFO	DB Port accessible internally	f	2025-11-28 13:27:17.455033+00	\N
200e51d6-26d0-443e-a7f2-5a8403aa2d16	INFO	DB Port accessible internally	f	2025-11-28 13:28:17.457615+00	\N
2a2d4f0d-056e-4c76-b165-f6bfe56c0944	INFO	DB Port accessible internally	f	2025-11-28 13:29:17.435037+00	\N
a1a6c362-8533-4ca2-bb77-06350bc8a99e	INFO	DB Port accessible internally	f	2025-11-28 13:30:17.46328+00	\N
b27733e5-a433-47dd-b1e6-d1ebf69bc77d	INFO	DB Port accessible internally	f	2025-11-28 13:31:17.427033+00	\N
dc19e5cd-58ed-43c1-a774-e9e47c3290e9	INFO	DB Port accessible internally	f	2025-11-28 13:32:17.488137+00	\N
f2cd6a23-5258-48d8-b08b-0fc0c85ecdae	INFO	DB Port accessible internally	f	2025-11-28 13:33:17.443391+00	\N
c5f86cfb-3c6a-426d-99ba-d0fae8f22be7	INFO	DB Port accessible internally	f	2025-11-28 13:34:17.46979+00	\N
cecef73f-2314-444a-a204-5991918df893	INFO	DB Port accessible internally	f	2025-11-28 13:35:17.440331+00	\N
d849f9b5-6097-45e8-b629-da6296bb0431	INFO	DB Port accessible internally	f	2025-11-28 13:36:17.489859+00	\N
42596958-a88f-45b0-b4c1-dd95a0ed25b7	INFO	DB Port accessible internally	f	2025-11-28 13:37:17.451465+00	\N
0f67f8a9-d803-48dc-be69-1f02b84dd5b2	INFO	DB Port accessible internally	f	2025-11-28 13:38:17.432712+00	\N
ec17e6d3-d70e-49c0-9711-3b7112ecb39c	INFO	DB Port accessible internally	f	2025-11-28 13:39:17.465951+00	\N
03f0b2d8-e7af-41c7-948a-4e9e12143c9f	INFO	DB Port accessible internally	f	2025-11-28 13:40:17.48467+00	\N
7a672578-451c-4b2c-afaf-33644f7c75d4	INFO	DB Port accessible internally	f	2025-11-28 13:41:17.47448+00	\N
4339092a-161f-474f-a2a2-3b74247fe57e	INFO	DB Port accessible internally	f	2025-11-28 13:42:17.491902+00	\N
de2e0de1-324a-4c7c-93f2-2786f9e5a669	INFO	DB Port accessible internally	f	2025-11-28 13:43:17.43809+00	\N
a50248e0-242a-4b87-8872-830171ee078c	INFO	DB Port accessible internally	f	2025-11-28 13:44:17.450817+00	\N
ac49f58a-4850-4033-8684-16becb25f5d0	INFO	DB Port accessible internally	f	2025-11-28 13:45:17.451311+00	\N
c1748f54-2ba7-4f6e-8144-fba360c64a19	INFO	DB Port accessible internally	f	2025-11-28 13:46:17.481879+00	\N
91b9f2c8-a85f-411a-aa0f-2d28f6289987	INFO	DB Port accessible internally	f	2025-11-28 13:47:17.486989+00	\N
594cbee4-654b-4442-bfb2-03823eeeca10	INFO	DB Port accessible internally	f	2025-11-28 13:48:17.452031+00	\N
305a8221-3dc2-4581-996f-882a5195b98a	INFO	DB Port accessible internally	f	2025-11-28 13:49:17.516166+00	\N
a5c096ee-45ab-4460-8d08-756573da341e	INFO	DB Port accessible internally	f	2025-11-28 13:50:17.473448+00	\N
5c8a0385-0244-4d0d-b791-c01cc28bfef0	INFO	DB Port accessible internally	f	2025-11-28 13:51:17.47241+00	\N
4a97111f-6cd8-437b-88ef-b041abd7b4c1	INFO	DB Port accessible internally	f	2025-11-28 13:52:17.438656+00	\N
68351b07-ba18-48dd-95d6-eb5f403cdb4b	INFO	DB Port accessible internally	f	2025-11-28 13:53:17.473173+00	\N
2f1984e5-1331-4c8f-920c-df01f959143b	INFO	DB Port accessible internally	f	2025-11-28 13:54:17.453489+00	\N
c7e05d4a-62ce-49d6-b48a-9f2fa97399b5	INFO	DB Port accessible internally	f	2025-11-28 13:55:17.461103+00	\N
48370cb4-f12d-4db2-898a-a769b15ad6fb	INFO	DB Port accessible internally	f	2025-11-28 13:56:17.4611+00	\N
d73710d1-cd14-4e44-a34c-1788547034d4	INFO	DB Port accessible internally	f	2025-11-28 13:57:17.460158+00	\N
694a0c87-c101-4ded-8ce5-c9b1f1d45c39	INFO	DB Port accessible internally	f	2025-11-28 13:58:17.45397+00	\N
15bbecf9-374a-490e-8fe8-019d0a664944	INFO	DB Port accessible internally	f	2025-11-28 13:59:17.485553+00	\N
d10ddcb6-c755-4e94-a339-630e7b2f4092	INFO	DB Port accessible internally	f	2025-11-28 14:00:17.484509+00	\N
31fedea9-256e-458e-a672-3a2094e1e911	INFO	DB Port accessible internally	f	2025-11-28 14:01:17.468238+00	\N
8e9220cc-d7fc-4899-a236-f47a48d5f7c9	INFO	DB Port accessible internally	f	2025-11-28 14:02:17.534134+00	\N
94831ff0-765b-4282-9556-9ad07663388e	INFO	DB Port accessible internally	f	2025-11-28 14:03:17.4507+00	\N
0a5febd6-9120-4b25-94cb-9b09721947a2	INFO	DB Port accessible internally	f	2025-11-28 14:04:17.451566+00	\N
30528990-522a-4a3c-ad19-10f1ef2a4e6e	INFO	DB Port accessible internally	f	2025-11-28 14:05:17.438242+00	\N
8be9983f-2f3b-4446-9a09-44954d335f8f	INFO	DB Port accessible internally	f	2025-11-28 14:06:17.470747+00	\N
a03232ce-1708-4a35-8605-b217d13b3c60	INFO	DB Port accessible internally	f	2025-11-28 14:07:17.45567+00	\N
d5b74c32-35f0-4aa2-9ec6-13cdd62bb415	INFO	DB Port accessible internally	f	2025-11-28 14:08:17.486424+00	\N
ed861263-6640-43a3-a7ee-456d325dae78	INFO	DB Port accessible internally	f	2025-11-28 14:09:17.462973+00	\N
1ac18747-cb1f-4e33-a954-caea4c0cd0f6	INFO	DB Port accessible internally	f	2025-11-28 14:10:17.504862+00	\N
b3d82714-e364-4430-83f5-dd5bccbb8675	INFO	DB Port accessible internally	f	2025-11-28 14:11:17.473+00	\N
c0a1373d-e421-4ba1-83f9-7001c28365ab	INFO	DB Port accessible internally	f	2025-11-28 14:12:17.470023+00	\N
1d2160cf-80ec-4529-8e47-dc6756e22bb0	INFO	DB Port accessible internally	f	2025-11-28 14:13:17.47321+00	\N
2cf7f565-7424-4a87-b9fe-52c89dfb416c	INFO	DB Port accessible internally	f	2025-11-28 14:14:17.46132+00	\N
f2a09403-a307-41ae-a509-302ed9909c73	INFO	DB Port accessible internally	f	2025-11-28 14:15:17.479557+00	\N
144bca04-5352-4df1-95e1-a3b6d6c75848	INFO	DB Port accessible internally	f	2025-11-28 14:16:17.471448+00	\N
bfe474c0-398b-461c-b9dd-790004b843a7	INFO	DB Port accessible internally	f	2025-11-28 14:17:17.450599+00	\N
0bf89e0d-88c3-44d0-8ef4-23c40073aca0	INFO	DB Port accessible internally	f	2025-11-28 14:18:17.493006+00	\N
f5be9316-ad2e-4530-bce0-f54c500190a4	INFO	DB Port accessible internally	f	2025-11-28 14:19:17.488137+00	\N
e00c2b79-b232-44b7-ae76-8ba158de4fa8	INFO	DB Port accessible internally	f	2025-11-28 14:20:17.482397+00	\N
b7af6e1a-0d82-46a4-912d-27408b25fac2	INFO	DB Port accessible internally	f	2025-11-28 14:21:17.456336+00	\N
015bbfa2-b1b7-4775-8b07-930cf8c85eb9	INFO	DB Port accessible internally	f	2025-11-28 14:22:17.452271+00	\N
e9ce864f-ad42-4e4b-aafd-fcb07b54377d	INFO	DB Port accessible internally	f	2025-11-28 14:23:17.45859+00	\N
1d78e215-0482-453c-b549-6981f959b296	INFO	DB Port accessible internally	f	2025-11-28 14:24:17.478204+00	\N
6514cf8b-d544-4ae7-8638-e889545696bd	INFO	DB Port accessible internally	f	2025-11-28 14:25:17.478998+00	\N
1a572771-13a3-4cb9-9fc3-f70072616f28	INFO	DB Port accessible internally	f	2025-11-28 14:26:17.460323+00	\N
09535ab0-5cc5-4b78-b3d5-5a1249f8cc07	INFO	DB Port accessible internally	f	2025-11-28 14:27:17.46782+00	\N
cc076b2c-45aa-444c-a7c6-56e459f2c990	INFO	DB Port accessible internally	f	2025-11-28 14:28:17.459692+00	\N
c05a7ddc-0338-47ee-8281-1ec8e750e360	INFO	DB Port accessible internally	f	2025-11-28 14:29:17.476138+00	\N
0774a526-93c2-4ecd-8b4a-bd0d2985d190	INFO	DB Port accessible internally	f	2025-11-28 14:30:17.443591+00	\N
e7e0011b-94ef-4b2a-bccc-0b5eea9cf994	INFO	DB Port accessible internally	f	2025-11-28 14:31:17.475478+00	\N
961cd4f7-9b5c-430d-89c6-23e7da713fb7	INFO	DB Port accessible internally	f	2025-11-28 14:32:17.467454+00	\N
41ab8c2c-e9fd-4d6a-91bb-49c6df3bd358	INFO	DB Port accessible internally	f	2025-11-28 14:33:17.487962+00	\N
92195421-d07d-46bf-8027-9a2886f8fa6e	INFO	DB Port accessible internally	f	2025-11-28 14:34:17.42917+00	\N
10be58fc-c295-47de-984b-5196f99a5803	INFO	DB Port accessible internally	f	2025-11-28 14:35:17.488036+00	\N
0fc3dbbe-4c24-497c-beca-fbc96dc78f64	INFO	DB Port accessible internally	f	2025-11-28 14:36:17.449766+00	\N
79ca9da4-930c-473d-85dd-8dec745b20f2	INFO	DB Port accessible internally	f	2025-11-28 14:37:17.475544+00	\N
f6ae0448-8daf-4b83-aebd-1ee6f8ec74f7	INFO	DB Port accessible internally	f	2025-11-28 14:38:17.510023+00	\N
fef73828-96ef-4e75-9616-5fe362ee0929	INFO	DB Port accessible internally	f	2025-11-28 14:39:17.492196+00	\N
2f966d4a-319f-44b8-8996-414bcc7896fc	INFO	DB Port accessible internally	f	2025-11-28 14:40:17.473792+00	\N
5484d50e-33b5-4144-b576-10bc4aed68bb	INFO	DB Port accessible internally	f	2025-11-28 14:41:17.482006+00	\N
3e442924-f953-4fdc-a2e0-0e4f1e1235a9	INFO	DB Port accessible internally	f	2025-11-28 14:42:17.452958+00	\N
b654344d-f51b-4575-a4de-71f2ea28a257	INFO	DB Port accessible internally	f	2025-11-28 14:43:17.443614+00	\N
d913e72e-8d3e-41c8-a0d7-bc16883f9dd8	INFO	DB Port accessible internally	f	2025-11-28 14:44:17.508894+00	\N
cb2e326c-fd5b-4ffe-bbf9-0fa31f0850b8	INFO	DB Port accessible internally	f	2025-11-28 14:45:17.487354+00	\N
31c25e5b-99c4-49a9-b7d5-9fbe1bed90e0	INFO	DB Port accessible internally	f	2025-11-28 14:46:17.504226+00	\N
54f637a5-65b0-44dc-a195-8c96e177426c	INFO	DB Port accessible internally	f	2025-11-28 14:47:17.464075+00	\N
4ea55a08-adf6-41bd-9720-ab8fa0cfeda2	INFO	DB Port accessible internally	f	2025-11-28 14:48:17.511823+00	\N
04e6972c-9726-4cbc-a7a2-a09eac77a31f	INFO	DB Port accessible internally	f	2025-11-28 14:49:17.48806+00	\N
ca401735-fd53-4740-965b-54b897984bc0	INFO	DB Port accessible internally	f	2025-11-28 14:50:17.469775+00	\N
dda400e5-99e6-49ff-99c1-fd2c3a7b28cc	INFO	DB Port accessible internally	f	2025-11-28 14:51:17.473142+00	\N
0370afad-9395-4f10-ade9-d5567c129613	INFO	DB Port accessible internally	f	2025-11-28 14:52:17.47778+00	\N
895ccd21-a47e-4c4f-a6c9-33ea08b16ab0	INFO	DB Port accessible internally	f	2025-11-28 14:53:17.435845+00	\N
e031433c-13b3-4052-8f53-e19aecad7014	INFO	DB Port accessible internally	f	2025-11-28 14:54:17.491991+00	\N
7e1e4cc8-a9b9-4831-9528-0ff9dfd4b719	INFO	DB Port accessible internally	f	2025-11-28 14:55:17.460595+00	\N
ec095b26-75ab-4e76-a619-fb0a1ca8a268	INFO	DB Port accessible internally	f	2025-11-28 14:56:17.475911+00	\N
fd10f968-4282-4075-874e-6cb676ecd44e	INFO	DB Port accessible internally	f	2025-11-28 14:57:17.452156+00	\N
8966252c-5914-46c9-82bf-f195a524c5f2	INFO	DB Port accessible internally	f	2025-11-28 14:58:17.451967+00	\N
cdf00f6a-7177-4b75-a629-3610cc0b38b0	INFO	DB Port accessible internally	f	2025-11-28 14:59:17.457663+00	\N
07e21403-ef8a-40ac-b6ef-f118fd8e52d5	INFO	DB Port accessible internally	f	2025-11-28 15:00:17.478863+00	\N
20a55bde-9e85-4813-87b9-80c6cb3af548	INFO	DB Port accessible internally	f	2025-11-28 15:01:17.47779+00	\N
421054e2-835b-4957-aadd-004f0b5e16e3	INFO	DB Port accessible internally	f	2025-11-28 15:02:17.47374+00	\N
5b4be8e2-3454-466f-b0cf-e26df6951c38	INFO	DB Port accessible internally	f	2025-11-28 15:03:17.488709+00	\N
404efc64-0491-4fad-a874-26e35c8877e8	INFO	DB Port accessible internally	f	2025-11-28 15:04:17.459909+00	\N
fa2a4999-8ebe-441c-a7cf-798815d56c09	INFO	DB Port accessible internally	f	2025-11-28 15:05:17.464599+00	\N
5d3e8ded-390e-46f6-a398-66d6470ddf75	INFO	DB Port accessible internally	f	2025-11-28 15:06:17.468332+00	\N
314aa542-6286-4442-8df8-a7c8e4509248	INFO	DB Port accessible internally	f	2025-11-28 15:07:17.455944+00	\N
dbf2c0b5-a1a9-42fc-92da-704148af3fb3	INFO	DB Port accessible internally	f	2025-11-28 15:08:17.462102+00	\N
dc07001a-ebb2-4582-b420-88230ca7b74a	INFO	DB Port accessible internally	f	2025-11-28 15:09:17.457382+00	\N
d4dda4d4-f57d-4ca3-a20a-6388381df4bf	INFO	DB Port accessible internally	f	2025-11-28 15:10:17.494801+00	\N
8f69ee24-91fb-41e3-a332-efbea758e96d	INFO	DB Port accessible internally	f	2025-11-28 15:11:17.450937+00	\N
d255737b-2b70-451e-9454-83eaf2a639a6	INFO	DB Port accessible internally	f	2025-11-28 15:12:17.438826+00	\N
0ed4a512-f125-4cb2-825e-7497310bff19	INFO	DB Port accessible internally	f	2025-11-28 15:13:17.461575+00	\N
f0055110-4cee-44d1-9ba0-d5baf0c90dfd	INFO	DB Port accessible internally	f	2025-11-28 15:14:17.501484+00	\N
9769a16e-ef42-47e0-a8c3-01ce4e87505e	INFO	DB Port accessible internally	f	2025-11-28 15:15:17.439487+00	\N
31f22834-2c61-4c8c-8b0b-ef1a6d13bc2d	INFO	DB Port accessible internally	f	2025-11-28 15:16:17.447554+00	\N
4438f17d-718a-4f42-be41-657334487e5f	INFO	DB Port accessible internally	f	2025-11-28 15:17:17.477207+00	\N
c1dd0736-c3e9-4a6d-b7fc-a6aa440cbb82	INFO	DB Port accessible internally	f	2025-11-28 15:18:17.454328+00	\N
79767204-8332-4d1a-882b-b4ca11592f87	INFO	DB Port accessible internally	f	2025-11-28 15:19:17.468982+00	\N
b750ce96-178c-491d-86a0-9c59e11fd7b1	INFO	DB Port accessible internally	f	2025-11-28 15:20:17.458872+00	\N
a0d329a3-5f04-4306-904c-5ee8ff7acc86	INFO	DB Port accessible internally	f	2025-11-28 15:21:17.452838+00	\N
461d0237-9a50-4ed1-a068-112d436f1ff0	INFO	DB Port accessible internally	f	2025-11-28 15:22:17.47689+00	\N
7769f7bd-2a88-44ae-a785-4e91d6d64ca7	INFO	DB Port accessible internally	f	2025-11-28 15:23:17.464417+00	\N
372f5d02-dcad-4176-a497-49389681771e	INFO	DB Port accessible internally	f	2025-11-28 15:24:17.452414+00	\N
a2bd4a8f-d9db-4498-a8f3-be64b89bf023	INFO	DB Port accessible internally	f	2025-11-28 15:25:17.456677+00	\N
8606896c-98da-4d66-9676-24ddcee8d31f	INFO	DB Port accessible internally	f	2025-11-28 15:26:17.482282+00	\N
59bafcd5-fc2e-4182-9a37-9f1b98410ed6	INFO	DB Port accessible internally	f	2025-11-28 15:27:17.490697+00	\N
13de7499-6066-49c9-a282-4f93d9619602	INFO	DB Port accessible internally	f	2025-11-28 15:28:17.461127+00	\N
738623c7-b012-43c9-bf8b-1718e6071d04	INFO	DB Port accessible internally	f	2025-11-28 15:29:17.431298+00	\N
8cf4d917-6003-401e-8b0a-d40e69a12e75	INFO	DB Port accessible internally	f	2025-11-28 15:30:17.4612+00	\N
49481d23-82ff-4af0-96ba-9c6c4389cf5b	INFO	DB Port accessible internally	f	2025-11-28 15:31:17.440343+00	\N
249e23d1-5bef-466d-86ec-1b3942d56395	INFO	DB Port accessible internally	f	2025-11-28 15:32:17.465857+00	\N
22973985-8124-4898-83f4-d283601486fe	INFO	DB Port accessible internally	f	2025-11-28 15:33:17.454231+00	\N
7c99c01a-6c55-45f0-890d-829a93ae1161	INFO	DB Port accessible internally	f	2025-11-28 15:34:17.495142+00	\N
79c6f3f6-e4c2-4310-b4bb-32ccdc30eecd	INFO	DB Port accessible internally	f	2025-11-28 15:35:17.492139+00	\N
a215d1ed-bb68-4f90-9cf8-92cbf620b053	INFO	DB Port accessible internally	f	2025-11-28 15:36:17.465312+00	\N
6b40cdab-b651-469a-a917-dd27fed64905	INFO	DB Port accessible internally	f	2025-11-28 15:37:17.430469+00	\N
42032526-a526-4f0a-89ff-c1adee0ecd35	INFO	DB Port accessible internally	f	2025-11-28 15:38:17.454933+00	\N
9ba70e8d-e03e-4b4c-b59f-7d1889e94fcc	INFO	DB Port accessible internally	f	2025-11-28 15:39:17.4352+00	\N
810998a6-8885-4467-b3ab-eee40922f2d3	INFO	DB Port accessible internally	f	2025-11-28 15:40:17.501925+00	\N
c2972e2c-e41d-45e7-b2c6-c42e1d697162	INFO	DB Port accessible internally	f	2025-11-28 15:41:17.440706+00	\N
762ed625-e561-45f7-9828-ee4911583560	INFO	DB Port accessible internally	f	2025-11-28 15:42:17.475234+00	\N
f1b70755-5801-4f87-a8be-fb593fb25aa5	INFO	DB Port accessible internally	f	2025-11-28 15:43:17.441673+00	\N
95aa089e-f0e7-411c-b424-ec85f8680b1d	INFO	DB Port accessible internally	f	2025-11-28 15:44:17.487652+00	\N
ee6c43e7-deb8-4bef-885c-9c3efd8746c6	INFO	DB Port accessible internally	f	2025-11-28 15:45:17.445312+00	\N
8e3084ed-65b1-4133-9ab7-34aff01265e2	INFO	DB Port accessible internally	f	2025-11-28 15:46:17.442351+00	\N
d441e814-2153-471b-946b-c7a6eb64943e	INFO	DB Port accessible internally	f	2025-11-28 15:47:17.445206+00	\N
a4fd799a-f7c3-4948-8eb3-bf38ba322544	INFO	DB Port accessible internally	f	2025-11-28 15:48:17.49422+00	\N
fc43ac1a-69cf-4238-98fb-e8bdefcef9b5	INFO	DB Port accessible internally	f	2025-11-28 15:49:17.450089+00	\N
d32036ab-fe1b-444a-b0fa-e9dc131613a2	INFO	DB Port accessible internally	f	2025-11-28 15:50:17.45359+00	\N
b506af18-0ee7-4ece-9dbc-45f8a8b6f80b	INFO	DB Port accessible internally	f	2025-11-28 15:51:17.464004+00	\N
26f3a489-54b3-4800-a4c0-392852fb10f0	INFO	DB Port accessible internally	f	2025-11-28 15:52:17.481829+00	\N
d99a5056-0ce9-41ba-b19d-f9470cb44c51	INFO	DB Port accessible internally	f	2025-11-28 15:53:17.498858+00	\N
432331a4-64d6-4cf7-9518-2449c3066233	INFO	DB Port accessible internally	f	2025-11-28 15:54:17.514664+00	\N
bb2e16df-1bcd-45ce-a827-fddcab7686d3	INFO	DB Port accessible internally	f	2025-11-28 15:55:17.477064+00	\N
aa2cb51e-1d10-4acd-849b-b810474fd729	INFO	DB Port accessible internally	f	2025-11-28 15:56:17.484088+00	\N
93a330bc-b6eb-4d6c-afbc-b033c1587413	INFO	DB Port accessible internally	f	2025-11-28 15:57:17.461816+00	\N
b409f741-e062-47da-b4f8-fed793374fc9	INFO	DB Port accessible internally	f	2025-11-28 15:58:17.458119+00	\N
7b374394-5165-4737-9cc8-54b9869ca3fa	INFO	DB Port accessible internally	f	2025-11-28 15:59:17.48185+00	\N
01bf5871-9082-47ba-a84c-33a913607d9e	INFO	DB Port accessible internally	f	2025-11-28 16:00:17.489598+00	\N
423d926a-e701-43d0-96cb-b35f3bf26d68	INFO	DB Port accessible internally	f	2025-11-28 16:02:10.228892+00	\N
faa165c4-d29a-4702-904c-e3ddf54dd5df	INFO	DB Port accessible internally	f	2025-11-28 16:03:10.279843+00	\N
6e6680f5-80af-49ed-bcb5-1d1b1fab02be	INFO	DB Port accessible internally	f	2025-11-28 16:04:10.273209+00	\N
7c5002cf-afab-4099-bb1c-a9e959159ffa	INFO	DB Port accessible internally	f	2025-11-28 16:05:10.277279+00	\N
96935064-7dd3-4365-8060-47f065d83667	INFO	DB Port accessible internally	f	2025-11-28 16:06:10.279491+00	\N
f6823c34-3299-408f-ae4d-aa3cbd941858	INFO	DB Port accessible internally	f	2025-11-28 16:07:10.370842+00	\N
17a425b2-2389-4fb2-ac67-1901cfd5398e	INFO	DB Port accessible internally	f	2025-11-28 16:08:10.267054+00	\N
5cb43524-6b78-47c6-a49f-028cab9e2209	INFO	DB Port accessible internally	f	2025-11-28 16:09:10.233983+00	\N
f8f0bc77-9e8c-41bd-9b63-a7bf5041b205	INFO	DB Port accessible internally	f	2025-11-28 16:10:10.29037+00	\N
e8fcfad4-e1af-4a54-9868-1cb5c2c5572f	INFO	DB Port accessible internally	f	2025-11-28 16:11:10.279599+00	\N
9de705a9-1d06-4d0a-9c21-4de6cacb7340	INFO	DB Port accessible internally	f	2025-11-28 16:12:10.306407+00	\N
a6323a46-a4e5-4da8-aa40-9745f59b41fa	INFO	DB Port accessible internally	f	2025-11-28 16:58:00.006492+00	\N
cbbb420c-237e-4e50-9e00-95040a28004c	INFO	DB Port accessible internally	f	2025-11-28 16:59:00.017881+00	\N
03da0a78-1c13-418d-b717-58306a3309c5	INFO	DB Port accessible internally	f	2025-11-28 17:00:00.032357+00	\N
1b5ddca6-c7d1-40e1-bd24-1fc7aeef72ff	INFO	DB Port accessible internally	f	2025-11-28 17:01:00.053892+00	\N
b2c93355-a234-4401-8a9f-892bc057a4d4	INFO	DB Port accessible internally	f	2025-11-28 17:02:00.033318+00	\N
15d9657c-5ed6-476c-a79f-9a5456f5c0e8	INFO	DB Port accessible internally	f	2025-11-28 17:03:00.022098+00	\N
e803ef10-20f3-4cca-b853-29d5034b7fc9	INFO	DB Port accessible internally	f	2025-11-28 17:04:00.02098+00	\N
80e151c7-6d0f-440b-ae35-dfb1cb6d1cb8	INFO	DB Port accessible internally	f	2025-11-28 17:05:00.039402+00	\N
a69e941d-3274-4473-80fd-a8072fc08de7	INFO	DB Port accessible internally	f	2025-11-28 17:06:00.031463+00	\N
62481418-9532-488b-817a-5ba85313f517	INFO	DB Port accessible internally	f	2025-11-28 17:07:00.031032+00	\N
e0c6c17b-3fbe-4beb-adb5-123780854aaf	INFO	DB Port accessible internally	f	2025-11-28 17:08:00.037428+00	\N
187c5e81-e40e-469c-a546-2eef9a5da771	INFO	DB Port accessible internally	f	2025-11-28 17:09:00.007009+00	\N
192e879c-465f-4107-8ed2-2c5d4eec84cb	INFO	DB Port accessible internally	f	2025-11-28 17:10:00.045367+00	\N
234c17f4-f871-4324-a43f-013411cb3862	INFO	DB Port accessible internally	f	2025-11-28 17:11:00.037828+00	\N
1c8a0f25-3037-4e7b-bd97-9c6c14226be0	INFO	DB Port accessible internally	f	2025-11-28 17:12:01.324491+00	\N
0326e56f-de92-4215-9118-1bf8ecc0ee8e	INFO	DB Port accessible internally	f	2025-11-28 17:13:01.403006+00	\N
8665503f-f43f-4776-b21d-d3a378a81fa7	INFO	DB Port accessible internally	f	2025-11-28 17:14:01.372131+00	\N
f2369e1b-f680-4dae-bc6f-f7b8d0214eee	INFO	DB Port accessible internally	f	2025-11-28 17:15:01.370922+00	\N
8eb0266d-eca7-4271-830a-aa774f375aba	INFO	DB Port accessible internally	f	2025-11-28 17:16:01.36138+00	\N
a20e3dd7-7d32-4839-98a4-860c6263d14c	INFO	DB Port accessible internally	f	2025-11-28 17:17:01.373268+00	\N
e525bef6-101a-42bb-9f8d-c3e21caf3642	INFO	DB Port accessible internally	f	2025-11-28 17:18:01.375456+00	\N
06cd306b-e2ef-49dd-9e86-381a716b78b2	INFO	DB Port accessible internally	f	2025-11-28 17:19:01.353882+00	\N
c734c516-a048-4010-9951-b9426b255729	INFO	DB Port accessible internally	f	2025-11-28 17:20:01.374594+00	\N
9884f816-9c72-485f-9976-303f296cea1c	INFO	DB Port accessible internally	f	2025-11-28 17:21:01.372244+00	\N
d639a099-2163-4c1a-8f26-f27b4beca22f	INFO	DB Port accessible internally	f	2025-11-28 17:22:01.379525+00	\N
1916d2aa-3cb3-4547-8a21-7eedac1ad1d9	INFO	DB Port accessible internally	f	2025-11-28 17:22:58.789802+00	\N
c3ac462c-9b8e-4a51-b2b8-1453ca4602fe	INFO	DB Port accessible internally	f	2025-11-28 17:23:58.805994+00	\N
332bdee2-d31f-40b0-b03a-7f6052b925e3	INFO	DB Port accessible internally	f	2025-11-28 17:24:58.790212+00	\N
9dc4976c-c29e-40e3-8dc9-add2b9c02ed7	INFO	DB Port accessible internally	f	2025-11-28 17:25:58.833135+00	\N
54cac20b-e78d-4032-939a-a5a4e595c95e	INFO	DB Port accessible internally	f	2025-11-28 17:26:58.83485+00	\N
d1bdb066-c8a6-40c8-aba5-cfae5640bf80	INFO	DB Port accessible internally	f	2025-11-28 17:27:58.829254+00	\N
1d2b7b7e-1e67-4ffe-bc12-b733ea57977f	INFO	DB Port accessible internally	f	2025-11-28 17:28:58.794312+00	\N
ef92ac01-9d90-4534-84de-099f6938c022	INFO	DB Port accessible internally	f	2025-11-28 17:29:58.85735+00	\N
10364cbe-7012-439c-9fa0-892dce61c3c1	INFO	DB Port accessible internally	f	2025-11-28 17:30:58.828188+00	\N
af6eb77b-011d-4b49-986d-0ebcba03ce03	INFO	DB Port accessible internally	f	2025-11-28 17:31:58.834098+00	\N
4923e62d-fe1e-4c89-b12a-7eec3d5c486d	INFO	DB Port accessible internally	f	2025-11-29 12:24:07.656354+00	\N
547f974c-132d-47d5-9b0f-63ae2a9f561c	INFO	DB Port accessible internally	f	2025-11-29 12:25:07.647015+00	\N
0c1798d4-10c0-44c2-8e20-3d4416e4be25	INFO	DB Port accessible internally	f	2025-11-29 12:26:07.651886+00	\N
d0b745b7-615f-4ce4-9980-08867d78270a	INFO	DB Port accessible internally	f	2025-11-29 12:27:07.66239+00	\N
6428a339-e306-4e62-8b4d-764d2347449a	INFO	DB Port accessible internally	f	2025-11-29 12:28:07.725345+00	\N
e53b04c3-6aca-45af-8a41-3552d46a2c4f	INFO	DB Port accessible internally	f	2025-11-29 12:29:07.692356+00	\N
813ef1cc-f4c1-48dd-bcb1-20eb3db390e7	INFO	DB Port accessible internally	f	2025-11-29 12:30:07.65113+00	\N
ef0f0d2c-c89a-41e6-b264-95b3977cade2	INFO	DB Port accessible internally	f	2025-11-29 12:31:07.685697+00	\N
1c60fd88-267f-49d7-a98b-c0883c1de7fe	INFO	DB Port accessible internally	f	2025-11-29 12:32:07.670878+00	\N
e755f21a-b6b1-4242-87fc-cc92fe2b4cc5	INFO	DB Port accessible internally	f	2025-11-29 12:35:07.706879+00	\N
5b4adb96-3411-45e5-979e-166934e8f685	INFO	DB Port accessible internally	f	2025-11-29 12:37:07.678966+00	\N
b46ca5df-625c-4185-b617-afe66f9bd55d	INFO	DB Port accessible internally	f	2025-11-29 12:38:07.683909+00	\N
249a2211-fa1a-4e5a-b48b-04a080873142	INFO	DB Port accessible internally	f	2025-11-29 12:39:07.726079+00	\N
ae9d76f1-08b4-40dd-9045-2771b6e4144f	INFO	DB Port accessible internally	f	2025-11-29 12:40:07.706466+00	\N
bd5683ab-79b8-435f-9973-dda60acf50c4	INFO	DB Port accessible internally	f	2025-11-29 12:41:07.752159+00	\N
21df8740-b445-4ba2-ae37-3e06eb20f682	INFO	DB Port accessible internally	f	2025-11-29 12:42:07.710536+00	\N
a802bf4c-e2cf-477e-a349-3763ea1e793f	INFO	DB Port accessible internally	f	2025-11-29 12:43:07.664299+00	\N
99ee14ba-1abe-4306-9c73-3b549a22dac7	INFO	DB Port accessible internally	f	2025-11-29 12:44:07.689555+00	\N
3c3919ee-ce02-4a16-805b-c6e83537a14f	INFO	DB Port accessible internally	f	2025-11-29 12:45:07.658701+00	\N
c7d56b19-7bc5-4dd2-b21a-52c533c3cbad	INFO	DB Port accessible internally	f	2025-11-29 12:46:07.695224+00	\N
980d4429-a96e-40a6-86d2-132490acf961	INFO	DB Port accessible internally	f	2025-11-29 12:47:07.652358+00	\N
81faf835-2ef3-4e0f-9dcf-68e1bd249a62	INFO	DB Port accessible internally	f	2025-11-29 12:48:07.708504+00	\N
f6ca5fcc-0915-4226-913d-00045c84f537	INFO	DB Port accessible internally	f	2025-11-29 12:49:07.673048+00	\N
22d96400-54b1-424d-9fb4-0bd8f0c54101	INFO	DB Port accessible internally	f	2025-11-29 12:50:07.654573+00	\N
93b8d694-a0ed-4e81-b471-52af6627280e	INFO	DB Port accessible internally	f	2025-11-29 12:51:07.700685+00	\N
34b4cbb2-6ab2-4c05-b273-9e67a2ba692a	INFO	DB Port accessible internally	f	2025-11-29 12:52:07.652997+00	\N
ca8e5346-f1fd-4400-9be1-640f41ed36e8	INFO	DB Port accessible internally	f	2025-11-29 12:53:07.663105+00	\N
f8c963c0-c1dc-44b5-b79d-51028a8692f7	INFO	DB Port accessible internally	f	2025-11-29 12:54:07.644857+00	\N
8d7b904e-b06d-42ee-baaa-bbf0459f9f96	INFO	DB Port accessible internally	f	2025-11-29 12:55:07.683686+00	\N
94df1bc8-1b37-47c2-9eff-a8023797f146	INFO	DB Port accessible internally	f	2025-11-29 12:56:07.650113+00	\N
6aef40d6-5324-48bc-9dc8-2f46505d7b75	INFO	DB Port accessible internally	f	2025-11-29 12:57:07.71635+00	\N
880ebbce-f7c7-4d0e-bdd4-d3cc1e3ee255	INFO	DB Port accessible internally	f	2025-11-29 12:58:07.654492+00	\N
a6cf1463-e9b3-483e-9796-4cd122951ec8	INFO	DB Port accessible internally	f	2025-11-29 12:59:07.683277+00	\N
c43c01d6-1538-4ee8-9997-e2822d43bc54	INFO	DB Port accessible internally	f	2025-11-29 13:00:07.66661+00	\N
cc46a0e1-f6df-4fea-9ba3-4542196a3590	INFO	DB Port accessible internally	f	2025-11-29 13:01:07.674639+00	\N
e75e2394-fa0d-41e9-b854-7ae5d76651f3	INFO	DB Port accessible internally	f	2025-11-29 13:18:07.726577+00	\N
38e8c95a-bf46-4845-bb0a-9f33a3ee5745	INFO	DB Port accessible internally	f	2025-11-29 13:19:07.69152+00	\N
0de2e548-f137-4c8c-9114-ade7839f1717	INFO	DB Port accessible internally	f	2025-11-29 13:20:07.657621+00	\N
a2a763c2-1159-480c-a69a-669ee3ac15b5	INFO	DB Port accessible internally	f	2025-11-29 13:21:07.662403+00	\N
6d197cb2-17af-44b8-989c-9f3895cf4070	INFO	DB Port accessible internally	f	2025-11-29 13:22:07.675217+00	\N
769c408f-00ca-4c2b-a2e2-de04f8512974	INFO	DB Port accessible internally	f	2025-11-29 13:23:07.668292+00	\N
76396ef2-e214-417f-a8fa-4bc54a6ae91b	INFO	DB Port accessible internally	f	2025-11-29 13:24:07.688165+00	\N
d777cf8c-a0fb-443d-acff-0f138051e3ec	INFO	DB Port accessible internally	f	2025-11-29 13:25:07.677958+00	\N
a00d172f-be19-41a9-8736-882afe4fb368	INFO	DB Port accessible internally	f	2025-11-29 13:26:07.660021+00	\N
cfac122e-bef5-4304-ac1f-dd63fc57efda	INFO	DB Port accessible internally	f	2025-11-29 13:27:07.686301+00	\N
c126e6cb-57df-4713-952c-85021e190995	INFO	DB Port accessible internally	f	2025-11-29 13:28:07.694592+00	\N
1703842d-4339-4ac9-abc8-a5badb93dd85	INFO	DB Port accessible internally	f	2025-11-29 13:29:07.704649+00	\N
f5aa75d5-6345-4588-960d-9f0d95f05978	INFO	DB Port accessible internally	f	2025-11-29 13:30:07.6933+00	\N
bdf7ebfb-8ad8-4631-9ec8-36594e636666	INFO	DB Port accessible internally	f	2025-11-29 13:31:07.702462+00	\N
f801a0dd-0cdc-4de6-b0a9-a16ae4309f09	INFO	DB Port accessible internally	f	2025-11-29 13:32:07.651393+00	\N
964ce7d3-ffed-4413-99e0-4212a179acc5	INFO	DB Port accessible internally	f	2025-11-29 13:33:07.671095+00	\N
\.


--
-- Data for Name: system_metrics; Type: TABLE DATA; Schema: public; Owner: nas_user
--

COPY public.system_metrics (id, agent_id, cpu_usage, ram_usage, disk_usage, created_at) FROM stdin;
ea6ed4ad-6f8a-4345-bcb4-2bafcb8732ac	prod-agent	17.60	35.72	51.06	2025-11-27 20:11:28.212309+00
66617f33-4163-46a3-aca8-d1b317133d24	prod-agent	13.18	35.79	51.06	2025-11-27 20:11:38.211408+00
e03f4d9e-41e2-4fad-9056-a8c6cfb2ade9	prod-agent	9.50	35.89	51.06	2025-11-27 20:11:48.213014+00
5597667e-31dc-480e-9df1-7c645d581570	prod-agent	8.20	35.79	51.06	2025-11-27 20:11:58.220819+00
801d53f2-aa68-4ee3-8407-556eca6655d6	prod-agent	6.73	35.74	51.06	2025-11-27 20:12:08.208671+00
70369534-9a3b-446e-8661-fe4446730a94	prod-agent	10.38	35.87	51.06	2025-11-27 20:12:18.206262+00
9055f291-091a-457d-ab1f-084b15315f7a	prod-agent	10.75	35.88	51.06	2025-11-27 20:12:28.212195+00
76356a97-1a3f-4fee-bd3a-b803e1cce7ee	prod-agent	100.00	35.91	51.06	2025-11-27 20:12:52.877655+00
b7a3ada2-892c-4866-a5fb-b1e9fd315a87	prod-agent	12.17	35.59	51.06	2025-11-27 20:13:02.871869+00
2324de72-ab8f-4953-8e8f-39d03fd727bb	prod-agent	15.28	35.71	51.06	2025-11-27 20:13:12.872395+00
1b499293-3634-4267-90f2-8075f81d29f4	prod-agent	12.72	35.67	51.06	2025-11-27 20:13:22.86628+00
6ecfdfa8-916a-4970-8932-37209f59126c	prod-agent	10.57	35.71	51.06	2025-11-27 20:13:32.871857+00
b7f3e4d6-a418-4edb-a314-9d550cc9949f	prod-agent	18.31	35.79	51.06	2025-11-27 20:13:42.877155+00
02d6eb2c-9eee-40ba-8cf9-e1e1ebc7e259	prod-agent	12.47	35.61	51.06	2025-11-27 20:13:52.86497+00
982f3e65-6b6f-4621-884b-4a072a3769ba	prod-agent	16.47	35.68	51.06	2025-11-27 20:14:02.868894+00
8b686801-4d59-40c3-9839-b2be4de7c61c	prod-agent	21.03	35.97	51.06	2025-11-27 20:14:12.881525+00
f8084009-17dc-4dfd-a2c8-39f9a707f56a	prod-agent	28.49	35.92	51.06	2025-11-27 20:14:22.863873+00
fab55d00-8d84-4a6b-9a5b-719c1679e21f	prod-agent	13.04	35.86	51.06	2025-11-27 20:14:32.870962+00
67b4aec2-7266-4294-a9aa-4ddab7c63330	prod-agent	37.68	37.31	51.06	2025-11-27 20:14:42.872201+00
34b668a6-d757-4678-8b84-f1556270434c	prod-agent	54.90	39.58	51.06	2025-11-27 20:14:52.864609+00
4f50a879-a7ea-406a-ad58-70064caf15f1	prod-agent	51.84	40.41	51.06	2025-11-27 20:15:02.871984+00
e416eaac-945d-4f57-a9d6-3c38fa390861	prod-agent	12.89	36.52	51.06	2025-11-27 20:15:12.878996+00
53c20813-fc0c-403b-b56b-91abb5b4f893	prod-agent	9.39	36.07	51.06	2025-11-27 20:15:22.865006+00
aea72beb-4669-400f-b665-cadd1cc01db2	prod-agent	8.08	36.03	51.06	2025-11-27 20:15:32.871484+00
1c6055d0-3879-4586-a242-2c3299f94663	prod-agent	9.40	35.97	51.06	2025-11-27 20:15:42.872156+00
e331929e-e82f-4c00-a8e8-879b2c8038cc	prod-agent	13.82	35.99	51.06	2025-11-27 20:15:52.871751+00
cb02a755-4860-4da1-b293-1b9bd76a9838	prod-agent	7.30	36.07	51.06	2025-11-27 20:16:02.869562+00
adb0d3f5-3e46-4e41-90d4-9a9e52638603	prod-agent	9.90	35.90	51.06	2025-11-27 20:16:12.872371+00
293431d1-2d9a-482d-9e4f-248476941c9b	prod-agent	9.80	35.93	51.06	2025-11-27 20:16:22.874717+00
3c6df441-1259-4004-a8ba-985a0d631b61	prod-agent	11.98	35.88	51.06	2025-11-27 20:16:32.876514+00
4d573ee6-06cb-4cb6-9749-07279128fd32	prod-agent	8.35	35.95	51.06	2025-11-27 20:16:42.866604+00
f1f18372-c098-49ee-af8f-3fbec76fa709	prod-agent	16.60	35.68	51.06	2025-11-27 20:16:52.870031+00
e83e5ee3-7db6-4e02-9bb8-1de58130db12	prod-agent	8.41	35.62	51.06	2025-11-27 20:17:02.872175+00
85509593-3003-4b01-9155-3dcad5eb036e	prod-agent	10.46	35.65	51.06	2025-11-27 20:17:12.871232+00
4e305501-649b-4df8-8c3d-e5fb7cd179e0	prod-agent	11.42	35.60	51.06	2025-11-27 20:17:22.871425+00
471067ff-6460-4a91-a432-69da9ff23f20	prod-agent	10.50	35.63	51.06	2025-11-27 20:17:32.86562+00
148cf030-9fd1-4089-b3e7-9fe472629129	prod-agent	9.24	35.60	51.06	2025-11-27 20:17:42.871496+00
1fe15e6d-23d4-4269-b5c4-b371121ce060	prod-agent	13.48	35.65	51.06	2025-11-27 20:17:52.87077+00
9b999b5f-acb1-4290-bec9-1d47d3bbfe80	prod-agent	7.00	35.62	51.06	2025-11-27 20:18:02.872426+00
171447fa-37a6-49b6-9849-4e7b13a674d3	prod-agent	9.62	35.62	51.06	2025-11-27 20:18:12.865639+00
95f08ecd-8b2c-46b0-88fb-7f7782019b13	prod-agent	30.46	35.37	51.06	2025-11-27 20:18:22.86995+00
98f4af36-3dc8-4be2-8701-adbf9705f70d	prod-agent	5.85	35.40	51.06	2025-11-27 20:18:32.874114+00
d27f56ff-f847-41ec-8ff0-e36b62138de2	prod-agent	8.41	35.42	51.06	2025-11-27 20:18:42.87463+00
d2e2bedf-368d-434f-b807-fc1a88abf1bd	prod-agent	14.99	35.39	51.06	2025-11-27 20:18:52.865587+00
b989f63a-3a35-4c7e-9e04-589aa333e85b	prod-agent	7.39	35.24	51.06	2025-11-27 20:19:02.869155+00
0628c966-9a14-44a6-ae4d-722754666b68	prod-agent	10.10	35.26	51.06	2025-11-27 20:19:12.866027+00
b833f622-96ba-4cd2-8246-bd1bab49d11f	prod-agent	9.78	35.29	51.06	2025-11-27 20:19:22.868853+00
140a2420-4945-4f9d-a372-ebc8c91aa338	prod-agent	7.65	35.28	51.06	2025-11-27 20:19:32.872751+00
988b5e65-184c-4bea-bc5b-fec26e74b013	prod-agent	8.77	35.21	51.06	2025-11-27 20:19:42.870265+00
070a5215-bc7a-4219-bcb6-ac8b87ff8d98	prod-agent	12.53	35.30	51.06	2025-11-27 20:19:52.869697+00
bd7544d7-859b-40b5-af4e-0dd70a3f85bb	prod-agent	8.58	35.27	51.06	2025-11-27 20:20:02.871966+00
ea36b839-87d4-46ce-88bf-05c7fdaa556f	prod-agent	9.78	35.31	51.06	2025-11-27 20:20:12.873254+00
279003fa-6f0f-460b-a3fd-2312730493d5	prod-agent	9.68	35.30	51.06	2025-11-27 20:20:22.866098+00
b89e35fb-9dec-4a09-b039-3890086c3820	prod-agent	7.23	35.28	51.06	2025-11-27 20:20:32.920654+00
5b2f9b68-a47b-45da-aae6-1fc90bd0780c	prod-agent	9.25	35.32	51.06	2025-11-27 20:20:42.87266+00
2c77cd12-be28-4dcc-a428-0cdc18005348	prod-agent	9.14	35.29	51.06	2025-11-27 20:20:52.870666+00
da5f3b8a-74dc-4fa1-8fb6-b94d34f95483	prod-agent	12.62	35.30	51.06	2025-11-27 20:21:02.872511+00
6229bd9e-00e2-4ddd-9c8f-eb6b0c5e8559	prod-agent	9.12	16.03	51.11	2025-11-28 10:54:29.990337+00
737eda1f-0f25-4b3d-b4f7-3cdccd5f2cb5	prod-agent	13.23	15.26	51.11	2025-11-28 10:54:39.952211+00
01e94e91-1074-4207-b0d3-7a1673848523	prod-agent	19.62	15.05	51.11	2025-11-28 10:54:49.958432+00
eea244c8-c21d-4251-92d4-cebeab10fdc9	prod-agent	8.55	15.09	51.11	2025-11-28 10:54:59.962108+00
28c4fa52-758b-442f-9570-bb5f34d8631c	prod-agent	12.58	15.01	51.11	2025-11-28 10:55:09.959527+00
0bfb6b78-5128-4d87-b111-27150ca46f18	prod-agent	11.12	15.01	51.11	2025-11-28 10:55:19.958353+00
a65fe04b-e4ee-426f-9520-947ba2997d00	prod-agent	9.46	15.03	51.11	2025-11-28 10:55:29.958696+00
407b0f5a-27f0-4048-8cc4-9332af495705	prod-agent	8.80	15.07	51.11	2025-11-28 10:55:39.960553+00
a802d373-8047-4184-b196-8126281ba2d6	prod-agent	14.77	15.08	51.11	2025-11-28 10:55:49.957313+00
2789fd6d-484b-44bb-9de5-dfcc9980d924	prod-agent	7.13	15.08	51.11	2025-11-28 10:55:59.95345+00
d6b49db2-8775-4680-a776-5a80a60a2585	prod-agent	9.72	15.17	51.11	2025-11-28 10:56:09.958393+00
7054b388-dc77-42e5-a680-a5313714f1c7	prod-agent	10.35	15.15	51.11	2025-11-28 10:56:19.95815+00
450b852e-744d-445c-9ca3-77cf2d8ec6c1	prod-agent	8.10	15.09	51.11	2025-11-28 10:56:29.957098+00
eaf1f5f4-9f4b-4703-b0dc-f1e4fd8bdd79	prod-agent	10.43	15.02	51.11	2025-11-28 10:56:39.954352+00
21bf38bd-2a07-4303-bfa2-28bd418e35cd	prod-agent	15.72	15.11	51.11	2025-11-28 10:56:49.957992+00
62f0fe10-2aff-47e0-b3aa-e9324a3c4af4	prod-agent	8.03	15.07	51.11	2025-11-28 10:56:59.957648+00
e5354a05-7b31-4c8e-aacd-ba7cbefdb2ae	prod-agent	9.20	15.11	51.11	2025-11-28 10:57:09.958199+00
a66354c1-cc21-4e88-ae7e-b77a1d3a68c1	prod-agent	9.62	15.09	51.11	2025-11-28 10:57:19.961995+00
9d12f61b-9775-4466-8327-5d0240bcbfc9	prod-agent	6.53	15.05	51.11	2025-11-28 10:57:29.955576+00
f1c904ac-2985-4913-b092-cbbf64903ae8	prod-agent	9.16	15.10	51.11	2025-11-28 10:57:39.956045+00
014f4986-dda2-43e0-b9fa-89c53e403db2	prod-agent	15.10	15.20	51.11	2025-11-28 10:57:49.952821+00
5c061b77-7d5c-45c9-ac30-154f19da1f3c	prod-agent	6.75	15.07	51.11	2025-11-28 10:57:59.959239+00
8c17b9cd-906c-482f-8b5d-57cd845b353c	prod-agent	9.10	15.12	51.11	2025-11-28 10:58:09.952143+00
f9993703-6434-45c2-bdc7-d425b180c1d2	prod-agent	8.98	15.11	51.11	2025-11-28 10:58:19.953594+00
235663f5-6c1d-437c-9adf-08af5a99c7a5	prod-agent	7.85	15.12	51.11	2025-11-28 10:58:29.958545+00
f4ff75dc-8ab3-4f17-99ba-b24026873b7a	prod-agent	9.42	15.11	51.11	2025-11-28 10:58:39.961384+00
b78a3f3b-ab02-43f4-acd8-037d8d8effbf	prod-agent	11.08	15.07	51.11	2025-11-28 10:58:49.961242+00
76543e32-0158-489b-89c9-df5e238ec4cb	prod-agent	11.44	15.15	51.11	2025-11-28 10:58:59.95879+00
43ed58f2-6fa2-42c1-a26f-14a64e31d877	prod-agent	9.46	15.08	51.11	2025-11-28 10:59:09.95166+00
f1c33258-66d8-474b-b718-e1a819865e62	prod-agent	9.32	15.05	51.11	2025-11-28 10:59:19.958786+00
b0827fd0-018d-43cb-93b1-1d332b556161	prod-agent	7.06	15.08	51.11	2025-11-28 10:59:29.951991+00
1ec566fc-842d-4618-8ea8-970e3dea7fe1	prod-agent	27.03	15.45	51.11	2025-11-28 10:59:39.95716+00
90133ada-b070-4da7-891b-cd0f5021bc2a	prod-agent	37.23	15.57	51.11	2025-11-28 10:59:49.960248+00
5496056f-d490-4150-8f24-f3c816a8f927	prod-agent	38.20	15.69	51.11	2025-11-28 10:59:59.95724+00
1cd3daa3-ea95-4b3b-8fab-9d2634e856fe	prod-agent	34.65	16.17	51.11	2025-11-28 11:00:09.958066+00
98ca2c85-5fb3-4d5e-99fc-3fe3c0069f2f	prod-agent	31.64	16.22	51.11	2025-11-28 11:00:19.960408+00
1276685a-20c8-47fc-9abb-458b1883ef70	prod-agent	11.69	16.28	51.11	2025-11-28 11:00:29.955213+00
7fdc9a07-74c2-4d15-93e4-c76d61908311	prod-agent	9.87	16.29	51.11	2025-11-28 11:00:40.004893+00
e65fa528-222f-46af-945d-78db0be5f068	prod-agent	10.44	16.26	51.11	2025-11-28 11:00:49.958095+00
6ee65d64-987f-4a9d-b3a6-bf76351759b1	prod-agent	25.68	16.15	51.11	2025-11-28 11:00:59.956533+00
69b9c007-7929-4517-b67d-bf42e4398ac6	prod-agent	7.59	15.93	51.11	2025-11-28 11:01:09.956868+00
aaa2f2b8-9e2d-48a9-ad11-bf6b1aa1f338	prod-agent	9.31	15.88	51.11	2025-11-28 11:01:19.961662+00
08e2a1cb-722b-4253-a630-08ff3864b146	prod-agent	7.95	15.86	51.11	2025-11-28 11:01:29.959002+00
7014e7ae-e67f-4fe6-b326-d77e197d36a5	prod-agent	8.26	15.88	51.11	2025-11-28 11:01:39.95853+00
a78f2410-c9aa-405e-9161-2836a202b72a	prod-agent	8.74	15.92	51.11	2025-11-28 11:01:49.959109+00
99c1897e-96d0-4ffc-a975-2ea6715ccfe4	prod-agent	20.55	15.68	51.11	2025-11-28 11:01:59.95738+00
665a9272-8986-4d89-9a46-8d23507ab003	prod-agent	8.23	15.60	51.11	2025-11-28 11:02:09.959218+00
3a5c6161-6936-41d0-aa40-057b517938a5	prod-agent	8.83	15.49	51.11	2025-11-28 11:02:19.958195+00
275e3cad-916f-4b1d-a304-7f0c9024e287	prod-agent	8.03	15.55	51.11	2025-11-28 11:02:29.954608+00
894c591e-7d29-4461-8aef-d73addd89b88	prod-agent	8.77	15.45	51.11	2025-11-28 11:02:39.965064+00
c503f65c-ac1e-4e5a-a5f9-5befe1bd7eb9	prod-agent	7.85	15.51	51.11	2025-11-28 11:02:49.959126+00
b47a993c-a5ba-4c03-8a11-cb41bff854c2	prod-agent	12.55	15.53	51.11	2025-11-28 11:02:59.958838+00
586f9b9d-1aa4-4949-85a1-84eecf178346	prod-agent	10.48	15.55	51.11	2025-11-28 11:03:09.95119+00
1715bd7b-0067-46e4-9c67-3cdd0fbcff63	prod-agent	9.44	15.54	51.11	2025-11-28 11:03:19.960288+00
506e4807-72d5-4d43-a987-e0728a28f3fa	prod-agent	12.56	13.17	51.11	2025-11-28 11:03:29.970503+00
8f251fdc-e861-45fc-92ff-c7deda2fc704	prod-agent	19.65	14.43	51.11	2025-11-28 11:03:39.960285+00
d579d982-44be-4b3a-8eab-09ff660180c4	prod-agent	37.82	15.85	51.11	2025-11-28 11:03:49.959525+00
2e4fce10-f816-4551-8165-b8dae52f5e20	prod-agent	17.43	15.54	51.11	2025-11-28 11:03:59.960809+00
f8245087-dace-4bc9-9389-4dfb7c847bce	prod-agent	37.18	15.53	51.11	2025-11-28 11:04:09.950325+00
312823ac-0b2e-4a0e-9f0c-50bdf6e2cb67	prod-agent	8.14	15.54	51.11	2025-11-28 11:04:19.957897+00
ebddb72f-8ea9-4fbc-b0cb-ca9ebd4fe4d1	prod-agent	8.29	15.50	51.11	2025-11-28 11:04:29.963017+00
97a8da97-709e-4ed3-add0-fd7c6794571c	prod-agent	8.79	15.49	51.11	2025-11-28 11:04:39.955774+00
cd32a664-434d-4702-a017-392b5820f92a	prod-agent	15.55	14.67	51.11	2025-11-28 11:04:49.951628+00
08c69858-da10-46e5-a1a7-39a944cc61b8	prod-agent	7.87	14.58	51.11	2025-11-28 11:04:59.957038+00
1d7d2b77-e33c-4ed7-bc41-ebc8bc84cbe3	prod-agent	15.09	14.65	51.11	2025-11-28 11:05:09.955402+00
43de7731-f5a7-48f8-af1e-f355636a1e69	prod-agent	8.92	14.67	51.11	2025-11-28 11:05:19.959221+00
95b1fb0f-4c8b-4ad3-955e-0b1580eadfa1	prod-agent	8.03	14.67	51.11	2025-11-28 11:05:29.95959+00
15918432-07ea-4305-84c6-5593e55839d6	prod-agent	10.05	14.71	51.11	2025-11-28 11:05:39.952+00
5ebcc02b-40a9-481d-9ac1-c592f2aa566b	prod-agent	8.33	14.67	51.11	2025-11-28 11:05:50.010596+00
f1824342-9a5c-4131-af59-1221c278618e	prod-agent	8.71	14.70	51.11	2025-11-28 11:05:59.959636+00
c142da63-7df6-4f09-9e4c-88572bc0b652	prod-agent	14.47	14.75	51.11	2025-11-28 11:06:09.952649+00
aff6c9e2-5e11-4cde-ba42-3abbe120c489	prod-agent	8.52	14.69	51.11	2025-11-28 11:06:19.959148+00
d154476d-aa39-4766-a48b-2f039573c3c6	prod-agent	9.35	14.68	51.11	2025-11-28 11:06:29.957845+00
718b8458-c7c6-40ed-8ec5-36c45e331dd6	prod-agent	11.64	14.64	51.11	2025-11-28 11:06:39.951637+00
93dc060f-ead5-424c-8932-a567e1df9946	prod-agent	9.10	14.66	51.11	2025-11-28 11:06:49.956771+00
4a8d86d4-5bc5-4f6b-b455-4234a25a7e69	prod-agent	9.59	14.59	51.11	2025-11-28 11:06:59.95314+00
69ad4e15-074b-4f44-b48c-e49924bd9de0	prod-agent	15.43	14.65	51.11	2025-11-28 11:07:09.959463+00
e2be8e89-4f40-4962-ad78-d27040bb3aaa	prod-agent	7.44	14.62	51.11	2025-11-28 11:07:19.965214+00
ecc813a5-b24c-4601-9502-67ed7db07156	prod-agent	8.86	14.69	51.11	2025-11-28 11:07:29.955513+00
c62664c7-8f0a-408e-92bd-1fd936421d06	prod-agent	8.18	14.64	51.11	2025-11-28 11:07:39.959872+00
87cb6bd6-739a-4f64-8c04-16ac96654bf6	prod-agent	9.15	14.63	51.11	2025-11-28 11:07:49.958313+00
0a694324-2797-4b8e-8533-43cfd55c558e	prod-agent	7.82	14.64	51.11	2025-11-28 11:07:59.962014+00
d9a538d0-8abd-4db4-8fa9-92864d238d33	prod-agent	13.21	14.65	51.11	2025-11-28 11:08:09.956241+00
9ad7759a-8301-4679-8ad1-ff15f05cf7f9	prod-agent	9.83	14.64	51.11	2025-11-28 11:08:19.958363+00
8597b321-0218-4b32-9d13-a822401779b2	prod-agent	8.32	14.71	51.11	2025-11-28 11:08:29.952369+00
a8b27372-df46-4cd7-9949-1d3511413729	prod-agent	8.59	14.56	51.11	2025-11-28 11:08:39.952978+00
fcaba5ef-4d6c-4d81-8061-d0a2d6759b3c	prod-agent	9.22	14.54	51.11	2025-11-28 11:08:49.953267+00
cad67d4c-3aed-4cbc-a6e3-36dc89da3144	prod-agent	8.49	14.59	51.11	2025-11-28 11:08:59.957877+00
3fc66a94-e132-4dbf-9c54-11c1d06dcae5	prod-agent	9.24	14.59	51.11	2025-11-28 11:09:09.961576+00
3fbdd0b9-e090-46b8-9a2d-a7bf13ad56e7	prod-agent	13.76	14.64	51.11	2025-11-28 11:09:19.955537+00
e731ba9e-d3ae-4ccd-b190-2eceee482395	prod-agent	8.07	14.64	51.11	2025-11-28 11:09:29.96169+00
69af738d-0767-4de6-a718-c1420dbef760	prod-agent	8.42	14.57	51.11	2025-11-28 11:09:39.958568+00
9a4b7d3b-9252-43bc-ae21-d847952f874d	prod-agent	8.48	14.60	51.11	2025-11-28 11:09:49.95846+00
3abbdcb9-45f3-4715-9a7a-6fd39a8138e7	prod-agent	8.79	14.54	51.11	2025-11-28 11:09:59.961598+00
c40c90dc-fd05-4809-b7d6-d707db4100d6	prod-agent	10.40	14.57	51.11	2025-11-28 11:10:09.953188+00
4d590bd2-f6b1-476d-895f-0f5122345f2d	prod-agent	13.54	14.63	51.11	2025-11-28 11:10:19.957926+00
0325b4d1-ebb3-4b35-b7ca-5613907d2610	prod-agent	7.50	14.50	51.11	2025-11-28 11:10:29.956949+00
c6ab28dd-dc99-4e23-a052-ddb7cb9f0619	prod-agent	9.22	14.61	51.11	2025-11-28 11:10:39.955902+00
e7c1d9af-b1f4-4d52-bfcc-b2d6550a3b5b	prod-agent	8.65	14.58	51.11	2025-11-28 11:10:49.95227+00
5059191a-1110-4d8f-bb00-d547ffd38c90	prod-agent	8.44	14.54	51.11	2025-11-28 11:11:00.014599+00
888ba101-d082-41d7-9626-eaa17b3fbf2e	prod-agent	9.94	14.52	51.11	2025-11-28 11:11:09.960063+00
49a81f3d-484c-46fa-b4bc-d03a265546a7	prod-agent	14.47	14.59	51.11	2025-11-28 11:11:19.960035+00
0d934bf8-f05b-45d7-857a-7b88de3d2b8d	prod-agent	6.07	14.57	51.11	2025-11-28 11:11:29.959957+00
a7311b5f-7ef4-4909-9bb8-e8abfb2ae7c1	prod-agent	9.32	14.59	51.11	2025-11-28 11:11:39.958991+00
53e61350-2778-43d4-9073-3abd970997ae	prod-agent	7.23	14.63	51.11	2025-11-28 11:11:49.958569+00
4f43692f-7ab7-4cf4-8280-aed36d357895	prod-agent	8.33	14.65	51.11	2025-11-28 11:11:59.958647+00
6bf87d3d-31b2-45f9-aa4a-3e7060e147d7	prod-agent	9.39	14.61	51.11	2025-11-28 11:12:09.953325+00
37d89cdb-0839-46b1-b136-ea39bb2a44ff	prod-agent	14.77	14.63	51.11	2025-11-28 11:12:19.958983+00
f901aafd-3169-474f-bf64-8f503820b419	prod-agent	12.34	14.68	51.11	2025-11-28 11:12:29.959519+00
1c93d38f-9f04-46f3-9d53-72a3c7a25271	prod-agent	9.00	14.65	51.11	2025-11-28 11:12:39.957883+00
277f6ce1-fd86-4c78-b342-ea3130dc7402	prod-agent	8.61	14.62	51.11	2025-11-28 11:12:49.959576+00
de6f46ba-729d-406b-b141-a7dda091048f	prod-agent	10.08	14.70	51.11	2025-11-28 11:12:59.952469+00
7b101f3c-3cf5-4c21-8b13-22a7123528c4	prod-agent	8.72	14.70	51.11	2025-11-28 11:13:09.958049+00
c5c73998-633a-4cbb-9e5f-7e27c4090efa	prod-agent	8.98	14.63	51.11	2025-11-28 11:13:19.953278+00
fff6b537-af78-4b5a-af9a-fe3908099e89	prod-agent	12.72	14.63	51.11	2025-11-28 11:13:29.956848+00
54c09630-98cc-4c5e-8d0b-cc716a6f53ad	prod-agent	8.84	14.64	51.11	2025-11-28 11:13:39.958387+00
1c79c735-2176-468b-970d-ae56a0b98637	prod-agent	8.35	14.67	51.11	2025-11-28 11:13:49.952061+00
ef412da0-4d04-4a18-bcd6-f65c49cd0a50	prod-agent	7.77	14.62	51.11	2025-11-28 11:13:59.96599+00
0b9978e4-6921-4bcf-9ef9-44ac0e8404c1	prod-agent	11.05	14.62	51.11	2025-11-28 11:14:09.956807+00
3612cb32-ad88-469b-951b-ad08be54208a	prod-agent	9.28	14.55	51.11	2025-11-28 11:14:19.954438+00
0af47fbb-c21d-473e-b438-fc20bb26a500	prod-agent	14.89	14.59	51.11	2025-11-28 11:14:29.95196+00
25fd6135-1ddf-4d38-b721-1821399059ba	prod-agent	8.35	14.59	51.11	2025-11-28 11:14:39.958792+00
decfb77e-153c-488f-95e1-257c85d61e30	prod-agent	8.74	14.65	51.11	2025-11-28 11:14:49.95665+00
98b8f1c0-e854-42a8-90e5-b6b8ac6a0b78	prod-agent	9.38	14.61	51.11	2025-11-28 11:14:59.95751+00
25dbe0f1-af36-4651-848f-b84c7c4c4aed	prod-agent	8.98	14.66	51.11	2025-11-28 11:15:09.961691+00
003ccd3f-770d-4c11-8f71-181cebf16a39	prod-agent	8.82	14.61	51.11	2025-11-28 11:15:19.953527+00
3a05d79a-3f83-42c0-813d-417faedaa0e8	prod-agent	12.80	14.70	51.11	2025-11-28 11:15:29.952125+00
7ca9b169-6dd1-4638-a657-ae7d5e0938e8	prod-agent	9.16	14.68	51.11	2025-11-28 11:15:39.958256+00
46440b7b-9a71-4b81-9527-dc17a8d4c7de	prod-agent	8.77	14.63	51.11	2025-11-28 11:15:49.957377+00
ef63bfa4-aca3-4772-88a9-bdfa2ed9b5a5	prod-agent	17.88	15.04	51.11	2025-11-28 11:15:59.960089+00
b38fae66-c93a-4054-affd-947e0aab4a72	prod-agent	14.46	14.66	51.11	2025-11-28 11:16:09.976482+00
dac2cfab-c905-4c80-ba07-6c287d1ebd49	prod-agent	7.86	14.61	51.11	2025-11-28 11:16:19.954259+00
d0cc4549-8c46-46ca-8d22-664cc79c849e	prod-agent	14.50	14.56	51.11	2025-11-28 11:16:29.952588+00
08a7b6b8-8adc-4592-aec1-b1f5f4b6ba85	prod-agent	9.35	14.57	51.11	2025-11-28 11:16:39.959774+00
c5f406c5-5d0a-494e-ac54-1aca3853678f	prod-agent	9.35	14.65	51.11	2025-11-28 11:16:49.958406+00
3e6679e7-ce8d-4ab2-a4a1-9a395b76a17f	prod-agent	8.40	14.65	51.11	2025-11-28 11:16:59.956239+00
4d530034-5007-41a1-a9f4-1cea189ab6c1	prod-agent	8.68	14.65	51.11	2025-11-28 11:17:09.965853+00
91aac0b4-deca-4668-89e1-124f51bc2beb	prod-agent	8.00	14.61	51.11	2025-11-28 11:17:19.958285+00
f3454c15-0896-4577-a715-78150334e423	prod-agent	11.26	14.61	51.11	2025-11-28 11:17:29.952922+00
f31c4adb-6b5b-46c7-b76d-3a0261c6548d	prod-agent	10.71	14.64	51.11	2025-11-28 11:17:39.958017+00
3e3c8a84-9db3-4d6e-b9e4-c36994951532	prod-agent	7.75	14.58	51.11	2025-11-28 11:17:49.955263+00
a5799d9e-fd2c-4f3c-96c4-5f29da438d15	prod-agent	8.82	14.66	51.11	2025-11-28 11:17:59.959421+00
2e55540f-3bbf-4c70-8425-88e713841f70	prod-agent	11.19	14.58	51.11	2025-11-28 11:18:09.958019+00
f88ba54a-c525-41df-ab5a-38aba69339b2	prod-agent	7.41	14.65	51.11	2025-11-28 11:18:19.95883+00
2256fdf0-9955-4555-9bd2-f5de30efe594	prod-agent	9.16	14.66	51.11	2025-11-28 11:18:29.952555+00
d6031402-2b6e-4be2-ac2f-2d6bc3121bb4	prod-agent	12.99	14.58	51.11	2025-11-28 11:18:39.963477+00
a7e52f69-0bb2-4244-8fc3-385a95e35775	prod-agent	8.54	14.69	51.11	2025-11-28 11:18:49.951603+00
2c545f1c-c7a7-486c-bcb3-7a728a0c238c	prod-agent	8.90	14.63	51.11	2025-11-28 11:18:59.960145+00
34eef604-0685-4e48-a9de-afe974484719	prod-agent	8.28	14.58	51.11	2025-11-28 11:19:09.951565+00
898cb7a1-46fb-40bb-af56-5cea2ce510b5	prod-agent	8.11	14.63	51.11	2025-11-28 11:19:19.968192+00
7388d82c-8f9d-404c-bfc8-74c4164069be	prod-agent	8.21	14.68	51.11	2025-11-28 11:19:29.957842+00
c522485d-9d83-489b-a9ec-228b64a82bd3	prod-agent	15.06	14.69	51.11	2025-11-28 11:19:39.957417+00
afb664f1-d7a7-4087-b69c-38870f3d5d6a	prod-agent	7.67	14.67	51.11	2025-11-28 11:19:49.959284+00
2d4131bb-88d7-4b8b-a68d-fab89f34ed1b	prod-agent	8.31	14.70	51.11	2025-11-28 11:19:59.953948+00
ce01b3d8-7e61-4689-a249-81432589c912	prod-agent	10.25	14.73	51.11	2025-11-28 11:20:09.952188+00
6f13f77f-8346-48a5-ae6d-b1f70cae83d7	prod-agent	7.72	14.67	51.11	2025-11-28 11:20:19.958956+00
c58b073e-e79e-4be3-9fa6-c0b2a97d62e0	prod-agent	8.55	14.69	51.11	2025-11-28 11:20:29.95249+00
8326be12-8528-4181-b689-d739b455aaaa	prod-agent	14.09	14.72	51.11	2025-11-28 11:20:39.957637+00
8dcde7de-8cf1-4977-be55-82eb9fd5ef19	prod-agent	8.14	14.66	51.11	2025-11-28 11:20:49.959793+00
677a0f14-b3b8-41d1-b889-c7278c38a66c	prod-agent	7.85	14.61	51.11	2025-11-28 11:20:59.958277+00
38874625-bdaf-4387-8877-6caabb5c8116	prod-agent	9.92	14.69	51.11	2025-11-28 11:21:09.959479+00
18142f30-3e27-454e-bd96-fc00d4b49e0c	prod-agent	7.80	14.71	51.11	2025-11-28 11:21:20.025973+00
d6593daf-c034-434e-9f7e-ae77b4bf4a11	prod-agent	8.59	14.78	51.11	2025-11-28 11:21:29.953166+00
ef666634-62c6-4cd6-890b-24a8cbf629c8	prod-agent	15.25	14.73	51.11	2025-11-28 11:21:39.956443+00
671e8a4d-c656-426d-bf2c-ffe482c81d5d	prod-agent	7.84	14.70	51.11	2025-11-28 11:21:49.959066+00
338184b3-c24d-49c2-b248-5e537d222599	prod-agent	7.89	14.69	51.11	2025-11-28 11:21:59.952499+00
d93f9148-6c1b-4ee8-b32c-961e84b1488f	prod-agent	9.74	14.64	51.11	2025-11-28 11:22:09.961546+00
2f4afc1d-b65c-4248-bb69-2bff9cd68d47	prod-agent	8.90	14.65	51.11	2025-11-28 11:22:19.952078+00
3c875f92-77f8-4617-aa72-c7b93fe532d2	prod-agent	8.63	14.63	51.11	2025-11-28 11:22:29.958899+00
32f33fe6-26a0-4cfe-bfed-66ebd61d9f24	prod-agent	14.03	14.83	51.11	2025-11-28 11:22:39.963021+00
c6aeec97-2dc9-42a3-b47f-1967d11f6cc4	prod-agent	10.19	14.69	51.11	2025-11-28 11:22:49.957385+00
69206ec8-3146-4dd4-b7b0-53aacc19d40c	prod-agent	7.43	14.69	51.11	2025-11-28 11:22:59.961344+00
8023aa8c-67ef-4ce8-937d-05d9cccfb009	prod-agent	8.10	14.69	51.11	2025-11-28 11:23:09.958498+00
5a6b1613-2989-47a0-84cc-18819eab3351	prod-agent	7.95	14.57	51.11	2025-11-28 11:23:19.958501+00
d5254ef9-4c76-4eba-b704-fb86d4487144	prod-agent	8.40	14.65	51.11	2025-11-28 11:23:29.998448+00
3a34e568-0411-49e3-a1fa-cef9b5c8ee6d	prod-agent	8.84	14.70	51.11	2025-11-28 11:23:39.957193+00
ad8d53ef-efb8-48b6-b16d-c4e638660d56	prod-agent	13.40	14.68	51.11	2025-11-28 11:23:49.958959+00
aeff7339-b85a-413f-aa37-94e973e4fd3a	prod-agent	8.55	14.69	51.11	2025-11-28 11:23:59.959657+00
97d0e225-5365-40f2-bcef-1b9b0c2b0770	prod-agent	9.70	14.67	51.11	2025-11-28 11:24:09.951812+00
09264b52-84c9-45e0-b305-8590a65d0748	prod-agent	8.48	14.72	51.11	2025-11-28 11:24:19.952683+00
f7e2ef45-59cc-4976-9917-a2458c689bf1	prod-agent	8.62	14.74	51.11	2025-11-28 11:24:29.952808+00
72e75064-af23-4eaa-a456-a84379911dfa	prod-agent	8.98	14.74	51.11	2025-11-28 11:24:39.958581+00
1224ffc2-d0c0-4eb2-a624-740418090fb7	prod-agent	14.35	14.76	51.11	2025-11-28 11:24:49.962227+00
6c4de0b8-adb0-4c98-b1fe-ee8da78ababd	prod-agent	8.28	14.80	51.11	2025-11-28 11:24:59.9545+00
cb38bd2b-af68-4236-bff6-e0e38cdae1ff	prod-agent	8.37	14.73	51.11	2025-11-28 11:25:09.954559+00
cfbb778e-7d68-483e-8e61-e4091d38b9b5	prod-agent	8.00	14.75	51.11	2025-11-28 11:25:19.958569+00
d8824bfd-07c3-44cd-9522-de9c10c718e2	prod-agent	8.57	14.73	51.11	2025-11-28 11:25:29.952703+00
04101661-bae0-4ad2-9a08-ead0ff7f4baf	prod-agent	9.22	14.75	51.11	2025-11-28 11:25:39.989127+00
35f1c653-e4cb-4a81-8cb2-c2c02963c8ab	prod-agent	15.03	14.66	51.11	2025-11-28 11:25:49.953456+00
2ec64690-f929-4a26-a55a-c1ae2d00572f	prod-agent	7.30	14.67	51.11	2025-11-28 11:25:59.958965+00
26e6417c-d622-460f-8aa1-4c62a16dc632	prod-agent	9.87	14.67	51.11	2025-11-28 11:26:09.959347+00
9538e97e-c8ac-47b5-9e11-b1aae90cb8b8	prod-agent	8.79	14.70	51.11	2025-11-28 11:26:19.962634+00
f6c351b3-ca77-4b62-9ede-f1649aa686c9	prod-agent	8.82	14.70	51.11	2025-11-28 11:26:30.015232+00
841efbd7-2d01-4707-8279-aa1a3e16c092	prod-agent	9.26	14.72	51.11	2025-11-28 11:26:39.957917+00
1b0f8bee-cfcd-4647-9172-a159f4bd4431	prod-agent	14.63	14.74	51.11	2025-11-28 11:26:49.957786+00
934bafc1-502c-4ea2-aa5e-5cc12fe6f673	prod-agent	7.46	14.72	51.11	2025-11-28 11:26:59.958473+00
e72e2eff-a7f9-4ec8-85e0-baff32c4bdfa	prod-agent	8.58	14.75	51.11	2025-11-28 11:27:09.960497+00
1f32ca61-87b1-4e4b-9b4e-a304b391d391	prod-agent	11.78	14.77	51.11	2025-11-28 11:27:19.958621+00
08beed36-6514-435f-b810-37ff8f46df39	prod-agent	11.16	14.80	51.11	2025-11-28 11:27:29.953461+00
23964632-4cb5-4117-81bb-7d4b4fff6c83	prod-agent	12.68	14.83	51.11	2025-11-28 11:27:39.952789+00
7779664a-a1b0-4f21-a060-9bd30adb00c2	prod-agent	16.03	14.76	51.11	2025-11-28 11:27:49.954503+00
502da509-b972-4314-9e40-915b8e1da2d7	prod-agent	12.61	14.81	51.11	2025-11-28 11:27:59.955051+00
63e71cca-f06c-41f7-917d-12667169ce81	prod-agent	13.08	14.85	51.11	2025-11-28 11:28:09.961543+00
c60369f8-6c9d-436c-b328-661dacb912e4	prod-agent	12.92	14.87	51.11	2025-11-28 11:28:19.951836+00
4e34268c-2800-41b0-825c-d6c21dbd2f33	prod-agent	7.64	14.90	51.11	2025-11-28 11:28:29.963215+00
32d0823b-109d-46dc-8193-925bc3c550dc	prod-agent	9.04	14.89	51.11	2025-11-28 11:28:39.953582+00
9dc0c494-cae9-47d5-96be-575150b41b4c	prod-agent	10.22	14.90	51.11	2025-11-28 11:28:49.961147+00
fd394285-b0bb-4301-a663-b27d776d71b3	prod-agent	14.45	14.73	51.11	2025-11-28 11:28:59.959351+00
a17e60af-aa7a-424a-83ba-bd93a85a1d41	prod-agent	9.94	14.88	51.11	2025-11-28 11:29:09.953706+00
08721fd2-c79a-4a54-82c1-87540bfa28b6	prod-agent	11.65	14.90	51.11	2025-11-28 11:29:19.957588+00
3d7716b5-447a-46c5-be3a-f942a4448472	prod-agent	13.22	14.80	51.11	2025-11-28 11:29:29.958413+00
5ffc8c72-f289-448c-9f28-aeb76dc2d17a	prod-agent	13.04	14.83	51.11	2025-11-28 11:29:39.952117+00
28779c5f-ae6c-4f7a-87a6-c04d69753d68	prod-agent	12.23	14.74	51.11	2025-11-28 11:29:49.963779+00
fda4e1a7-bc82-4322-baff-0c828b53a666	prod-agent	18.96	14.81	51.11	2025-11-28 11:29:59.958131+00
cabf9ec0-db8f-4980-8b80-ef096f19a283	prod-agent	14.12	14.77	51.11	2025-11-28 11:30:09.958351+00
4c7ae77e-9544-4b83-8114-c997b2433d7c	prod-agent	12.49	14.73	51.11	2025-11-28 11:30:19.952835+00
81b26c08-9768-4e37-9724-224f5630cf83	prod-agent	11.45	14.84	51.11	2025-11-28 11:30:29.957599+00
b537f185-650b-405b-a880-33204eb2528c	prod-agent	8.51	14.91	51.11	2025-11-28 11:30:39.959116+00
4e9c9baa-f772-4334-a1f0-73a177d4f5cf	prod-agent	8.56	14.87	51.11	2025-11-28 11:30:49.958751+00
98eef863-c92e-439f-84fb-74f3a6dab692	prod-agent	36.52	15.69	51.11	2025-11-28 11:30:59.968565+00
8158d555-5852-49df-ab59-b6e80439fc8a	prod-agent	46.03	15.63	51.11	2025-11-28 11:31:09.956135+00
f44fbf8c-0e6e-4676-a676-905b028366c1	prod-agent	41.12	15.61	51.11	2025-11-28 11:31:19.955034+00
135ca811-9aec-4866-859b-885ef5bb3c3b	prod-agent	35.83	15.72	51.11	2025-11-28 11:31:29.95984+00
25674af7-d108-4ed4-8f05-80ae5780987f	prod-agent	47.39	16.11	51.11	2025-11-28 11:31:39.958507+00
d140b9f9-beb9-408d-adf4-48dcaef50a4f	prod-agent	30.09	16.76	52.27	2025-11-28 11:43:02.041893+00
89073a1c-54f2-48d7-a27d-49be5e36cac1	prod-agent	14.72	16.58	52.28	2025-11-28 11:43:12.03834+00
3049c391-03d1-4319-b5bd-5cefd6871641	prod-agent	21.82	16.79	52.28	2025-11-28 11:43:22.037117+00
b9dbbb70-c5d5-419a-b793-ffe2de10593e	prod-agent	12.86	16.70	52.28	2025-11-28 11:43:32.037262+00
6ddf530f-8c9d-4d23-8ce0-2a57530adade	prod-agent	13.43	16.69	52.28	2025-11-28 11:43:42.039009+00
87978dde-86cd-4895-968d-7b7bfa88f726	prod-agent	13.79	16.66	52.28	2025-11-28 11:43:52.032448+00
29063b14-4070-4d3b-824e-8bdc8e06826c	prod-agent	12.77	16.66	52.28	2025-11-28 11:44:02.03732+00
5356a32b-858c-4a73-9623-d66787cd96d2	prod-agent	11.84	16.39	52.28	2025-11-28 11:44:12.031081+00
2a971440-41cb-4fb5-b72b-b089ebce0740	prod-agent	12.54	16.56	52.28	2025-11-28 11:44:22.045907+00
49be9dbd-b2b9-44aa-aa32-1644dd08e9b8	prod-agent	11.16	16.67	52.28	2025-11-28 11:44:32.037345+00
0d31593f-dd0b-4f40-8229-599e7de50b46	prod-agent	8.18	16.49	52.28	2025-11-28 11:44:42.038557+00
2d3d0a48-bc5e-4309-b019-a0215dbc0302	prod-agent	8.66	16.53	52.28	2025-11-28 11:44:52.039383+00
7582b8b8-2643-4fd3-aab4-f6d67e94955e	prod-agent	8.13	16.49	52.28	2025-11-28 11:45:02.036682+00
fcf1a802-8adf-4988-a468-c050dace8cc5	prod-agent	8.33	16.50	52.28	2025-11-28 11:45:12.039575+00
6ac188ea-eed8-48ae-8a02-1e1395c1e016	prod-agent	8.90	16.51	52.28	2025-11-28 11:45:22.038427+00
80436136-6a16-4da2-815a-ca233045fb30	prod-agent	12.26	16.45	52.28	2025-11-28 11:45:32.031328+00
b0b642c7-4674-4b17-928f-e76067cf95b1	prod-agent	7.88	16.49	52.28	2025-11-28 11:45:42.037835+00
43ee741c-0f79-4d8d-85bb-fb36966ff925	prod-agent	9.22	16.52	52.28	2025-11-28 11:45:52.038627+00
3cd802e0-6255-4d91-b415-be1a36b2d824	prod-agent	15.79	16.85	52.28	2025-11-28 11:46:02.039425+00
52aa0b0a-9018-4a09-82d7-ed5286ed03dd	prod-agent	14.17	16.51	52.28	2025-11-28 11:46:12.030298+00
2a5838f0-57a3-4310-bedc-225b1f587a81	prod-agent	8.12	16.44	52.28	2025-11-28 11:46:22.036964+00
6e0a2f00-3bd0-4bf2-92e6-d9e776ef846f	prod-agent	14.11	16.33	52.28	2025-11-28 11:46:32.032615+00
af822863-4250-433e-9662-ca3a64b01d90	prod-agent	8.19	16.37	52.28	2025-11-28 11:46:42.039689+00
4ffba05b-c7cf-4bc0-b3a2-99874dfcd49c	prod-agent	8.87	16.40	52.28	2025-11-28 11:46:52.047232+00
5ee5d9ec-2b79-44b2-8fd3-080054d2df69	prod-agent	7.42	16.43	52.28	2025-11-28 11:47:02.034542+00
55f8ced1-e7c0-4fed-a9bf-991d6edae69e	prod-agent	8.70	16.53	52.28	2025-11-28 11:47:12.032543+00
03e06be9-6b82-4dfb-ba4e-62876b20edfc	prod-agent	9.78	16.35	52.28	2025-11-28 11:47:22.037533+00
2fe6c3ec-af7f-4e51-95f4-849ba19a0db5	prod-agent	13.99	16.42	52.28	2025-11-28 11:47:32.034807+00
555b7b96-1488-4cbf-b37f-32bd38a0aa0c	prod-agent	6.96	16.38	52.28	2025-11-28 11:47:42.035997+00
ab410016-1c22-4c32-bec6-47e9b7463c94	prod-agent	9.55	16.45	52.28	2025-11-28 11:47:52.041046+00
ab3e7f19-9e6b-4b6e-91e0-8a0732f61bb9	prod-agent	8.23	16.39	52.28	2025-11-28 11:48:02.038917+00
b007918b-d0dd-412b-b261-79066c9cffb1	prod-agent	8.32	16.45	52.28	2025-11-28 11:48:12.032409+00
e81b9657-75fb-40d6-9272-0b575a703170	prod-agent	10.05	16.43	52.28	2025-11-28 11:48:22.038808+00
18b96c81-66ea-4e8c-9756-1d898ee269d2	prod-agent	12.99	16.29	52.28	2025-11-28 11:48:32.090993+00
f75f127e-d58e-4769-8696-f04274c06a3d	prod-agent	7.44	16.40	52.28	2025-11-28 11:48:42.031912+00
ee39c5a4-4b35-4659-a470-0d215fee5b8d	prod-agent	9.06	16.39	52.28	2025-11-28 11:48:52.039094+00
b800ac64-9616-4984-baf8-7411ce0deed4	prod-agent	7.59	16.39	52.28	2025-11-28 11:49:02.039596+00
b395a74b-466e-4854-a6ed-267b442933e7	prod-agent	7.84	16.44	52.28	2025-11-28 11:49:12.040997+00
ff0b218a-51af-461a-88be-0a9794f22794	prod-agent	9.81	16.48	52.28	2025-11-28 11:49:22.039195+00
a471179a-05ba-4b32-bcad-83ccedf2b780	prod-agent	10.84	16.45	52.28	2025-11-28 11:49:32.052599+00
fb3b3088-385e-4559-aaff-f4a7dae8ceb9	prod-agent	10.18	16.37	52.28	2025-11-28 11:49:42.037864+00
99b89087-8c26-48f8-ac12-480d070b5f9e	prod-agent	9.61	16.50	52.28	2025-11-28 11:49:52.039796+00
e3662433-bd94-4b06-8dd4-e4d020ddeb8e	prod-agent	7.37	16.48	52.28	2025-11-28 11:50:02.038564+00
93ab1e0d-d0d8-4ee8-80b2-6cccba719ab7	prod-agent	9.45	16.51	52.28	2025-11-28 11:50:12.031758+00
20b49638-2406-4d46-a257-09417654d2ec	prod-agent	9.24	16.56	52.28	2025-11-28 11:50:22.048639+00
25763ddb-3696-43b8-a654-2d558927fbc3	prod-agent	14.54	16.41	52.28	2025-11-28 11:50:32.043451+00
3c8b224d-5b3e-483d-8ef9-38c5c267e594	prod-agent	18.18	16.39	52.28	2025-11-28 11:50:42.037394+00
24d9c832-bf00-4222-a0ec-5e64f593175d	prod-agent	12.68	16.36	52.28	2025-11-28 11:50:52.04211+00
4a42d603-6cb7-4bc7-a01f-f0e37890dfde	prod-agent	11.15	16.33	52.28	2025-11-28 11:51:02.036674+00
670defbd-501a-469b-aed8-167a8c9cc0e0	prod-agent	13.86	16.39	52.28	2025-11-28 11:51:12.031296+00
075aec47-a0f7-49df-b886-0a6532071f61	prod-agent	12.44	16.39	52.28	2025-11-28 11:51:22.035989+00
0bbcfbab-cc1c-40d7-8028-1178fdb17e72	prod-agent	12.20	16.32	52.28	2025-11-28 11:51:32.039642+00
01282139-68ac-4a21-9a85-0f214e8c3e83	prod-agent	16.89	16.41	52.28	2025-11-28 11:51:42.036291+00
b27cfe16-d064-4804-bf8d-ecb749de755d	prod-agent	12.01	16.31	52.28	2025-11-28 11:51:52.03404+00
14570fd5-5b60-46e7-bd67-2bce482cb3ac	prod-agent	26.85	16.58	52.28	2025-11-28 11:52:02.041367+00
a9f10783-20a1-48ca-ba05-388643c9a2f1	prod-agent	92.77	19.79	52.29	2025-11-28 11:52:12.041621+00
e82ad755-3212-497b-a580-8c8e302c2fa6	prod-agent	99.94	20.63	52.30	2025-11-28 11:52:22.039683+00
7360c4ec-ea8f-43b2-a404-988e66b42ed9	prod-agent	99.84	21.72	52.32	2025-11-28 11:52:32.041081+00
ffa4a426-ee6f-4c38-987b-5458c59a284c	prod-agent	98.92	21.84	52.33	2025-11-28 11:52:42.040793+00
1900e2f7-a6cf-4e1b-ae74-a258f2d11454	prod-agent	96.71	21.53	52.33	2025-11-28 11:52:52.039409+00
779766bd-0bc6-4bdc-ac40-5e40ff77dfd7	prod-agent	65.46	21.48	52.33	2025-11-28 11:53:02.035538+00
5a4cd499-bf85-402d-abdd-7089b32166d7	prod-agent	66.37	22.13	52.33	2025-11-28 11:53:12.039181+00
69f5f678-7c5d-4cfd-b398-dee579cf5175	prod-agent	71.47	24.43	52.33	2025-11-28 11:53:22.043046+00
1da3a100-d783-4048-b24b-6ff9af0bdbfa	prod-agent	91.38	27.15	52.33	2025-11-28 11:53:32.055843+00
4d6c19df-e60a-4413-ae1f-7d34cdac348b	prod-agent	99.69	29.89	52.33	2025-11-28 11:53:42.093202+00
3197a245-5ed9-44a9-988d-54207d395497	prod-agent	99.79	31.96	52.33	2025-11-28 11:53:52.04524+00
18ffd0bd-ce7b-493c-a590-cd882f0d4410	prod-agent	60.91	19.62	52.33	2025-11-28 11:54:02.039426+00
8fa2b5e4-846f-41a1-8241-0880ffddaf03	prod-agent	73.24	22.85	52.33	2025-11-28 11:54:12.040812+00
1571e0d7-92ea-4a09-8264-6e55e7434944	prod-agent	65.28	22.24	52.33	2025-11-28 11:54:22.036873+00
6a43c585-fccc-4cc6-ae3a-76cbbcbca682	prod-agent	67.41	23.04	52.33	2025-11-28 11:54:32.041933+00
f6fa6e58-7cce-4d08-a541-1c223ace669f	prod-agent	66.98	24.07	52.33	2025-11-28 11:54:42.04065+00
23370c3f-083c-441f-9068-307d438f46e0	prod-agent	73.70	25.78	52.33	2025-11-28 11:54:52.046607+00
68487f37-a078-4b9f-bc91-c5ed346c3cb5	prod-agent	99.67	30.67	52.33	2025-11-28 11:55:02.044079+00
21847fb9-fb1a-48d4-ad32-d82cd2f9e4bf	prod-agent	99.59	33.14	52.33	2025-11-28 11:55:12.042576+00
bbffb192-b604-4367-b92e-514a365bcc8d	prod-agent	75.67	20.45	52.38	2025-11-28 11:55:22.033486+00
9f8aa698-f0ae-4fc6-ba2c-d1fc05c1a3ee	prod-agent	64.09	20.01	52.38	2025-11-28 11:55:32.035067+00
033d2b6b-35dc-4c0f-9b45-f0137a4d5907	prod-agent	9.34	19.13	52.38	2025-11-28 11:55:42.033405+00
96aa8234-ef3d-430f-aeb3-ed8ff3b8cf5b	prod-agent	15.30	19.07	52.38	2025-11-28 11:55:52.032089+00
9c85d652-9732-4780-b5f6-8e642d5b5f08	prod-agent	7.49	19.02	52.38	2025-11-28 11:56:02.032435+00
92e6064d-b989-4ba9-9238-371357da3267	prod-agent	9.74	18.93	52.38	2025-11-28 11:56:12.037055+00
e0b2f23f-2e3c-4f01-8060-95074203cf66	prod-agent	9.79	18.99	52.38	2025-11-28 11:56:22.031648+00
41e72d3e-0419-4f20-86bf-a0b3e9488a03	prod-agent	9.96	19.06	52.38	2025-11-28 11:56:32.038327+00
1818c729-5d78-4431-a3d9-bb331cdb999c	prod-agent	8.68	19.04	52.38	2025-11-28 11:56:42.038971+00
ec1b2509-e3c9-47e2-ad90-5d2b52343a1c	prod-agent	14.91	19.08	52.38	2025-11-28 11:56:52.033894+00
d2261630-f40d-4c95-80a7-38a0ed6046bc	prod-agent	11.59	18.06	52.38	2025-11-28 11:57:02.033095+00
680be887-2ff0-4443-afa6-d01c1704df4d	prod-agent	17.36	18.31	52.38	2025-11-28 11:57:12.037781+00
81cc3d99-da28-4541-96d7-793f7282c579	prod-agent	9.49	18.28	52.38	2025-11-28 11:57:22.035582+00
5c17f190-1809-41bc-90cd-d0c4dcd4cd16	prod-agent	9.21	18.31	52.38	2025-11-28 11:57:32.04492+00
80a313da-572f-4abf-8788-b1122026c229	prod-agent	100.00	18.55	52.49	2025-11-28 11:58:46.701188+00
68248c5e-afc1-44ca-b7ea-cab86b12d9cd	prod-agent	29.77	18.65	52.49	2025-11-28 11:59:14.204827+00
667bfb4d-aa49-46b4-962d-88ee96300903	prod-agent	9.34	18.57	52.49	2025-11-28 11:59:24.195479+00
92545116-5dda-44c4-b0a8-0bea1f82e439	prod-agent	12.60	18.46	52.49	2025-11-28 11:59:34.194338+00
e865770b-25f8-448c-be3f-fa34c8c5d4b5	prod-agent	14.73	18.52	52.49	2025-11-28 11:59:44.194633+00
32ba7314-5acf-4602-ab96-4c8371b5a0c9	prod-agent	10.75	18.54	52.49	2025-11-28 11:59:54.196769+00
9d523e35-8049-4033-8f1a-c04e1214f864	prod-agent	16.71	18.46	52.49	2025-11-28 12:00:04.197807+00
4524f5c5-1134-49bf-978b-adec3a9e42fd	prod-agent	12.14	18.51	52.49	2025-11-28 12:00:14.206703+00
90db9e78-9494-48d0-bbad-911446609c87	prod-agent	14.31	18.51	52.49	2025-11-28 12:00:24.193022+00
5d0640da-d851-4af7-adf0-8f8ccc4d7b95	prod-agent	8.63	18.54	52.49	2025-11-28 12:00:34.200944+00
026fddf9-abde-44c1-8642-a483686d9cae	prod-agent	8.20	18.51	52.49	2025-11-28 12:00:44.194881+00
21863785-d0e5-4785-953c-77ace8b1303b	prod-agent	8.29	18.49	52.49	2025-11-28 12:00:54.199513+00
92b8a34d-4888-42a7-b4af-620b730ce4a6	prod-agent	23.13	18.91	52.49	2025-11-28 12:01:04.214197+00
90062f00-6e56-443c-afbd-4b40349e0079	prod-agent	11.75	18.56	52.49	2025-11-28 12:01:14.199489+00
3310275e-1fd8-4b7f-8566-e2b9612fda77	prod-agent	8.33	18.51	52.49	2025-11-28 12:01:24.203086+00
224599bc-8d8f-4aea-b6cc-5bc3c9f0c867	prod-agent	8.48	18.51	52.49	2025-11-28 12:01:34.193923+00
46fc99d8-b27d-4ced-a92c-1a39ab743753	prod-agent	9.34	18.48	52.49	2025-11-28 12:01:44.194516+00
cf4bdd0c-5e8c-4e3c-bc2d-971ed93a817a	prod-agent	8.83	18.46	52.49	2025-11-28 12:01:54.197417+00
1a363897-9ff6-4f73-a1ec-d83c7e47e2c1	prod-agent	15.10	18.46	52.49	2025-11-28 12:02:04.199377+00
ea155fff-1110-4bf6-9099-fc951a286ea5	prod-agent	8.64	18.47	52.49	2025-11-28 12:02:14.202199+00
74299e5b-d08d-4555-9586-8f4973a2e84a	prod-agent	12.30	18.50	52.49	2025-11-28 12:02:24.199662+00
633d7127-9982-47bc-8d2e-13f3f2d39f1b	prod-agent	11.31	18.53	52.49	2025-11-28 12:02:34.206405+00
82339ea9-e131-4bf9-b312-231950667e4b	prod-agent	11.00	18.50	52.49	2025-11-28 12:02:44.194839+00
ac1b7bde-bf88-437f-b299-023843c3a225	prod-agent	12.36	18.49	52.49	2025-11-28 12:02:54.201786+00
efdfe6aa-7974-4136-946e-e7638e28da40	prod-agent	13.46	18.54	52.49	2025-11-28 12:03:04.200338+00
e82d8101-fed9-4724-b1fa-a654e220f527	prod-agent	7.91	18.52	52.49	2025-11-28 12:03:14.195756+00
46907d33-e92f-4af7-b8af-37c03b55f8ab	prod-agent	9.65	18.53	52.49	2025-11-28 12:03:24.193366+00
71b40f7f-b3fc-44fc-b975-7d29105982fe	prod-agent	9.16	18.54	52.49	2025-11-28 12:03:34.197916+00
4779527f-a22c-44c6-84a5-379756d5cc51	prod-agent	12.68	18.35	52.49	2025-11-28 12:03:44.201806+00
1760f4ee-cde6-4ab8-91a8-285b30109eee	prod-agent	11.77	18.55	52.49	2025-11-28 12:03:54.195621+00
be37b984-06c3-472a-a5c3-20d6260a0e2d	prod-agent	18.30	19.00	52.49	2025-11-28 12:04:04.308467+00
267f3a0e-9885-454b-849a-4c4147237e5e	prod-agent	21.35	18.63	52.49	2025-11-28 12:04:14.202313+00
88681b79-5fbb-4cd3-8094-b9f72a67b2f4	prod-agent	8.37	18.50	52.49	2025-11-28 12:04:24.194527+00
d33ea034-82ca-498a-978e-ea7e16e96178	prod-agent	11.20	18.53	52.49	2025-11-28 12:04:34.20046+00
a2e5bae0-4fd5-4455-ac3c-196df8ed139d	prod-agent	10.67	18.50	52.49	2025-11-28 12:04:44.198972+00
a75baa19-ab8a-478b-881e-720d143a4b25	prod-agent	11.78	18.53	52.49	2025-11-28 12:04:54.197668+00
2aa0551d-06f6-4e89-9c69-e1a8bd82634e	prod-agent	11.18	18.48	52.49	2025-11-28 12:05:04.195196+00
695d5159-d631-4050-8856-921140989a1e	prod-agent	17.37	18.51	52.49	2025-11-28 12:05:14.199298+00
5341caae-8f37-4468-9a68-f9eb3ed050b8	prod-agent	12.18	18.49	52.49	2025-11-28 12:05:24.200417+00
de3fc930-827d-4e45-afad-43b0d2904085	prod-agent	10.32	18.50	52.49	2025-11-28 12:05:34.199332+00
c2daa9bd-9373-47d0-a630-bf833a1b828c	prod-agent	13.52	18.52	52.49	2025-11-28 12:05:44.19999+00
a16c7e84-d510-41bd-9fa1-a05884151ec9	prod-agent	11.51	18.53	52.49	2025-11-28 12:05:54.198952+00
b931bbb6-f982-4e1f-b8a7-ba5dc9c95e47	prod-agent	9.78	18.52	52.49	2025-11-28 12:06:04.195524+00
8e63a300-9f22-47dc-9d4f-5be9263e21b9	prod-agent	16.55	18.52	52.49	2025-11-28 12:06:14.203407+00
1b36431d-667b-457d-8cd6-3498352724b0	prod-agent	10.55	18.55	52.49	2025-11-28 12:06:24.194468+00
dd7e6ea1-3753-4025-b6b5-b81d95ea9b56	prod-agent	11.67	18.52	52.49	2025-11-28 12:06:34.194353+00
463d6ad7-a879-420d-848f-a25f2d24bef7	prod-agent	10.64	18.51	52.49	2025-11-28 12:06:44.200844+00
916dca3e-19e9-41ec-9108-1544ef6b838a	prod-agent	12.02	18.46	52.49	2025-11-28 12:06:54.201116+00
1feef627-daf3-48a3-9b66-1a7075c00a34	prod-agent	12.99	18.50	52.49	2025-11-28 12:07:04.200905+00
03d7c6d7-7077-4238-9f05-779e77579c74	prod-agent	18.03	18.59	52.49	2025-11-28 12:07:14.199687+00
97ae6ca9-e4bc-4d85-b2ab-538069dea1b3	prod-agent	12.98	18.49	52.49	2025-11-28 12:07:24.193526+00
917bb236-8095-441c-a5ac-2b71af0e3096	prod-agent	11.43	18.51	52.49	2025-11-28 12:07:34.194015+00
0d9a188e-597c-464d-9806-a7e21f4fb155	prod-agent	11.16	18.50	52.49	2025-11-28 12:07:44.20162+00
fcc56e30-4f8d-46e1-8128-2f4f3ee2df9b	prod-agent	9.76	18.51	52.49	2025-11-28 12:07:54.203237+00
29b88d35-93b7-4114-822f-919ab5b5f836	prod-agent	13.94	18.55	52.49	2025-11-28 12:08:04.204776+00
c714a6f3-4ee5-481e-99c2-051033b3a9f4	prod-agent	15.55	18.55	52.49	2025-11-28 12:08:14.20079+00
04be9515-1f36-4a49-ac59-43e90689f57b	prod-agent	10.07	18.50	52.49	2025-11-28 12:08:24.195+00
aa036234-97a3-449f-a983-ffa4020e17f3	prod-agent	12.44	18.56	52.49	2025-11-28 12:08:34.194022+00
e2dcbf60-2caf-46b9-b5bf-140630d60006	prod-agent	9.17	18.38	52.49	2025-11-28 12:08:44.199662+00
5aa0a112-053a-4f2d-b36b-3bbb0fea51e6	prod-agent	9.02	18.55	52.49	2025-11-28 12:08:54.196001+00
216721d8-4bea-4e24-ad13-3e1973846a2a	prod-agent	9.96	18.54	52.49	2025-11-28 12:09:04.193941+00
23978eed-86f5-4812-8782-1eae8a42eb47	prod-agent	13.53	18.56	52.49	2025-11-28 12:09:14.273624+00
02696d05-7df8-44ad-8870-01d81996dcf6	prod-agent	13.16	18.71	52.49	2025-11-28 12:09:24.194666+00
8a44d332-8cdf-4fa0-92e0-8dba40b9ed4e	prod-agent	37.11	20.70	52.49	2025-11-28 12:09:34.200743+00
b25e0b0a-bc13-413e-a9ae-58e920ba3369	prod-agent	52.06	22.72	52.49	2025-11-28 12:09:44.195449+00
534c5314-711a-4d31-9ad4-d951b3cb892f	prod-agent	43.18	18.81	52.49	2025-11-28 12:09:54.198275+00
24e9f4f1-f579-48ee-9673-c22de8329c45	prod-agent	7.05	18.57	52.49	2025-11-28 12:10:04.194354+00
1af9af79-aa4d-43b2-aedb-336841d4f623	prod-agent	9.96	18.58	52.49	2025-11-28 12:10:14.20436+00
32127120-3f09-4426-9feb-cbef8e6aa20c	prod-agent	17.81	18.53	52.49	2025-11-28 12:10:24.198288+00
ee5858b2-4c52-4ff2-abdd-f54df9b3a53b	prod-agent	7.35	18.56	52.49	2025-11-28 12:10:34.203515+00
13aebfa3-e42d-43ae-a756-06910410d442	prod-agent	8.75	18.57	52.49	2025-11-28 12:10:44.194423+00
9e69e56b-76c4-46f1-a96c-3081fea2c919	prod-agent	8.06	18.55	52.49	2025-11-28 12:10:54.196116+00
d858c955-c7c1-48b7-b032-95dbd981fc34	prod-agent	9.12	18.49	52.49	2025-11-28 12:11:04.201133+00
2d72e3e0-331a-4efa-baca-f0e6e376ae4c	prod-agent	8.76	18.58	52.49	2025-11-28 12:11:14.199701+00
8105f7d7-d8e4-4ad3-89f1-9bd6a502324d	prod-agent	12.64	18.55	52.49	2025-11-28 12:11:24.199538+00
e12a3842-4d89-4e85-a2df-5bb0ce80d057	prod-agent	8.89	18.54	52.49	2025-11-28 12:11:34.198267+00
397043e0-b258-446c-8570-ff052b7fdb56	prod-agent	11.32	18.55	52.49	2025-11-28 12:11:44.195475+00
a2e3a741-c7bf-45bc-aeda-be7e5207b73f	prod-agent	7.43	18.45	52.49	2025-11-28 12:11:54.19478+00
cf763e73-17d8-4fed-b422-9f9649e3d05c	prod-agent	9.07	18.51	52.49	2025-11-28 12:12:04.199877+00
00d214a0-0bdd-4814-996f-e810e5fe7832	prod-agent	8.09	18.52	52.49	2025-11-28 12:12:14.195527+00
d283fe63-4a5b-4b33-8843-4271dac1e102	prod-agent	11.48	18.48	52.49	2025-11-28 12:12:24.200087+00
bef63818-a88e-4549-b104-5d17656055ca	prod-agent	9.07	18.52	52.49	2025-11-28 12:12:34.200595+00
ff5595dc-a0bd-4f71-b7b9-a28d5c32233b	prod-agent	9.00	18.51	52.49	2025-11-28 12:12:44.202701+00
fb0f5aa4-5c99-40e4-ab2e-f8f31288fd31	prod-agent	7.72	18.46	52.49	2025-11-28 12:12:54.202847+00
1a2d8555-d6b2-4a4c-87fb-fec6d1d07e65	prod-agent	9.23	18.45	52.49	2025-11-28 12:13:04.19445+00
15658904-b917-4ab0-a967-1fa5a1bd5225	prod-agent	10.15	18.37	52.49	2025-11-28 12:13:14.201016+00
9057e750-5335-490c-a143-a64d288886df	prod-agent	14.84	18.44	52.49	2025-11-28 12:13:24.200772+00
ed43b172-299a-413f-b2c4-71c3e1774077	prod-agent	8.21	18.31	52.49	2025-11-28 12:13:34.200258+00
99c7af69-4482-43b5-9154-300d65c4dad5	prod-agent	8.77	18.42	52.49	2025-11-28 12:13:44.196118+00
16d75d1c-d49a-4594-8a29-dc4030e1eb67	prod-agent	7.61	18.38	52.49	2025-11-28 12:13:54.198228+00
7c6f8405-4f77-4ac5-82fe-b55432744b1a	prod-agent	9.49	18.38	52.49	2025-11-28 12:14:04.194632+00
181e563c-fb0a-41fe-8085-685c7cf6080f	prod-agent	8.58	18.43	52.49	2025-11-28 12:14:14.201275+00
1f8e4225-e779-4ef1-85b7-d7a3d1051e4c	prod-agent	11.22	18.45	52.49	2025-11-28 12:14:24.202317+00
33840d33-c1bd-4a4b-839c-9f7db7d756fb	prod-agent	10.08	18.38	52.49	2025-11-28 12:14:34.19522+00
eea8a570-62db-44cd-9282-476f3e3d5a91	prod-agent	8.75	18.37	52.49	2025-11-28 12:14:44.195467+00
cc403e97-4caa-4a75-80c2-8b302437eb31	prod-agent	7.82	18.32	52.49	2025-11-28 12:14:54.196664+00
6946d185-024a-4085-b280-4d5218e2763f	prod-agent	9.22	18.28	52.49	2025-11-28 12:15:04.20112+00
eb7d316d-5362-493a-a93f-c8bbcc5a29b7	prod-agent	8.01	18.27	52.49	2025-11-28 12:15:14.199877+00
ffab37d3-6d6c-414a-9f39-e06cc6b4acc8	prod-agent	8.54	18.26	52.49	2025-11-28 12:15:24.199649+00
84614ac9-7b7a-49f4-9527-b843cbf39512	prod-agent	13.37	18.39	52.49	2025-11-28 12:15:34.199694+00
66470716-92a8-4518-8767-93695c21a856	prod-agent	9.49	18.38	52.49	2025-11-28 12:15:44.203169+00
3c9add75-8e0f-44a3-aa31-d6b50c190673	prod-agent	7.72	18.33	52.49	2025-11-28 12:15:54.194362+00
e48a44ad-d69d-475e-a2b7-df570d2808e0	prod-agent	14.29	18.49	52.49	2025-11-28 12:16:04.207498+00
0f8f7c95-c0be-4451-b01e-28b5f4468fd3	prod-agent	16.71	18.39	52.49	2025-11-28 12:16:14.199588+00
e14217c8-8d47-4303-a0be-207dbd18418b	prod-agent	7.64	18.38	52.49	2025-11-28 12:16:24.200595+00
08a48c78-5d3d-4a01-bd7a-c5167f660747	prod-agent	15.69	18.35	52.49	2025-11-28 12:16:34.206662+00
7c9b981d-c475-45f2-8aa8-c92f9d96cd72	prod-agent	8.22	18.43	52.49	2025-11-28 12:16:44.20241+00
f7fc8f21-82f8-4b68-9e9b-2435f4c0df2e	prod-agent	6.40	18.41	52.49	2025-11-28 12:16:54.19414+00
d23f5d9e-9537-4cdb-bb95-2d3f32d1d229	prod-agent	9.27	18.36	52.49	2025-11-28 12:17:04.199798+00
c5ed6322-3b50-4bf6-a124-22ee834e4fe1	prod-agent	8.60	18.38	52.49	2025-11-28 12:17:14.19455+00
be4bae7c-c8c9-41ed-8a7d-91a7a195d13d	prod-agent	8.50	18.40	52.49	2025-11-28 12:17:24.194953+00
73d0d8ea-aeb8-422d-99f9-b17f7eb6f6f5	prod-agent	14.36	18.45	52.49	2025-11-28 12:17:34.193447+00
26131b1e-172f-4ec9-9876-9e169468be80	prod-agent	7.96	18.41	52.49	2025-11-28 12:17:44.20729+00
4339403e-a61e-4f3b-a7ff-81fa357bafb3	prod-agent	7.94	18.41	52.49	2025-11-28 12:17:54.194601+00
b717c89b-f8b4-4cdd-9afe-15e587964167	prod-agent	9.19	18.43	52.49	2025-11-28 12:18:04.195777+00
1e1f9bdf-a9fe-400e-bca8-f72e1266e60b	prod-agent	7.63	18.49	52.49	2025-11-28 12:18:14.204705+00
ac8631a8-0983-4870-8572-69ee3166a912	prod-agent	8.39	18.43	52.49	2025-11-28 12:18:24.193876+00
8297cc23-9c94-4e82-b976-fdbfa77dff28	prod-agent	14.99	18.56	52.49	2025-11-28 12:18:34.204489+00
4c521493-2004-4b70-aae2-216491308d27	prod-agent	8.70	18.44	52.49	2025-11-28 12:18:44.202974+00
ba01f9c2-68bf-4d99-9f41-811ef6c27b11	prod-agent	7.75	18.47	52.49	2025-11-28 12:18:54.194461+00
bee63f09-3b44-4d45-990d-460c8d4b2021	prod-agent	8.20	18.45	52.49	2025-11-28 12:19:04.194367+00
cf3457fe-ce17-466f-b3d2-4255923aabaf	prod-agent	8.71	18.40	52.49	2025-11-28 12:19:14.200558+00
d705c77d-b7c1-4509-8933-bf4385ee3727	prod-agent	8.37	18.45	52.49	2025-11-28 12:19:24.25169+00
72eb1e4f-1930-4164-84f5-cfa90b8c954b	prod-agent	10.12	18.48	52.49	2025-11-28 12:19:34.198466+00
e605875e-f9b6-4d28-ab70-7be35e62584e	prod-agent	11.87	18.48	52.49	2025-11-28 12:19:44.193899+00
b16f2126-78c1-4c65-b70e-56a356b6fa32	prod-agent	7.33	18.47	52.49	2025-11-28 12:19:54.196136+00
d0cca30a-d65d-4714-ac48-812d97cb0798	prod-agent	8.26	18.54	52.49	2025-11-28 12:20:04.202572+00
b78cc1e2-0a74-454b-817a-632ce617ef3e	prod-agent	9.12	18.51	52.49	2025-11-28 12:20:14.19465+00
6e20ccfb-1c48-4fdc-a499-672d5f4605da	prod-agent	8.85	18.47	52.49	2025-11-28 12:20:24.200116+00
45215ab5-12cc-44a7-a462-18a18b55febb	prod-agent	9.40	18.54	52.49	2025-11-28 12:20:34.194386+00
e5e69167-9cee-4ff4-b28e-bfa28c383df6	prod-agent	13.57	18.54	52.49	2025-11-28 12:20:44.192711+00
ee060dc7-eabb-433a-9718-1319d6fd0959	prod-agent	8.13	18.56	52.49	2025-11-28 12:20:54.194301+00
925c2758-9de2-4c43-8e3d-463a0b232ce4	prod-agent	9.30	18.54	52.49	2025-11-28 12:21:04.203153+00
3a83a8d5-95c0-4931-8817-854f2d93fdd6	prod-agent	7.77	18.58	52.49	2025-11-28 12:21:14.202981+00
ac7b4f8b-bc0a-4157-98b1-91aca1dcb57f	prod-agent	8.43	18.49	52.49	2025-11-28 12:21:24.200781+00
d75a4f61-2ad8-40fd-affb-c35d26600b45	prod-agent	9.11	18.53	52.49	2025-11-28 12:21:34.200252+00
038ccd6d-8507-47c3-bd7b-5b4912da3542	prod-agent	13.87	18.47	52.49	2025-11-28 12:21:44.20035+00
b0e0bf04-a2c7-474f-b7c7-1c4e471cec7b	prod-agent	7.16	18.49	52.49	2025-11-28 12:21:54.20726+00
7b76289d-4e49-4213-8d8a-5ad9aa9fc556	prod-agent	7.55	18.48	52.49	2025-11-28 12:22:04.2055+00
54f68596-6a2b-4004-b827-50f4d2c68e4b	prod-agent	9.79	18.47	52.49	2025-11-28 12:22:14.194616+00
7941ff62-b34d-4950-9ca3-a95cae6ae620	prod-agent	8.07	18.57	52.49	2025-11-28 12:22:24.194971+00
5b866500-a932-4100-874a-b24d3a50cae1	prod-agent	8.40	18.39	52.49	2025-11-28 12:22:34.1943+00
194df6b3-1728-4fae-a6c6-e7063ed777c2	prod-agent	14.38	18.53	52.49	2025-11-28 12:22:44.199393+00
3cdce80a-6793-4635-a98f-79a3f6607f6b	prod-agent	7.04	18.51	52.49	2025-11-28 12:22:54.204897+00
3c894bf6-1dce-4450-a4e1-7a55472556d6	prod-agent	8.54	18.53	52.49	2025-11-28 12:23:04.194626+00
4865d51c-bdc0-40df-8fd2-5e12939a5ec6	prod-agent	9.29	18.46	52.49	2025-11-28 12:23:14.204267+00
9e6431f2-cb24-4fff-bc0f-b1201317d26a	prod-agent	7.85	18.50	52.49	2025-11-28 12:23:24.200951+00
720ae01b-00e5-46a5-b02c-006596479beb	prod-agent	8.16	18.51	52.49	2025-11-28 12:23:34.194139+00
100786ab-c679-4d48-90a7-b5dfca9be2ff	prod-agent	13.98	18.49	52.49	2025-11-28 12:23:44.195536+00
ccd6dbc4-d245-4cbf-87ae-01f2391aa036	prod-agent	8.05	18.39	52.49	2025-11-28 12:23:54.201788+00
b262db67-8b5b-450d-880f-46e7e10ded5f	prod-agent	7.25	18.42	52.49	2025-11-28 12:24:04.194138+00
599a5d05-b1fd-4fd9-b4cf-fb201d1ad28b	prod-agent	8.28	18.39	52.49	2025-11-28 12:24:14.20032+00
3de85a50-63f5-4a7f-8b2a-ff8f6feb7d80	prod-agent	8.75	18.45	52.49	2025-11-28 12:24:24.198029+00
1c4ed92b-899c-4e64-b1b5-9a7d54323029	prod-agent	9.03	18.31	52.49	2025-11-28 12:24:34.258132+00
24736e50-4da5-4aff-89af-dc6270dea38d	prod-agent	9.87	18.38	52.49	2025-11-28 12:24:44.200279+00
2fb7858f-4fc0-4d69-9c62-5e95e4e3ccaf	prod-agent	12.56	18.40	52.49	2025-11-28 12:24:54.201146+00
2aa234ee-862a-4416-9ceb-c4636ab701c8	prod-agent	8.39	18.41	52.49	2025-11-28 12:25:04.202478+00
17292af6-a757-4a8a-8aa9-8b2d38bb669b	prod-agent	8.40	18.37	52.49	2025-11-28 12:25:14.194072+00
b060f5a6-26ce-4722-8b15-5fc83f085258	prod-agent	9.20	18.32	52.49	2025-11-28 12:25:24.201701+00
ec3b4814-2340-4db1-99d6-015224dac868	prod-agent	8.42	18.43	52.49	2025-11-28 12:25:34.197+00
361f2b00-4520-4769-8938-cf2e269c5052	prod-agent	8.53	18.35	52.49	2025-11-28 12:25:44.19851+00
04d775d4-cd72-4aed-8b0f-c0d1c8b37631	prod-agent	13.80	18.43	52.49	2025-11-28 12:25:54.203245+00
75dc3e50-4357-4aca-bd3c-fd36c96057e5	prod-agent	7.22	18.45	52.49	2025-11-28 12:26:04.195995+00
235bd228-26c4-4ff2-9b86-cdf6fbcf343b	prod-agent	9.68	18.45	52.49	2025-11-28 12:26:14.201517+00
e034a378-5f4d-43a2-b77f-e2db915c841f	prod-agent	9.26	18.45	52.49	2025-11-28 12:26:24.203313+00
3f0f748b-0804-4382-9f5c-d1dff4f8b504	prod-agent	7.42	18.43	52.49	2025-11-28 12:26:34.194258+00
09d2faa9-0f03-4da8-b702-7f33ab8c551a	prod-agent	9.19	18.42	52.49	2025-11-28 12:26:44.194309+00
bdd7b034-fb47-4139-a0de-e8dc48bf6fe5	prod-agent	15.16	18.49	52.49	2025-11-28 12:26:54.19941+00
68fd9e8b-3d1f-4966-a1d6-9d465754fea3	prod-agent	6.66	18.41	52.49	2025-11-28 12:27:04.196618+00
ed1523bc-b4d9-415b-9480-57c95aefed71	prod-agent	9.33	18.50	52.49	2025-11-28 12:27:14.195014+00
3e8feba6-a77b-4199-89b1-9052f36a69be	prod-agent	9.51	18.52	52.49	2025-11-28 12:27:24.194786+00
195263a4-6278-4875-9e7f-14bde0ed4e74	prod-agent	7.77	18.44	52.49	2025-11-28 12:27:34.199956+00
316af9ee-0f1e-4008-988d-f185d8f9ac42	prod-agent	10.06	18.50	52.49	2025-11-28 12:27:44.194566+00
e86a6a89-0edf-4a89-833a-fb21e9f0aba4	prod-agent	12.03	18.52	52.49	2025-11-28 12:27:54.206929+00
06a97b11-720c-4005-ad89-ac18d26e6f45	prod-agent	6.85	18.45	52.49	2025-11-28 12:28:04.197228+00
79e7ba85-b6bf-4f87-b799-694461b98126	prod-agent	10.24	18.45	52.49	2025-11-28 12:28:14.194908+00
99bc8e9b-8d6b-4782-a617-9ff3ed1ef319	prod-agent	9.01	18.42	52.49	2025-11-28 12:28:24.200572+00
4b38d934-5110-4ed8-acc1-a941443d1a62	prod-agent	8.04	18.44	52.49	2025-11-28 12:28:34.200692+00
a3dff8d1-a0e0-4342-a289-ad005d3a2fc5	prod-agent	10.19	18.49	52.49	2025-11-28 12:28:44.199277+00
aabdd119-f73d-4a2b-bf3b-25c4ae05bc65	prod-agent	13.28	18.51	52.49	2025-11-28 12:28:54.195526+00
53eeb09c-a26d-4787-a53c-46f5da0988bb	prod-agent	8.27	18.42	52.49	2025-11-28 12:29:04.208164+00
86a73136-ef5f-4d1b-b675-b6a10399dc70	prod-agent	9.14	18.49	52.49	2025-11-28 12:29:14.194525+00
b7df5549-d74a-44b4-9813-5f6be8e55c07	prod-agent	9.47	18.48	52.49	2025-11-28 12:29:24.198254+00
027b43d3-6b59-4b7c-ad72-9bc9c6062931	prod-agent	7.71	18.45	52.49	2025-11-28 12:29:34.196594+00
bd9f6a6f-b834-4fe1-a6d0-0891fec3759c	prod-agent	9.04	18.48	52.49	2025-11-28 12:29:44.254066+00
a17141fd-cc28-404f-8d5b-bb2557dcf22a	prod-agent	7.76	18.52	52.49	2025-11-28 12:29:54.203324+00
45f32337-095e-4b85-82e3-233982b51438	prod-agent	12.20	18.47	52.49	2025-11-28 12:30:04.19586+00
9dada046-78c0-44d3-b3a4-8b85a612ced9	prod-agent	8.79	18.51	52.49	2025-11-28 12:30:14.200888+00
1e7e64a8-9a22-41a5-82bb-eb7fed8bba08	prod-agent	9.28	18.47	52.49	2025-11-28 12:30:24.195004+00
0d1a4afd-23df-4e99-b2b7-564aa7cc9d62	prod-agent	7.72	18.52	52.49	2025-11-28 12:30:34.19428+00
620e3976-2999-4e9d-9669-26925b5c6086	prod-agent	8.77	18.44	52.49	2025-11-28 12:30:44.193635+00
2536be12-b37d-4472-bd18-d395f8f745f2	prod-agent	8.94	18.50	52.49	2025-11-28 12:30:54.204381+00
319c13a7-5694-4576-8143-50b45218ddb0	prod-agent	12.97	18.48	52.49	2025-11-28 12:31:04.193637+00
9a2c1e72-29b4-47e9-82ab-f8e32c1d620e	prod-agent	23.60	18.59	52.49	2025-11-28 12:31:14.19714+00
2a94f21e-d397-4631-b5d9-46a2bc1b50bc	prod-agent	7.71	18.50	52.49	2025-11-28 12:31:24.201083+00
b88d26bf-e97b-4af6-9e77-2ab98b96e737	prod-agent	7.32	18.47	52.49	2025-11-28 12:31:34.194143+00
0771f99d-d3e9-425d-bc52-263d1b58d2f1	prod-agent	9.86	18.54	52.49	2025-11-28 12:31:44.200091+00
41b1e91d-b04e-4d3b-9e52-907930d4832f	prod-agent	6.91	18.47	52.49	2025-11-28 12:31:54.195317+00
0d36af25-bca2-409e-a165-3447c2e038b8	prod-agent	14.11	18.57	52.49	2025-11-28 12:32:04.290129+00
0b708a84-0edc-4a80-a195-4e37aecc08fa	prod-agent	8.04	18.48	52.49	2025-11-28 12:32:14.195004+00
ee13c5e2-8f9c-452c-b11c-d3f8b03bc16c	prod-agent	9.12	18.47	52.49	2025-11-28 12:32:24.20021+00
b471ee9c-1388-4630-966d-ca9e0869c1eb	prod-agent	7.53	18.46	52.49	2025-11-28 12:32:34.20185+00
ab16782f-b6a3-4179-adc7-7092908c4388	prod-agent	9.12	18.43	52.49	2025-11-28 12:32:44.193999+00
346bf47c-0491-4011-b778-6b849990c8c0	prod-agent	8.60	18.48	52.49	2025-11-28 12:32:54.203307+00
309d8642-4699-4c97-acf6-37f13d708441	prod-agent	13.70	18.52	52.49	2025-11-28 12:33:04.200382+00
86c967a7-ebd2-440d-9575-6c7b2f34ba1d	prod-agent	8.49	18.46	52.49	2025-11-28 12:33:14.202359+00
366e67dc-bbec-466f-b338-4957c2113fe2	prod-agent	9.02	18.46	52.49	2025-11-28 12:33:24.19398+00
5d205ce1-ba65-498c-a0d4-90e77cdf39a4	prod-agent	7.57	18.49	52.49	2025-11-28 12:33:34.198401+00
3c99267a-ff34-4930-acd7-a6163cd80527	prod-agent	9.86	18.48	52.49	2025-11-28 12:33:44.200564+00
b01e1ccc-31d6-47d6-968a-384753f19191	prod-agent	7.45	18.53	52.49	2025-11-28 12:33:54.194362+00
4a49e970-37cd-4e2b-8096-5e868af879af	prod-agent	19.07	19.03	52.49	2025-11-28 12:34:04.202404+00
86da23d1-b840-4590-bda6-4f4e7d018d89	prod-agent	21.31	18.67	52.49	2025-11-28 12:34:14.193092+00
bfc1d1a9-6164-4abb-b17b-197ce9e6c313	prod-agent	8.53	18.60	52.49	2025-11-28 12:34:24.208383+00
d28752ba-4d17-402e-b744-82ec308111eb	prod-agent	6.85	18.58	52.49	2025-11-28 12:34:34.199625+00
df57b604-d8de-4aa7-a71d-2ef630c3398a	prod-agent	8.96	18.55	52.49	2025-11-28 12:34:44.195226+00
b5da6dcb-0126-4807-8f25-0b4b698cac5c	prod-agent	8.85	18.61	52.49	2025-11-28 12:34:54.252679+00
aba16a93-91e1-4734-b722-5dd4f3e11d69	prod-agent	8.00	18.56	52.49	2025-11-28 12:35:04.196902+00
a347103b-50f4-4c8a-95e5-a926ec0e1b69	prod-agent	14.56	18.59	52.49	2025-11-28 12:35:14.193371+00
d4452327-3ca6-4d70-be01-0c9ad8506568	prod-agent	9.32	18.69	52.49	2025-11-28 12:35:24.200116+00
306948c8-bdb3-4897-9e06-92c44706ad8b	prod-agent	7.87	18.61	52.49	2025-11-28 12:35:34.196959+00
5abca549-7112-4fba-a10a-63490cff0cda	prod-agent	10.05	18.58	52.49	2025-11-28 12:35:44.20228+00
fd6d3cbe-d95f-41ec-91a0-b1b7499db2d7	prod-agent	7.40	18.58	52.49	2025-11-28 12:35:54.194323+00
cef4c423-530a-44c5-93d0-6f81c7acd870	prod-agent	7.08	18.59	52.49	2025-11-28 12:36:04.205558+00
338451ee-1846-434d-add1-84f5b7cccf02	prod-agent	14.93	18.61	52.49	2025-11-28 12:36:14.195392+00
40f8ac81-be35-4b75-bc9a-41877bb6e692	prod-agent	8.25	18.57	52.49	2025-11-28 12:36:24.194084+00
a8d8a98c-b870-4919-87df-65393f542e62	prod-agent	7.95	18.60	52.49	2025-11-28 12:36:34.200682+00
22dbb580-569f-4837-b0d2-6314dee5e375	prod-agent	9.26	18.49	52.49	2025-11-28 12:36:44.203196+00
5ee67bda-2440-4746-87f4-442fa5e33bae	prod-agent	8.47	18.52	52.49	2025-11-28 12:36:54.201199+00
0dba38a7-fed4-4226-9814-077574ba5fad	prod-agent	7.71	18.44	52.49	2025-11-28 12:37:04.200045+00
d80af763-20ac-461e-9667-7950290e8176	prod-agent	14.30	18.59	52.49	2025-11-28 12:37:14.203637+00
95b767fa-5ead-42cb-b277-10097a2b3ba0	prod-agent	8.55	18.50	52.49	2025-11-28 12:37:24.194701+00
4aac6bfe-616b-4a2c-946d-9ae8fd3d0932	prod-agent	7.57	18.54	52.49	2025-11-28 12:37:34.200709+00
ab806479-dceb-4718-84ce-5ca35519d447	prod-agent	9.23	18.52	52.49	2025-11-28 12:37:44.20711+00
561aecf7-d327-4b7f-9cb0-5521bcbbac46	prod-agent	7.97	18.50	52.49	2025-11-28 12:37:54.195591+00
c0eab385-ebf1-4c72-9509-3eaa1de8f57f	prod-agent	7.11	18.46	52.49	2025-11-28 12:38:04.193874+00
db5b781a-2d88-4d9f-bbfa-db2981d2915b	prod-agent	15.96	18.54	52.49	2025-11-28 12:38:14.200163+00
0b26ff85-af69-475c-9999-49ce9b097507	prod-agent	8.46	18.52	52.49	2025-11-28 12:38:24.194802+00
e8d1293f-63b9-40aa-89b9-d58cd7dd879d	prod-agent	7.74	18.42	52.49	2025-11-28 12:38:34.19795+00
e39ebd57-1129-486f-b00b-d65dbd8c4d27	prod-agent	9.66	18.50	52.49	2025-11-28 12:38:44.194443+00
8acb6e18-8613-4821-a717-b858df90e0e3	prod-agent	7.72	18.44	52.49	2025-11-28 12:38:54.194864+00
c5773348-cbc0-416b-8e60-f6973055f64b	prod-agent	7.26	18.50	52.49	2025-11-28 12:39:04.198677+00
f66e2d9c-47de-4f0a-92f9-f9cc1590d0ca	prod-agent	12.60	18.65	52.49	2025-11-28 12:39:14.200071+00
78afab73-4ec5-4651-97b3-591b0d50efcb	prod-agent	10.15	18.48	52.49	2025-11-28 12:39:24.194054+00
fbda9353-ed4d-4243-8d50-46d263a2fe4d	prod-agent	7.26	18.48	52.49	2025-11-28 12:39:34.194235+00
033e9efc-998c-430c-8e82-19813f53e5e1	prod-agent	9.47	18.54	52.49	2025-11-28 12:39:44.200069+00
78080887-77a4-48b1-924c-98049f6a3587	prod-agent	7.69	18.50	52.49	2025-11-28 12:39:54.203185+00
132c3cf3-aadb-44a1-88e8-4c0241612f26	prod-agent	8.01	18.47	52.49	2025-11-28 12:40:04.262904+00
1901b670-7c21-4e6e-9abb-bd922061e97d	prod-agent	10.00	18.44	52.49	2025-11-28 12:40:14.193856+00
021e3abd-3382-4885-bb6f-196805c06e61	prod-agent	14.03	18.55	52.49	2025-11-28 12:40:24.193128+00
a4e254a0-3ec9-41be-9b84-576302b94fbb	prod-agent	7.44	18.52	52.49	2025-11-28 12:40:34.203552+00
50419f3c-f682-4400-8b21-b252e6f54c4a	prod-agent	10.15	18.52	52.49	2025-11-28 12:40:44.195157+00
b1df29e1-2b0b-4f2d-9a30-ecee19a5c24e	prod-agent	8.17	18.43	52.49	2025-11-28 12:40:54.199782+00
9b8ac6b8-7ec4-4fd5-9b83-4b8662ff5cd1	prod-agent	7.05	18.47	52.49	2025-11-28 12:41:04.20146+00
168f71cd-115f-4bbf-8cbb-6778d34ea250	prod-agent	9.69	18.47	52.49	2025-11-28 12:41:14.203214+00
e05b65f3-e8ed-4f0b-a2e4-e1a5b9035b1b	prod-agent	13.55	18.42	52.49	2025-11-28 12:41:24.199387+00
abb02e54-0e5c-448a-86ae-ffcc6361be60	prod-agent	7.13	18.40	52.49	2025-11-28 12:41:34.200595+00
1bf1ceb1-2969-459e-a45c-d4b54562979f	prod-agent	9.29	18.46	52.49	2025-11-28 12:41:44.194313+00
c9bed0d5-b9a1-4e8f-93dd-06a210223070	prod-agent	7.76	18.49	52.49	2025-11-28 12:41:54.196889+00
aab06cef-f31a-4db6-bb2b-d18ca83534cd	prod-agent	6.87	18.43	52.49	2025-11-28 12:42:04.199018+00
ffa11a5b-0cf6-4f39-88a0-ef3417560943	prod-agent	8.97	18.44	52.49	2025-11-28 12:42:14.200788+00
ef8ff31a-c239-4ce9-998b-067240aa254b	prod-agent	14.80	18.57	52.49	2025-11-28 12:42:24.194227+00
958a0ecd-c6ec-4b54-83dd-6915a19908b1	prod-agent	6.70	18.53	52.49	2025-11-28 12:42:34.20123+00
bdc7e2cd-41d8-403d-b991-ff2bf68c577b	prod-agent	9.49	18.50	52.49	2025-11-28 12:42:44.195106+00
c9fd4664-5e66-42ac-a8d0-078176a2b52e	prod-agent	8.16	18.47	52.49	2025-11-28 12:42:54.194913+00
ba842718-10a6-4c92-bcb7-a7dea3428916	prod-agent	7.81	18.53	52.49	2025-11-28 12:43:04.20046+00
77bd0881-9856-48f6-a956-36f7864e6611	prod-agent	10.04	18.57	52.49	2025-11-28 12:43:14.199665+00
a1efcda0-84e2-49b0-9ee1-96f7fac5f2f6	prod-agent	13.08	18.63	52.49	2025-11-28 12:43:24.199563+00
ca0883df-c8e3-425b-a224-0613e15ba5f6	prod-agent	6.04	18.62	52.49	2025-11-28 12:43:34.2014+00
cbbe0a60-7d6f-489e-bd0b-461812696e91	prod-agent	9.60	18.55	52.49	2025-11-28 12:43:44.200812+00
989bc0fa-c0d9-4fd4-ba17-cd34520a1c8d	prod-agent	7.67	18.50	52.49	2025-11-28 12:43:54.201803+00
1fa9194d-8477-42ad-bc1b-d7a107b58fc5	prod-agent	7.41	18.58	52.49	2025-11-28 12:44:04.194824+00
1ca1f58c-b1b2-4040-9f4a-f3479dc57cd2	prod-agent	9.36	18.59	52.49	2025-11-28 12:44:14.199633+00
42b6875b-d87a-4a1c-9ed1-ac17c310a3f4	prod-agent	12.73	18.65	52.49	2025-11-28 12:44:24.196643+00
e4efb429-3618-4741-8c44-5a0954a16595	prod-agent	9.54	18.61	52.49	2025-11-28 12:44:34.194271+00
c05c46ef-4a71-44d0-867c-da925193ba40	prod-agent	8.88	18.58	52.49	2025-11-28 12:44:44.194712+00
63fc4506-2581-40b6-823a-2fed91867cf6	prod-agent	9.62	18.53	52.49	2025-11-28 12:44:54.202166+00
e925d526-0eba-455c-832c-d65b17f12e03	prod-agent	7.67	18.57	52.49	2025-11-28 12:45:04.199857+00
70a0a288-d7b9-45b0-82e5-cb8c60dbf206	prod-agent	8.83	18.53	52.49	2025-11-28 12:45:14.251616+00
805d3d7e-e0cd-4952-bfce-d97f0f2d8601	prod-agent	9.67	18.58	52.49	2025-11-28 12:45:24.197687+00
8fa1ecf5-12eb-405a-849b-c1ea0b5c3bec	prod-agent	10.97	18.62	52.49	2025-11-28 12:45:34.194104+00
98295526-400f-4943-8e0d-bf323f37bfeb	prod-agent	9.06	18.51	52.49	2025-11-28 12:45:44.199594+00
006b7c6f-dd86-4a6a-bb06-43774104cff1	prod-agent	9.13	18.50	52.49	2025-11-28 12:45:54.200609+00
b7acd4c0-9acf-428d-bbd3-a24b471227f3	prod-agent	8.18	18.62	52.49	2025-11-28 12:46:04.194104+00
65dd4bd3-1887-4b56-ae75-a882fb10e122	prod-agent	26.92	18.72	52.49	2025-11-28 12:46:14.201288+00
3f92d42b-5a53-463f-ac2b-7513a0b11745	prod-agent	7.27	18.56	52.49	2025-11-28 12:46:24.196665+00
b9d7a5ae-4897-4691-83df-ae41a53236f8	prod-agent	15.34	18.57	52.49	2025-11-28 12:46:34.199079+00
2cbbc8a5-6b6a-44b1-b154-0f1d981d8fda	prod-agent	7.91	18.49	52.49	2025-11-28 12:46:44.193861+00
dc716c7a-a812-4bcc-a2f5-8906f8012b18	prod-agent	10.46	18.51	52.49	2025-11-28 12:46:54.203501+00
6890fb88-203e-4e39-a76d-94bb7efa0fd0	prod-agent	7.70	18.50	52.49	2025-11-28 12:47:04.19639+00
a85c7815-5041-4389-b69d-fbebed993971	prod-agent	8.06	18.50	52.49	2025-11-28 12:47:14.203955+00
88211cc4-d154-45e5-adee-1f783f905645	prod-agent	10.49	18.46	52.49	2025-11-28 12:47:24.19589+00
75e0be3f-0f89-42ef-b5f3-767513476603	prod-agent	12.32	18.55	52.49	2025-11-28 12:47:34.199414+00
5c8328c6-fc9a-4bf7-a394-9156ca022483	prod-agent	7.00	18.55	52.49	2025-11-28 12:47:44.203233+00
38d8dca5-d716-426f-978e-b78fcbb03d0a	prod-agent	9.18	18.56	52.49	2025-11-28 12:47:54.193804+00
1151f9e5-ff42-4067-8947-e4b3229dbb3b	prod-agent	6.43	18.57	52.49	2025-11-28 12:48:04.210847+00
d6546722-0e2c-403e-8203-f3a048e22927	prod-agent	8.27	18.61	52.49	2025-11-28 12:48:14.20028+00
aba045d9-ef05-468b-8ac0-d22f7ac131e0	prod-agent	10.75	18.55	52.49	2025-11-28 12:48:24.196192+00
9e089701-468f-4f94-b6a5-e9e55a65f05e	prod-agent	13.65	18.59	52.49	2025-11-28 12:48:34.201271+00
0a7debf5-85bd-4074-898c-0ad838da44a6	prod-agent	7.41	18.61	52.49	2025-11-28 12:48:44.195044+00
5fc1928f-8f8a-426d-aef3-dbe89c3d5c46	prod-agent	9.62	18.59	52.49	2025-11-28 12:48:54.19993+00
e41fc377-3858-4b0d-8732-8cfad3661ead	prod-agent	7.46	18.55	52.49	2025-11-28 12:49:04.200611+00
a50aa22b-031c-4552-b36d-fb1137a324c3	prod-agent	8.05	18.57	52.49	2025-11-28 12:49:14.201207+00
82a122de-2e04-4460-91c7-9f862d7d307e	prod-agent	9.62	18.63	52.49	2025-11-28 12:49:24.198691+00
82e71b53-c4dd-4084-a5c5-9ca778072728	prod-agent	8.70	18.59	52.49	2025-11-28 12:49:34.195781+00
b23af5ad-ed8c-4303-aa7d-a9c2d8b17887	prod-agent	12.03	18.62	52.49	2025-11-28 12:49:44.198232+00
5a9b1c87-0e86-4cef-8bc1-6ee629ea1652	prod-agent	9.20	18.58	52.49	2025-11-28 12:49:54.196286+00
2cbfd07e-93c5-4254-a59a-31d9536ec95b	prod-agent	8.02	18.60	52.49	2025-11-28 12:50:04.20005+00
48e74a87-7bc7-4b6f-ad18-def069b7d4a0	prod-agent	8.68	18.56	52.49	2025-11-28 12:50:14.194476+00
f72d072f-2fb7-4b1e-8845-0759e89114fd	prod-agent	9.98	18.63	52.49	2025-11-28 12:50:24.204579+00
2a20100e-38ba-43e5-8fab-fcf6836b5440	prod-agent	8.41	18.59	52.49	2025-11-28 12:50:34.200254+00
d5b20b44-11c6-4ae6-a58c-367c7893c66d	prod-agent	13.15	18.61	52.49	2025-11-28 12:50:44.198609+00
b1de98bc-d22e-45c4-9187-0d742b7fbbc1	prod-agent	9.08	18.57	52.49	2025-11-28 12:50:54.195176+00
8eef98a6-f32e-4438-a902-a7f0a9269bbc	prod-agent	8.12	18.62	52.49	2025-11-28 12:51:04.194295+00
a8b536c7-5b53-4a4a-b9c0-471b71249e28	prod-agent	8.27	18.61	52.49	2025-11-28 12:51:14.200108+00
1802473e-f1fc-4ed4-8cfd-ef9e0a2dd7d3	prod-agent	9.60	18.60	52.49	2025-11-28 12:51:24.200812+00
2ed52177-c9d7-465c-9e25-2e9527e104b0	prod-agent	7.62	18.59	52.49	2025-11-28 12:51:34.19653+00
e74fafa7-46d8-4a27-bbdb-12443a1312b7	prod-agent	14.55	18.60	52.49	2025-11-28 12:51:44.200254+00
20d0fe3b-7f36-4930-89f3-ee43f6602773	prod-agent	9.15	18.61	52.49	2025-11-28 12:51:54.201422+00
4ec13880-7f54-4911-ab8a-550c615fbab5	prod-agent	8.53	18.58	52.49	2025-11-28 12:52:04.203288+00
02a6e2a5-2de2-4e4b-8ca3-60a36d39a2fd	prod-agent	7.83	18.63	52.49	2025-11-28 12:52:14.198922+00
44594e57-cc87-4cc6-9bac-6a7a64b908bf	prod-agent	9.43	18.59	52.49	2025-11-28 12:52:24.194479+00
fc7d94d9-d8ed-4dbe-a97e-a6cf83342c42	prod-agent	8.09	18.62	52.49	2025-11-28 12:52:34.194399+00
384fe705-979b-47a7-874b-0e83857a208d	prod-agent	13.96	18.73	52.49	2025-11-28 12:52:44.194515+00
dc3079f8-7481-42b5-8d45-ec77f33fe92c	prod-agent	8.94	18.63	52.49	2025-11-28 12:52:54.195276+00
dd64a627-b9e3-4e1b-b220-bae4e91f5f4d	prod-agent	8.28	18.57	52.49	2025-11-28 12:53:04.202604+00
390fe5aa-f71c-40ee-b693-aa27732bb034	prod-agent	8.13	18.52	52.49	2025-11-28 12:53:14.194885+00
891ffd49-4968-40be-906b-917bda17b808	prod-agent	9.12	18.46	52.49	2025-11-28 12:53:24.200646+00
6a4c63fe-628e-4348-8635-0ff2b2cd2836	prod-agent	8.34	18.47	52.49	2025-11-28 12:53:34.200484+00
f34db36f-df02-422a-9eec-2b23f4e85450	prod-agent	12.86	18.53	52.49	2025-11-28 12:53:44.201911+00
4d627a40-8cf7-4ca0-ad22-8427012906b0	prod-agent	9.26	18.46	52.49	2025-11-28 12:53:54.206834+00
b274698e-3ace-4610-86c4-523619d0c084	prod-agent	7.75	18.39	52.49	2025-11-28 12:54:04.194604+00
c36bc26f-8047-460b-9853-55d844af298a	prod-agent	7.96	18.44	52.49	2025-11-28 12:54:14.199853+00
8f69dc4f-dd3f-4b25-b424-74918c1b6c6f	prod-agent	10.05	18.48	52.49	2025-11-28 12:54:24.20017+00
2b48049a-c660-4676-85f2-6d7466330688	prod-agent	7.26	18.48	52.49	2025-11-28 12:54:34.197929+00
bfc1e803-9b4d-4868-83e9-596a577d0da2	prod-agent	8.13	18.50	52.49	2025-11-28 12:54:44.198664+00
a29beee5-0851-4c51-bfb1-02ba7694aa89	prod-agent	14.67	18.48	52.49	2025-11-28 12:54:54.200314+00
200242e4-4fc4-4e0d-8a2a-5459eae4c77e	prod-agent	8.42	18.55	52.49	2025-11-28 12:55:04.193672+00
6e9cde38-ef9c-4a0e-b036-fcec26826063	prod-agent	7.41	18.55	52.49	2025-11-28 12:55:14.204518+00
e69da5f7-da97-432e-991e-2658d2f6cf65	prod-agent	8.86	18.52	52.49	2025-11-28 12:55:24.244923+00
ab09da3a-ae0a-4863-b8de-ccf7f8ad499b	prod-agent	8.17	18.58	52.49	2025-11-28 12:55:34.194522+00
c9128725-248f-4e91-b5a8-36ce8688d9f8	prod-agent	7.96	18.56	52.49	2025-11-28 12:55:44.203523+00
114184fb-e60e-4488-8f01-d8d7ad72d0de	prod-agent	15.35	18.62	52.49	2025-11-28 12:55:54.196861+00
5a704b26-3d08-47d2-8fb9-838a1fca8040	prod-agent	7.51	18.53	52.49	2025-11-28 12:56:04.193989+00
5ebeb8d9-712e-4868-837b-58d72655dbe6	prod-agent	8.38	18.48	52.49	2025-11-28 12:56:14.201016+00
076de176-2918-48df-b9ca-97da5fc9b8b6	prod-agent	10.40	18.44	52.49	2025-11-28 12:56:24.194292+00
ad75d635-4a5d-4a4f-8423-6ade4f79364b	prod-agent	8.87	18.38	52.49	2025-11-28 12:56:34.203469+00
e2974bbd-d315-4d66-bf9c-d817097de877	prod-agent	6.58	18.32	52.49	2025-11-28 12:56:44.194334+00
5afbc4bc-7ae5-45a3-b36f-6f386d1c6134	prod-agent	14.66	18.46	52.49	2025-11-28 12:56:54.199196+00
35e11073-acf9-414f-bcd7-cf7eaf61e5d2	prod-agent	7.95	18.37	52.49	2025-11-28 12:57:04.199581+00
9ac16b6c-ce4d-45d2-83c2-ff22a1be5ca9	prod-agent	7.66	18.41	52.49	2025-11-28 12:57:14.208096+00
851d9249-d574-4399-9b11-5d48fd407c4e	prod-agent	9.10	18.37	52.49	2025-11-28 12:57:24.194009+00
8e18a293-a21e-45e7-b381-68d59a81df21	prod-agent	8.04	18.39	52.49	2025-11-28 12:57:34.199648+00
fb8e6473-88ab-41b3-824c-0896366f1853	prod-agent	7.75	18.38	52.49	2025-11-28 12:57:44.200665+00
e48d6f4f-545e-4f5f-accd-f1a66f20e5ea	prod-agent	15.14	18.58	52.49	2025-11-28 12:57:54.202584+00
6d8c9593-385f-4340-a112-4313bf9e7478	prod-agent	7.21	18.48	52.49	2025-11-28 12:58:04.195977+00
88d2daf3-84b4-4787-a8d6-162c5c8a3cb9	prod-agent	8.22	18.50	52.49	2025-11-28 12:58:14.200564+00
5861f860-1d7b-4bae-9490-fd61b4e8776f	prod-agent	10.18	18.47	52.49	2025-11-28 12:58:24.201162+00
97ef6149-732b-4cd3-b00a-82e42d1579bf	prod-agent	7.75	18.49	52.49	2025-11-28 12:58:34.203656+00
e1c4dea7-106d-406f-a550-c75ba7c38439	prod-agent	7.92	18.44	52.49	2025-11-28 12:58:44.200378+00
03d95626-de6a-4100-a320-a442ef43bfd5	prod-agent	12.58	18.57	52.49	2025-11-28 12:58:54.204825+00
9d643fd3-f3c8-49c7-8351-2b607625ecc0	prod-agent	8.77	18.49	52.49	2025-11-28 12:59:04.193811+00
3802eccf-6fca-4ac4-9b87-ef3d289c8150	prod-agent	6.95	18.52	52.49	2025-11-28 12:59:14.195504+00
453e7f0b-c26d-4c53-a51c-3d957eb0b5b7	prod-agent	9.94	18.55	52.49	2025-11-28 12:59:24.19498+00
997b06eb-11b0-4a1c-9b08-56238aa56aa4	prod-agent	7.95	18.50	52.49	2025-11-28 12:59:34.19548+00
807dbd21-7013-413b-8f6b-2ce68d3ccee8	prod-agent	8.61	18.58	52.49	2025-11-28 12:59:44.198146+00
76184523-b5cc-430f-b93a-285dfddbe627	prod-agent	8.85	18.54	52.49	2025-11-28 12:59:54.202808+00
afc538f6-1efc-4a9a-9bb5-f6b3c983193a	prod-agent	12.25	18.55	52.49	2025-11-28 13:00:04.193329+00
1ed2438d-0878-4551-ab48-1a8829ec5b9d	prod-agent	10.19	19.02	52.49	2025-11-28 13:00:14.19368+00
770c3ba3-9049-4247-9c48-94c6ba96fd46	prod-agent	9.56	19.05	52.49	2025-11-28 13:00:24.197645+00
6ba9a36b-7e46-4c5b-b73b-83778c42ad29	prod-agent	8.08	18.97	52.49	2025-11-28 13:00:34.252813+00
c3277960-cd35-48a5-8c6d-4da91d12bd33	prod-agent	8.18	18.99	52.49	2025-11-28 13:00:44.194709+00
48bb739a-d765-4a09-9829-052261c62e3e	prod-agent	9.01	19.05	52.49	2025-11-28 13:00:54.203698+00
a87e05ec-c940-462f-aa15-e8609d074160	prod-agent	13.97	18.95	52.49	2025-11-28 13:01:04.199328+00
d9e572cc-52ff-4252-92d7-3b66e3544dd8	prod-agent	23.02	19.35	52.49	2025-11-28 13:01:14.201222+00
91a4c052-84cc-4d9c-8bf4-6d055fee6ffd	prod-agent	6.48	18.98	52.49	2025-11-28 13:01:24.203487+00
236abe2d-3049-4fa3-a23c-4f19124720ee	prod-agent	8.68	19.00	52.49	2025-11-28 13:01:34.194788+00
467a36a4-6bd5-4add-9a12-c1bde8b063e0	prod-agent	8.09	18.98	52.49	2025-11-28 13:01:44.193413+00
8e6c63b4-f448-45b9-b10c-b443344a8bdf	prod-agent	9.63	19.00	52.49	2025-11-28 13:01:54.199214+00
a5be18e3-7c9b-41ac-afe3-fd52f297a460	prod-agent	13.98	19.10	52.49	2025-11-28 13:02:04.194203+00
db038318-932b-45ac-a62e-ebc625b03462	prod-agent	7.27	18.89	52.49	2025-11-28 13:02:14.194274+00
e66bc24e-afcc-43c5-a931-17cfbe1a59f6	prod-agent	9.70	18.98	52.49	2025-11-28 13:02:24.200394+00
26f9fc73-7888-4906-8c8b-c9393f7cca36	prod-agent	9.04	18.98	52.49	2025-11-28 13:02:34.194184+00
40c54147-8e4c-46b4-80ed-8e3d3ba74f79	prod-agent	8.30	18.83	52.49	2025-11-28 13:02:44.194611+00
d2f619a3-8889-42d1-a918-1e0d19581f83	prod-agent	8.60	18.83	52.49	2025-11-28 13:02:54.201522+00
80d2c48e-222c-447f-a1bd-3d0c82b4ad13	prod-agent	14.74	18.84	52.49	2025-11-28 13:03:04.196111+00
c8494dfd-9de2-4ae6-8493-e579b6899932	prod-agent	6.05	18.85	52.49	2025-11-28 13:03:14.200964+00
1dbc517e-741d-47ac-aa3f-312780282c6b	prod-agent	9.09	18.84	52.49	2025-11-28 13:03:24.205783+00
7108defb-6892-42de-a69a-fc80258a6161	prod-agent	8.98	18.81	52.49	2025-11-28 13:03:34.200237+00
4fc8decf-47ed-4075-8ac1-78eb7e4c23b7	prod-agent	8.01	18.77	52.49	2025-11-28 13:03:44.194406+00
6c8f4e28-b170-47ff-82d6-868e880e1ad2	prod-agent	9.05	18.81	52.49	2025-11-28 13:03:54.195049+00
d1d84e0e-8f68-4a07-8d67-84c5bef903d7	prod-agent	18.53	19.29	52.49	2025-11-28 13:04:04.201021+00
d04f4c3a-0c41-4b1c-b054-1cf9cedc566b	prod-agent	21.10	18.55	52.49	2025-11-28 13:04:14.191889+00
b36e179b-4abb-43ae-bea2-6d3a9f481872	prod-agent	9.09	18.50	52.49	2025-11-28 13:04:24.200553+00
1bfa7680-bc7e-4c93-9192-760b28170b1f	prod-agent	9.28	18.47	52.49	2025-11-28 13:04:34.201863+00
4239feb2-6c07-48e5-962d-f42b3a680b78	prod-agent	7.87	18.50	52.49	2025-11-28 13:04:44.196877+00
09eb6356-b988-4e96-b2b1-f41525c8d5c6	prod-agent	8.50	18.48	52.49	2025-11-28 13:04:54.196899+00
8a6c2670-63e3-46bb-ad16-689aa919001d	prod-agent	8.63	18.51	52.49	2025-11-28 13:05:04.201313+00
a6ec5f58-c443-4d27-bb3f-b58334ecb8c4	prod-agent	12.97	18.51	52.49	2025-11-28 13:05:14.201037+00
226bc436-7fd5-4fd5-925c-f12e0838852e	prod-agent	9.31	18.48	52.49	2025-11-28 13:05:24.199519+00
d59483ad-76bb-44ef-be1e-6267e15b49f5	prod-agent	7.94	18.50	52.49	2025-11-28 13:05:34.203874+00
1c6a2d8d-92f3-4671-b1b4-e8619f4d6433	prod-agent	7.95	18.47	52.49	2025-11-28 13:05:44.257093+00
1bd1b95c-6e62-43e8-8968-e9ae32b5d9db	prod-agent	8.79	18.51	52.49	2025-11-28 13:05:54.194785+00
b4a4be46-285c-4571-9816-da7635202207	prod-agent	8.82	18.52	52.49	2025-11-28 13:06:04.201678+00
c6a3d228-8cae-4312-a77f-46d54245abc1	prod-agent	13.40	18.53	52.49	2025-11-28 13:06:14.193584+00
b01eba33-6c7e-44a9-9f86-e0ec7eb00ed0	prod-agent	9.13	18.48	52.49	2025-11-28 13:06:24.200249+00
23c78605-33e0-48a8-a36f-83633db7c252	prod-agent	9.01	18.51	52.49	2025-11-28 13:06:34.199947+00
3889379c-366b-42f5-918e-387161c16a32	prod-agent	7.93	18.54	52.49	2025-11-28 13:06:44.21085+00
c52f97d1-b7e8-4404-a67f-fe5abdc947b8	prod-agent	8.73	18.50	52.49	2025-11-28 13:06:54.20339+00
f1e0ca30-b0a9-427e-9c92-798f901ffc2d	prod-agent	7.62	18.57	52.49	2025-11-28 13:07:04.200939+00
c5e8e3b1-3208-4f4a-a447-429e7886d1b8	prod-agent	13.13	18.62	52.49	2025-11-28 13:07:14.194305+00
4e39bb99-f825-4c99-99a9-77b3ab3fb0d4	prod-agent	8.96	18.56	52.49	2025-11-28 13:07:24.199644+00
96e60ab6-e594-4839-8f3d-194b35451f57	prod-agent	9.25	18.56	52.49	2025-11-28 13:07:34.200639+00
ebf57a10-6a37-44ac-88db-26e85434ec5a	prod-agent	7.46	18.56	52.49	2025-11-28 13:07:44.203691+00
9bce7ab6-d69b-48db-afe5-ec193e1d8edc	prod-agent	8.45	18.57	52.49	2025-11-28 13:07:54.201123+00
a8638104-a625-4cbf-bf4c-7d37bbbd72da	prod-agent	9.28	18.60	52.49	2025-11-28 13:08:04.200019+00
f168d99d-a043-4eab-9472-0a8c3a0ef945	prod-agent	14.51	18.63	52.49	2025-11-28 13:08:14.194653+00
0a9253d0-c4dc-4fdc-9b26-61a230c6726d	prod-agent	8.39	18.52	52.49	2025-11-28 13:08:24.20045+00
55fb4d78-324f-40a8-83fc-8a3d06c07c92	prod-agent	9.50	18.56	52.49	2025-11-28 13:08:34.199949+00
74c55829-5b81-4aaa-843c-d4ba04b32a8b	prod-agent	7.41	18.58	52.49	2025-11-28 13:08:44.195265+00
e1529c6b-2858-4d4f-a16e-48b653922510	prod-agent	8.24	18.59	52.49	2025-11-28 13:08:54.20967+00
eaa447dc-041e-4426-895e-837446335490	prod-agent	8.50	18.57	52.49	2025-11-28 13:09:04.197676+00
109ccf86-5e8b-4cb5-aa9d-42460bcdae08	prod-agent	11.34	18.58	52.49	2025-11-28 13:09:14.199507+00
06079d36-6216-4559-94bf-1d97b8e91342	prod-agent	10.29	18.59	52.49	2025-11-28 13:09:24.199684+00
afba25c2-1be0-4f3f-8859-21196b36733d	prod-agent	9.27	18.56	52.49	2025-11-28 13:09:34.199985+00
7750c7d1-3bcf-470a-a2d1-bed0e8211a6a	prod-agent	7.71	18.62	52.49	2025-11-28 13:09:44.202807+00
99086463-d134-4375-b81c-75ae68ffc780	prod-agent	8.80	18.53	52.49	2025-11-28 13:09:54.200976+00
033e08c4-18f8-4471-93df-b05391efd7f4	prod-agent	8.80	18.56	52.49	2025-11-28 13:10:04.198389+00
c1714f8c-e941-4395-b547-7b6da9f33afd	prod-agent	8.22	18.52	52.49	2025-11-28 13:10:14.195304+00
81b92f13-a83c-4db7-bc20-9212b3916b04	prod-agent	13.91	18.43	52.49	2025-11-28 13:10:24.203276+00
36533efe-d3c2-493e-b311-e1e9a5b11232	prod-agent	8.72	18.49	52.49	2025-11-28 13:10:34.194512+00
c3bb136a-5ee6-4218-84f2-e8ee4732b220	prod-agent	7.21	18.46	52.49	2025-11-28 13:10:44.200697+00
a66d75ce-7a76-4e03-bb8b-a9f51c89d548	prod-agent	7.51	18.46	52.49	2025-11-28 13:10:54.291428+00
045ba894-bb7e-4654-8e7b-5ca0d57cdf42	prod-agent	9.54	18.45	52.49	2025-11-28 13:11:04.194988+00
7637d82e-9717-45b1-a887-ef2f7c6b4053	prod-agent	8.39	18.46	52.49	2025-11-28 13:11:14.195008+00
40d5af78-4fc2-47d7-92ec-227cd05f67c3	prod-agent	13.80	18.47	52.49	2025-11-28 13:11:24.193554+00
92e84265-d073-43a2-a0eb-4d99acb6fb64	prod-agent	9.76	18.43	52.49	2025-11-28 13:11:34.199328+00
07e8a558-f6bd-437c-aa79-f6ca46ae0a83	prod-agent	7.15	18.41	52.49	2025-11-28 13:11:44.198425+00
96727beb-76c1-4193-985d-0ea1fdc6b264	prod-agent	7.79	18.45	52.49	2025-11-28 13:11:54.195769+00
4ffb6216-351a-4c6d-8e2f-a4908d3ff3e8	prod-agent	9.12	18.53	52.49	2025-11-28 13:12:04.199467+00
02ba03c8-ac94-4b01-92e0-cac6988844b9	prod-agent	7.74	18.54	52.49	2025-11-28 13:12:14.207509+00
f3937df1-1af7-404f-90fc-4697f88dd2b2	prod-agent	14.12	18.46	52.49	2025-11-28 13:12:24.200094+00
db175728-da0d-4e0e-b430-f6f4f945c92e	prod-agent	8.92	18.49	52.49	2025-11-28 13:12:34.204911+00
11bdcbb1-aecf-497e-8f06-22967ad0645d	prod-agent	8.18	18.49	52.49	2025-11-28 13:12:44.194926+00
f7f3903c-7a8d-4d11-b05b-ab1b39cd467d	prod-agent	7.00	18.51	52.49	2025-11-28 13:12:54.202473+00
1ca3ecc3-761d-478d-b671-6db65235c579	prod-agent	10.07	18.49	52.49	2025-11-28 13:13:04.194542+00
b34dafd2-6586-4a72-8474-149d3cc668db	prod-agent	8.09	18.48	52.49	2025-11-28 13:13:14.200583+00
2b0a8c65-60af-4767-a30c-edf0fb08cb00	prod-agent	14.70	18.57	52.49	2025-11-28 13:13:24.194466+00
df4d4f17-1c5e-4f09-a4c5-0494ad23e978	prod-agent	9.77	18.53	52.49	2025-11-28 13:13:34.200798+00
55967a46-56cb-4879-ad81-a5212ea8cc6a	prod-agent	7.28	18.53	52.49	2025-11-28 13:13:44.198588+00
50be47ff-4c3c-4924-a496-6693cf0e6b36	prod-agent	8.15	18.54	52.49	2025-11-28 13:13:54.194466+00
6f36f97d-2c46-4e27-bbbc-8936d08e49c2	prod-agent	9.81	18.54	52.49	2025-11-28 13:14:04.200089+00
d6700348-355f-490e-952d-f9dacc74e8e6	prod-agent	7.26	18.56	52.49	2025-11-28 13:14:14.194008+00
e1038192-06c8-40ec-a337-f6b9b0adf157	prod-agent	9.01	18.56	52.49	2025-11-28 13:14:24.204008+00
2053467c-80b9-4dea-bfb1-121aea90324a	prod-agent	13.42	18.53	52.49	2025-11-28 13:14:34.193189+00
b79bda44-dac8-43ee-87b0-32d9e7dd88f0	prod-agent	7.23	18.54	52.49	2025-11-28 13:14:44.194223+00
d61ea1ee-52bb-4f4a-a85d-f75178b0bd7a	prod-agent	7.07	18.53	52.49	2025-11-28 13:14:54.202502+00
550ead94-55b3-4c2b-88f1-f82a90c2e849	prod-agent	9.96	18.50	52.49	2025-11-28 13:15:04.196051+00
88115efc-ed90-4f94-9450-81684766fa4d	prod-agent	7.59	18.47	52.49	2025-11-28 13:15:14.199033+00
71c4a6aa-2649-49f5-b436-02e6f71bfe62	prod-agent	8.46	18.53	52.49	2025-11-28 13:15:24.202935+00
e5a66f27-935b-4177-9668-e6ade7fc47d2	prod-agent	14.39	18.64	52.49	2025-11-28 13:15:34.197905+00
543d23f5-93ad-4387-af15-f4e7045fa039	prod-agent	7.97	18.54	52.49	2025-11-28 13:15:44.200458+00
5564bba6-92f1-4c60-bf7b-2ffd6655c0ab	prod-agent	8.12	18.49	52.49	2025-11-28 13:15:54.195776+00
599a451e-da9e-4a1d-95af-f1ce10405fbc	prod-agent	9.76	18.61	52.49	2025-11-28 13:16:04.262448+00
0dc26f91-0dff-42c5-9150-de44c5196638	prod-agent	19.46	19.10	52.49	2025-11-28 13:16:14.199771+00
229e2b06-ab13-40c1-b17d-a74de5a94a35	prod-agent	8.69	18.56	52.49	2025-11-28 13:16:24.192852+00
7787cde3-4869-46da-bbd1-adb5e83978c7	prod-agent	14.54	18.58	52.49	2025-11-28 13:16:34.206473+00
0e173da5-5ae2-4c1a-aa80-28422dcb36b3	prod-agent	8.74	18.56	52.49	2025-11-28 13:16:44.195053+00
2f9f0e85-fc24-4111-b8ef-696486266917	prod-agent	5.84	18.48	52.49	2025-11-28 13:16:54.199183+00
f1115db4-9dc9-4b94-8942-7a90168bc391	prod-agent	9.56	18.57	52.49	2025-11-28 13:17:04.203205+00
808d4c65-a099-4566-b23f-e7273b0ba86d	prod-agent	9.12	18.61	52.49	2025-11-28 13:17:14.203706+00
1f54e22b-f5af-4ad6-a02f-85e8c5a671d9	prod-agent	8.20	18.57	52.49	2025-11-28 13:17:24.20071+00
9282194e-f049-4362-90f8-e9a92fc4189b	prod-agent	16.78	18.59	52.49	2025-11-28 13:17:34.193988+00
51759488-6786-4562-97f9-55ca3c210f8d	prod-agent	8.06	18.61	52.49	2025-11-28 13:17:44.194918+00
b03f22a0-4697-43d4-83e5-118dfd7940cd	prod-agent	7.88	18.55	52.49	2025-11-28 13:17:54.203482+00
17e1cc7f-aae1-4599-b51d-99196fb37102	prod-agent	8.47	18.57	52.49	2025-11-28 13:18:04.20081+00
f556e69f-b7ae-474a-bebc-f41f04fefd20	prod-agent	8.94	18.54	52.49	2025-11-28 13:18:14.200184+00
a9c93670-9752-49af-b6fe-8e4cb114608e	prod-agent	8.62	18.49	52.49	2025-11-28 13:18:24.194623+00
32c67261-f9cd-4f78-b670-8969e285a501	prod-agent	12.89	18.70	52.49	2025-11-28 13:18:34.204043+00
be456295-e8d7-4b63-8e4f-610e17243fd4	prod-agent	9.35	18.56	52.49	2025-11-28 13:18:44.19929+00
822bdc92-7f3c-4f62-9514-e638f670aee8	prod-agent	6.22	18.57	52.49	2025-11-28 13:18:54.201136+00
ecf7de21-9f9f-485e-b647-c99aeb2cb4d3	prod-agent	9.25	18.60	52.49	2025-11-28 13:19:04.200884+00
6ccd51b4-74da-4f28-be80-0265160e3986	prod-agent	8.77	18.64	52.49	2025-11-28 13:19:14.194649+00
e3721f70-45f7-443f-8b75-7f63a8d00203	prod-agent	8.65	18.58	52.49	2025-11-28 13:19:24.200772+00
3b76d3bc-e201-461e-bb0f-5b99f8e2e118	prod-agent	9.11	18.55	52.49	2025-11-28 13:19:34.200655+00
601301a9-4e87-4e4a-a887-bb57ecae6bc6	prod-agent	14.28	18.60	52.49	2025-11-28 13:19:44.200278+00
4363a0f9-f90f-496e-a8d9-cb0eb9ae9ef4	prod-agent	8.16	18.60	52.49	2025-11-28 13:19:54.198775+00
6c53a192-545b-4586-b396-832254c9fe7e	prod-agent	8.62	18.58	52.49	2025-11-28 13:20:04.196263+00
e05fe178-0917-49f1-a4aa-93cf1d454a18	prod-agent	9.47	18.57	52.49	2025-11-28 13:20:14.205386+00
df7bec62-8663-4de0-9219-d1ea80b62ce6	prod-agent	8.29	18.54	52.49	2025-11-28 13:20:24.193616+00
878223c1-4425-404a-b6d3-3ef790d301d6	prod-agent	7.51	18.59	52.49	2025-11-28 13:20:34.205561+00
321914c9-744d-4b19-9d4a-e7c3e85ad222	prod-agent	15.07	18.40	52.49	2025-11-28 13:20:44.195539+00
c8792dc7-0c9d-4a2c-aeb1-ccbe74f2d34f	prod-agent	7.27	18.45	52.49	2025-11-28 13:20:54.194121+00
2f48d699-afd0-484d-b192-22e381fc4002	prod-agent	7.86	18.47	52.49	2025-11-28 13:21:04.204343+00
ef7299bf-0f42-4ab7-9de4-d6cb6ee4d17c	prod-agent	9.26	18.45	52.49	2025-11-28 13:21:14.251561+00
451a0d22-2a91-4035-a94a-7594f8d9ef69	prod-agent	8.58	18.48	52.49	2025-11-28 13:21:24.197532+00
dcb815b9-1b12-4601-8526-17e3421e9e21	prod-agent	8.54	18.45	52.49	2025-11-28 13:21:34.194366+00
2a0fbc73-e113-403d-98d0-397c8023115a	prod-agent	15.21	18.55	52.49	2025-11-28 13:21:44.20019+00
04284514-99f6-4f19-80c8-24d131ef8b9f	prod-agent	7.82	18.51	52.49	2025-11-28 13:21:54.194246+00
9f54ec9e-109b-48d5-8f0c-a0024f023a5f	prod-agent	8.44	18.46	52.49	2025-11-28 13:22:04.199995+00
1d30cec8-9a0a-41dd-9ad0-a36cf42b23a7	prod-agent	9.51	18.48	52.49	2025-11-28 13:22:14.199814+00
e703b7fa-338d-4dd0-8b19-3fd07ad67112	prod-agent	8.46	18.44	52.49	2025-11-28 13:22:24.194401+00
a74a1314-afb9-40c0-a2c6-a9770b3d4820	prod-agent	7.39	18.48	52.49	2025-11-28 13:22:34.198439+00
ff3671c1-9aa2-4619-8902-f043e2635009	prod-agent	14.15	18.54	52.49	2025-11-28 13:22:44.193374+00
7b8c0f50-f076-4f18-bf52-e76c4e2d0cea	prod-agent	7.26	18.52	52.49	2025-11-28 13:22:54.199132+00
8efd8875-0296-4503-8028-5d7b56b13540	prod-agent	8.44	18.44	52.49	2025-11-28 13:23:04.194044+00
20cc32f7-58d3-4209-a541-650b9b7fde67	prod-agent	9.27	18.56	52.49	2025-11-28 13:23:14.204144+00
4a817cf5-409b-42a0-aabd-8c518e80d804	prod-agent	8.39	18.52	52.49	2025-11-28 13:23:24.201512+00
b7c2aa3a-674e-46a7-8aa9-d6ed933f8499	prod-agent	8.27	18.58	52.49	2025-11-28 13:23:34.20026+00
a0e5f72e-bcfb-4c23-83f5-9d68987ac1e9	prod-agent	13.38	18.68	52.49	2025-11-28 13:23:44.202955+00
e14c5ce7-7428-457b-9017-1301d0028bde	prod-agent	9.50	18.55	52.49	2025-11-28 13:23:54.199581+00
159cd411-66c6-4360-a8bd-069c5faf4ae1	prod-agent	7.49	18.48	52.49	2025-11-28 13:24:04.204548+00
d0d6a17f-1493-4c5b-994b-f492ae58301d	prod-agent	9.26	18.54	52.49	2025-11-28 13:24:14.200396+00
858a44a5-1c4a-414c-9c31-4ddd8f4c0e5b	prod-agent	8.61	18.55	52.49	2025-11-28 13:24:24.194861+00
046fc717-005a-43bf-87f2-3d12ad67daca	prod-agent	7.36	18.49	52.49	2025-11-28 13:24:34.200484+00
0f22aa08-9233-415b-9083-00d0c60178ea	prod-agent	9.69	18.58	52.49	2025-11-28 13:24:44.19966+00
ee7d45ae-4fc8-423a-8a72-64be7937855e	prod-agent	13.49	18.57	52.49	2025-11-28 13:24:54.20129+00
fb09321b-dfca-4e7e-8e77-333650d1fb8a	prod-agent	8.70	18.52	52.49	2025-11-28 13:25:04.200515+00
23def5b8-35d4-48cd-b54f-f991dfede51a	prod-agent	9.76	18.54	52.49	2025-11-28 13:25:14.195914+00
d06a8074-5312-410b-8799-ddc45862a698	prod-agent	8.82	18.53	52.49	2025-11-28 13:25:24.200035+00
1810aa9c-02dc-4d8a-a0d5-70a17cce1171	prod-agent	7.75	18.58	52.49	2025-11-28 13:25:34.194225+00
d67c6b35-e6f6-48ea-bf58-06049a6ba047	prod-agent	9.09	18.58	52.49	2025-11-28 13:25:44.193678+00
e9e01e52-4dd3-437e-b2ab-255e9414ad59	prod-agent	14.35	18.60	52.49	2025-11-28 13:25:54.197388+00
0c33f373-b7df-4165-9b00-f9c825b7023c	prod-agent	8.01	18.57	52.49	2025-11-28 13:26:04.194184+00
0fd5fd82-df93-4cce-a189-8580e0b74b77	prod-agent	10.29	18.56	52.49	2025-11-28 13:26:14.199959+00
020f6559-8cf8-4cb4-8fa8-2054d280fcfd	prod-agent	8.34	18.54	52.49	2025-11-28 13:26:24.209342+00
67d4727e-f6ec-434b-97ea-c5a53420af6b	prod-agent	7.84	18.58	52.49	2025-11-28 13:26:34.199875+00
79e9a672-bc04-4ac7-b4ee-fe169a026a79	prod-agent	8.74	18.59	52.49	2025-11-28 13:26:44.194118+00
217c9417-c3c3-457a-be58-67dfb821c6af	prod-agent	13.57	18.65	52.49	2025-11-28 13:26:54.202641+00
cae29df1-8ff1-4cfb-b49a-a9f6e8f9683f	prod-agent	7.97	18.52	52.49	2025-11-28 13:27:04.203132+00
7627b086-8c24-4fe2-9b2b-460d7b68c9d3	prod-agent	9.57	18.53	52.49	2025-11-28 13:27:14.20187+00
d1f94413-f8d2-4918-a5e1-e5d0de0ab67a	prod-agent	8.41	18.59	52.49	2025-11-28 13:27:24.201349+00
7a4aadb0-5c5c-4d0c-9908-a0404ac25e16	prod-agent	8.20	18.52	52.49	2025-11-28 13:27:34.210077+00
f4f1d41c-18d9-4a37-9f8b-f9ad269691c1	prod-agent	9.07	18.62	52.49	2025-11-28 13:27:44.202403+00
abd03931-fa5c-4990-8ec8-2265fd0f9a21	prod-agent	13.56	18.42	52.49	2025-11-28 13:27:54.200834+00
12f2f83b-ded7-4a61-a122-caaa890d154c	prod-agent	8.01	18.40	52.49	2025-11-28 13:28:04.194192+00
c72c001e-0b8a-4c8e-84a2-a4d774ced71b	prod-agent	9.75	18.39	52.49	2025-11-28 13:28:14.199193+00
f6b953af-cd52-4d07-9419-c70243defd11	prod-agent	8.35	18.46	52.49	2025-11-28 13:28:24.206554+00
2b2b5499-e9e4-41c8-8baf-12aa687343e4	prod-agent	7.09	18.38	52.49	2025-11-28 13:28:34.194451+00
57ef26b0-a042-4071-b22a-49e35f23aca2	prod-agent	9.43	18.35	52.49	2025-11-28 13:28:44.194533+00
bf659636-4da0-428c-8cd8-930822275d08	prod-agent	10.90	18.40	52.49	2025-11-28 13:28:54.226259+00
eaba46fb-4b4e-4828-8e52-5fd02f9bbca2	prod-agent	11.01	18.33	52.49	2025-11-28 13:29:04.193792+00
61d71dab-b012-4ba0-b207-d40b95544d9f	prod-agent	9.93	18.33	52.49	2025-11-28 13:29:14.19501+00
9be0a277-412a-4d93-8a37-b5619b3958f9	prod-agent	8.61	18.34	52.49	2025-11-28 13:29:24.195034+00
874a2144-6287-4165-bb81-d980ae8c08d6	prod-agent	8.41	18.41	52.49	2025-11-28 13:29:34.199974+00
440393c3-62a5-4f1b-8e88-5955351e0390	prod-agent	9.37	18.40	52.49	2025-11-28 13:29:44.194796+00
d8f9a772-91d1-4f41-afd2-cad0dc0d521e	prod-agent	7.30	18.34	52.49	2025-11-28 13:29:54.210445+00
d125e2c3-b01b-44fe-9135-5f3b08ab0c0f	prod-agent	13.18	18.41	52.49	2025-11-28 13:30:04.194318+00
a2ec62ef-8ace-4be4-878e-38fcdb8afa87	prod-agent	10.06	18.34	52.49	2025-11-28 13:30:14.200854+00
32432cf6-f2fd-48d7-a518-9ec5d7e01736	prod-agent	7.50	18.42	52.49	2025-11-28 13:30:24.200094+00
47f7a0fd-d203-4037-bb56-7ec71885cff3	prod-agent	8.29	18.38	52.49	2025-11-28 13:30:34.196586+00
919cc575-0d17-4dae-9e8e-58beb681c373	prod-agent	10.02	18.53	52.49	2025-11-28 13:30:44.203065+00
cdb122ce-16d2-44eb-9aeb-78076f9e5ed4	prod-agent	7.58	18.57	52.49	2025-11-28 13:30:54.201478+00
3860f2c9-1a7a-41a8-bbd3-4c4707032bc3	prod-agent	13.81	18.56	52.49	2025-11-28 13:31:04.195091+00
44ec6009-3ec7-4f54-954f-684d6fa46683	prod-agent	18.92	19.10	52.49	2025-11-28 13:31:14.199177+00
ff4724a7-2912-4a57-a4b9-3575cc8a3c63	prod-agent	12.11	18.56	52.49	2025-11-28 13:31:24.229518+00
10d2c980-4217-4f2c-94d2-69b1cec888b0	prod-agent	7.69	18.58	52.49	2025-11-28 13:31:34.200714+00
905c0ba5-b4b0-4f5b-9c8f-12ce47f2bb12	prod-agent	9.57	18.57	52.49	2025-11-28 13:31:44.195225+00
151a0af7-123e-4bc1-98c9-0ad55e024b44	prod-agent	7.97	18.58	52.49	2025-11-28 13:31:54.201101+00
88d5eb39-4632-4b15-87c8-11c9fd301f20	prod-agent	14.34	18.62	52.49	2025-11-28 13:32:04.199913+00
889fd869-7083-490f-8498-4b7dfe6c5ff1	prod-agent	8.55	18.52	52.49	2025-11-28 13:32:14.197107+00
9ce88f72-e3bc-4356-8227-286c3e32bfaf	prod-agent	7.59	18.50	52.49	2025-11-28 13:32:24.196662+00
0fe2c2b0-6224-4188-8c79-419821d3a466	prod-agent	7.66	18.52	52.49	2025-11-28 13:32:34.200528+00
d1a9515e-d57d-4dbe-bf0c-32f6d53ad679	prod-agent	9.56	18.51	52.49	2025-11-28 13:32:44.194223+00
d5b81641-c303-4110-b0fe-5b38b0340d9d	prod-agent	9.34	18.51	52.49	2025-11-28 13:32:54.194137+00
582d433c-35f6-482f-9f50-8eb44f0a2ed5	prod-agent	12.95	18.64	52.49	2025-11-28 13:33:04.195073+00
67ddb28a-2422-483a-965b-77f96715c82e	prod-agent	9.40	18.57	52.49	2025-11-28 13:33:14.200625+00
96e1674e-f24a-4776-8c6c-9ad6406f61e1	prod-agent	8.28	18.54	52.49	2025-11-28 13:33:24.195667+00
a0e5e84c-0e1c-4bb3-8456-8cdd7846df96	prod-agent	7.93	18.49	52.49	2025-11-28 13:33:34.195368+00
f9248e1c-98e3-4fc9-a0f1-06cea73531ed	prod-agent	9.71	18.56	52.49	2025-11-28 13:33:44.193112+00
1faff531-e066-4fc0-bf51-03d98da4b133	prod-agent	7.73	18.55	52.49	2025-11-28 13:33:54.198199+00
976558ce-ae88-44de-9d4e-360c209b5091	prod-agent	16.88	18.88	52.49	2025-11-28 13:34:04.208136+00
d32020af-5671-4c6f-914b-9a4bf4fffc07	prod-agent	18.83	18.62	52.49	2025-11-28 13:34:14.197877+00
b1de7d7f-6441-4ddb-ba8d-13d3845ad6ec	prod-agent	6.52	18.58	52.49	2025-11-28 13:34:24.194846+00
b31e8038-cbae-4615-9317-7091715c85c6	prod-agent	8.49	18.59	52.49	2025-11-28 13:34:34.209035+00
852489cc-0a2a-4a1c-9489-0e5218fc09e4	prod-agent	9.94	18.54	52.49	2025-11-28 13:34:44.195251+00
58c32569-b6aa-4ab6-ab39-4962ca1c1769	prod-agent	8.15	18.53	52.49	2025-11-28 13:34:54.19705+00
79ec93da-ea2b-4704-83f0-3850d56793bf	prod-agent	7.60	18.53	52.49	2025-11-28 13:35:04.194567+00
fd79bc42-896f-4a05-9acf-e2ae41d59554	prod-agent	16.56	18.58	52.49	2025-11-28 13:35:14.200013+00
e05e4ffc-42cf-43fb-8378-90d52faf6e83	prod-agent	8.10	18.56	52.49	2025-11-28 13:35:24.21158+00
de4a49d9-ab4a-4ffc-b80b-8581ca087a60	prod-agent	8.00	18.59	52.49	2025-11-28 13:35:34.203878+00
ed5f6155-97df-4e5f-a944-0c4cce5b11c5	prod-agent	9.57	18.54	52.49	2025-11-28 13:35:44.203272+00
db627702-30e8-48df-b9ac-e5d8dda78bee	prod-agent	7.59	18.62	52.49	2025-11-28 13:35:54.202295+00
758e2a02-78ba-4322-a0ab-45014e2ead5c	prod-agent	7.25	18.59	52.49	2025-11-28 13:36:04.19428+00
c15c78f6-0700-4a71-b10f-56175430c2c1	prod-agent	14.94	18.59	52.49	2025-11-28 13:36:14.199211+00
495135de-0082-4ead-9134-3085bca1fd2a	prod-agent	8.16	18.57	52.49	2025-11-28 13:36:24.194578+00
7a2166c8-197a-4fec-82c9-203763fa12dd	prod-agent	6.94	18.55	52.49	2025-11-28 13:36:34.262192+00
52e7a058-30dc-4a6f-a92f-b6eaa82c5f2f	prod-agent	10.19	18.56	52.49	2025-11-28 13:36:44.194082+00
21c941a2-de27-4280-b404-38447f22ac7e	prod-agent	7.94	18.58	52.49	2025-11-28 13:36:54.201074+00
d40fea68-64b0-4e92-a52f-701a857f0e02	prod-agent	7.90	18.59	52.49	2025-11-28 13:37:04.193993+00
d3976054-9b3e-471f-b502-6dc18b03fe98	prod-agent	16.82	18.75	52.49	2025-11-28 13:37:14.201142+00
09266571-b901-4b70-8403-ce16ffdc790e	prod-agent	7.70	18.63	52.49	2025-11-28 13:37:24.197671+00
0829a3f9-1c78-48c9-bf6e-1147f678df97	prod-agent	7.85	18.64	52.49	2025-11-28 13:37:34.195946+00
e775509a-8674-4996-b1c4-52c27c2afc23	prod-agent	10.33	18.64	52.49	2025-11-28 13:37:44.200601+00
3a4058a4-98b5-490f-8b0e-4d387b2339bc	prod-agent	8.02	18.62	52.49	2025-11-28 13:37:54.195423+00
17fe47dc-941a-4363-8130-0264f31db513	prod-agent	7.18	18.63	52.49	2025-11-28 13:38:04.204491+00
9a827fbd-c39d-4001-bd35-40e84a90d9a6	prod-agent	12.96	18.75	52.49	2025-11-28 13:38:14.202807+00
c032778f-0aab-4ebb-9efb-6b83ecf604ea	prod-agent	9.72	18.65	52.49	2025-11-28 13:38:24.199353+00
e38b09c8-d48e-48ab-835f-185ba77d1287	prod-agent	7.51	18.68	52.49	2025-11-28 13:38:34.206831+00
7e0349d5-f873-4b54-a09b-2be1ea20a80b	prod-agent	9.26	18.64	52.49	2025-11-28 13:38:44.195546+00
62e09ddd-d48d-45a2-8ac8-40d5e0da389f	prod-agent	8.10	18.63	52.49	2025-11-28 13:38:54.196043+00
dcb481d9-f275-4cfe-927e-2e2ac09dac10	prod-agent	7.59	18.62	52.49	2025-11-28 13:39:04.198815+00
70dcb97d-8be0-4d22-b170-3230ad4a8cb5	prod-agent	9.74	18.66	52.49	2025-11-28 13:39:14.194444+00
2a4ea1d1-6e4f-481d-8875-88092687e42e	prod-agent	14.12	18.65	52.49	2025-11-28 13:39:24.201025+00
6b035a5c-68dc-4812-8d5e-0f379d145dc5	prod-agent	7.93	18.61	52.49	2025-11-28 13:39:34.197342+00
a1809482-74f9-43ef-8694-f1f647edee79	prod-agent	9.33	18.61	52.49	2025-11-28 13:39:44.19522+00
7892675f-ee2f-4ea1-b7bd-1950105552ad	prod-agent	8.90	18.66	52.49	2025-11-28 13:39:54.198871+00
66332a58-3c9d-4a85-a99d-d330dc7f2cb9	prod-agent	26.11	20.89	52.49	2025-11-28 13:40:04.197273+00
652661ba-b81d-4c35-b5c1-2c24087451cc	prod-agent	7.15	18.71	52.49	2025-11-28 13:40:14.193788+00
608b1da1-e303-43f8-912a-5df2bd90dca2	prod-agent	13.46	18.60	52.49	2025-11-28 13:40:24.201839+00
76ecffb4-488a-4b58-aceb-4c0a2c91c667	prod-agent	7.27	18.58	52.49	2025-11-28 13:40:34.194209+00
177757a0-48e1-41ff-8bfa-8d470b8481d4	prod-agent	9.20	18.58	52.49	2025-11-28 13:40:44.201615+00
ce3629e7-410d-4dfb-8562-2d54a7f7e443	prod-agent	8.68	18.58	52.49	2025-11-28 13:40:54.201221+00
cb29b74d-f0b6-4fcc-a501-73f6ceb59cf6	prod-agent	7.59	18.61	52.49	2025-11-28 13:41:04.194067+00
fc9376a2-46cf-4978-a2ce-d5e7268dbe3f	prod-agent	8.26	18.65	52.49	2025-11-28 13:41:14.198976+00
053c0961-4f4f-4620-81af-f0bbe20049df	prod-agent	15.00	18.66	52.49	2025-11-28 13:41:24.199657+00
2a90daef-848f-41d0-9a92-91f751228101	prod-agent	7.13	18.67	52.49	2025-11-28 13:41:34.200519+00
cccba1e1-0f7b-45aa-94cb-2b992f9e4d3d	prod-agent	8.69	18.59	52.49	2025-11-28 13:41:44.243972+00
50fc6423-7d62-42e4-a5a5-0abe8be94bc9	prod-agent	9.04	18.60	52.49	2025-11-28 13:41:54.207565+00
545760e4-1568-421c-9332-cc0a48770ada	prod-agent	7.30	18.61	52.49	2025-11-28 13:42:04.203218+00
b7018ae1-c308-4c5c-ba80-b78a8f78874b	prod-agent	7.82	18.55	52.49	2025-11-28 13:42:14.194522+00
ce0ae010-3817-49ca-944b-f1889340d8c9	prod-agent	15.14	18.64	52.49	2025-11-28 13:42:24.199847+00
fc10a56c-167b-4873-9986-61b7a6bc1994	prod-agent	6.92	18.58	52.49	2025-11-28 13:42:34.197345+00
91a200e4-42f6-4f52-9e17-dae9808dbca4	prod-agent	8.61	18.67	52.49	2025-11-28 13:42:44.201034+00
10f350c0-1464-439a-9fce-11c919941345	prod-agent	9.01	18.65	52.49	2025-11-28 13:42:54.199592+00
74a50c19-5042-4deb-82a8-5068ee404a2e	prod-agent	7.20	18.58	52.49	2025-11-28 13:43:04.194598+00
522c0292-8b3d-4b8e-9751-ba92c4c57f9c	prod-agent	8.01	18.56	52.49	2025-11-28 13:43:14.206792+00
8e120c3f-ebbe-496d-9a5f-6ec1755a161b	prod-agent	13.22	18.73	52.49	2025-11-28 13:43:24.197436+00
3db48af6-c159-4b57-9a83-7ea414a14cd0	prod-agent	8.94	18.60	52.49	2025-11-28 13:43:34.204061+00
dfac4c07-0535-4214-844d-53c558680ccd	prod-agent	8.71	18.67	52.49	2025-11-28 13:43:44.195518+00
3de8aaee-4df9-455e-9a05-b2cc4635015a	prod-agent	9.82	18.65	52.49	2025-11-28 13:43:54.199918+00
b88b2fc3-bdf1-4bec-ac3d-6ef664c797b5	prod-agent	7.48	18.59	52.49	2025-11-28 13:44:04.200333+00
bb66ce4f-eb2d-4336-8d01-1d2fc7e8ca5f	prod-agent	8.23	18.62	52.49	2025-11-28 13:44:14.194469+00
9c9b3d21-a3d2-44b9-a4ab-fe36d9819087	prod-agent	9.65	18.55	52.49	2025-11-28 13:44:24.199977+00
f0b81051-90f6-4a59-9f7e-0128cc397cdb	prod-agent	12.36	18.65	52.49	2025-11-28 13:44:34.196044+00
00e05405-2f0b-466b-bc58-6e516d98f299	prod-agent	9.05	18.65	52.49	2025-11-28 13:44:44.198899+00
6d106b7f-0950-4bf9-81d8-df5075db16be	prod-agent	9.38	18.58	52.49	2025-11-28 13:44:54.200713+00
b573ad78-ba19-4acb-bb81-078955bcfe95	prod-agent	7.85	18.65	52.49	2025-11-28 13:45:04.194639+00
f8e0b992-346f-4086-ba64-5575d883da5e	prod-agent	8.89	18.66	52.49	2025-11-28 13:45:14.200037+00
c48d8e3c-a488-4c7d-b511-9a57243a90c9	prod-agent	8.80	18.63	52.49	2025-11-28 13:45:24.195471+00
6596383c-a16a-4f10-816d-854bd221b157	prod-agent	12.77	18.66	52.49	2025-11-28 13:45:34.193214+00
2d5bd9a1-44a4-466a-9312-ded6ce9ad8fc	prod-agent	8.08	18.65	52.49	2025-11-28 13:45:44.195084+00
363eff4b-49e6-482e-9499-7fe2d7cb97e6	prod-agent	9.91	18.61	52.49	2025-11-28 13:45:54.194676+00
0e05dbb0-4d5f-4a9a-b450-bee5d350c792	prod-agent	6.88	18.62	52.49	2025-11-28 13:46:04.208868+00
19d875a5-972a-44a7-84b2-4e9604508575	prod-agent	13.17	19.03	52.49	2025-11-28 13:46:14.200893+00
256793bd-49b6-4a9d-833e-fabe5aa86cfe	prod-agent	18.09	18.65	52.49	2025-11-28 13:46:24.193399+00
8c069b14-805d-4f33-a386-bc0484da478f	prod-agent	12.43	18.61	52.49	2025-11-28 13:46:34.199965+00
28ed42fe-6349-4f75-9784-ffd7e2bd8c9d	prod-agent	8.21	18.67	52.49	2025-11-28 13:46:44.19445+00
8eb36633-65a1-461b-af83-0dd302585d5a	prod-agent	9.66	18.61	52.49	2025-11-28 13:46:54.251867+00
010a4591-8d7d-4669-b623-026a082135ed	prod-agent	8.10	18.62	52.49	2025-11-28 13:47:04.199027+00
b47139cc-d3b8-4a62-bb68-537fac05f7c5	prod-agent	8.28	18.64	52.49	2025-11-28 13:47:14.199621+00
b7a516a2-9725-4289-bea7-97d802f3ff73	prod-agent	10.29	18.62	52.49	2025-11-28 13:47:24.20054+00
10cf9479-d008-4f44-89ae-4f8152d27bc5	prod-agent	13.52	18.73	52.49	2025-11-28 13:47:34.194255+00
4582a26a-958c-4f4e-9694-2aa0e0e6931a	prod-agent	8.85	18.53	52.49	2025-11-28 13:47:44.198012+00
497ba698-9aa8-4fc4-a8b9-845bcd126312	prod-agent	9.16	18.46	52.49	2025-11-28 13:47:54.20545+00
2af143f0-48e1-4879-abb4-3397344a7ea8	prod-agent	6.84	18.52	52.49	2025-11-28 13:48:04.203099+00
8446261a-35a0-424b-bd28-80e5641fb926	prod-agent	9.07	18.52	52.49	2025-11-28 13:48:14.204003+00
7811c55c-1529-4334-89c4-43aa22b44c9a	prod-agent	9.92	18.51	52.49	2025-11-28 13:48:24.197759+00
3a8e8448-3d1d-4c4e-91ed-dccdc48c15ff	prod-agent	11.04	18.51	52.49	2025-11-28 13:48:34.197125+00
af141cdd-bfa2-45e8-8c05-c7890e299278	prod-agent	11.42	18.47	52.49	2025-11-28 13:48:44.199276+00
be5d59c0-5d45-4424-a520-28c83e007e10	prod-agent	9.36	18.52	52.49	2025-11-28 13:48:54.195222+00
7b15eda6-871c-46b5-8cf5-d1dddb8463dd	prod-agent	7.06	18.50	52.49	2025-11-28 13:49:04.194937+00
fd94dfd8-4c02-4a06-9dce-34b6f90b1b1c	prod-agent	8.81	18.43	52.49	2025-11-28 13:49:14.193983+00
4210c6f9-e1f3-476e-9fce-03c98d796345	prod-agent	10.27	18.48	52.49	2025-11-28 13:49:24.199595+00
811bc627-a8ae-439f-a087-bdda6947ff68	prod-agent	6.96	18.48	52.49	2025-11-28 13:49:34.194465+00
6a97dd5a-8232-45fc-b196-d07b92234012	prod-agent	13.62	18.55	52.49	2025-11-28 13:49:44.196087+00
e7405919-75e6-48bd-9575-978c42ebeeb4	prod-agent	8.30	18.63	52.49	2025-11-28 13:49:54.197638+00
2a3ecf3d-96b7-49c8-bde9-9fb016249d3d	prod-agent	5.75	18.53	52.49	2025-11-28 13:50:04.202983+00
789097ea-f24f-4a6f-a322-b6e529fa5728	prod-agent	10.33	18.48	52.49	2025-11-28 13:50:14.195395+00
91e916f5-8cc0-4104-8944-888ed4b8521f	prod-agent	10.03	18.48	52.49	2025-11-28 13:50:24.199885+00
b6a514c2-1fb9-4dd2-8bd5-b7240a0f4743	prod-agent	7.25	18.51	52.49	2025-11-28 13:50:34.194234+00
c586c594-40f2-4e19-a465-ad10c357c9f3	prod-agent	14.32	18.56	52.49	2025-11-28 13:50:44.198255+00
1d8ee356-8a23-4623-9dfb-39afaab82164	prod-agent	9.26	18.60	52.49	2025-11-28 13:50:54.200581+00
513eb2dc-ac49-434b-bc56-5a4cf5e34fa9	prod-agent	7.19	18.57	52.49	2025-11-28 13:51:04.200309+00
03d786d5-7006-4fd2-aab2-8e62e71c6ad8	prod-agent	9.36	18.55	52.49	2025-11-28 13:51:14.193637+00
eee8421f-dab8-4e75-9472-4aee14d325a0	prod-agent	9.65	18.52	52.49	2025-11-28 13:51:24.200179+00
e61d5d0c-6fcd-4a5b-ab97-6d0f237772d6	prod-agent	7.22	18.59	52.49	2025-11-28 13:51:34.194523+00
6a51ce53-609c-47ba-b396-71f989588856	prod-agent	13.32	18.56	52.49	2025-11-28 13:51:44.201375+00
23d1f7dd-804a-49d2-a6e6-1e6592701265	prod-agent	7.90	18.48	52.49	2025-11-28 13:51:54.201162+00
03dab958-ef6e-4ed9-ac0a-3cce54832a76	prod-agent	6.85	18.44	52.49	2025-11-28 13:52:04.252526+00
0dfe49a1-15cd-4f1c-8d8f-e9be38fd3eff	prod-agent	8.96	18.55	52.49	2025-11-28 13:52:14.194894+00
c44378eb-090f-4522-92dd-b96eaa7df536	prod-agent	10.22	18.51	52.49	2025-11-28 13:52:24.193978+00
fce2eb15-c994-460e-8c8d-a4fcf37f79c1	prod-agent	7.05	18.51	52.49	2025-11-28 13:52:34.19993+00
0fa8e2ae-50e1-45ac-8a0e-0400bf664d1b	prod-agent	14.04	18.54	52.49	2025-11-28 13:52:44.202011+00
cf5d76c3-e6db-4228-9c9e-ba31e15f1eab	prod-agent	9.03	18.50	52.49	2025-11-28 13:52:54.20054+00
ad650f6c-6542-4ff3-b83a-333ca8d6c1ee	prod-agent	7.57	18.50	52.49	2025-11-28 13:53:04.200006+00
b43bf30f-71d0-441a-952e-1323f5074c77	prod-agent	8.59	18.52	52.49	2025-11-28 13:53:14.194192+00
c6203a3d-5756-42ff-b59b-9f38cec69c2c	prod-agent	9.84	18.54	52.49	2025-11-28 13:53:24.205958+00
c53bf93e-c2e4-427c-a70c-daaefaf7b2d4	prod-agent	7.21	18.60	52.49	2025-11-28 13:53:34.194065+00
2e3352a6-f436-4e3c-a950-42e0dab9bda2	prod-agent	11.67	18.65	52.49	2025-11-28 13:53:44.201052+00
e2a59950-af3e-425b-bad4-059195363a30	prod-agent	9.63	18.53	52.49	2025-11-28 13:53:54.203492+00
a934da82-d5ed-420a-aa27-8926eeb9bbdc	prod-agent	7.40	18.54	52.49	2025-11-28 13:54:04.199309+00
48847fa9-8194-4c50-8576-c0f4c241f14c	prod-agent	8.60	18.63	52.49	2025-11-28 13:54:14.199997+00
cac20dbe-d5f9-4568-852a-ce730975ff55	prod-agent	9.96	18.59	52.49	2025-11-28 13:54:24.199711+00
fc7c63ed-88e0-4423-a199-fae5b7000965	prod-agent	7.49	18.61	52.49	2025-11-28 13:54:34.197514+00
016c6866-b962-4238-aa98-8116caa9ba3d	prod-agent	8.86	18.62	52.49	2025-11-28 13:54:44.200877+00
0eea60ac-5d9a-4a2d-8b2b-5db16bf98f66	prod-agent	14.49	18.55	52.49	2025-11-28 13:54:54.19461+00
dbf69b1c-24e6-47f0-8bce-0d6244a4d5a3	prod-agent	7.27	18.61	52.49	2025-11-28 13:55:04.198785+00
3e3ae2e5-648d-4976-85f4-91e4015e8bbd	prod-agent	8.53	18.56	52.49	2025-11-28 13:55:14.200333+00
d585fe2b-5424-4857-9c0d-5a4ff02e9030	prod-agent	9.31	18.66	52.49	2025-11-28 13:55:24.193756+00
ae252066-a976-4417-8064-e107ef9f7bf3	prod-agent	7.59	18.56	52.49	2025-11-28 13:55:34.200379+00
997292df-5cae-4c2a-9091-c0addc810ce0	prod-agent	8.84	18.61	52.49	2025-11-28 13:55:44.216086+00
80d1bc49-7f06-4e3d-84fe-730069bcb5d9	prod-agent	14.58	18.56	52.49	2025-11-28 13:55:54.199711+00
29b6e312-3194-42a1-b950-dc6fcb891eca	prod-agent	7.14	18.59	52.49	2025-11-28 13:56:04.200095+00
dda5dfdc-20d3-4ac2-bd48-1e4bf13ca00a	prod-agent	9.19	18.60	52.49	2025-11-28 13:56:14.201134+00
d7aca758-18a6-4cd4-bd1e-bbca36952345	prod-agent	8.58	18.62	52.49	2025-11-28 13:56:24.206307+00
42647480-275e-48d8-8882-58bff7fdba84	prod-agent	8.25	18.61	52.49	2025-11-28 13:56:34.194536+00
08c15c3e-396d-4bb5-9987-68aca8fe539c	prod-agent	8.22	18.67	52.49	2025-11-28 13:56:44.199411+00
a9d81c9b-e6a0-4a98-a40e-46a0ab9415b8	prod-agent	14.97	18.65	52.49	2025-11-28 13:56:54.194963+00
7e85efcd-cf85-4b94-9d81-9e26ff5561c1	prod-agent	7.28	18.56	52.49	2025-11-28 13:57:04.200479+00
20169124-05a1-49a0-a082-121b02223cbf	prod-agent	8.47	18.56	52.49	2025-11-28 13:57:14.251399+00
b24a412a-fd6a-4ebe-8e72-76645c9a080d	prod-agent	9.53	18.63	52.49	2025-11-28 13:57:24.200205+00
ca35a8f3-04a2-46d3-a0d8-fe5ac21f3e3d	prod-agent	8.05	18.61	52.49	2025-11-28 13:57:34.198778+00
01facf38-a053-4e72-bd25-fa815a8bbd8c	prod-agent	9.23	18.60	52.49	2025-11-28 13:57:44.203524+00
feda0481-3d1c-473c-8169-7f90847337df	prod-agent	13.85	18.65	52.49	2025-11-28 13:57:54.196082+00
475030a9-56bb-4c2d-b0bb-e9d9e38a091a	prod-agent	7.82	18.65	52.49	2025-11-28 13:58:04.199914+00
c3bcbe0c-d3fb-4d33-bcd9-ca886e78bfc6	prod-agent	8.68	18.61	52.49	2025-11-28 13:58:14.199999+00
82b80bc8-4b31-4473-86a6-3f76c3d3d808	prod-agent	9.38	18.61	52.49	2025-11-28 13:58:24.201031+00
ba944fc1-263a-4995-93cb-d59e7dff1a84	prod-agent	7.72	18.60	52.49	2025-11-28 13:58:34.194239+00
0caf7171-a794-4602-9dbb-d2b2694cc74a	prod-agent	8.05	18.64	52.49	2025-11-28 13:58:44.19911+00
9b8df0c6-d5d7-4209-a3c4-819cf4b771bd	prod-agent	12.44	18.73	52.49	2025-11-28 13:58:54.19598+00
e6fd909d-c990-4936-aaaa-827592d4e524	prod-agent	9.39	18.60	52.49	2025-11-28 13:59:04.194269+00
756e9c50-d5b0-4469-994b-375adb76604b	prod-agent	9.16	18.61	52.49	2025-11-28 13:59:14.202287+00
a3818cc5-de61-489e-9bc9-8e70d22bce8c	prod-agent	8.80	18.66	52.49	2025-11-28 13:59:24.193991+00
54fb5a5f-662d-459d-b05e-600a3968b0e9	prod-agent	8.21	18.70	52.49	2025-11-28 13:59:34.200384+00
aecbb140-aa15-4a35-8c40-1d2bf9fc0f73	prod-agent	7.53	18.63	52.49	2025-11-28 13:59:44.210613+00
0525ca01-a38b-48e2-950e-480086142f9c	prod-agent	9.04	18.68	52.49	2025-11-28 13:59:54.207382+00
3d2b461f-3bae-43a1-bb32-13a99336ff2d	prod-agent	13.39	18.70	52.49	2025-11-28 14:00:04.203243+00
b1a5c44a-f51e-49be-b576-e6f53c1ae4ce	prod-agent	6.83	18.60	52.49	2025-11-28 14:00:14.199933+00
e8e81f33-1fbc-418f-befe-106ff15f3bea	prod-agent	9.71	18.68	52.49	2025-11-28 14:00:24.193633+00
46558200-52d6-47d1-a585-2f2b2322e181	prod-agent	7.97	18.60	52.49	2025-11-28 14:00:34.198939+00
3093e1c5-ee4e-41f5-80c9-d1f969cf6f7a	prod-agent	8.60	18.69	52.49	2025-11-28 14:00:44.194155+00
4dab98e7-6ed5-4c34-80a6-0f754b3939f8	prod-agent	8.49	18.78	52.49	2025-11-28 14:00:54.202961+00
ae2890ce-2983-4e5b-8e13-9944a0ee7254	prod-agent	13.94	18.66	52.49	2025-11-28 14:01:04.199494+00
27e1e5d6-8cd5-46fb-9462-ca9c7ea834e4	prod-agent	7.59	18.60	52.49	2025-11-28 14:01:14.19456+00
1eb5bac0-01ef-4a6b-a6d8-b9eff71802e8	prod-agent	24.48	18.79	52.49	2025-11-28 14:01:24.197524+00
c550914f-9b20-416a-9201-e0a0b2621ad9	prod-agent	6.32	18.66	52.49	2025-11-28 14:01:34.194477+00
0b47bd2e-198c-4be7-8bb4-47a414c2afe9	prod-agent	8.14	18.65	52.49	2025-11-28 14:01:44.202058+00
7fd4bd89-72c1-453e-874e-eaf7c6e87e83	prod-agent	8.49	18.64	52.49	2025-11-28 14:01:54.213399+00
96f42cb1-3a1a-4255-a575-65d0bb78d4f0	prod-agent	14.35	18.77	52.49	2025-11-28 14:02:04.194023+00
8a8a3fc5-7b6c-4158-bb05-85407e0cc3ce	prod-agent	7.79	18.64	52.49	2025-11-28 14:02:14.199478+00
2c71982e-35b6-4a3e-8823-e22649633c3c	prod-agent	10.43	18.66	52.49	2025-11-28 14:02:24.201444+00
2fae31f6-d259-4baa-80e6-f3144f3afcd1	prod-agent	7.61	18.62	52.49	2025-11-28 14:02:34.194486+00
1b4e925e-e743-4a62-b5ee-ab40ae20861a	prod-agent	8.15	18.63	52.49	2025-11-28 14:02:44.194608+00
5eaf43e2-6ad5-4f65-9563-8bf75dd3b137	prod-agent	9.60	18.67	52.49	2025-11-28 14:02:54.204867+00
b9712339-0290-4c8e-887b-4c4c471c7338	prod-agent	13.69	18.61	52.49	2025-11-28 14:03:04.19412+00
91cf568b-b2b2-4931-b0a8-c4f0ec69f83e	prod-agent	6.56	18.63	52.49	2025-11-28 14:03:14.200636+00
dc22190c-fe0e-48fd-a0fc-33636e100946	prod-agent	9.31	18.61	52.49	2025-11-28 14:03:24.19956+00
5027534c-c999-4256-8bb3-1ad64b94444f	prod-agent	8.19	18.62	52.49	2025-11-28 14:03:34.204205+00
d606bdd6-e8a6-466c-adcc-bcc244582419	prod-agent	7.53	18.61	52.49	2025-11-28 14:03:44.201294+00
49810731-8ea5-45aa-b9d9-ff5bf5ba3230	prod-agent	9.05	18.62	52.49	2025-11-28 14:03:54.201162+00
758e6973-f874-462f-b858-ce4a9d5ce4a6	prod-agent	18.89	19.04	52.49	2025-11-28 14:04:04.20229+00
34ede2e7-a059-4f7f-9dd9-88e9a4b0c040	prod-agent	22.67	18.68	52.49	2025-11-28 14:04:14.200006+00
c2bca7f5-1eef-4820-9b55-1e9437dcbc5c	prod-agent	9.09	18.64	52.49	2025-11-28 14:04:24.200259+00
734648ec-9a6e-4802-baf3-915cecc57bfb	prod-agent	7.87	18.69	52.49	2025-11-28 14:04:34.200435+00
e0b74aff-3ca8-4662-a4f5-f39e1dcecbbc	prod-agent	8.14	18.68	52.49	2025-11-28 14:04:44.195061+00
ab907812-5346-48f0-83bc-a6a4ffad2071	prod-agent	9.80	18.72	52.49	2025-11-28 14:04:54.202181+00
af5df4dc-3b2b-485a-bc53-453bdf7d20de	prod-agent	7.96	18.66	52.49	2025-11-28 14:05:04.193361+00
656ce34b-ae23-426b-b146-9f4dd59a8ffc	prod-agent	14.14	18.70	52.49	2025-11-28 14:05:14.199343+00
b948d4ee-5500-4145-974b-42931c78c274	prod-agent	10.46	18.71	52.49	2025-11-28 14:05:24.198516+00
b5a21bf9-582b-41dc-afea-887be5e8a3d5	prod-agent	6.71	18.62	52.49	2025-11-28 14:05:34.194092+00
7a2477c0-305f-4725-98e3-a866b9ed1b23	prod-agent	7.12	18.67	52.49	2025-11-28 14:05:44.205112+00
2018d84f-1517-4217-b571-5b784d857a16	prod-agent	9.31	18.67	52.49	2025-11-28 14:05:54.198621+00
e5318e52-32e9-459f-852f-c40c3031e3b8	prod-agent	8.13	18.67	52.49	2025-11-28 14:06:04.201429+00
797c1342-608e-41f0-84e1-18a526c30c1b	prod-agent	13.54	18.49	52.49	2025-11-28 14:06:14.199689+00
e52f4f79-e905-47ed-a031-c83152faf820	prod-agent	10.04	18.53	52.49	2025-11-28 14:06:24.200168+00
82ffe7f2-b62a-4f3a-b5c7-68c2dce7089e	prod-agent	8.59	18.50	52.49	2025-11-28 14:06:34.210947+00
bf4c8e27-1612-4a84-8bb4-76ad978fdea7	prod-agent	9.61	18.47	52.49	2025-11-28 14:06:44.199378+00
9d04b8ad-d663-43f4-9066-d3cbd5a80e98	prod-agent	7.54	18.44	52.49	2025-11-28 14:06:54.198854+00
e2f20317-5b81-4a77-8093-90ed2cf5a055	prod-agent	9.40	18.55	52.49	2025-11-28 14:07:04.194456+00
6d996c53-e32a-4c12-aef8-ef656d9a218e	prod-agent	14.04	18.63	52.49	2025-11-28 14:07:14.201821+00
17969b61-2d77-419b-be31-fff4e089b12f	prod-agent	8.45	18.58	52.49	2025-11-28 14:07:24.252091+00
1721788e-13f4-4bd9-a027-5a8858d10039	prod-agent	9.37	18.59	52.49	2025-11-28 14:07:34.195493+00
aaa52037-6253-4467-9ed0-f75b0bc75bc1	prod-agent	6.05	18.57	52.49	2025-11-28 14:07:44.199442+00
243c654c-9862-4afb-9436-85e8e4bf6b91	prod-agent	8.78	18.63	52.49	2025-11-28 14:07:54.195587+00
4a057b90-cff4-47f0-9aff-718b028362a1	prod-agent	9.46	18.65	52.49	2025-11-28 14:08:04.21128+00
874f20c3-8347-4d92-b61b-e6ca514b2cd1	prod-agent	11.52	18.58	52.49	2025-11-28 14:08:14.206395+00
7dc5a79f-41b3-4e5d-b5fa-84caf5c51f3d	prod-agent	10.56	18.65	52.49	2025-11-28 14:08:24.202796+00
cd60f835-de0a-4f0a-ae15-ea5cf2b48214	prod-agent	9.48	18.63	52.49	2025-11-28 14:08:34.205132+00
4dd6a2ce-97c1-4ec9-869f-ed09a51dd2fa	prod-agent	7.97	18.49	52.49	2025-11-28 14:08:44.195365+00
a5d02c9b-e78b-457d-a94a-088cab6bc5a6	prod-agent	8.77	18.48	52.49	2025-11-28 14:08:54.200064+00
accef17e-2abf-4c77-aa6e-888b08c578bf	prod-agent	9.40	18.51	52.49	2025-11-28 14:09:04.196007+00
e2348cd3-3820-4fa5-a8ab-3b2016d333e0	prod-agent	7.60	18.45	52.49	2025-11-28 14:09:14.200091+00
e039a338-8127-43f3-a488-13e6de56d27e	prod-agent	13.76	18.51	52.49	2025-11-28 14:09:24.199626+00
59c5e5af-80c2-44d3-9734-dd95f93d2193	prod-agent	8.20	18.48	52.49	2025-11-28 14:09:34.196686+00
33b91d70-3596-47f5-aa13-b6e938271ed5	prod-agent	8.63	18.42	52.49	2025-11-28 14:09:44.194473+00
9f4af7e8-55a8-4537-8630-98a08cc48b2f	prod-agent	8.17	18.43	52.49	2025-11-28 14:09:54.198433+00
294ad666-337d-49a2-8b2a-0df33d026f15	prod-agent	9.07	18.51	52.49	2025-11-28 14:10:04.201152+00
370a7ab7-78c6-4711-a819-41296443149a	prod-agent	8.32	18.44	52.49	2025-11-28 14:10:14.207175+00
0d5c7a6c-ced7-4004-b4da-e3cce9311e06	prod-agent	14.67	18.52	52.49	2025-11-28 14:10:24.199461+00
572ecd54-790f-431a-8890-f887618fbd82	prod-agent	8.69	18.54	52.49	2025-11-28 14:10:34.194619+00
e221c143-65a3-4998-adaf-15ecb93f0a1e	prod-agent	8.31	18.46	52.49	2025-11-28 14:10:44.203055+00
d0c89e09-241a-4cac-9d92-bee2d57502ca	prod-agent	8.49	18.46	52.49	2025-11-28 14:10:54.19455+00
a136d9f5-1deb-44e0-8dd7-2f816e3f353f	prod-agent	9.29	18.47	52.49	2025-11-28 14:11:04.207474+00
d400af6f-bfb9-4f1a-b65c-22a30f38fac3	prod-agent	7.58	18.43	52.49	2025-11-28 14:11:14.200133+00
60f99125-6735-41d7-885a-5ae9cea9d1d1	prod-agent	15.33	18.50	52.49	2025-11-28 14:11:24.197498+00
20817a23-a964-400a-a92d-52752de03f3f	prod-agent	7.33	18.47	52.49	2025-11-28 14:11:34.217823+00
c4119cad-2380-4d75-a79e-4754751f2671	prod-agent	7.39	18.46	52.49	2025-11-28 14:11:44.199258+00
e65e6634-a692-42db-afb7-9dad51d453f2	prod-agent	9.10	18.44	52.49	2025-11-28 14:11:54.196893+00
e38776c9-a600-4614-9612-d40185a8e8f0	prod-agent	8.40	18.44	52.49	2025-11-28 14:12:04.203282+00
f9e10789-8d2b-489d-a20d-9fcc7b6c1bc3	prod-agent	8.25	18.55	52.49	2025-11-28 14:12:14.200441+00
34e518ca-1f88-4b3a-9c7a-3b81abba313b	prod-agent	15.04	18.50	52.49	2025-11-28 14:12:24.193581+00
b5d13fc1-12f3-4c0c-b69b-545620c55385	prod-agent	8.64	18.47	52.49	2025-11-28 14:12:34.255809+00
7d040b57-742e-4065-9da9-5f3f55a16ffd	prod-agent	7.57	18.51	52.49	2025-11-28 14:12:44.195112+00
782b8273-3a5d-4201-a0fe-0ad3d2966ba6	prod-agent	8.76	18.56	52.49	2025-11-28 14:12:54.202178+00
2a3a825d-028c-40d0-ae7f-a580f3ed79b4	prod-agent	9.34	18.52	52.49	2025-11-28 14:13:04.200002+00
fd54a47a-07e3-4432-9ab5-9576e0094cbd	prod-agent	7.80	18.55	52.49	2025-11-28 14:13:14.200071+00
89f76c81-2709-489b-adce-7fd522c77076	prod-agent	11.96	18.56	52.49	2025-11-28 14:13:24.222438+00
0c9b9abc-13b0-4881-a767-725f243d2479	prod-agent	10.58	18.57	52.49	2025-11-28 14:13:34.202283+00
23b2ee5f-7bd8-479d-970e-09befd57b06b	prod-agent	8.68	18.62	52.51	2025-11-28 14:13:44.195029+00
171a1e34-abfd-4abc-82f6-c86798f32002	prod-agent	8.76	18.61	52.51	2025-11-28 14:13:54.20221+00
a41abc0c-3a5d-4f5b-9612-3a4099113c70	prod-agent	8.64	18.62	52.51	2025-11-28 14:14:04.195317+00
4028f3e5-434f-45c3-96a5-6df231ea3b96	prod-agent	8.05	18.62	52.51	2025-11-28 14:14:14.204528+00
dcc6bede-dfda-4769-b692-33bd2ad22a4d	prod-agent	8.61	18.58	52.51	2025-11-28 14:14:24.204302+00
5c737ddd-6e1c-4fb8-99a5-6a122fa483a2	prod-agent	14.36	18.54	52.51	2025-11-28 14:14:34.198972+00
16f9eb7e-f1d5-4a83-be44-20dd1bbffdb9	prod-agent	7.75	18.55	52.51	2025-11-28 14:14:44.201005+00
4e1be1ed-b9e2-44b4-95b9-cca4e4ae0a35	prod-agent	8.97	18.56	52.51	2025-11-28 14:14:54.194294+00
f2ab8028-9d1a-42f5-9b8d-9de49fa7e393	prod-agent	8.59	18.53	52.51	2025-11-28 14:15:04.2031+00
096573fd-70a1-4c7b-a581-f1902a2ae603	prod-agent	9.19	18.56	52.51	2025-11-28 14:15:14.2007+00
e5eaa72c-c535-47bc-a820-4ecca8ae472e	prod-agent	9.33	18.64	52.51	2025-11-28 14:15:24.194799+00
5872a8fb-43a6-455a-972e-e5df34e1e47c	prod-agent	13.26	18.61	52.51	2025-11-28 14:15:34.192998+00
277ca5d1-66e1-45e0-8c76-722b2c9a24e1	prod-agent	8.16	18.57	52.51	2025-11-28 14:15:44.199561+00
44870bb3-2588-448f-ac69-f90ae0347152	prod-agent	8.50	18.57	52.51	2025-11-28 14:15:54.193926+00
de5efd9c-71a3-4f38-b738-134334cfd6d9	prod-agent	8.34	18.55	52.51	2025-11-28 14:16:04.198167+00
fca948f2-aaf9-4f64-95c9-88f4a0cd3726	prod-agent	8.88	18.58	52.51	2025-11-28 14:16:14.193939+00
856039d6-8684-48c7-86b3-e378de5033cb	prod-agent	26.12	18.77	52.51	2025-11-28 14:16:24.191143+00
5517e265-d0a4-4088-936c-21c16dd301df	prod-agent	10.82	18.63	52.51	2025-11-28 14:16:34.199477+00
450ac5c4-8d7d-4765-8f9f-0e55b797d2d3	prod-agent	8.24	18.57	52.51	2025-11-28 14:16:44.199833+00
94d1332f-74b2-4c96-bb0e-0c9e1783adbd	prod-agent	8.41	18.60	52.51	2025-11-28 14:16:54.200025+00
e5b4523d-ff29-4aa0-ad0f-e7f89c91c4c5	prod-agent	8.58	18.57	52.51	2025-11-28 14:17:04.200366+00
a114b338-79e3-42be-a130-888cd1366900	prod-agent	9.15	18.67	52.51	2025-11-28 14:17:14.202892+00
da1d0e4f-1e54-425f-978b-a8899c35881f	prod-agent	9.51	18.69	52.51	2025-11-28 14:17:24.193684+00
e8dcd1c5-18e0-4d78-8da7-747909f8ce2d	prod-agent	11.33	18.68	52.51	2025-11-28 14:17:34.201916+00
e188655b-5be3-4d9e-8ee2-2da6b21f4ce5	prod-agent	7.59	18.60	52.51	2025-11-28 14:17:44.244204+00
464a0e57-9b22-434f-9555-4d88b5c69498	prod-agent	8.08	18.60	52.51	2025-11-28 14:17:54.198954+00
aa1a4a82-5172-411a-80aa-b433dd0911d1	prod-agent	8.16	18.64	52.51	2025-11-28 14:18:04.199057+00
a2ac886d-aa17-45aa-8940-ca768f2b63ac	prod-agent	8.58	18.65	52.51	2025-11-28 14:18:14.200578+00
672c1667-f611-443f-a33b-74b20afe9961	prod-agent	9.63	18.62	52.51	2025-11-28 14:18:24.193841+00
ab3a4a22-1803-4103-82bb-4601662d1602	prod-agent	9.68	18.63	52.51	2025-11-28 14:18:34.201273+00
b5ee0f29-ddce-41c8-bedd-72f36070bcd1	prod-agent	12.74	18.69	52.51	2025-11-28 14:18:44.200536+00
cc3b56f2-8650-4863-9218-f22aca4065d6	prod-agent	8.06	18.64	52.51	2025-11-28 14:18:54.194319+00
100af92f-3aad-49b5-ab95-c1bdc9b6c05d	prod-agent	8.57	18.63	52.51	2025-11-28 14:19:04.199898+00
6d8c516a-a4a5-43da-8a45-7ee0853596f8	prod-agent	8.43	18.59	52.51	2025-11-28 14:19:14.203044+00
7e2be22e-f707-4411-80a7-2c4da474bdc8	prod-agent	9.16	18.65	52.51	2025-11-28 14:19:24.194545+00
a6293f8a-159d-4899-92d8-09657509ed9b	prod-agent	7.47	18.62	52.51	2025-11-28 14:19:34.201094+00
95a5b76c-2cb4-4658-8ffd-8ce07b20a51e	prod-agent	14.58	18.61	52.51	2025-11-28 14:19:44.192884+00
4d1bab72-3b76-4637-92cb-3cc07f1100d5	prod-agent	7.43	18.67	52.51	2025-11-28 14:19:54.198847+00
297b1d88-51a3-4e46-b482-119a30f9244e	prod-agent	8.99	18.58	52.51	2025-11-28 14:20:04.205545+00
00eeb741-ab31-4dac-9761-5e43616ff5f6	prod-agent	9.34	18.65	52.51	2025-11-28 14:20:14.200456+00
9a6b97db-0e36-43b9-ab5d-7946fa4954e6	prod-agent	9.57	18.61	52.51	2025-11-28 14:20:24.196105+00
d4d9faac-e9b1-43b2-922b-73ba955ba475	prod-agent	8.51	18.62	52.51	2025-11-28 14:20:34.202182+00
785a7850-8619-4190-ab11-7cac0b3d9858	prod-agent	13.92	18.65	52.51	2025-11-28 14:20:44.199267+00
2872e353-ec2b-40c3-b101-0c61312ecddc	prod-agent	8.21	18.60	52.51	2025-11-28 14:20:54.194765+00
48e74014-1274-47d5-9c3d-284c404b5d65	prod-agent	7.13	18.56	52.51	2025-11-28 14:21:04.202712+00
e9d30778-1c15-4fe8-96ae-3c6a4cad0b59	prod-agent	8.30	18.65	52.51	2025-11-28 14:21:14.198291+00
a5a2d1e0-0f8f-43b6-baf1-232d28960a66	prod-agent	9.21	18.61	52.51	2025-11-28 14:21:24.208281+00
c87d2020-6bc0-49a9-985f-38c0b439be87	prod-agent	6.78	18.58	52.51	2025-11-28 14:21:34.195003+00
399568c2-aea2-477d-81c4-eb436aaee11f	prod-agent	15.55	18.79	52.51	2025-11-28 14:21:44.197704+00
e0a75532-fbab-4c47-a6bb-6812494d51ff	prod-agent	7.35	18.62	52.51	2025-11-28 14:21:54.195698+00
b069ec9f-367a-4157-b62f-8dd4f0cec382	prod-agent	8.14	18.59	52.51	2025-11-28 14:22:04.204318+00
5fe1c1d0-1325-4444-a3ec-28ef903ef227	prod-agent	8.76	18.64	52.51	2025-11-28 14:22:14.200322+00
1935f3d8-5ba0-4589-9886-2569f23a9fc9	prod-agent	9.67	18.64	52.51	2025-11-28 14:22:24.209592+00
7fc7b254-f9b5-4414-bbb7-25cf19fd2f60	prod-agent	8.59	18.68	52.51	2025-11-28 14:22:34.199802+00
784f89ae-f05a-4c5c-9eb1-291c564e75b8	prod-agent	13.67	18.70	52.51	2025-11-28 14:22:44.201198+00
ffb7c88b-b011-4b9b-a69a-a80f815dbc4e	prod-agent	8.60	18.56	52.51	2025-11-28 14:22:54.249032+00
6ebe4774-24db-405e-a3b4-d45558ff0157	prod-agent	8.46	18.60	52.51	2025-11-28 14:23:04.20059+00
c05ac69a-8fc0-48d0-8925-353a9dccb17b	prod-agent	7.53	18.56	52.51	2025-11-28 14:23:14.200256+00
126863e8-9f5e-4393-b16d-575de064edd1	prod-agent	8.08	18.65	52.51	2025-11-28 14:23:24.224296+00
609b2911-f407-4d73-a1d4-097d0f726c2e	prod-agent	7.95	18.60	52.51	2025-11-28 14:23:34.194979+00
388beeae-11b9-4f40-a589-6a4fa954e2b1	prod-agent	8.78	18.56	52.51	2025-11-28 14:23:44.202872+00
85c3c213-ed07-48d0-ae9e-1b0a145e3367	prod-agent	13.35	18.64	52.51	2025-11-28 14:23:54.194851+00
01f355d0-2d84-483f-970d-b0e054e98d8a	prod-agent	8.77	18.61	52.51	2025-11-28 14:24:04.200218+00
ff0d8aa4-d037-4717-99e2-41650ad9d969	prod-agent	8.24	18.60	52.51	2025-11-28 14:24:14.202798+00
f3c103eb-7da1-4eb1-80e2-0bc5a4f4a0fa	prod-agent	7.97	18.59	52.51	2025-11-28 14:24:24.200043+00
f2aed219-d82b-49bf-8716-f915bf980fa7	prod-agent	9.31	18.60	52.51	2025-11-28 14:24:34.194219+00
f593d3af-d9cd-407d-959b-f279150c5565	prod-agent	8.67	18.57	52.51	2025-11-28 14:24:44.202093+00
4a63c39b-4e69-40f2-b8cd-83fec97c4b95	prod-agent	13.79	18.60	52.51	2025-11-28 14:24:54.203556+00
4e11a7d2-5556-4f7d-9e76-2710aa6f35f8	prod-agent	9.09	18.59	52.51	2025-11-28 14:25:04.201023+00
4d142a0e-aec5-457b-8c21-875af755e2a6	prod-agent	8.95	18.58	52.51	2025-11-28 14:25:14.194913+00
5029a6a2-cfac-4c0e-ac1c-735cdfcbeb7c	prod-agent	6.90	18.55	52.51	2025-11-28 14:25:24.200134+00
bf96fc68-e3a2-4bd3-b0fc-4dd293c94799	prod-agent	9.30	18.55	52.51	2025-11-28 14:25:34.202962+00
a267f4b1-ab63-4c24-ab93-e9bdb06e9d07	prod-agent	8.01	18.59	52.51	2025-11-28 14:25:44.202913+00
31b196ca-87b1-465e-a5d5-8da23ac28c26	prod-agent	13.88	18.61	52.51	2025-11-28 14:25:54.196014+00
e575a72d-d141-4dd6-9580-ddf4cf5c9886	prod-agent	8.55	18.61	52.51	2025-11-28 14:26:04.198613+00
33ebfe80-dd08-438b-87c6-49b9d300fce1	prod-agent	8.93	18.50	52.51	2025-11-28 14:26:14.200483+00
62e51077-a50f-4913-a987-3763bed82dc9	prod-agent	8.77	18.51	52.51	2025-11-28 14:26:24.214022+00
f178588d-a110-40ee-a028-6d924a7f9eee	prod-agent	8.88	18.50	52.51	2025-11-28 14:26:34.194355+00
7af5ec31-ac3a-4a8e-95fd-4dd3c1b4927a	prod-agent	8.69	18.57	52.51	2025-11-28 14:26:44.202923+00
b2b8f1b4-8533-4b15-8731-722463440142	prod-agent	13.85	18.62	52.51	2025-11-28 14:26:54.194325+00
f0d21d4b-9dc0-426c-8267-340b0c649e8a	prod-agent	8.05	18.61	52.51	2025-11-28 14:27:04.203196+00
34bd6acc-5e21-4718-9471-20d8f2624f0f	prod-agent	8.31	18.58	52.51	2025-11-28 14:27:14.196791+00
cd705628-0027-4458-8bf2-2aee01753a6b	prod-agent	8.17	18.56	52.51	2025-11-28 14:27:24.19599+00
17fe39d5-b77f-40fe-967d-1b5871d0bf0d	prod-agent	7.56	18.58	52.51	2025-11-28 14:27:34.207762+00
b54d8460-6cd0-4125-892b-84640b736b36	prod-agent	9.24	18.44	52.51	2025-11-28 14:27:44.200127+00
fc528dc9-60bd-4aba-8648-6d7cb9c1f327	prod-agent	11.16	18.42	52.51	2025-11-28 14:27:54.200536+00
c9e979d5-e809-46c8-a614-3e48840f7826	prod-agent	10.46	18.41	52.51	2025-11-28 14:28:04.241012+00
e890435f-819b-40d5-b311-a182234734ec	prod-agent	8.68	18.48	52.51	2025-11-28 14:28:14.20068+00
ca9e78c2-b4d0-4496-ad00-7dccc878b2fd	prod-agent	7.90	18.48	52.51	2025-11-28 14:28:24.197269+00
16d41dc3-e45d-4507-abf1-0f311469fab0	prod-agent	10.04	18.54	52.51	2025-11-28 14:28:34.196602+00
958df912-4d53-45d9-b2c4-0221af3ae1e4	prod-agent	8.48	18.51	52.51	2025-11-28 14:28:44.194422+00
791b22c0-06c7-4d82-a46f-da4da8f7b664	prod-agent	7.99	18.59	52.51	2025-11-28 14:28:54.202684+00
f1be6ff2-d468-4eb7-ad15-303be2a8ede5	prod-agent	14.43	18.55	52.51	2025-11-28 14:29:04.196632+00
219fb2a0-e886-4947-bb9f-f27de0e27121	prod-agent	7.34	18.56	52.51	2025-11-28 14:29:14.194618+00
e05c8cbb-93a6-49df-9118-6bea07dfbff5	prod-agent	7.98	18.58	52.51	2025-11-28 14:29:24.201689+00
4bd3c1f3-ab8a-4106-97f9-1ad2c19fc370	prod-agent	8.96	18.50	52.51	2025-11-28 14:29:34.194869+00
585a59d1-7cb6-4d65-aca1-b95bdb74981a	prod-agent	7.97	18.56	52.51	2025-11-28 14:29:44.197414+00
67b284ed-a033-412b-8a42-54be2d89c80f	prod-agent	7.28	18.51	52.51	2025-11-28 14:29:54.200954+00
2f2e9873-9b9e-4738-9456-db29dc5a33f5	prod-agent	14.81	18.54	52.51	2025-11-28 14:30:04.202083+00
d97093ec-69ad-432b-b73e-ae39aec2145f	prod-agent	8.87	18.62	52.51	2025-11-28 14:30:14.194326+00
6e4c253e-5444-4490-ab37-848df3df21a9	prod-agent	8.55	18.60	52.51	2025-11-28 14:30:24.200609+00
ad49f368-1161-47f7-bad0-eb8f9c7f4cb3	prod-agent	9.46	18.62	52.51	2025-11-28 14:30:34.199935+00
58f24c18-65ba-45a2-a719-13d1f3fbe05f	prod-agent	8.94	18.57	52.51	2025-11-28 14:30:44.194463+00
b53ce17d-6316-40e0-813b-79cb1cf8e984	prod-agent	7.94	18.56	52.51	2025-11-28 14:30:54.194129+00
197713ac-5f28-4f81-93d4-edc995062aad	prod-agent	15.20	18.57	52.51	2025-11-28 14:31:04.194271+00
d5b8350f-b65b-4c22-8e49-3321181820c5	prod-agent	6.77	18.55	52.51	2025-11-28 14:31:14.195058+00
6d7b4721-fe19-4cda-9c3f-4bddeddf6a5f	prod-agent	26.07	19.01	52.51	2025-11-28 14:31:24.196089+00
a0cdd133-08ce-4fe3-81aa-bdfc88548a8d	prod-agent	6.91	18.53	52.51	2025-11-28 14:31:34.20033+00
02132314-0f40-40d7-88af-64adf39e7867	prod-agent	9.10	18.61	52.51	2025-11-28 14:31:44.20268+00
266b0697-cb5f-4853-95e4-99d4e1672127	prod-agent	7.78	18.56	52.51	2025-11-28 14:31:54.194508+00
54efdef6-0d76-43c9-b5c6-1326717baae7	prod-agent	14.24	18.70	52.51	2025-11-28 14:32:04.199318+00
667386bd-1b71-451b-adb2-c9d3dc06350b	prod-agent	8.08	18.61	52.51	2025-11-28 14:32:14.202533+00
ad21cb33-44d9-4a94-9e35-95c78a7a0b13	prod-agent	8.67	18.55	52.51	2025-11-28 14:32:24.194493+00
425979c0-87bb-4e82-9310-76ca973bdbe8	prod-agent	9.17	18.59	52.51	2025-11-28 14:32:34.20492+00
f25edeaa-7ea6-4de1-8b37-56d197e1716d	prod-agent	8.30	18.56	52.51	2025-11-28 14:32:44.202376+00
957bbc66-614e-4f12-9be4-8b19936c23d4	prod-agent	8.77	18.56	52.51	2025-11-28 14:32:54.194449+00
12130ff9-abc9-4ab1-8b2f-bb9613ea2f52	prod-agent	13.19	18.72	52.51	2025-11-28 14:33:04.194635+00
2acac19b-a904-48f5-8519-cf13b745c89b	prod-agent	8.41	18.41	52.51	2025-11-28 14:33:14.266464+00
11a843b5-787a-4394-b730-23e91670ac6f	prod-agent	9.34	18.41	52.51	2025-11-28 14:33:24.200633+00
560956fd-59c8-4214-ac5d-9d3b03dae5af	prod-agent	7.71	18.41	52.51	2025-11-28 14:33:34.205625+00
c5b26e69-9254-4669-9a7f-dd1e9e7f945b	prod-agent	8.52	18.39	52.51	2025-11-28 14:33:44.194685+00
942e1e0d-7af8-4d39-88e6-cddb6392b965	prod-agent	8.47	18.31	52.51	2025-11-28 14:33:54.20719+00
47270af2-4024-466c-b7ea-838a608dc426	prod-agent	15.91	18.76	52.51	2025-11-28 14:34:04.195366+00
5f924033-61e4-441e-93e1-5a0739b405ca	prod-agent	24.02	18.57	52.51	2025-11-28 14:34:14.195142+00
dc42a7d3-8234-4d25-bb9a-52c80ac03bb0	prod-agent	7.61	18.54	52.51	2025-11-28 14:34:24.195418+00
35e744db-cebd-46de-b151-5ac5f9a23cd7	prod-agent	8.77	18.42	52.51	2025-11-28 14:34:34.203067+00
aff34610-4162-4115-87d6-f9b0cd062b53	prod-agent	8.83	18.39	52.51	2025-11-28 14:34:44.204198+00
33715b71-0eab-468e-8e05-ebc29f8b54c6	prod-agent	8.91	18.37	52.51	2025-11-28 14:34:54.203809+00
b9630836-3f7f-401b-9ad1-cd54474ae03b	prod-agent	8.22	18.40	52.51	2025-11-28 14:35:04.200389+00
7f96a49c-5dbd-41c5-8f2f-ba075bde5ccd	prod-agent	13.21	18.54	52.52	2025-11-28 14:35:14.196551+00
45a530fe-9389-4461-8753-23f5dacda21d	prod-agent	8.52	18.58	52.52	2025-11-28 14:35:24.194025+00
c45f90fe-47d3-4286-9a04-ae9232966f3e	prod-agent	8.13	18.54	52.52	2025-11-28 14:35:34.203209+00
6a3b003d-03cb-4802-aa87-22eb43e8e138	prod-agent	8.24	18.54	52.52	2025-11-28 14:35:44.19412+00
7afa5a9f-f9a3-448b-a252-e1fa8e58f482	prod-agent	8.68	18.60	52.52	2025-11-28 14:35:54.19529+00
e4ce5281-e6d4-4654-9a92-0596b2265664	prod-agent	8.83	18.53	52.52	2025-11-28 14:36:04.199791+00
7488ba0e-8137-4350-8f63-de111533a6f8	prod-agent	14.44	18.62	52.52	2025-11-28 14:36:14.193913+00
8c05bcba-f34b-4920-8e1d-1c4c2b1778a3	prod-agent	8.74	18.54	52.52	2025-11-28 14:36:24.200494+00
1274e784-2fb0-4291-8e5c-3753dbbc0c51	prod-agent	8.16	18.58	52.52	2025-11-28 14:36:34.1957+00
adfe34ff-b9d9-4fbd-a4e2-8d1728c151f9	prod-agent	9.33	18.58	52.52	2025-11-28 14:36:44.195093+00
beb3d3f7-f466-4ad0-9626-7bb9440df888	prod-agent	8.35	18.58	52.52	2025-11-28 14:36:54.194792+00
e03eb94a-0023-4212-8f9f-13e06dee050c	prod-agent	8.67	18.57	52.52	2025-11-28 14:37:04.195196+00
41e698fe-edb6-454c-9963-66d9513cc9f3	prod-agent	14.31	18.60	52.52	2025-11-28 14:37:14.202535+00
e167255b-e9e0-4831-8b22-6a6525ac96ee	prod-agent	8.18	18.57	52.52	2025-11-28 14:37:24.19996+00
076971b9-4548-4777-a477-a6454a6f28d6	prod-agent	8.41	18.58	52.52	2025-11-28 14:37:34.193826+00
9c00b4b2-0ab4-4d70-a90c-2590d43bcbf4	prod-agent	8.82	18.57	52.52	2025-11-28 14:37:44.202295+00
917492f3-5dc1-4a8e-99fd-26091ba6b7c5	prod-agent	8.60	18.58	52.52	2025-11-28 14:37:54.194703+00
7f437953-0a5c-4c69-bfe6-332374d48d99	prod-agent	8.70	18.57	52.52	2025-11-28 14:38:04.200442+00
96bba947-3763-4f4c-a347-43953e3ceb89	prod-agent	12.71	18.68	52.52	2025-11-28 14:38:14.200228+00
2d039dbd-a0ba-40a4-9606-ecbcee8ce5e5	prod-agent	11.35	18.59	52.52	2025-11-28 14:38:24.204337+00
37807914-06cf-4000-9db5-7287e6a91e59	prod-agent	8.58	18.60	52.52	2025-11-28 14:38:34.194253+00
a8750c94-2d7a-4d0f-8413-914f8967f18b	prod-agent	8.05	18.62	52.52	2025-11-28 14:38:44.20621+00
a778a974-4403-48fa-8fa5-c59f62aa1557	prod-agent	8.95	18.67	52.52	2025-11-28 14:38:54.200251+00
aafea20a-5f64-4375-8c25-c7e2a97e3008	prod-agent	7.81	18.61	52.52	2025-11-28 14:39:04.200702+00
33d0eaf7-421b-4f6b-ac29-ba6461fee7bd	prod-agent	7.85	18.61	52.52	2025-11-28 14:39:14.195813+00
29ac8033-6013-492b-9cc4-9604d8a5b3ee	prod-agent	15.58	18.68	52.52	2025-11-28 14:39:24.193142+00
ab0461f3-2da0-4e76-871d-e8d5c707d899	prod-agent	7.98	18.66	52.52	2025-11-28 14:39:34.199636+00
25ed864a-df2d-47eb-a797-3359fffd0b09	prod-agent	9.45	18.66	52.52	2025-11-28 14:39:44.19416+00
58c103c0-1a52-4e69-9d3f-0482c4ce9576	prod-agent	8.59	18.68	52.52	2025-11-28 14:39:54.19989+00
3eb9ee15-a056-4878-8961-99f8524b5934	prod-agent	9.41	18.65	52.52	2025-11-28 14:40:04.19622+00
7669054c-bc5a-41f6-968c-970abaec559c	prod-agent	8.69	18.67	52.52	2025-11-28 14:40:14.193574+00
18d873d3-fc2e-4a3e-bc48-9d1e84487bca	prod-agent	13.57	18.61	52.52	2025-11-28 14:40:24.199591+00
4be2a923-b9cd-4709-8186-dbe22913c710	prod-agent	8.57	18.63	52.52	2025-11-28 14:40:34.194062+00
accd78c8-d976-4237-b9e4-8ed5a5397fe9	prod-agent	8.35	18.64	52.52	2025-11-28 14:40:44.198785+00
25ae5588-adfd-4681-9351-dd90b1527e30	prod-agent	7.98	18.61	52.52	2025-11-28 14:40:54.202871+00
ddad5fff-cc3b-42b6-8b52-54b02c200aa8	prod-agent	8.23	18.68	52.52	2025-11-28 14:41:04.195809+00
45a227e4-e389-4b55-9da8-1dbe795db91d	prod-agent	7.55	18.67	52.52	2025-11-28 14:41:14.20466+00
6d9bbf73-fb3e-45de-8ae3-b11c507722d6	prod-agent	14.95	18.76	52.52	2025-11-28 14:41:24.192863+00
5236988d-5385-4eed-a11c-9c10f88ae428	prod-agent	7.27	18.65	52.52	2025-11-28 14:41:34.194056+00
744fa3ee-4982-4b7b-bf09-9fc81a04f41b	prod-agent	8.22	18.63	52.52	2025-11-28 14:41:44.194815+00
64c9f5a3-6052-4d90-a676-3f11cccdb6a8	prod-agent	8.12	18.72	52.52	2025-11-28 14:41:54.200126+00
8fab29f3-7ef5-4b23-86c0-baf2bbb338e2	prod-agent	8.72	18.66	52.52	2025-11-28 14:42:04.199804+00
d2520dfc-e3fe-4f60-a490-583d6cc8329e	prod-agent	14.67	18.92	52.52	2025-11-28 14:42:14.200223+00
50299f69-b5ab-4add-a694-054a2600ed97	prod-agent	12.37	18.66	52.52	2025-11-28 14:42:24.199988+00
99d059fb-e324-4b11-9c25-f2d68ccb3f7e	prod-agent	7.83	18.70	52.52	2025-11-28 14:42:34.200452+00
3f9f209b-03a6-417a-a828-0a1315710782	prod-agent	8.79	18.65	52.52	2025-11-28 14:42:44.194017+00
4ba218ff-7db9-48f2-a07c-b56184c5afcd	prod-agent	8.85	18.65	52.52	2025-11-28 14:42:54.201759+00
1b6ffd43-f283-4374-837e-86cc1e186c4c	prod-agent	7.13	18.68	52.52	2025-11-28 14:43:04.207297+00
93e5b1f3-0d31-4352-8961-b766c6b18506	prod-agent	8.99	18.70	52.52	2025-11-28 14:43:14.201014+00
c79203e2-739b-4cf4-89f8-f2fa49dfb0df	prod-agent	11.74	18.67	52.52	2025-11-28 14:43:24.319345+00
650f9b2a-89ef-4823-8807-40d1e703cde6	prod-agent	11.16	18.70	52.52	2025-11-28 14:43:34.200346+00
47ba453e-93a1-4f27-9f55-6b58fad23179	prod-agent	8.59	18.74	52.52	2025-11-28 14:43:44.200787+00
573bdd30-bcf0-4cde-ae84-7ea2cc5c9ff4	prod-agent	8.38	18.68	52.52	2025-11-28 14:43:54.200439+00
99d1fb35-aaed-40c3-acde-ae7b571b57d6	prod-agent	8.79	18.68	52.52	2025-11-28 14:44:04.195036+00
47f35803-9bbb-41b5-8cba-6f1d7018576d	prod-agent	8.54	18.71	52.52	2025-11-28 14:44:14.20026+00
395fab3d-a920-40ca-b369-2b5d2daa62b3	prod-agent	9.47	18.72	52.52	2025-11-28 14:44:24.194705+00
3cce7b6e-31b3-467f-a111-5ea0f6688085	prod-agent	14.31	18.75	52.52	2025-11-28 14:44:34.193207+00
f178b15a-4a1d-4769-a744-5163ff378212	prod-agent	8.74	18.71	52.52	2025-11-28 14:44:44.199243+00
6b4b62b9-0f2c-4a75-a7fe-bd881d39d0af	prod-agent	8.18	18.68	52.52	2025-11-28 14:44:54.204058+00
dca0ba8a-a7ca-4c77-ba57-5a2177090268	prod-agent	7.71	18.66	52.52	2025-11-28 14:45:04.201284+00
8c962a46-17eb-47b2-a36c-f5d29f01a401	prod-agent	8.43	18.73	52.52	2025-11-28 14:45:14.203562+00
a60023e9-814d-40e2-ad45-a0c6d7dd7ca6	prod-agent	8.41	18.74	52.52	2025-11-28 14:45:24.198485+00
a575421a-07a3-442a-8dec-eba064924ebe	prod-agent	13.71	18.68	52.52	2025-11-28 14:45:34.199469+00
265f5776-f346-49b0-a51b-f6cc3852e947	prod-agent	7.80	18.64	52.52	2025-11-28 14:45:44.194843+00
df2ec2fc-e7c8-471d-a5d3-af14d0a1b3fa	prod-agent	8.60	18.69	52.52	2025-11-28 14:45:54.199462+00
a03f1b67-c7a0-42c1-b7cc-e23455f06db4	prod-agent	9.32	18.70	52.52	2025-11-28 14:46:04.20053+00
07571456-1293-4b21-9149-e1fd106c19a7	prod-agent	8.55	18.72	52.52	2025-11-28 14:46:14.200291+00
11a79fa9-6971-4eb6-b914-2f0c3c727794	prod-agent	22.21	19.19	52.52	2025-11-28 14:46:24.198174+00
6f4f59bd-3556-4910-a509-9d25a552b0d1	prod-agent	15.81	18.70	52.52	2025-11-28 14:46:34.202258+00
caba1e66-c583-47bd-94f8-c2319a7a16de	prod-agent	8.94	18.68	52.52	2025-11-28 14:46:44.193813+00
c07918d0-c045-4103-af9e-36d8c2e2cd96	prod-agent	8.61	18.66	52.52	2025-11-28 14:46:54.204615+00
31171a64-ce9b-4896-b583-3c06dbf9f31a	prod-agent	7.08	18.64	52.52	2025-11-28 14:47:04.202134+00
13cc9296-3eef-4baa-954e-bd26f67a5ad5	prod-agent	8.63	18.69	52.52	2025-11-28 14:47:14.197694+00
53104774-56d6-4e3d-b866-fdc3956a32de	prod-agent	9.31	18.74	52.52	2025-11-28 14:47:24.194688+00
6377a395-239d-461b-b832-6085ca4fdf20	prod-agent	11.78	18.75	52.52	2025-11-28 14:47:34.199616+00
b04387c7-678a-490c-a397-2e2caaea5237	prod-agent	9.38	18.75	52.52	2025-11-28 14:47:44.192913+00
e4e2929e-a83e-4061-84ab-d3b844a4a97e	prod-agent	8.50	18.74	52.52	2025-11-28 14:47:54.194325+00
eecfdb8a-7b30-4f1c-9bd3-3e5ca5924809	prod-agent	8.01	18.68	52.52	2025-11-28 14:48:04.201309+00
0f32dd37-bcdb-42f7-a5d1-092fcb6755c2	prod-agent	9.13	18.76	52.52	2025-11-28 14:48:14.19761+00
c23fd9dd-3de4-4ecc-b083-fe09825cb2fe	prod-agent	9.34	18.70	52.52	2025-11-28 14:48:24.194207+00
e81e2fce-e753-4bf4-9fdd-f45dbd4fa87b	prod-agent	11.00	18.71	52.52	2025-11-28 14:48:34.267852+00
6b71ac01-3fb3-46e4-bc96-040e1ccb9b6b	prod-agent	11.09	18.71	52.52	2025-11-28 14:48:44.200952+00
88cd05a3-c357-4a5c-a039-13de2930ded9	prod-agent	8.51	18.65	52.52	2025-11-28 14:48:54.219004+00
548b79fd-b944-4fe3-920e-1df025a56b04	prod-agent	6.62	18.62	52.52	2025-11-28 14:49:04.200248+00
396442a9-aecb-4b0d-97b7-de43a9a14f30	prod-agent	9.23	18.69	52.52	2025-11-28 14:49:14.193755+00
17d6254d-d1f8-4d60-8fb9-943b2c9d6afa	prod-agent	9.44	18.66	52.52	2025-11-28 14:49:24.194613+00
35beb6e2-e24a-4125-818f-28cf796d299d	prod-agent	7.43	18.67	52.52	2025-11-28 14:49:34.200469+00
d72178a2-679d-47fd-8a0b-1bc5702ca736	prod-agent	14.66	18.74	52.52	2025-11-28 14:49:44.194281+00
8d95e24e-6ea8-40cb-8068-400d906309d2	prod-agent	7.45	18.69	52.52	2025-11-28 14:49:54.201058+00
b788b0ec-ac42-4797-95a5-818be34630d9	prod-agent	8.02	18.70	52.52	2025-11-28 14:50:04.200413+00
f34f8219-74ff-4a61-a864-1bfe50a2b8dc	prod-agent	9.93	18.74	52.52	2025-11-28 14:50:14.194609+00
60e7dc82-cdd1-426d-a917-c345279f624e	prod-agent	8.72	18.72	52.52	2025-11-28 14:50:24.202935+00
8570d123-1e52-49b6-9eaa-0e8e6191d2f4	prod-agent	7.77	18.68	52.52	2025-11-28 14:50:34.200302+00
4a7c67ea-58e3-40e2-95d1-92d856ac905d	prod-agent	15.11	18.75	52.52	2025-11-28 14:50:44.202412+00
05d33fce-ea73-40ad-b6d6-6a9f778087d5	prod-agent	7.25	18.72	52.52	2025-11-28 14:50:54.194854+00
f8b4f9a5-2f84-4851-960c-1854595a7086	prod-agent	7.23	18.68	52.52	2025-11-28 14:51:04.208395+00
6979bf6c-99a6-40b0-abc4-845d21b3896a	prod-agent	9.32	18.72	52.52	2025-11-28 14:51:14.199649+00
868e35d0-b109-4a49-8751-00b0a93b5b7d	prod-agent	8.28	18.73	52.52	2025-11-28 14:51:24.194453+00
8549047d-2d90-4fb7-9948-fa7513ef112a	prod-agent	8.48	18.70	52.52	2025-11-28 14:51:34.202826+00
0faac5d3-1988-46ce-a368-3edaea353bd2	prod-agent	15.13	18.78	52.52	2025-11-28 14:51:44.199269+00
3a0edd69-1e3c-4045-bfa2-119d1ad2a68b	prod-agent	7.00	18.68	52.52	2025-11-28 14:51:54.195438+00
9906a53f-a8ae-4018-8b52-6db7a39e6b42	prod-agent	8.09	18.72	52.52	2025-11-28 14:52:04.204143+00
c338f485-e23c-4e35-a88e-5bf02d8e374e	prod-agent	10.29	18.77	52.52	2025-11-28 14:52:14.198322+00
4c6cc6ec-9051-4216-b16e-56812ec0d558	prod-agent	8.72	18.68	52.52	2025-11-28 14:52:24.194277+00
f9683d8f-c2a2-4df1-bc32-fb2d8dca1134	prod-agent	8.63	18.71	52.52	2025-11-28 14:52:34.194604+00
204980e0-f533-43f6-98aa-d9db011477de	prod-agent	15.63	18.82	52.52	2025-11-28 14:52:44.204315+00
c13b89dd-f81b-43e9-903d-144ffe205eba	prod-agent	5.91	18.66	52.52	2025-11-28 14:52:54.205319+00
26005129-9aca-428e-8436-1528cfb67ba3	prod-agent	8.29	18.70	52.52	2025-11-28 14:53:04.206388+00
831cd737-cb93-4208-b7d0-78efb59dfeb7	prod-agent	9.92	18.65	52.52	2025-11-28 14:53:14.206393+00
89dc3588-7517-429a-9436-395ca98d3d81	prod-agent	8.36	18.63	52.52	2025-11-28 14:53:24.202628+00
2a295757-a9bb-49b5-9314-77a0165ddfe1	prod-agent	7.85	18.61	52.52	2025-11-28 14:53:34.194332+00
649452a9-45c2-4baf-b796-340cd5e78caa	prod-agent	11.27	18.70	52.52	2025-11-28 14:53:44.294128+00
6f448e22-a883-4b7e-8c86-77348645b2ab	prod-agent	10.51	18.66	52.52	2025-11-28 14:53:54.200396+00
ec8a399c-c93a-4d61-b95c-5c34bf1b16d8	prod-agent	8.65	18.71	52.52	2025-11-28 14:54:04.201199+00
40c18de6-1b06-4e2a-877e-fc7c8f1cd038	prod-agent	8.85	18.65	52.52	2025-11-28 14:54:14.201542+00
440b54a9-de7b-4e6a-a6c3-908ba14ef840	prod-agent	8.21	18.65	52.52	2025-11-28 14:54:24.203048+00
c0ef1eb9-1a67-4dd3-ae44-25b3c3a937eb	prod-agent	8.40	18.66	52.52	2025-11-28 14:54:34.193697+00
320330e5-e25f-48ef-ae31-155f9013ba6a	prod-agent	9.46	18.72	52.52	2025-11-28 14:54:44.200845+00
cf27c3e4-2c78-4b40-bf4e-f3bd3c6e1b3c	prod-agent	12.51	18.73	52.52	2025-11-28 14:54:54.201348+00
55b65a8b-6a02-4362-82ae-4de7b4dafce7	prod-agent	8.32	18.67	52.52	2025-11-28 14:55:04.19544+00
b1785c1f-d3e8-4fba-a7a3-25a07b943d8a	prod-agent	9.55	18.71	52.52	2025-11-28 14:55:14.203341+00
dae8e8f7-287b-40c3-83a2-059a094c930c	prod-agent	8.68	18.67	52.52	2025-11-28 14:55:24.200343+00
847c09a8-935e-4842-8273-5d008177b9b1	prod-agent	8.77	18.69	52.52	2025-11-28 14:55:34.199913+00
cc316a3a-67f0-4ca8-a03c-5f74a32fe30c	prod-agent	10.09	18.70	52.52	2025-11-28 14:55:44.203777+00
644219d8-7572-48e3-abb7-64e3962d15b0	prod-agent	12.31	18.73	52.52	2025-11-28 14:55:54.198014+00
2fccdd45-a705-413c-93e8-eab2e56e29a7	prod-agent	7.64	18.75	52.52	2025-11-28 14:56:04.203197+00
e78494c5-f0e0-42d9-8d15-8ed2745d0c52	prod-agent	9.28	18.75	52.52	2025-11-28 14:56:14.193653+00
3bc163ed-ed3b-48aa-8014-2eedc7d6ada5	prod-agent	7.94	18.66	52.52	2025-11-28 14:56:24.201756+00
ee047e3f-f4f7-47b8-86ab-5c9222a476ea	prod-agent	7.75	18.74	52.52	2025-11-28 14:56:34.203346+00
207e6be7-9761-432c-8ae8-7e5e56b9028e	prod-agent	9.72	18.68	52.52	2025-11-28 14:56:44.195701+00
e0d651eb-5526-4dd9-a1a4-fd6d12714ccb	prod-agent	11.65	18.76	52.52	2025-11-28 14:56:54.199443+00
fec8ac5b-07ee-4d14-9702-80388366e97f	prod-agent	7.90	18.77	52.52	2025-11-28 14:57:04.194313+00
dcb081dd-b499-455d-be64-5ef1e48ffc4b	prod-agent	9.95	18.62	52.52	2025-11-28 14:57:14.195876+00
280eb40b-3a55-4b8f-817c-2c87a64b5754	prod-agent	8.29	18.62	52.52	2025-11-28 14:57:24.203087+00
e940153a-a76c-46d6-be34-9ee9fc11f022	prod-agent	8.32	18.64	52.52	2025-11-28 14:57:34.194699+00
51c75bb4-2b6a-47ee-85ad-1111e40276b0	prod-agent	8.78	18.69	52.52	2025-11-28 14:57:44.200644+00
8ea3e420-3703-46cf-a5d7-58be5a8cc464	prod-agent	13.45	18.68	52.52	2025-11-28 14:57:54.202248+00
685adb3a-5fcf-464e-b90e-525e8f0ae5eb	prod-agent	7.71	18.68	52.52	2025-11-28 14:58:04.199779+00
25b4847f-d29f-4f56-b2bd-00d1a252bf37	prod-agent	9.01	18.61	52.52	2025-11-28 14:58:14.200058+00
4996d84d-6d93-49f3-beff-e064664d5fa5	prod-agent	8.69	18.65	52.52	2025-11-28 14:58:24.20027+00
a69931e1-fbb1-4d30-9e64-215d5805bf68	prod-agent	7.98	18.61	52.52	2025-11-28 14:58:34.203197+00
e4750908-44cc-4446-aea4-6911dde32e0c	prod-agent	8.38	18.61	52.52	2025-11-28 14:58:44.202127+00
5cd83eae-ab92-4a83-bdc0-361a61e44ea1	prod-agent	9.16	18.62	52.52	2025-11-28 14:58:54.254666+00
35775af1-cdfc-40c3-bfcb-ea8269ec2ee7	prod-agent	12.63	18.62	52.52	2025-11-28 14:59:04.194404+00
8b2abcb4-f209-42b5-8cfc-a90702095348	prod-agent	9.01	18.70	52.52	2025-11-28 14:59:14.20487+00
573dbe4f-bb96-40a4-9573-75e6f0bf410c	prod-agent	9.10	18.73	52.52	2025-11-28 14:59:24.2001+00
5b52b13c-de75-465c-8d44-c16ac42957fb	prod-agent	8.97	18.74	52.52	2025-11-28 14:59:34.200184+00
f4441745-a875-4e15-b389-2c219bb16db8	prod-agent	8.94	18.69	52.52	2025-11-28 14:59:44.200406+00
478cbd67-e3b8-43fd-b88e-3993e3deacf2	prod-agent	8.72	18.67	52.52	2025-11-28 14:59:54.198236+00
e8b8cd8d-7d59-40bb-8eec-cc23a06ac5be	prod-agent	14.58	18.69	52.52	2025-11-28 15:00:04.19779+00
63f5b03e-dee5-41dd-a26c-9000b9863c37	prod-agent	10.05	18.90	52.52	2025-11-28 15:00:14.200116+00
15835f8e-a337-4926-a44f-17517a10c35c	prod-agent	9.07	18.80	52.52	2025-11-28 15:00:24.193263+00
658eeffb-1947-4911-8c91-e36afc60de8c	prod-agent	7.82	18.84	52.52	2025-11-28 15:00:34.203876+00
48bfa60c-1928-4d91-bd33-ee9008eae50b	prod-agent	8.48	18.85	52.52	2025-11-28 15:00:44.203032+00
26e36fca-17ec-4ae8-b554-91f119b50809	prod-agent	8.24	18.87	52.52	2025-11-28 15:00:54.19487+00
75e3d6b6-54dd-485f-b664-ae7f92c5d6ac	prod-agent	15.15	18.83	52.52	2025-11-28 15:01:04.199493+00
4eeb6313-9eea-4ba7-b496-fe8b9a05dc27	prod-agent	8.24	19.04	52.52	2025-11-28 15:01:14.194445+00
d3534c49-3a69-49c2-a1b2-28a45fb0393e	prod-agent	20.20	19.44	52.52	2025-11-28 15:01:24.200083+00
5d761cc0-d123-4164-b1fc-72f5ac838a2c	prod-agent	10.57	19.07	52.52	2025-11-28 15:01:34.192531+00
dc9bbb8e-859d-4d42-a906-06c6347f4686	prod-agent	8.55	19.00	52.52	2025-11-28 15:01:44.198416+00
c3bd3d63-3177-4eea-b0e9-e41472690893	prod-agent	8.64	19.05	52.52	2025-11-28 15:01:54.194081+00
ee486153-5936-49ab-ad6d-76baec153043	prod-agent	14.39	19.01	52.52	2025-11-28 15:02:04.199788+00
dcb003d0-40f2-49f7-9ac5-adf444d3d7e1	prod-agent	8.48	18.71	52.52	2025-11-28 15:02:14.193922+00
93248a06-82d3-4f23-8269-1041ba676b7f	prod-agent	8.58	18.67	52.52	2025-11-28 15:02:24.194554+00
90f95ca1-c6c3-4dbd-9695-176a24dfa381	prod-agent	7.39	18.69	52.52	2025-11-28 15:02:34.199574+00
2a966172-a2cd-40bf-be7f-4ee3aaccbca3	prod-agent	7.78	18.70	52.52	2025-11-28 15:02:44.194383+00
056ceb23-a7f5-4e49-9703-55fe496ed9cc	prod-agent	7.74	18.65	52.52	2025-11-28 15:02:54.202535+00
0873b626-2cdc-4b1e-9353-caf092d6ceb3	prod-agent	12.03	18.80	52.52	2025-11-28 15:03:04.199407+00
4c98676e-ce10-4d0e-af9e-3bf67bd07842	prod-agent	9.35	18.78	52.52	2025-11-28 15:03:14.200769+00
c22b761a-f6f4-4316-bc0b-a9d220992436	prod-agent	9.52	18.69	52.52	2025-11-28 15:03:24.194789+00
4d5a233c-11c1-4410-8d02-1afad55a8b93	prod-agent	8.27	18.67	52.52	2025-11-28 15:03:34.200224+00
b1ef1086-0be1-4d8d-ae72-766703fa9066	prod-agent	8.82	18.67	52.52	2025-11-28 15:03:44.194023+00
6f9d45f6-79a4-4ebf-8a1a-e2718326f56d	prod-agent	8.65	18.71	52.52	2025-11-28 15:03:54.199375+00
a3c5a84e-21f2-428f-b8a4-6f82a31df16a	prod-agent	15.59	19.03	52.52	2025-11-28 15:04:04.262165+00
279b6eab-abd9-4ce7-b1c4-00f97d920291	prod-agent	22.93	18.67	52.52	2025-11-28 15:04:14.196942+00
1787719b-ff44-4245-a2fd-11c6b3ec33da	prod-agent	7.82	18.58	52.52	2025-11-28 15:04:24.203507+00
cbfae6a2-b4f7-4190-9c4b-b15c0ea6f1bc	prod-agent	7.54	18.64	52.52	2025-11-28 15:04:34.194967+00
facfd57a-1e74-4673-9bdf-6602c3bba17d	prod-agent	8.45	18.58	52.52	2025-11-28 15:04:44.203942+00
fe25feaa-9687-470f-ad11-fee026b233d6	prod-agent	8.89	18.62	52.52	2025-11-28 15:04:54.201081+00
4e59ede7-73c9-4d33-90a0-18952cfb98ac	prod-agent	8.31	18.64	52.52	2025-11-28 15:05:04.194839+00
96f3c18b-b733-449c-9732-750a7e24ae7a	prod-agent	12.88	18.66	52.52	2025-11-28 15:05:14.198776+00
ab4c5f03-b1d2-4cc0-af4a-10965881fb7e	prod-agent	8.60	18.62	52.52	2025-11-28 15:05:24.198075+00
2b7510a9-ea85-4171-a35d-949e047017a9	prod-agent	8.02	18.60	52.52	2025-11-28 15:05:34.200349+00
53b8a231-f22a-4ba9-864f-ac8b7763056b	prod-agent	9.02	18.63	52.52	2025-11-28 15:05:44.196937+00
4b3a827b-3bfe-4bdf-b20b-cf76f442f11d	prod-agent	8.35	18.66	52.52	2025-11-28 15:05:54.201691+00
83d03c17-057c-4ceb-bc51-d3890eafbd66	prod-agent	8.86	18.66	52.52	2025-11-28 15:06:04.194372+00
28bcdb3a-8ede-457b-9b7d-eec27f9fd45f	prod-agent	13.68	18.74	52.52	2025-11-28 15:06:14.19984+00
a1353e42-7ee7-4c59-bd88-f29f88048cf1	prod-agent	8.67	18.70	52.52	2025-11-28 15:06:24.194466+00
62e8a840-a8df-46f2-b0df-87cc88b4c474	prod-agent	7.91	18.68	52.52	2025-11-28 15:06:34.196048+00
4d63056d-fb48-4acc-9dc8-04c390247773	prod-agent	9.45	18.52	52.52	2025-11-28 15:06:44.194384+00
b01219c7-3318-4c1e-8948-d298011dacb7	prod-agent	8.14	18.55	52.52	2025-11-28 15:06:54.200513+00
c2073709-f1a2-4f7c-abdb-a3dc19a54ba2	prod-agent	8.77	18.59	52.52	2025-11-28 15:07:04.194697+00
06a5d142-5d4a-49b6-8f43-41c2e1eecffd	prod-agent	13.91	18.67	52.52	2025-11-28 15:07:14.200553+00
221d6235-33b0-48ba-9d73-a407539cba24	prod-agent	8.06	18.54	52.52	2025-11-28 15:07:24.19764+00
74ff84ea-2bde-4555-9251-1145a4ae5b59	prod-agent	8.39	18.53	52.52	2025-11-28 15:07:34.208977+00
a97c06f1-3f1c-4c68-897c-8ba4dff7da5d	prod-agent	8.62	18.51	52.52	2025-11-28 15:07:44.202074+00
4db320e5-ab7d-4d65-bfac-3d6cbd18ee9e	prod-agent	7.72	18.53	52.52	2025-11-28 15:07:54.202799+00
ed4ff7fb-1d5f-4012-a9ce-f2983db837a6	prod-agent	8.84	18.54	52.52	2025-11-28 15:08:04.195022+00
8508c6b2-487f-40d0-8498-51b4086585cb	prod-agent	12.91	18.54	52.52	2025-11-28 15:08:14.194765+00
1c258c43-c9e0-42e2-9fd8-a2ef522b6cb5	prod-agent	9.61	18.53	52.52	2025-11-28 15:08:24.200797+00
c57cb02c-59cb-4e6a-999d-238cf892bfbe	prod-agent	7.32	18.48	52.52	2025-11-28 15:08:34.195068+00
20eb4de6-c394-4120-911b-988d35ab3ffb	prod-agent	9.02	18.49	52.52	2025-11-28 15:08:44.194609+00
f09848c9-a6f6-44bd-bbb2-4f184f1bdee7	prod-agent	8.27	18.43	52.52	2025-11-28 15:08:54.200152+00
25285cee-2714-4171-abb2-a177ef728fea	prod-agent	7.54	18.45	52.52	2025-11-28 15:09:04.200089+00
0105ee38-2824-4cf1-ab4a-878593d2ea72	prod-agent	10.15	18.38	52.52	2025-11-28 15:09:14.250964+00
48f3d01c-bcd0-4613-9fac-33c40a9627db	prod-agent	15.14	18.47	52.52	2025-11-28 15:09:24.202166+00
e7deb518-f24f-4499-8765-afc7711aa852	prod-agent	7.39	18.50	52.52	2025-11-28 15:09:34.200244+00
dd05fb7e-9550-4db2-b113-61ddd6809b00	prod-agent	9.70	18.48	52.52	2025-11-28 15:09:44.195938+00
83995f89-ce39-4298-b8a9-b3fbf2980e7b	prod-agent	7.42	18.49	52.52	2025-11-28 15:09:54.20588+00
f318d6e7-e814-48f8-a872-46a06e919248	prod-agent	8.60	18.50	52.52	2025-11-28 15:10:04.201222+00
9e3cd722-eda0-474a-96d5-4d1797eb0f90	prod-agent	9.13	18.57	52.52	2025-11-28 15:10:14.194263+00
48500f07-01dd-4f71-afd4-2c4b8ae80ebc	prod-agent	13.88	18.55	52.52	2025-11-28 15:10:24.199847+00
90f634fe-edaa-4b15-adbd-110cddd30684	prod-agent	7.09	18.52	52.52	2025-11-28 15:10:34.200433+00
52af039a-515e-4e9a-bef7-468c9593295b	prod-agent	9.68	18.50	52.52	2025-11-28 15:10:44.200519+00
8fbe7761-c28b-4503-84d3-845b6658e5a1	prod-agent	8.65	18.56	52.52	2025-11-28 15:10:54.20486+00
656cdb86-838a-4ba9-b572-9d7601053c12	prod-agent	7.28	18.51	52.52	2025-11-28 15:11:04.207964+00
a03d8251-ef5a-43ee-afd8-880860ca44f3	prod-agent	8.76	18.46	52.52	2025-11-28 15:11:14.198987+00
2d8bb148-9ab4-4191-bfc3-f9888fc584f4	prod-agent	16.40	18.52	52.52	2025-11-28 15:11:24.197078+00
27a6f3f0-e3e9-4515-8046-c04ab510d52e	prod-agent	6.97	18.55	52.52	2025-11-28 15:11:34.194641+00
89379b92-df78-40a5-9de6-c17e083e88d1	prod-agent	8.81	18.53	52.52	2025-11-28 15:11:44.202576+00
635cbdfb-0acf-48c9-b9ce-52833c4c010e	prod-agent	8.89	18.54	52.52	2025-11-28 15:11:54.194408+00
babfdcdc-c054-46a0-8c5c-8e6870140b8d	prod-agent	7.89	18.57	52.52	2025-11-28 15:12:04.194621+00
5bab0163-c1ec-456b-bbf2-625cf4187822	prod-agent	7.76	18.53	52.52	2025-11-28 15:12:14.199609+00
6b597f63-32a4-4118-8c43-58a94f85bc8a	prod-agent	14.58	18.53	52.52	2025-11-28 15:12:24.194469+00
f546df2d-d24a-4fd7-bee1-40ec931f88f2	prod-agent	6.90	18.56	52.52	2025-11-28 15:12:34.199047+00
d9a8f56b-8535-4592-a4ce-2c039f4b255e	prod-agent	8.59	18.62	52.52	2025-11-28 15:12:44.199192+00
6f70baf9-3940-4681-ba57-5d2c49a84ddd	prod-agent	8.93	18.62	52.52	2025-11-28 15:12:54.204165+00
df6db8cd-66de-45e0-9bc5-b3f10f555d6a	prod-agent	7.77	18.62	52.52	2025-11-28 15:13:04.202815+00
16d5943f-a684-4180-8999-8a8f5105e383	prod-agent	8.90	18.62	52.52	2025-11-28 15:13:14.199931+00
b884a46a-7636-4443-898e-cba052ce69ca	prod-agent	13.41	18.50	52.52	2025-11-28 15:13:24.220022+00
15a0910e-6e0d-46b2-8968-f74b25d7891c	prod-agent	9.24	18.49	52.52	2025-11-28 15:13:34.206654+00
2f097665-e07d-43c8-bdd5-707500a58460	prod-agent	8.90	18.56	52.52	2025-11-28 15:13:44.196643+00
26a9cdb4-ee5d-4cf3-a4c6-484f23218cc0	prod-agent	9.29	18.61	52.52	2025-11-28 15:13:54.2033+00
0695c8c7-f2f8-4b2c-b27e-b76a5a1762a8	prod-agent	8.19	18.61	52.52	2025-11-28 15:14:04.202098+00
e72c6eb7-79fd-4fb7-8944-a40f492552af	prod-agent	7.48	18.61	52.52	2025-11-28 15:14:14.20193+00
bf844c61-739d-44ec-aa90-28b944166f35	prod-agent	9.43	18.68	52.52	2025-11-28 15:14:24.201851+00
c1eeaa59-a3cb-44fa-80da-8e666d032b97	prod-agent	13.12	18.70	52.52	2025-11-28 15:14:34.197692+00
01ae5b51-aee3-4238-b61d-984e4b414ff5	prod-agent	8.46	18.62	52.52	2025-11-28 15:14:44.200551+00
b7da8793-4dd0-4d34-9885-2e8eee8fe3d3	prod-agent	9.48	18.68	52.52	2025-11-28 15:14:54.19401+00
03503778-6d99-4920-a361-0491864897e4	prod-agent	8.24	18.67	52.52	2025-11-28 15:15:04.199451+00
3d252e61-3e7c-440f-b0db-6a737667d3aa	prod-agent	8.69	18.63	52.52	2025-11-28 15:15:14.204015+00
c8dea9a6-6d0a-4928-abbe-62ddf980d2c5	prod-agent	10.24	18.69	52.52	2025-11-28 15:15:24.199591+00
1b81e00b-f903-490a-92f8-27acdea2ccd0	prod-agent	11.99	18.69	52.52	2025-11-28 15:15:34.19559+00
dd35f772-3853-4734-9fcf-a947e4d7e4f5	prod-agent	9.58	18.66	52.52	2025-11-28 15:15:44.197489+00
dd1aad40-fa22-4c06-905d-4636b5894477	prod-agent	9.39	18.68	52.52	2025-11-28 15:15:54.199623+00
ccfce2af-e751-4e11-bbcf-e73de2ce2a23	prod-agent	8.13	18.65	52.52	2025-11-28 15:16:04.1948+00
c8490c41-124b-49cc-8210-9956508b5221	prod-agent	8.15	18.59	52.52	2025-11-28 15:16:14.201182+00
f9c01907-107c-4659-b9af-a1714a914781	prod-agent	20.33	18.93	52.52	2025-11-28 15:16:24.195688+00
e553e920-1bdb-4704-8a50-2f38dec584d7	prod-agent	16.28	18.64	52.52	2025-11-28 15:16:34.200775+00
b34d1508-71ff-4d3c-b71b-d1b94be7ff18	prod-agent	7.34	18.63	52.52	2025-11-28 15:16:44.194018+00
f9987b0a-b995-4f3a-b972-6e6cbfe6eea3	prod-agent	9.22	18.59	52.52	2025-11-28 15:16:54.197954+00
5f557964-aafb-4ff4-8b20-7eff9103a67c	prod-agent	7.28	18.57	52.52	2025-11-28 15:17:04.199628+00
f716ddce-3c7b-4cca-b5f9-be1b3c972d00	prod-agent	8.93	18.61	52.52	2025-11-28 15:17:14.201144+00
f211ae54-63f1-4435-ad7c-95675c709f32	prod-agent	9.21	18.57	52.52	2025-11-28 15:17:24.198978+00
34e14298-4889-4d74-90ae-3dcb1bc3f2ea	prod-agent	13.73	18.64	52.52	2025-11-28 15:17:34.197883+00
22fd8064-0c7e-4335-bdf5-997b6ec88803	prod-agent	7.56	18.73	52.52	2025-11-28 15:17:44.215222+00
bf18d0ac-3401-427a-8677-d22c0443b7a4	prod-agent	9.59	18.67	52.52	2025-11-28 15:17:54.199944+00
6aaace98-a689-4f50-adc2-47f6589c793d	prod-agent	7.79	18.69	52.52	2025-11-28 15:18:04.199631+00
5ebbf527-3ff2-4bde-b0d4-2b3315d24ee6	prod-agent	8.00	18.50	52.52	2025-11-28 15:18:14.194351+00
c34986bf-11a6-429d-b717-b762aec06462	prod-agent	9.93	18.51	52.52	2025-11-28 15:18:24.203052+00
483a8425-73f4-4efe-8c9e-43bf72f178fd	prod-agent	11.11	18.59	52.52	2025-11-28 15:18:34.200029+00
0a065fcc-d8e0-49ae-ab28-99045b897028	prod-agent	9.95	18.58	52.52	2025-11-28 15:18:44.202385+00
7051c33d-5434-44b2-bfd0-e450e5eda907	prod-agent	8.01	18.57	52.52	2025-11-28 15:18:54.193949+00
a38f7cd3-dc49-4bc4-b172-9d0db6cfbf5e	prod-agent	7.85	18.63	52.52	2025-11-28 15:19:04.203218+00
6b389d59-2120-4580-8183-dd45346f28e8	prod-agent	9.43	18.66	52.52	2025-11-28 15:19:14.202962+00
c37536e9-52ad-4f8a-93d3-7aabf0bf226f	prod-agent	10.24	18.61	52.52	2025-11-28 15:19:24.25184+00
604dccb8-b632-4d85-9a6c-28e4bf7d3c42	prod-agent	7.98	18.66	52.52	2025-11-28 15:19:34.20396+00
6b96903f-517e-4aeb-aef5-504b735b4e8f	prod-agent	14.39	18.70	52.52	2025-11-28 15:19:44.207505+00
e6325d88-4ff4-45bf-bef6-12ea280605ba	prod-agent	9.18	18.71	52.52	2025-11-28 15:19:54.200799+00
8414c93c-c1a2-436e-b258-70efea5f42c8	prod-agent	7.87	18.67	52.52	2025-11-28 15:20:04.203507+00
bce7d6bf-7c95-4aa2-9f11-4dc8595b1147	prod-agent	8.51	18.53	52.52	2025-11-28 15:20:14.203087+00
a3fcdfd2-00d6-44ef-b1a6-59f63a11bbd6	prod-agent	10.30	18.49	52.52	2025-11-28 15:20:24.203607+00
db7126f4-d736-4bca-96e8-f4baaa7ec19e	prod-agent	7.02	18.54	52.52	2025-11-28 15:20:34.194917+00
b1624a4b-1a8f-445e-a70e-c4d6f0cebbb5	prod-agent	15.42	18.61	52.52	2025-11-28 15:20:44.199313+00
9729ffa2-7723-4e83-a5a5-01272a4fd2c2	prod-agent	8.56	18.60	52.52	2025-11-28 15:20:54.203044+00
14b26564-ff33-434f-b136-3c8a7a442033	prod-agent	6.37	18.60	52.52	2025-11-28 15:21:04.20045+00
7668829b-362e-4156-b525-258a0fa109c9	prod-agent	9.13	18.57	52.52	2025-11-28 15:21:14.194262+00
245f7abc-9892-4afb-9602-72b04ce7fdad	prod-agent	9.25	18.58	52.52	2025-11-28 15:21:24.198779+00
41b34543-d186-44ad-9484-5a7af425ed13	prod-agent	7.55	18.56	52.52	2025-11-28 15:21:34.194424+00
792bb63c-38e0-4092-9c76-981f1455ff0f	prod-agent	14.35	18.65	52.52	2025-11-28 15:21:44.193215+00
0efbe37d-2e31-443f-9ed3-abf6e09b7b8b	prod-agent	9.01	18.57	52.52	2025-11-28 15:21:54.200435+00
96725b1c-589c-4730-b9ae-ae91896026c0	prod-agent	7.44	18.66	52.52	2025-11-28 15:22:04.199202+00
5eb17bac-89da-4143-a800-39914e3d6b30	prod-agent	7.66	18.65	52.52	2025-11-28 15:22:14.194277+00
0274ebe5-c48e-4e43-8332-aa190768e6ad	prod-agent	10.30	18.64	52.52	2025-11-28 15:22:24.193856+00
4239b4e7-766e-494c-bcbc-f615d6f4642d	prod-agent	6.83	18.66	52.52	2025-11-28 15:22:34.199875+00
6fb04308-0c94-453c-b073-c97b8ce54b55	prod-agent	14.78	18.67	52.52	2025-11-28 15:22:44.203491+00
e68f7e62-b2d1-491a-bfa9-188bb0d644b1	prod-agent	8.46	18.52	52.52	2025-11-28 15:22:54.201466+00
0dbebbe4-d253-4398-962b-a81e0d322464	prod-agent	7.18	18.54	52.52	2025-11-28 15:23:04.19649+00
0b1a7766-0cc6-423e-8e79-f85813e698c6	prod-agent	8.17	18.47	52.52	2025-11-28 15:23:14.201376+00
adb0ca7b-05fb-46ac-9474-7e4ac75e38fe	prod-agent	8.95	18.61	52.52	2025-11-28 15:23:24.198867+00
1cfdbc1b-30ef-4e56-bc16-5b39002522c9	prod-agent	7.03	18.62	52.52	2025-11-28 15:23:34.200778+00
0e91fe23-d899-4a89-ba45-371e2e8fc006	prod-agent	12.73	18.70	52.52	2025-11-28 15:23:44.194745+00
559735f3-0e25-44cb-8aad-ba8d66111a86	prod-agent	10.44	18.62	52.52	2025-11-28 15:23:54.199289+00
5b7a1967-6cb8-4300-a416-bf230f83e4f1	prod-agent	7.46	18.67	52.52	2025-11-28 15:24:04.199457+00
78ee8098-b3d9-44eb-9c25-8003c1f2393a	prod-agent	8.13	18.54	52.52	2025-11-28 15:24:14.195198+00
e66c689b-52c4-4c65-926d-15a715e58c26	prod-agent	9.43	18.58	52.52	2025-11-28 15:24:24.200499+00
346b3db0-d5d1-4e8e-996e-f8ac5966906a	prod-agent	7.51	18.55	52.52	2025-11-28 15:24:34.258564+00
6dd86eed-8c9e-467e-a820-5261c5663a15	prod-agent	8.74	18.57	52.52	2025-11-28 15:24:44.194496+00
305fa893-dee6-44ff-b742-3f6062a05118	prod-agent	14.59	18.63	52.52	2025-11-28 15:24:54.198789+00
c9103227-ae34-42fe-a0fc-99c04a68be48	prod-agent	7.95	18.63	52.52	2025-11-28 15:25:04.195959+00
5914a46c-a98e-405f-93fb-ab10c429696b	prod-agent	9.35	18.65	52.52	2025-11-28 15:25:14.203204+00
b65400ed-383e-4d0a-917e-4c7f8f720fb8	prod-agent	9.80	18.68	52.52	2025-11-28 15:25:24.193688+00
bb7c9133-a1f5-4b2c-a368-cef2a705e8bf	prod-agent	6.85	18.64	52.52	2025-11-28 15:25:34.199393+00
dbcfbe5d-f9d1-44de-8ed1-4b761148f95c	prod-agent	8.73	18.66	52.52	2025-11-28 15:25:44.196607+00
21d615f9-2c97-4f2a-9710-93a46396b0b6	prod-agent	14.83	18.61	52.52	2025-11-28 15:25:54.199012+00
4a0bc488-5db6-4178-9705-9020b9c3d870	prod-agent	9.64	18.66	52.52	2025-11-28 15:26:04.215492+00
cbee306d-c205-44e7-bc8d-407045c03d48	prod-agent	6.96	18.68	52.52	2025-11-28 15:26:14.203196+00
855d957c-7e41-41a2-a096-8797e780642f	prod-agent	10.37	18.61	52.52	2025-11-28 15:26:24.19431+00
2f9fd8ec-0f79-4249-aecf-556460870b2d	prod-agent	8.11	18.69	52.52	2025-11-28 15:26:34.196002+00
9c374014-1eeb-49e0-8dea-56c89f8a5b4e	prod-agent	8.28	18.68	52.52	2025-11-28 15:26:44.201487+00
15f995d2-2fa2-4a9b-9ae2-ae068f07e5f9	prod-agent	15.07	18.74	52.52	2025-11-28 15:26:54.200129+00
513f5046-98f9-4aab-9736-7a329ebb905f	prod-agent	7.78	18.68	52.52	2025-11-28 15:27:04.19547+00
b6dc0162-4e40-4eb4-a25d-4816bdabb2fc	prod-agent	8.10	18.64	52.52	2025-11-28 15:27:14.194153+00
2fa7b197-8e94-4201-a1ce-df908b813245	prod-agent	10.23	18.64	52.52	2025-11-28 15:27:24.208597+00
5dfa62f7-65b2-4a06-bb58-2973675475b5	prod-agent	7.88	18.63	52.52	2025-11-28 15:27:34.194941+00
80b0e187-58a3-4780-bdcc-1b2c2b276e74	prod-agent	7.39	18.70	52.52	2025-11-28 15:27:44.205805+00
0ac137bd-37a6-4381-b4ab-34cd744e2ac3	prod-agent	13.89	18.74	52.52	2025-11-28 15:27:54.208105+00
7b6f6c89-739f-4558-9e50-dcca746ea259	prod-agent	8.86	18.68	52.52	2025-11-28 15:28:04.200658+00
d0122df9-baf4-4d52-94ab-c3880062727b	prod-agent	7.21	18.69	52.52	2025-11-28 15:28:14.202884+00
17167684-0630-4ec2-93bf-e1ef8a00a3fd	prod-agent	9.91	18.61	52.52	2025-11-28 15:28:24.199386+00
a5c2123e-0304-47b9-ac77-88ffd50c5b01	prod-agent	8.06	18.67	52.52	2025-11-28 15:28:34.194468+00
18ec3903-59cf-46ee-ac4b-c35af62e2672	prod-agent	8.40	18.72	52.52	2025-11-28 15:28:44.196523+00
1f715c42-a651-4529-af99-4f76b67d1d88	prod-agent	9.23	18.64	52.52	2025-11-28 15:28:54.199377+00
00c3dd71-4a08-49cb-8390-7d5ed6db58e2	prod-agent	13.35	18.65	52.52	2025-11-28 15:29:04.200365+00
38d16e66-ae26-4ab3-b39c-459d55fa6b0c	prod-agent	6.97	18.67	52.52	2025-11-28 15:29:14.211+00
7ab225ac-5b47-438a-b025-1382222095d3	prod-agent	9.26	18.62	52.52	2025-11-28 15:29:24.193285+00
d990da31-315f-4a30-98d4-c8a45c5cf62e	prod-agent	8.23	18.68	52.52	2025-11-28 15:29:34.194502+00
f840ab68-f56f-4b74-b82a-c45952490017	prod-agent	7.39	18.68	52.52	2025-11-28 15:29:44.250453+00
c78dee81-c903-4c04-b914-dc836be6c62b	prod-agent	9.38	18.69	52.52	2025-11-28 15:29:54.203355+00
0eca5476-b8c3-4f88-8385-2479686381f4	prod-agent	12.88	18.77	52.52	2025-11-28 15:30:04.200275+00
cc18c121-8274-4bb6-b3c0-2598a9440909	prod-agent	7.34	18.69	52.52	2025-11-28 15:30:14.200887+00
e3b742f6-3d9e-4390-a7ca-c1c4a07abb2f	prod-agent	9.66	18.67	52.52	2025-11-28 15:30:24.194934+00
22a41eb9-0096-4817-a9c9-bfb57ce287ef	prod-agent	8.08	18.71	52.52	2025-11-28 15:30:34.200339+00
7cd12bd5-077c-4221-b04a-e2ddcc6a1e73	prod-agent	8.29	18.68	52.52	2025-11-28 15:30:44.196597+00
42889c45-e2eb-49c7-a9b7-7d01bd57df8d	prod-agent	9.51	18.61	52.52	2025-11-28 15:30:54.202475+00
b1aaa584-5bbf-4739-aaa2-19a6191608bc	prod-agent	14.99	18.67	52.52	2025-11-28 15:31:04.193359+00
a389ea93-96a9-403b-a04a-fa00a54cb1f5	prod-agent	6.95	18.72	52.52	2025-11-28 15:31:14.194085+00
d39a7040-ba0f-4323-8041-901b3038eb76	prod-agent	10.36	18.70	52.52	2025-11-28 15:31:24.194555+00
3031f368-2994-494c-9150-d5cb0a5c1d5f	prod-agent	22.88	18.79	52.52	2025-11-28 15:31:34.191783+00
7e2702e4-10d8-49af-b300-c7b14240b736	prod-agent	6.86	18.67	52.52	2025-11-28 15:31:44.203538+00
f38ea6b0-820b-4ab9-ab40-3fb173ae89a4	prod-agent	9.07	18.67	52.52	2025-11-28 15:31:54.20443+00
b59f5666-34aa-4f3b-b011-3fc971b756e3	prod-agent	13.09	18.64	52.52	2025-11-28 15:32:04.193563+00
70e0645d-fb30-4f7e-8727-25d9b3a52cf5	prod-agent	6.43	18.71	52.52	2025-11-28 15:32:14.19894+00
4d660ba1-cc78-411c-b977-b47e5e9cdd49	prod-agent	9.94	18.62	52.52	2025-11-28 15:32:24.199166+00
33e552f3-8d64-4cb1-876c-ac4d3068f5dc	prod-agent	8.18	18.68	52.52	2025-11-28 15:32:34.194019+00
ef7f4a17-f15d-4ad8-b86e-bc01e4383dfd	prod-agent	7.83	18.67	52.52	2025-11-28 15:32:44.199273+00
6685b3e5-93c1-4590-916e-8feec8ce2b8d	prod-agent	9.28	18.70	52.52	2025-11-28 15:32:54.201438+00
0ff43dd6-758c-4c53-8056-142a9b7be126	prod-agent	13.16	18.68	52.52	2025-11-28 15:33:04.205254+00
19f70c8f-11ad-4ace-bc9f-050f427dac17	prod-agent	8.86	18.69	52.52	2025-11-28 15:33:14.197186+00
379da526-fad7-44b2-b908-3f87410fbb66	prod-agent	10.20	18.68	52.52	2025-11-28 15:33:24.19551+00
03abe8c4-8e4d-4cb2-ba80-e32a2c031061	prod-agent	8.14	18.71	52.52	2025-11-28 15:33:34.200625+00
6b95acae-7079-457a-8629-d7d7999607e1	prod-agent	8.20	18.68	52.52	2025-11-28 15:33:44.202328+00
7b5051d7-f2ed-4d58-9c0c-b0b6fca1143b	prod-agent	8.44	18.68	52.52	2025-11-28 15:33:54.203694+00
b7fb7f34-8b8c-4903-966a-d9f236d0362f	prod-agent	14.68	19.06	52.52	2025-11-28 15:34:04.20455+00
44df6161-0cfb-45eb-a65e-c02ca2d6ccbc	prod-agent	26.83	18.59	52.52	2025-11-28 15:34:14.200964+00
936b02df-acda-4ac0-8963-96f978e07cfc	prod-agent	9.62	18.54	52.52	2025-11-28 15:34:24.196192+00
ba87f225-578b-462e-a82b-d27e7d4beedf	prod-agent	7.35	18.61	52.52	2025-11-28 15:34:34.194455+00
11fa2b90-3010-466c-8526-8e5162f215ac	prod-agent	8.22	18.59	52.52	2025-11-28 15:34:44.199189+00
8b4aa2c1-3868-44f3-a3fa-c854fcb37d9a	prod-agent	9.54	18.62	52.52	2025-11-28 15:34:54.258862+00
ca46ab15-33a4-4a01-a8bc-65a72f4a913c	prod-agent	8.85	18.59	52.52	2025-11-28 15:35:04.1948+00
01b4399a-35dd-4950-96d8-6035cebcf74a	prod-agent	12.37	18.61	52.52	2025-11-28 15:35:14.202203+00
02941822-5d02-4eb8-9679-4ac0c477df43	prod-agent	9.50	18.63	52.52	2025-11-28 15:35:24.199844+00
5c1ec14f-6fa3-4030-8d76-062dc08853b6	prod-agent	8.66	18.56	52.52	2025-11-28 15:35:34.196525+00
9836d63f-2b69-48dd-ac1b-5bbb3984eb87	prod-agent	7.85	18.64	52.52	2025-11-28 15:35:44.194048+00
4eb1531b-8c3e-457e-95c2-fb4f39237c87	prod-agent	8.28	18.71	52.52	2025-11-28 15:35:54.20016+00
6647ba3d-05e1-4ddf-98d7-8a23a27fd7a6	prod-agent	8.40	18.69	52.52	2025-11-28 15:36:04.200566+00
78c90a54-06c1-4027-8fac-d378be409fe9	prod-agent	14.45	18.64	52.52	2025-11-28 15:36:14.19976+00
bb7ae872-d179-4313-b69e-a73ae0f37b51	prod-agent	9.18	18.61	52.52	2025-11-28 15:36:24.199954+00
45d48d7b-1d81-4007-ab79-38290eab33e7	prod-agent	9.32	18.63	52.52	2025-11-28 15:36:34.194995+00
cf0aae21-81f1-4215-be68-5c1ddc70ed50	prod-agent	6.54	18.61	52.52	2025-11-28 15:36:44.19881+00
92798876-73bd-4921-b77d-0963d72e88ab	prod-agent	8.41	18.64	52.52	2025-11-28 15:36:54.199282+00
086485ed-867d-4ff6-9fe8-f92152f7b102	prod-agent	8.64	18.67	52.52	2025-11-28 15:37:04.194945+00
2a25ef36-b771-47f4-8511-ccd755ece7bb	prod-agent	13.78	18.67	52.52	2025-11-28 15:37:14.194057+00
ef7783c2-70e9-448d-af0f-ae9a715d86fc	prod-agent	8.48	18.69	52.52	2025-11-28 15:37:24.194035+00
684920d5-ade3-45f6-a957-935c1b350e11	prod-agent	9.02	18.63	52.52	2025-11-28 15:37:34.199202+00
e5e099ca-e5b2-49f6-b285-8731b747dc2e	prod-agent	7.90	18.65	52.52	2025-11-28 15:37:44.199592+00
25141ad0-ee4e-4d5d-9580-724d96a3896e	prod-agent	7.25	18.71	52.52	2025-11-28 15:37:54.194489+00
b1a27646-9799-4016-b8c1-bfd3e38a4aed	prod-agent	9.35	18.70	52.52	2025-11-28 15:38:04.206981+00
145871ca-16c3-4c41-807a-a3bae3832ce8	prod-agent	12.09	18.76	52.52	2025-11-28 15:38:14.200598+00
fca42f25-603e-4917-b3de-a674a60bc133	prod-agent	9.88	18.69	52.52	2025-11-28 15:38:24.199836+00
f67ecbe3-67fd-4f8f-9392-c8b3709d1b70	prod-agent	9.14	18.69	52.52	2025-11-28 15:38:34.194356+00
01ce786d-29d0-4a3c-9789-4d6c5f520957	prod-agent	7.84	18.65	52.52	2025-11-28 15:38:44.195512+00
a5653f68-fe5a-4883-a824-f7abb919dd47	prod-agent	8.06	18.66	52.52	2025-11-28 15:38:54.198688+00
f98a8783-0905-4093-a95f-a04ed3ee7ec1	prod-agent	9.05	18.63	52.52	2025-11-28 15:39:04.1943+00
52b15dca-66be-42af-9df8-92beb936b7f0	prod-agent	8.35	18.63	52.52	2025-11-28 15:39:14.200049+00
c048a78c-00a7-48d4-be8f-397fe6ff06aa	prod-agent	15.09	18.65	52.52	2025-11-28 15:39:24.193709+00
92afdc7c-c257-4828-acf7-25fd6fa25e60	prod-agent	8.33	18.65	52.52	2025-11-28 15:39:34.19492+00
f0a82713-1878-49f9-ac23-353e8f26081f	prod-agent	8.05	18.67	52.52	2025-11-28 15:39:44.194282+00
7593d96a-f700-45ac-ba20-cf5740713687	prod-agent	7.21	18.67	52.52	2025-11-28 15:39:54.197293+00
1555d84c-96f2-4e4d-95f9-9e13fda19653	prod-agent	9.24	18.62	52.52	2025-11-28 15:40:04.260489+00
79581b12-721a-4388-8381-4cf4d704e1cd	prod-agent	8.98	18.74	52.52	2025-11-28 15:40:14.203515+00
18469839-9ac0-4bbf-875a-626055cf35df	prod-agent	15.67	18.66	52.52	2025-11-28 15:40:24.193171+00
a2e92f5b-0da6-4389-b9a5-21152ce786ba	prod-agent	8.17	18.66	52.52	2025-11-28 15:40:34.203109+00
f199df4a-e6ee-45b2-b219-8bcaa10821df	prod-agent	8.47	18.69	52.52	2025-11-28 15:40:44.195762+00
ff9df60d-9ef6-4fe9-8995-0ccf918e9b56	prod-agent	8.74	18.72	52.52	2025-11-28 15:40:54.194415+00
b8b908f5-4e72-4e3c-8fdb-dc34fda331d4	prod-agent	8.86	18.72	52.52	2025-11-28 15:41:04.201007+00
ce2c3049-1882-4940-84fc-51a9f58f2b13	prod-agent	6.82	18.69	52.52	2025-11-28 15:41:14.194557+00
2a11de9b-c7f1-47a1-87ea-3ad7d839a7c5	prod-agent	15.09	18.74	52.52	2025-11-28 15:41:24.194134+00
dd5000e7-05ee-4f78-94a6-2fe012114483	prod-agent	9.09	18.73	52.52	2025-11-28 15:41:34.19428+00
8ccae6c0-d7d9-4cbc-aa4c-c56c30304198	prod-agent	7.73	18.64	52.52	2025-11-28 15:41:44.199713+00
965ebc76-d8e7-4cb5-850a-26a7a9e5b57d	prod-agent	7.83	18.65	52.52	2025-11-28 15:41:54.211178+00
e9257704-b6f6-491e-a63a-a6c47bcc311d	prod-agent	9.03	18.66	52.52	2025-11-28 15:42:04.200103+00
ed4c053d-527b-4d58-81e4-cdee536cf033	prod-agent	8.19	18.64	52.52	2025-11-28 15:42:14.194165+00
99bc63cb-8239-440c-96af-11bbede81f6f	prod-agent	15.20	18.66	52.52	2025-11-28 15:42:24.195604+00
7c897b6b-912f-483d-a8e2-ecd0eccd693d	prod-agent	7.78	18.53	52.52	2025-11-28 15:42:34.199552+00
acebd3e7-4d67-487d-a2d7-5a31c10fdb7d	prod-agent	7.70	18.55	52.52	2025-11-28 15:42:44.201134+00
9e749234-f2b2-4998-a4dd-40308f5cdd26	prod-agent	8.58	18.51	52.52	2025-11-28 15:42:54.202551+00
108a9fdd-e8df-49f5-905b-bb62bf06fbd1	prod-agent	8.72	18.61	52.52	2025-11-28 15:43:04.1974+00
275c5a4b-124f-4507-a85d-ad212b435999	prod-agent	8.23	18.59	52.52	2025-11-28 15:43:14.199938+00
cdd9258f-65db-4d5a-97bc-fa1557e2c9e9	prod-agent	13.07	18.77	52.52	2025-11-28 15:43:24.200639+00
2b2198b5-f3d5-4566-9b05-f852e2952ea4	prod-agent	10.10	18.58	52.52	2025-11-28 15:43:34.203083+00
7d1aa17f-581c-4a86-b841-b1e5c7fcebf6	prod-agent	8.05	18.54	52.52	2025-11-28 15:43:44.20132+00
410611d3-c020-49c3-a2a7-2d1f203d860d	prod-agent	7.24	18.56	52.52	2025-11-28 15:43:54.201051+00
6d85dc6d-4c60-4698-9c18-dccf011f9940	prod-agent	9.68	18.59	52.52	2025-11-28 15:44:04.194102+00
84bbc0e3-8285-4e83-9d08-a525a35a8ecd	prod-agent	8.34	18.55	52.52	2025-11-28 15:44:14.198199+00
65d98cf8-1619-483e-9878-9fa88b8dd7f3	prod-agent	8.40	18.57	52.52	2025-11-28 15:44:24.201546+00
70a7469a-7b4a-4114-84fb-9b0bf9c8d530	prod-agent	14.07	18.58	52.52	2025-11-28 15:44:34.193983+00
3dd4f9d4-f248-468c-8339-5bdcab0e16b0	prod-agent	7.79	18.51	52.52	2025-11-28 15:44:44.19507+00
b797985c-37e0-4490-9b56-9407a58af889	prod-agent	7.66	18.55	52.52	2025-11-28 15:44:54.201228+00
1b812192-3c6e-4dd3-a484-4560777f3db9	prod-agent	10.32	18.58	52.52	2025-11-28 15:45:04.194917+00
bf6cf955-101a-45d0-8024-5f2307339241	prod-agent	8.37	18.55	52.52	2025-11-28 15:45:14.264819+00
d5d3b457-e5a4-4919-a828-96e89dbb66fb	prod-agent	8.82	18.50	52.52	2025-11-28 15:45:24.195489+00
d9cddc0e-6cf1-4b92-8195-a0e15848aa4f	prod-agent	15.09	18.64	52.52	2025-11-28 15:45:34.193232+00
9b0762ae-d376-404a-8e8c-e5693f1ab15d	prod-agent	6.26	18.59	52.52	2025-11-28 15:45:44.201635+00
cf3ada89-da35-49ea-beef-2c03a0201503	prod-agent	7.40	18.56	52.52	2025-11-28 15:45:54.19419+00
9ebd1546-7f9c-4364-b462-63208455964e	prod-agent	10.05	18.60	52.52	2025-11-28 15:46:04.212978+00
1d79488d-3e47-4179-825b-a041cf82c822	prod-agent	7.84	18.58	52.52	2025-11-28 15:46:14.200486+00
8b857682-49f2-4fbc-a769-2d1dfd5ec030	prod-agent	8.68	18.63	52.52	2025-11-28 15:46:24.201862+00
0b2624ed-5a2c-4455-84b6-b53d59bf2870	prod-agent	31.45	18.79	52.52	2025-11-28 15:46:34.197407+00
e8f725ad-75e1-4620-99ec-10db348286f8	prod-agent	5.69	18.67	52.52	2025-11-28 15:46:44.194712+00
8324c7d2-4f97-489b-829a-1270508b116f	prod-agent	7.93	18.60	52.52	2025-11-28 15:46:54.201243+00
de1235c9-ccda-4a2c-8f5a-6082689e87ad	prod-agent	9.99	18.63	52.52	2025-11-28 15:47:04.200016+00
23dfad82-fe83-4be9-8e64-548858d8f5a7	prod-agent	6.92	18.60	52.52	2025-11-28 15:47:14.200254+00
70628b70-bd40-4652-9e1a-44b86d31026a	prod-agent	8.20	18.61	52.52	2025-11-28 15:47:24.200503+00
bb2b074f-00a6-488f-881f-adf4ced7775c	prod-agent	14.65	18.64	52.52	2025-11-28 15:47:34.194666+00
6f047281-8e5c-4116-a5f6-b92efbead090	prod-agent	6.41	18.60	52.52	2025-11-28 15:47:44.202809+00
541bc656-5ca4-4b0b-b6c7-3f131cc9fa1f	prod-agent	8.53	18.63	52.52	2025-11-28 15:47:54.199797+00
89c81143-6b41-4780-aece-aac24c939344	prod-agent	9.33	18.67	52.52	2025-11-28 15:48:04.193431+00
fcb6f1e4-e761-4bf2-9cb9-92e3ff75a71c	prod-agent	7.76	18.63	52.52	2025-11-28 15:48:14.207863+00
6650340f-5b97-4971-b7c1-fc3295191a35	prod-agent	9.20	18.58	52.52	2025-11-28 15:48:24.193674+00
e824820d-d5bc-41f8-8446-79e64061f07b	prod-agent	13.36	18.66	52.52	2025-11-28 15:48:34.199371+00
b7c3febc-21a7-4e51-8037-424739bf025b	prod-agent	10.32	18.64	52.52	2025-11-28 15:48:44.193316+00
2c8f0b94-d5ba-4d6f-8eb9-0e0dba41ed19	prod-agent	8.05	18.68	52.52	2025-11-28 15:48:54.194413+00
da4653f6-5317-4e47-9ac4-277cb8fa3503	prod-agent	10.31	18.66	52.52	2025-11-28 15:49:04.199789+00
cfe7b199-ed2e-485c-998f-a4b95c800825	prod-agent	11.14	18.63	52.52	2025-11-28 15:49:14.200299+00
47c70ef0-f1ef-4440-82fe-3b4f3935183a	prod-agent	8.75	18.68	52.52	2025-11-28 15:49:24.195344+00
da528143-425e-4e83-b60d-6730d7f64ef7	prod-agent	10.43	18.70	52.52	2025-11-28 15:49:34.199919+00
e7693841-41c6-4c20-90d8-061c79375a4b	prod-agent	12.13	18.69	52.52	2025-11-28 15:49:44.194049+00
9e84a47a-49b6-419b-85ef-1ef10edc8d32	prod-agent	7.89	18.68	52.52	2025-11-28 15:49:54.200162+00
c805d9e0-5326-46b1-a371-cb0074533323	prod-agent	9.36	18.59	52.52	2025-11-28 15:50:04.202846+00
e43f0aa5-fdad-4a80-a0b3-2ec64856fb1b	prod-agent	7.86	18.66	52.52	2025-11-28 15:50:14.200288+00
8f6ce42d-12dd-4683-bd35-3232231dcb6b	prod-agent	9.04	18.57	52.52	2025-11-28 15:50:24.194977+00
31d5797c-2849-44da-b45a-0439c5a1f46a	prod-agent	10.40	18.55	52.52	2025-11-28 15:50:34.201497+00
0292f7df-c70f-4814-9e1f-a242ff13dc3b	prod-agent	12.34	18.61	52.52	2025-11-28 15:50:44.20002+00
c7c370e3-30ab-4ddb-9c07-4a9712b1238f	prod-agent	7.55	18.60	52.52	2025-11-28 15:50:54.204032+00
a2989764-8e29-4b5a-8d4d-b6cf25ebb9f6	prod-agent	9.90	18.70	52.52	2025-11-28 15:51:04.193337+00
78c4eb36-f333-48e5-a58d-c4eade33c753	prod-agent	7.71	18.55	52.52	2025-11-28 15:51:14.19932+00
2a239599-7555-4991-b529-c0697959b50d	prod-agent	8.22	18.58	52.52	2025-11-28 15:51:24.194484+00
fb73d117-69c8-4b67-bee4-d07b49563d4c	prod-agent	10.03	18.71	52.52	2025-11-28 15:51:34.204004+00
83df302b-8274-4e27-ae2f-f8b56fe7f2f9	prod-agent	12.26	18.64	52.52	2025-11-28 15:51:44.193871+00
c43c6daa-4633-4f5b-a04d-7f399acf5c2a	prod-agent	6.75	18.59	52.52	2025-11-28 15:51:54.200541+00
ec00d3a9-07c0-4fca-8a0a-52875efa2175	prod-agent	9.84	18.68	52.52	2025-11-28 15:52:04.203078+00
ff1f76e4-14a1-45ac-b47e-36ccf227b600	prod-agent	27.21	19.18	52.52	2025-11-28 15:52:14.204531+00
742bb952-dc29-41e7-8c85-13cb5966a4bb	prod-agent	37.50	19.14	52.52	2025-11-28 15:52:24.193946+00
b4635d23-92cb-477e-b618-04d1267a8ff3	prod-agent	36.04	19.40	52.52	2025-11-28 15:52:34.199529+00
63899783-797f-49a4-8b05-f154a163b2f0	prod-agent	40.47	19.67	52.52	2025-11-28 15:52:44.19293+00
27717e00-b8a0-4fef-8706-e9503f89bd19	prod-agent	34.21	19.50	52.52	2025-11-28 15:52:54.197614+00
4d09cd31-e2fd-4d5d-97d3-abb8dfc241a5	prod-agent	39.41	19.60	52.52	2025-11-28 15:53:04.204468+00
39c703ba-b3ff-4512-81bb-61dd66d537fa	prod-agent	36.06	19.81	52.52	2025-11-28 15:53:14.19625+00
c7325be0-1fe3-4a05-b221-a06db331fb9f	prod-agent	38.24	19.86	52.52	2025-11-28 15:53:24.199312+00
8ce23e9d-5fca-4ccb-b301-f1f9a092eeb4	prod-agent	38.21	19.88	52.52	2025-11-28 15:53:34.193594+00
00025363-2d38-4d5a-86d1-b0d50a0488f9	prod-agent	37.45	19.79	52.52	2025-11-28 15:53:44.211355+00
7ea37c76-b7c6-4e62-9552-7fa590360e87	prod-agent	39.87	19.73	52.52	2025-11-28 15:53:54.200669+00
04f544fc-b4e7-4ca7-81bd-90c37e400946	prod-agent	39.04	19.84	52.52	2025-11-28 15:54:04.208911+00
6520a3f1-6f9e-4eb2-96d7-1a0b99c2bb38	prod-agent	37.51	19.75	52.52	2025-11-28 15:54:14.199978+00
87d2acf1-d6c0-4fe1-8e6e-06066294258e	prod-agent	40.89	19.72	52.52	2025-11-28 15:54:24.195908+00
d3b1ffbc-714a-4038-b26c-3d83e04f22d6	prod-agent	37.95	19.81	52.52	2025-11-28 15:54:34.198859+00
dcb04241-18ae-4ae9-a43d-67940147667f	prod-agent	36.59	19.78	52.52	2025-11-28 15:54:44.201873+00
dc354220-9ce8-40d2-a534-38c35edfc195	prod-agent	34.14	19.71	52.52	2025-11-28 15:54:54.206495+00
1d78e4ef-6417-4c23-b293-6fc493de4bda	prod-agent	34.60	19.66	52.52	2025-11-28 15:55:04.193793+00
0b929352-56f2-4012-94f5-ef3fe6eec9e9	prod-agent	30.13	19.62	52.52	2025-11-28 15:55:14.200497+00
c825967a-68f5-430e-ab69-0cc48523a07d	prod-agent	35.26	19.57	52.52	2025-11-28 15:55:24.201207+00
5baa9cec-8c0e-4cb7-a00a-0933d86b839f	prod-agent	33.96	19.67	52.52	2025-11-28 15:55:34.195346+00
8f4536ce-7dc0-4e18-9bda-1bdbb5246096	prod-agent	29.99	19.62	52.52	2025-11-28 15:55:44.206328+00
7002d642-a484-4e2d-9bfe-1c130068c31d	prod-agent	36.58	19.61	52.52	2025-11-28 15:55:54.196889+00
9d837341-5989-45ca-9b2c-d2d32d2dd4d8	prod-agent	30.36	19.65	52.52	2025-11-28 15:56:04.192767+00
7baa3158-7630-4491-b6e1-e9860c6b7bb5	prod-agent	34.65	19.62	52.52	2025-11-28 15:56:14.196154+00
f7bc1cb3-b850-4f6d-9f3b-a5dbffe415c9	prod-agent	5.73	19.57	52.52	2025-11-28 15:56:24.197905+00
b084bf5b-918e-4e59-bbfe-959c16bcdf49	prod-agent	9.71	19.55	52.52	2025-11-28 15:56:34.197586+00
13a37834-154d-4772-84dd-44db2bd74bfe	prod-agent	9.86	19.53	52.52	2025-11-28 15:56:44.207326+00
d297c09f-cc19-41d8-be95-d8099224b25b	prod-agent	13.64	19.59	52.52	2025-11-28 15:56:54.194957+00
8915bc48-6796-46fd-95b5-b77a0245f423	prod-agent	9.35	19.50	52.52	2025-11-28 15:57:04.19437+00
6e3b73ee-0fbe-4dec-b47b-1fcef8a2f18f	prod-agent	9.23	19.54	52.52	2025-11-28 15:57:14.19883+00
1aab4483-47e4-4997-9aee-eb31d42d1ddc	prod-agent	11.28	19.48	52.52	2025-11-28 15:57:24.199575+00
dd811971-ecce-49b8-912a-8efbf8147e1b	prod-agent	9.10	19.48	52.52	2025-11-28 15:57:34.200399+00
2e18c8e0-191c-4e88-9d65-3a2d3b364abe	prod-agent	8.04	19.56	52.52	2025-11-28 15:57:44.202958+00
d7659e11-319f-48dd-ac03-9563905d4ce7	prod-agent	18.81	18.59	52.52	2025-11-28 15:57:54.201535+00
f44f459a-677a-486a-b7ac-d390f17303a7	prod-agent	8.12	18.47	52.52	2025-11-28 15:58:04.194764+00
60673203-4f32-483f-8d79-5dac6e044803	prod-agent	8.34	18.55	52.52	2025-11-28 15:58:14.200085+00
dade1d02-a981-40de-9205-71ac23a2adcc	prod-agent	7.91	18.56	52.52	2025-11-28 15:58:24.196328+00
03470f8d-bf92-463f-b1a7-275fa122ea63	prod-agent	8.91	18.48	52.52	2025-11-28 15:58:34.197014+00
482445b3-6c24-4830-89ca-2874caff087c	prod-agent	9.23	18.47	52.52	2025-11-28 15:58:44.201299+00
755ec47c-b32d-47e4-9f4f-656f0d4c0b1c	prod-agent	9.75	18.46	52.52	2025-11-28 15:58:54.215385+00
d9c04ac5-3f54-47c8-95c2-8a245bfe185f	prod-agent	12.93	18.56	52.52	2025-11-28 15:59:04.199971+00
3f631ec5-7a35-4ac1-9345-13242fca75e9	prod-agent	9.31	18.57	52.52	2025-11-28 15:59:14.201055+00
399723dc-6ec0-4b8f-ac5a-0ab4a77da155	prod-agent	8.75	18.55	52.52	2025-11-28 15:59:24.257091+00
72ca7df8-b524-4895-b23a-1cd70e69012f	prod-agent	9.44	18.48	52.52	2025-11-28 15:59:34.203658+00
f442d988-976c-4fe2-9cbb-52de72476e21	prod-agent	8.14	18.57	52.52	2025-11-28 15:59:44.194135+00
2ecf5832-957b-4818-a756-ab611adced9d	prod-agent	8.09	18.58	52.52	2025-11-28 15:59:54.200349+00
30b180c5-6196-4d95-8cb1-76acb83da90d	prod-agent	13.55	18.63	52.52	2025-11-28 16:00:04.199656+00
274ed8eb-c29a-4458-9495-ceb70d283e7a	prod-agent	10.31	18.59	52.52	2025-11-28 16:00:14.208938+00
2c1d2e8e-9d48-4ac8-98d8-f0900d55d0af	prod-agent	7.82	18.58	52.52	2025-11-28 16:00:24.194607+00
764bf145-25de-4210-90d8-b12b1a4ca5a7	prod-agent	30.27	19.03	52.62	2025-11-28 16:02:04.898662+00
bec2868b-6353-4eed-8dac-642da23d7328	prod-agent	13.34	18.85	52.62	2025-11-28 16:02:14.894608+00
89aaeb80-ebb9-48bb-b8e9-fc9dd4a5c31f	prod-agent	13.05	18.93	52.62	2025-11-28 16:02:24.890871+00
37c5afa0-4abb-47e2-93b3-2a566b7357ec	prod-agent	11.94	19.11	52.62	2025-11-28 16:02:34.893639+00
35f38352-aa83-4956-807f-8e0b017fd994	prod-agent	11.09	18.80	52.62	2025-11-28 16:02:44.887714+00
37daadae-cfdc-4337-9979-e18cb8bc9c65	prod-agent	12.93	18.78	52.62	2025-11-28 16:02:54.891285+00
8a78df2e-2678-4c2c-b7dc-757c5c936b4e	prod-agent	19.08	19.05	52.62	2025-11-28 16:03:04.891625+00
5ea3b8e4-7486-47cf-b1e7-c680c223af65	prod-agent	11.89	18.59	52.62	2025-11-28 16:03:14.88677+00
33f8b682-b613-45d3-ab20-d3d8c72826f9	prod-agent	8.10	18.66	52.62	2025-11-28 16:03:24.892615+00
ffd63e26-8a8a-44e1-afde-652dbe4f1859	prod-agent	8.56	18.74	52.62	2025-11-28 16:03:34.894336+00
c8460640-b481-4976-8577-7b2d2b7dac73	prod-agent	9.73	18.78	52.62	2025-11-28 16:03:44.892567+00
0086de33-5304-4be5-80b1-84a2013e4903	prod-agent	11.32	18.75	52.62	2025-11-28 16:03:54.892501+00
37e858da-44b5-4feb-bcc2-0b217d800506	prod-agent	21.71	19.10	52.62	2025-11-28 16:04:04.88875+00
cc459df2-a81f-45b7-b87e-17970cb014a2	prod-agent	14.67	18.69	52.62	2025-11-28 16:04:14.89246+00
ad10e24a-868a-44ac-8fdd-537a6dfc8f51	prod-agent	7.34	18.71	52.62	2025-11-28 16:04:24.89469+00
6ce7df79-7945-440c-9b73-73c6fe842fbc	prod-agent	8.68	18.62	52.62	2025-11-28 16:04:34.892632+00
48e2c100-b5d7-4d72-aed1-4dfa3351ef43	prod-agent	9.01	18.66	52.62	2025-11-28 16:04:44.896259+00
99191f6c-4784-41a5-b549-c285e53e9368	prod-agent	7.00	18.63	52.62	2025-11-28 16:04:54.893242+00
a5322428-527f-4709-b418-0087b028a183	prod-agent	9.79	18.64	52.62	2025-11-28 16:05:04.897791+00
67a4c038-8ae5-4723-9fcc-97be4db05e61	prod-agent	14.46	18.69	52.62	2025-11-28 16:05:14.891943+00
ca4d7151-4bef-43b1-9a74-cef17fb160d8	prod-agent	7.66	18.67	52.62	2025-11-28 16:05:24.892524+00
f85758ef-8716-4cf9-b196-29c5369da819	prod-agent	7.92	18.67	52.62	2025-11-28 16:05:34.892838+00
ff5c7e1f-27c9-4dfc-92eb-9cd928c66287	prod-agent	9.53	18.70	52.62	2025-11-28 16:05:44.892282+00
56ae7a27-7927-4c8e-ad8e-c4c2c8f2ec63	prod-agent	7.64	18.67	52.62	2025-11-28 16:05:54.890429+00
a8c90d4d-5a81-4587-a46c-d2e60ab47bd5	prod-agent	8.34	18.65	52.62	2025-11-28 16:06:04.897697+00
42485558-a96c-4f1c-8632-1493c36a635b	prod-agent	14.18	18.74	52.62	2025-11-28 16:06:14.895599+00
0ace7e50-a357-40b9-9166-11675e3e66d5	prod-agent	6.98	18.64	52.62	2025-11-28 16:06:24.892692+00
d8c92176-fcbb-47bc-bf9a-a7b234d029a0	prod-agent	8.60	18.67	52.62	2025-11-28 16:06:34.900558+00
5e61e6eb-0065-45fb-8571-2010330e0935	prod-agent	9.73	18.69	52.62	2025-11-28 16:06:44.892996+00
5caa4248-d87f-4bb5-be7e-04674f248c59	prod-agent	7.78	18.66	52.62	2025-11-28 16:06:54.886673+00
6387ff91-1c19-4936-bb53-8883d956cc7f	prod-agent	7.18	18.70	52.62	2025-11-28 16:07:04.95055+00
1e5d9693-68a8-4e90-81d5-f9f9d1af9dfa	prod-agent	16.27	18.71	52.62	2025-11-28 16:07:14.893489+00
37a34519-003f-4083-bd6d-25df47317479	prod-agent	6.57	18.69	52.62	2025-11-28 16:07:24.893278+00
b222a750-e8df-4861-9b4f-f481b806ec52	prod-agent	7.89	18.72	52.62	2025-11-28 16:07:34.892582+00
300386fd-bedc-4018-8bb5-5301e7c79911	prod-agent	10.10	18.68	52.62	2025-11-28 16:07:44.887037+00
8f942f75-f4ba-4c1f-bb3d-fdb1525635e3	prod-agent	7.33	18.73	52.62	2025-11-28 16:07:54.891025+00
a87b209b-d2d4-4d14-a830-839643b40849	prod-agent	8.61	18.61	52.62	2025-11-28 16:08:04.89042+00
1b9c2c3a-dee1-486d-9758-ba12784840e8	prod-agent	15.67	18.56	52.62	2025-11-28 16:08:14.895488+00
28c40a20-ecbf-4c14-8deb-24711f75e25e	prod-agent	7.23	18.55	52.62	2025-11-28 16:08:24.886293+00
c5210680-5d8f-4c80-8e6d-589402158a76	prod-agent	7.94	18.62	52.62	2025-11-28 16:08:34.891153+00
fa9cf406-bb1a-4768-84da-823caa4ee79d	prod-agent	9.44	18.62	52.62	2025-11-28 16:08:44.893056+00
3431283f-0df7-426a-83fd-36852757b9f4	prod-agent	7.68	18.58	52.62	2025-11-28 16:08:54.890461+00
5206301f-d86d-4d5b-ad3e-951e819fb1c7	prod-agent	8.12	18.68	52.62	2025-11-28 16:09:04.886447+00
f622f33d-ef93-4a84-ac06-0e16df058699	prod-agent	12.61	18.63	52.62	2025-11-28 16:09:14.890107+00
2569dd78-5f09-4877-997e-ae3124232b38	prod-agent	9.25	18.62	52.62	2025-11-28 16:09:24.893563+00
790d3ad8-d2de-4e83-857a-f4c7047525d6	prod-agent	8.14	18.56	52.62	2025-11-28 16:09:34.893878+00
f74b0318-8f02-4531-aa8e-d05563de94ef	prod-agent	10.33	18.63	52.62	2025-11-28 16:09:44.893874+00
afe3c4c8-14e4-4e17-a0a7-6273d9c4cce9	prod-agent	7.28	18.64	52.62	2025-11-28 16:09:54.892608+00
87737fe9-1758-486b-b482-867e3ee77405	prod-agent	8.80	18.59	52.62	2025-11-28 16:10:04.893248+00
a3a13b0e-bd94-41a0-8aee-476def95f998	prod-agent	9.54	18.55	52.62	2025-11-28 16:10:14.891646+00
9bd24a16-59f2-40e5-acfd-96cfaeba9e3b	prod-agent	12.48	18.65	52.62	2025-11-28 16:10:24.892303+00
8dcd9d65-c96e-4c7c-808c-3076ea48c73d	prod-agent	7.56	18.66	52.62	2025-11-28 16:10:34.890615+00
1c3d16cb-bcfb-4591-a92b-fe99122adf16	prod-agent	9.99	18.61	52.62	2025-11-28 16:10:44.893309+00
6c2bd1fd-273f-494e-ac37-9907b8515437	prod-agent	7.44	18.62	52.62	2025-11-28 16:10:54.891975+00
554f3ab5-b2f9-41cf-bcd0-8b7f8153d614	prod-agent	8.38	18.70	52.62	2025-11-28 16:11:04.886616+00
4d7163fb-0fe4-4c5e-b141-c04e39cf5c62	prod-agent	9.44	18.62	52.62	2025-11-28 16:11:14.891245+00
fd62c524-5474-4f18-99b5-f11a53351ed9	prod-agent	13.80	18.67	52.62	2025-11-28 16:11:24.892492+00
5175a207-22e6-4cab-9461-0844f09f647b	prod-agent	7.38	18.76	52.62	2025-11-28 16:11:34.893855+00
d7bfd876-4600-4767-8756-1682b03358a5	prod-agent	9.39	18.74	52.62	2025-11-28 16:11:44.894694+00
6d1fcb52-2f3f-47c4-96fe-2e59e44924fc	prod-agent	7.78	18.68	52.62	2025-11-28 16:11:54.890646+00
022c957a-2005-409a-a6a5-c31973fbb2d3	prod-agent	8.54	18.72	52.62	2025-11-28 16:12:04.894032+00
ecbdaab0-e067-42f1-a153-42cbdc954fc2	prod-agent	9.30	18.69	52.62	2025-11-28 16:12:14.893663+00
1ad4ca64-15fc-48c3-9104-c10751b692da	prod-agent	13.71	18.60	52.62	2025-11-28 16:12:24.889355+00
c009cc3e-717e-4b1d-8664-e6a9f59c1850	prod-agent	33.49	13.43	52.63	2025-11-28 16:57:54.739396+00
f8363126-aa84-46f2-a7b7-2ed73cf50948	prod-agent	17.77	12.97	52.63	2025-11-28 16:58:04.586254+00
18111722-c629-4fdb-8304-dd8c6cebea77	prod-agent	18.27	12.97	52.63	2025-11-28 16:58:14.588721+00
93b76c29-51d5-4dc4-8bc3-4790651d2031	prod-agent	17.05	13.02	52.63	2025-11-28 16:58:24.596424+00
a88bc472-ac85-4bba-86a8-4fc5212756ed	prod-agent	15.94	12.78	52.63	2025-11-28 16:58:34.583314+00
c638b719-3697-4656-bbfd-68e264721418	prod-agent	12.08	12.90	52.63	2025-11-28 16:58:44.589599+00
dcd4c77e-ce7f-4aaa-a158-7dc2bb0422f5	prod-agent	16.86	12.91	52.63	2025-11-28 16:58:54.591487+00
370740bb-4b9d-4b21-9481-10bd70039c7f	prod-agent	11.50	12.89	52.63	2025-11-28 16:59:04.585363+00
21435372-e8ca-4b48-97f5-5ca2263458b3	prod-agent	8.72	12.95	52.63	2025-11-28 16:59:14.589533+00
bc21f7ee-f02f-490f-bd32-5fa8c256405e	prod-agent	8.03	13.00	52.63	2025-11-28 16:59:24.588315+00
4c89a2c9-ca0d-4a5d-9d65-ce7d1a81c746	prod-agent	9.56	12.95	52.63	2025-11-28 16:59:34.588618+00
45117dfb-a53c-4281-98f6-ae3a276ad717	prod-agent	9.15	12.93	52.63	2025-11-28 16:59:44.58981+00
5bfe73f2-9713-4e35-872d-1a9f882408bc	prod-agent	7.90	12.99	52.63	2025-11-28 16:59:54.58826+00
ee15146b-116a-4fbb-a1f1-291ffbb1293f	prod-agent	16.53	13.33	52.63	2025-11-28 17:00:04.582736+00
e323c958-c706-4523-a229-23d98f05e8be	prod-agent	7.92	13.32	52.63	2025-11-28 17:00:14.594183+00
32c0916a-4ff9-4797-a55d-c35589a7e6ad	prod-agent	8.50	13.31	52.63	2025-11-28 17:00:24.584404+00
352f5158-6b9d-48f2-bbd6-e7682f90f965	prod-agent	10.57	13.37	52.63	2025-11-28 17:00:34.589924+00
f5fd20e5-2c48-41dc-8017-28cd88e60118	prod-agent	8.42	13.28	52.63	2025-11-28 17:00:44.588769+00
de1707ba-20f9-4e98-bac2-3e105c637496	prod-agent	9.35	13.38	52.63	2025-11-28 17:00:54.593128+00
668df009-c4da-4db9-92f5-f7ccc275328a	prod-agent	29.59	13.49	52.63	2025-11-28 17:01:04.588699+00
874959b9-2e0a-49d9-b168-b377ed5b2e93	prod-agent	5.75	13.42	52.63	2025-11-28 17:01:14.59047+00
879dc54e-4bd4-49ee-8ce3-3868fe2a0c6d	prod-agent	9.52	13.37	52.63	2025-11-28 17:01:24.583982+00
dd54b712-ea97-4650-9aa1-4807332a1d77	prod-agent	10.42	13.38	52.63	2025-11-28 17:01:34.584878+00
5f7e4ba2-4286-47d3-82a7-bcd84acf04ec	prod-agent	11.86	13.23	52.63	2025-11-28 17:01:44.590648+00
8b38759e-0a5f-4fb9-b3ba-1a602b318966	prod-agent	13.64	13.38	52.63	2025-11-28 17:01:54.584213+00
0c54b3b8-947c-43c2-adf1-9c1bd95b2eff	prod-agent	17.72	13.47	52.63	2025-11-28 17:02:04.592933+00
82d13ff8-6e86-4b52-94e9-6e40b1d6cb05	prod-agent	11.72	13.40	52.63	2025-11-28 17:02:14.590118+00
1b056f81-aa69-4012-833e-dd6c4058c5a8	prod-agent	10.16	13.35	52.63	2025-11-28 17:02:24.590023+00
d7e1ae74-9537-4cbc-b5f4-00daa534ea03	prod-agent	8.97	13.35	52.63	2025-11-28 17:02:34.588624+00
b7aac56f-8178-4ec0-ac75-bb0d5907273a	prod-agent	9.65	13.30	52.63	2025-11-28 17:02:44.590813+00
9c274fbd-85c8-4586-b534-a8d443041053	prod-agent	13.73	13.40	52.63	2025-11-28 17:02:54.590737+00
2d9ca067-8530-42e7-a1ff-3fd9ce3d54e5	prod-agent	18.14	13.43	52.63	2025-11-28 17:03:04.588876+00
6131786b-7089-4a56-9cca-d1c09ce1d0ae	prod-agent	7.56	13.34	52.63	2025-11-28 17:03:14.584751+00
cfcaa339-1703-46ef-a4e2-31b75d0f5867	prod-agent	9.74	13.39	52.63	2025-11-28 17:03:24.591326+00
e64669fe-3270-4c80-a9fb-4b6f8e285cbb	prod-agent	10.63	13.44	52.63	2025-11-28 17:03:34.596349+00
642d1164-d9d5-4e4c-9923-58303f63c77f	prod-agent	12.90	13.39	52.63	2025-11-28 17:03:44.588818+00
a29d2c0b-c707-4126-a5b8-2cfe9d06be4d	prod-agent	12.43	13.45	52.63	2025-11-28 17:03:54.595454+00
fa9da993-83ec-4072-8d5b-478123b6a73c	prod-agent	10.40	13.21	52.63	2025-11-28 17:04:04.594715+00
e0533c4e-fa04-4324-8aa1-1ce080f5b2ea	prod-agent	11.61	13.00	52.63	2025-11-28 17:04:14.589275+00
b3a9b9eb-0a46-47da-853e-2a827d59a1af	prod-agent	7.76	13.04	52.63	2025-11-28 17:04:24.593967+00
476c8024-4207-4e95-9d8d-bcf4466f2f4e	prod-agent	9.62	13.05	52.63	2025-11-28 17:04:34.590884+00
4b9c3905-8036-4a60-86e8-56e7777b62a5	prod-agent	9.56	13.06	52.63	2025-11-28 17:04:44.583607+00
7f57a1bd-94ba-4c8a-b7c4-5badf316b980	prod-agent	12.29	12.98	52.63	2025-11-28 17:04:54.589543+00
9025fe77-c91a-46fc-8ef0-1d93bb714de8	prod-agent	9.68	13.02	52.63	2025-11-28 17:05:04.585709+00
7b18e3da-7e59-4ab9-b2b0-dc850dc53592	prod-agent	16.41	13.14	52.63	2025-11-28 17:05:14.594518+00
46af0a6d-89b3-455b-9b0c-4487797c12b8	prod-agent	8.92	13.09	52.63	2025-11-28 17:05:24.589929+00
3eb01466-85b5-477f-9de6-e7291f90fcbe	prod-agent	9.12	13.01	52.63	2025-11-28 17:05:34.588561+00
de4b9e34-5521-4a22-99e3-157c6e3bf6c1	prod-agent	12.50	13.08	52.63	2025-11-28 17:05:44.584854+00
41192c59-cb5f-4dca-b188-a5e8f63e2eec	prod-agent	11.47	13.29	52.63	2025-11-28 17:05:54.589725+00
c1be345c-8245-40b1-ba61-d1d257c7d222	prod-agent	8.76	13.17	52.63	2025-11-28 17:06:04.593622+00
b5cda168-1f88-4a16-a447-603122279292	prod-agent	14.25	13.13	52.63	2025-11-28 17:06:14.589519+00
857a1135-b5ec-401a-b70d-2cbfccbeaa59	prod-agent	7.49	13.04	52.63	2025-11-28 17:06:24.589902+00
02406de7-c414-4f42-91cf-c174df59aba4	prod-agent	8.63	13.05	52.63	2025-11-28 17:06:34.590422+00
931e0363-7078-4e98-9c46-ab275f03b163	prod-agent	8.71	13.11	52.63	2025-11-28 17:06:44.587097+00
0218653b-9fa8-4ba0-944b-c38fef4b7079	prod-agent	7.77	13.07	52.63	2025-11-28 17:06:54.589722+00
a60bca7f-4710-474c-9b61-262dabbdd683	prod-agent	8.88	13.08	52.63	2025-11-28 17:07:04.589726+00
f9255bc5-9a45-43ac-86d6-eb868e691a47	prod-agent	18.30	13.10	52.63	2025-11-28 17:07:14.58961+00
4b614fdd-ff74-4e79-a431-cc29f34fbda3	prod-agent	10.49	13.00	52.63	2025-11-28 17:07:24.589471+00
c1fe9ece-b414-4ac0-8a4e-ac0153d0ffa2	prod-agent	12.06	13.02	52.63	2025-11-28 17:07:34.601163+00
4e65e8e9-0723-4c03-8c9c-3a211b31e58f	prod-agent	14.26	13.22	52.63	2025-11-28 17:07:44.590897+00
24aa4fc2-ccff-4aeb-aa27-48658b1785dd	prod-agent	30.84	12.88	52.63	2025-11-28 17:07:54.611669+00
ba4d0de7-5a8e-451e-a05b-0ea2a6983326	prod-agent	10.53	13.03	52.63	2025-11-28 17:08:04.589913+00
6a9cab29-0b61-499f-b050-11332e745787	prod-agent	30.92	14.93	52.63	2025-11-28 17:08:14.590157+00
e0a61c17-2701-48b8-980c-d7da1cbfa749	prod-agent	49.14	15.63	52.63	2025-11-28 17:08:24.589366+00
5b6fbfac-24ca-440b-ba5f-68bd3510c62f	prod-agent	51.66	15.48	52.63	2025-11-28 17:08:34.582387+00
41bd744a-66a0-4bf1-8206-9babc12520f2	prod-agent	27.26	15.52	52.63	2025-11-28 17:08:44.582404+00
33d961e3-1d92-4016-9476-f3e94e319ada	prod-agent	9.19	15.50	52.63	2025-11-28 17:08:54.590468+00
184a82e3-2641-4191-bb65-e43e1fd7f6e3	prod-agent	11.39	15.45	52.63	2025-11-28 17:09:04.588713+00
a7927e3e-ef9a-4542-8c70-15acddd0fb7b	prod-agent	13.36	15.39	52.63	2025-11-28 17:09:14.594176+00
e55dd331-cee2-46ab-8ce1-298e24bcd02c	prod-agent	13.92	15.10	52.63	2025-11-28 17:09:24.588379+00
e6245f5d-40d0-4867-965d-798d3c1bd93b	prod-agent	11.71	15.05	52.63	2025-11-28 17:09:34.589808+00
6f487666-2298-441f-96ec-44856311b2fb	prod-agent	13.39	15.09	52.63	2025-11-28 17:09:44.590588+00
67177c78-eb10-467d-9e5e-c07cb15adaca	prod-agent	11.81	15.03	52.63	2025-11-28 17:09:54.587608+00
cda3d3ad-54c8-40d3-a75a-ffcaf31c35ea	prod-agent	7.89	15.05	52.63	2025-11-28 17:10:04.593383+00
43badee2-151b-4c1d-8f8c-9b34b464506b	prod-agent	7.97	15.04	52.63	2025-11-28 17:10:14.583824+00
32056aa1-5954-4dff-a2a6-ae499cb57317	prod-agent	18.74	14.82	52.63	2025-11-28 17:10:24.588295+00
f745eb8a-2850-4395-9f21-45b83a40d07f	prod-agent	7.64	14.77	52.63	2025-11-28 17:10:34.595325+00
54d3ffc5-e2bf-4e63-b389-45ca5814c7f1	prod-agent	9.61	14.72	52.63	2025-11-28 17:10:44.589134+00
60b16bc3-8627-4356-842f-52023533f3de	prod-agent	9.07	14.82	52.63	2025-11-28 17:10:54.584602+00
137516db-241d-4c33-a53e-b63874bab7ed	prod-agent	10.37	14.75	52.63	2025-11-28 17:11:04.590388+00
6d1f8248-021b-4160-8217-a25e0b7ea950	prod-agent	8.33	14.62	52.63	2025-11-28 17:11:14.588911+00
f871257b-c09d-4699-8524-4c29571c9e4b	prod-agent	14.51	14.66	52.63	2025-11-28 17:11:24.587906+00
c097961d-c5a1-4ae9-9619-4f21147f6e9e	prod-agent	100.00	15.00	52.63	2025-11-28 17:11:46.137641+00
d8a48065-7177-4fd6-aa04-ea6c65124ffa	prod-agent	20.20	14.71	52.63	2025-11-28 17:11:56.102265+00
0121f23a-7cfd-4b6e-a192-756a81e73a69	prod-agent	8.63	14.73	52.63	2025-11-28 17:12:06.099496+00
d31f8637-a4d2-44dc-a077-05e471db4b92	prod-agent	12.07	14.84	52.63	2025-11-28 17:12:16.106559+00
4360a5d1-70ef-42d5-9904-f2e6223d0482	prod-agent	13.98	14.20	52.63	2025-11-28 17:12:26.097321+00
546adeee-7874-4636-98ed-f54706ca77ef	prod-agent	8.86	14.19	52.63	2025-11-28 17:12:36.108601+00
e880a23c-c256-4e79-96ad-379a7fa66111	prod-agent	7.71	14.21	52.63	2025-11-28 17:12:46.103769+00
b80d9045-306c-4993-9176-043fe1acb997	prod-agent	7.31	14.21	52.63	2025-11-28 17:12:56.107301+00
16682274-89e4-4bf8-af47-c6c64abdfaca	prod-agent	16.67	14.39	52.63	2025-11-28 17:13:06.1049+00
3ecaff9f-51ff-412c-81ca-b4ca142773cf	prod-agent	11.30	14.40	52.63	2025-11-28 17:13:16.105443+00
2a39ab52-acc0-45b4-b0a6-4ba7530b6eb8	prod-agent	25.91	14.67	52.63	2025-11-28 17:13:26.101969+00
9357076b-e06d-46f3-87a5-c17de8c0d6a8	prod-agent	7.26	14.68	52.63	2025-11-28 17:13:36.103794+00
c53fc2b4-ec4f-408b-ae50-a8ddb5f432e5	prod-agent	9.30	14.60	52.63	2025-11-28 17:13:46.1078+00
fd5298f2-51f4-4b24-a7c5-d37eb6804d1c	prod-agent	8.58	14.66	52.63	2025-11-28 17:13:56.099538+00
5996d0f3-4f8d-4286-a16a-d7ed519794c5	prod-agent	9.51	14.67	52.63	2025-11-28 17:14:06.099124+00
0a5e243b-c325-4891-957c-a79a92be51e5	prod-agent	24.18	14.96	52.63	2025-11-28 17:14:16.099126+00
b2683dda-713a-4d4c-abc1-3f477db53452	prod-agent	43.06	15.45	52.63	2025-11-28 17:14:26.103976+00
62f092e7-0c75-40cf-8a2f-9bb30f5ba79a	prod-agent	31.94	15.42	52.63	2025-11-28 17:14:36.107532+00
d8e53d8e-9b38-4baa-8014-8607633d83af	prod-agent	32.46	15.43	52.63	2025-11-28 17:14:46.107111+00
989e137f-e587-4abb-a833-ab87288520f9	prod-agent	40.55	15.58	52.63	2025-11-28 17:14:56.101881+00
f915a955-5724-4990-af43-475294456d95	prod-agent	36.39	15.57	52.63	2025-11-28 17:15:06.097927+00
0ac5002b-b111-4730-a39e-3861ced4fbf0	prod-agent	40.47	15.72	52.63	2025-11-28 17:15:16.09627+00
35c7eae6-822e-4578-92ff-1d7549aa7deb	prod-agent	33.78	15.64	52.63	2025-11-28 17:15:26.110218+00
6aac5781-d15b-4efc-8952-40a9de99dac6	prod-agent	42.82	15.84	52.63	2025-11-28 17:15:36.102444+00
28093195-0b15-4458-8cc1-77a75b01e913	prod-agent	37.79	15.82	52.63	2025-11-28 17:15:46.09825+00
65775b2d-b3c5-4e3d-9df9-72f2c16997ab	prod-agent	33.23	15.79	52.63	2025-11-28 17:15:56.099145+00
2a7d1d71-1675-4c12-adcd-10ac5dcd3827	prod-agent	51.67	15.94	52.63	2025-11-28 17:16:06.096206+00
e5ba1f1d-d2f9-453c-aa79-2f8fb01ca170	prod-agent	33.02	15.84	52.63	2025-11-28 17:16:16.101009+00
170e4f26-f447-48e6-80cf-688bd43bfa12	prod-agent	51.07	17.30	52.63	2025-11-28 17:16:26.103517+00
ab5a0b48-95ad-40e4-888d-b4a482ff4fa1	prod-agent	72.51	19.00	52.63	2025-11-28 17:16:36.106145+00
bd2c774c-1a41-4b7e-9a93-28a8fb863bb1	prod-agent	66.23	20.56	52.63	2025-11-28 17:16:46.147642+00
ad0297d9-1018-42df-b6e6-f6d0fa55cbf2	prod-agent	46.14	15.85	52.63	2025-11-28 17:16:56.098769+00
2eae5231-d9bf-4879-b61f-1d7c2e48284a	prod-agent	30.43	15.94	52.63	2025-11-28 17:17:06.103053+00
2fa684ac-8384-4b0d-8715-cf3100388c39	prod-agent	20.23	15.99	52.63	2025-11-28 17:17:16.106681+00
9027655b-88c4-4cd9-b3f4-fe9f734bc011	prod-agent	7.87	15.99	52.63	2025-11-28 17:17:26.104521+00
980fc2f4-017c-409d-8c2c-68a1164704ca	prod-agent	14.04	16.04	52.63	2025-11-28 17:17:36.103322+00
c5816926-c2fd-4528-a5f7-e4fc7dac9f94	prod-agent	8.18	15.96	52.63	2025-11-28 17:17:46.108138+00
47897930-5e0f-4cf9-8b81-f7929ec34ce1	prod-agent	8.95	16.05	52.63	2025-11-28 17:17:56.104125+00
e665d09d-8ec9-40a2-a89a-79a2712907af	prod-agent	8.52	16.02	52.63	2025-11-28 17:18:06.103405+00
973bf299-086f-48bf-ba44-577bed0470c2	prod-agent	7.65	16.00	52.63	2025-11-28 17:18:16.104689+00
8b4fb948-69ac-4e7b-93dd-fc36c56b4ac3	prod-agent	9.43	15.69	52.63	2025-11-28 17:18:26.100194+00
19c570f9-dab8-4353-b922-0acab36d7216	prod-agent	38.94	16.12	52.63	2025-11-28 17:18:36.102769+00
9d1f28de-7ef6-4efe-9b6d-338009781c8d	prod-agent	31.46	16.01	52.63	2025-11-28 17:18:46.095649+00
861dab88-798d-4d65-b9c0-e4915e8f5df4	prod-agent	55.20	17.94	52.63	2025-11-28 17:18:56.112843+00
1dd68cd0-374a-4e00-8623-caceda55ab80	prod-agent	69.35	19.86	52.63	2025-11-28 17:19:06.108318+00
0b6d1014-d6a9-4429-b60e-4ad9c8eefb63	prod-agent	64.22	16.06	52.63	2025-11-28 17:19:16.105141+00
50b0a68a-3be7-489f-8dec-df07710c4939	prod-agent	28.85	15.96	52.63	2025-11-28 17:19:26.098706+00
46a85567-bfee-467b-b1ea-4be6e7d5ba0d	prod-agent	26.37	16.17	52.63	2025-11-28 17:19:36.10014+00
27fc6c6c-86d4-4c90-bcee-a8014bb79654	prod-agent	6.21	16.17	52.63	2025-11-28 17:19:46.099257+00
9b4bee17-6d06-4ff4-9f5a-016c7c7570b4	prod-agent	8.99	16.10	52.63	2025-11-28 17:19:56.099365+00
1b7ed2e3-33d7-45cd-afc1-b429e86e6f12	prod-agent	10.55	16.13	52.63	2025-11-28 17:20:06.103801+00
0c2d2140-42fc-491b-8631-7adb98cb6938	prod-agent	36.16	16.22	52.63	2025-11-28 17:20:16.102786+00
a713bc3f-8514-44b4-bada-0698828bda79	prod-agent	37.35	16.15	52.63	2025-11-28 17:20:26.102141+00
36087432-1988-4df4-b807-2fc2a15c7df8	prod-agent	35.33	16.14	52.63	2025-11-28 17:20:36.114617+00
51548aa9-0bd3-4e3a-affa-7086f5668cd0	prod-agent	38.31	16.01	52.64	2025-11-28 17:20:46.099599+00
4f3fc019-3f88-4d84-aa1c-8a7b4f899bec	prod-agent	38.45	16.07	52.64	2025-11-28 17:20:56.10285+00
1e151489-54fc-41b5-aaf7-1e35de727c0d	prod-agent	34.77	16.02	52.64	2025-11-28 17:21:06.102485+00
91af5b08-b2a6-42e9-bed7-f65daee378fb	prod-agent	28.50	16.11	52.64	2025-11-28 17:21:16.097898+00
68b89c4f-2458-4c2f-bb7a-5fc49b884953	prod-agent	19.91	16.04	52.64	2025-11-28 17:21:26.111612+00
07e5d217-9b7e-4632-b232-f65a6d65fff6	prod-agent	8.18	16.05	52.64	2025-11-28 17:21:36.102276+00
79d6530c-2270-4dc1-b8de-4b5b4f9fb60d	prod-agent	15.57	16.08	52.64	2025-11-28 17:21:46.098889+00
345a9458-2f17-48d6-83ee-163dacef8f82	prod-agent	9.79	16.12	52.64	2025-11-28 17:21:56.157303+00
1a46c415-90aa-467f-b826-0cf8839a5db0	prod-agent	9.90	16.13	52.64	2025-11-28 17:22:06.098745+00
dd7be4c6-f1b3-490b-8dd0-52ea4dda7cc8	prod-agent	8.31	16.07	52.64	2025-11-28 17:22:16.103013+00
e7e96d7c-6566-4ff7-a117-c49d2dba4e7d	prod-agent	100.00	16.12	52.64	2025-11-28 17:22:43.381762+00
ac739c25-d8d1-4acf-be84-ea677983796e	prod-agent	22.07	15.60	52.64	2025-11-28 17:22:53.369848+00
d5d763fa-34fc-42dd-a61a-63f9c13d97dd	prod-agent	9.82	15.31	52.64	2025-11-28 17:23:03.369128+00
599bfeff-8cc0-49ff-8108-0cdcbd027ea5	prod-agent	8.74	15.25	52.64	2025-11-28 17:23:13.371145+00
1ccf816e-bb83-45d3-972a-d23c852f8ae5	prod-agent	11.73	15.30	52.64	2025-11-28 17:23:23.370771+00
a26323e9-dd82-495c-89e6-7b3e301b0723	prod-agent	8.65	15.26	52.64	2025-11-28 17:23:33.376016+00
d4b8bd6b-7588-400e-918f-32dc4bb10baa	prod-agent	11.24	15.37	52.64	2025-11-28 17:23:43.371488+00
9d2219db-8001-42f5-a48e-fdb1c780ea7e	prod-agent	10.49	15.37	52.64	2025-11-28 17:23:53.366499+00
090e1eb5-3ae3-4180-b94f-6c1bea84641d	prod-agent	7.77	15.37	52.64	2025-11-28 17:24:03.370908+00
2db9604d-a158-481d-9a0c-749c27662481	prod-agent	22.45	15.46	52.64	2025-11-28 17:24:13.370897+00
99541889-4e3d-40c5-9c17-b0638e4a0114	prod-agent	35.86	15.97	52.64	2025-11-28 17:24:23.370243+00
f34aa7aa-7cba-4e41-baf6-11308e1ace36	prod-agent	39.04	16.02	52.64	2025-11-28 17:24:33.366776+00
ca66d088-4462-4d6c-a202-d7cf0a884843	prod-agent	37.25	16.04	52.64	2025-11-28 17:24:43.368222+00
9d063c95-a1e7-49ad-9f5d-e3120b379dd9	prod-agent	41.57	16.12	52.64	2025-11-28 17:24:53.365733+00
87e21480-30a9-4e0c-a452-4aa26dcf1ad4	prod-agent	39.29	16.52	52.64	2025-11-28 17:25:03.370806+00
8226b0f1-141f-41d9-8ca4-56bb49e6d94b	prod-agent	64.89	18.33	52.64	2025-11-28 17:25:13.374152+00
af029252-210f-4d44-a5e2-3c1b7ab0d923	prod-agent	64.97	20.22	52.64	2025-11-28 17:25:23.375239+00
a38ed9e0-d507-41ec-a9d9-f5292d12b912	prod-agent	56.84	16.04	52.64	2025-11-28 17:25:33.362663+00
8fb0ce4d-111f-4a3c-bdcc-fc4101f0d451	prod-agent	28.09	16.09	52.64	2025-11-28 17:25:43.370053+00
0ab3a9b5-b748-44c2-ba43-62b15bf90208	prod-agent	33.54	16.16	52.64	2025-11-28 17:25:53.368791+00
af530a38-c4a5-467a-adec-cca87746361f	prod-agent	6.44	16.20	52.64	2025-11-28 17:26:03.371849+00
95e938e1-b7de-430f-82c2-e18ec2976d40	prod-agent	9.03	16.11	52.64	2025-11-28 17:26:13.370366+00
457d281f-dc43-43e5-b2f5-38593faddc98	prod-agent	9.30	16.12	52.64	2025-11-28 17:26:23.370839+00
bdf0f8d3-041e-46dd-8855-848d004fb980	prod-agent	7.60	16.18	52.64	2025-11-28 17:26:33.371797+00
496e7403-7c45-4001-8359-36c3d72d52e7	prod-agent	10.02	15.79	52.64	2025-11-28 17:26:43.370983+00
c098a414-1b04-42ca-a6fe-e70b45baf72d	prod-agent	25.50	16.03	52.64	2025-11-28 17:26:53.365015+00
7fd51588-3a4e-4204-9c33-ad2da59557cc	prod-agent	24.85	16.09	52.64	2025-11-28 17:27:03.376332+00
7ea065ee-013b-419e-8904-269b6895e689	prod-agent	33.74	16.16	52.64	2025-11-28 17:27:13.373474+00
3f209d6c-2cdb-402c-9665-92c702857d63	prod-agent	37.28	16.21	52.64	2025-11-28 17:27:23.366373+00
f114a708-a60c-4076-bc85-264b2e1ba461	prod-agent	29.35	16.17	52.64	2025-11-28 17:27:33.371939+00
488507b4-2c36-4bc9-bbd7-e336e521e85f	prod-agent	54.34	16.62	52.65	2025-11-28 17:27:43.366249+00
2517dbe9-7adc-4b3f-b1c8-b7612b65785f	prod-agent	46.60	17.04	52.70	2025-11-28 17:27:53.368946+00
678aa405-96be-4318-b218-bb520f18e9e6	prod-agent	43.90	16.77	52.75	2025-11-28 17:28:03.370722+00
f552e615-9716-4555-9096-2ca7c8b10cf6	prod-agent	52.26	17.93	52.75	2025-11-28 17:28:13.374812+00
0db8d3b9-1d9d-4c93-86a4-00df578b62cd	prod-agent	72.15	19.88	52.75	2025-11-28 17:28:23.374669+00
14997f2e-afc7-4768-ad10-8589ecb42e1e	prod-agent	69.16	21.04	52.75	2025-11-28 17:28:33.405618+00
c96a0170-2836-4550-9764-1ed28df3fb0a	prod-agent	40.49	16.63	52.75	2025-11-28 17:28:43.374217+00
71592f1b-4db3-494a-adbe-5defea80d3c9	prod-agent	44.47	16.64	52.75	2025-11-28 17:28:53.368592+00
b4c1d146-dd24-472f-b6a3-f01d49734c9d	prod-agent	26.55	16.73	52.75	2025-11-28 17:29:03.363823+00
cdbea652-a300-4408-b3fd-552491458ec2	prod-agent	27.18	16.70	52.75	2025-11-28 17:29:13.370876+00
208f2d62-b390-44dd-b45a-d15a78448ab5	prod-agent	20.46	16.75	52.75	2025-11-28 17:29:23.362778+00
b2b1dd74-e947-494b-abeb-2cb51a8a0f7a	prod-agent	7.43	16.71	52.75	2025-11-28 17:29:33.371488+00
c620b717-999a-40c7-8f60-db04946d3c99	prod-agent	11.97	16.67	52.75	2025-11-28 17:29:43.369782+00
f3c01e87-2d86-46b7-9912-c68d3b8e0ad7	prod-agent	8.25	16.63	52.75	2025-11-28 17:29:53.373775+00
b913dd8d-4566-4f68-b93b-75f48b5e0256	prod-agent	16.71	16.63	52.75	2025-11-28 17:30:03.363882+00
af957b5a-c46d-46ac-b844-1efbdb77a7e6	prod-agent	12.67	16.72	52.75	2025-11-28 17:30:13.369296+00
8f03ac7b-d850-418d-ae59-c9e12157c6a2	prod-agent	8.90	16.72	52.75	2025-11-28 17:30:23.37096+00
04828cb1-d221-41f5-be41-701ff5428756	prod-agent	7.34	16.72	52.75	2025-11-28 17:30:33.371053+00
1d8c5c35-0669-4443-b018-dda6f7bb0b39	prod-agent	9.28	16.78	52.75	2025-11-28 17:30:43.371122+00
6b88b845-d29b-417f-a521-c4885be3fcc3	prod-agent	13.28	16.04	52.75	2025-11-28 17:30:53.379889+00
d13e290e-9db8-4426-8b97-b6f56c6fa483	prod-agent	23.78	16.37	52.75	2025-11-28 17:31:03.370135+00
6b282bb7-b41b-4fa1-a527-e1030d17e8a1	prod-agent	11.99	15.96	52.75	2025-11-28 17:31:13.364025+00
11994958-9fbb-450c-9cbf-2c4b23b42c75	prod-agent	8.00	15.91	52.75	2025-11-28 17:31:23.36464+00
6c118e5e-a45e-4e20-8f6c-02c21199e7e0	prod-agent	7.49	15.88	52.75	2025-11-28 17:31:33.37035+00
f079330c-6f1b-405d-bd71-73da74bf2cc3	prod-agent	9.27	15.84	52.75	2025-11-28 17:31:43.365207+00
087d91c7-18b5-4059-b52c-cab57e14059f	prod-agent	12.60	15.38	52.75	2025-11-28 17:31:53.370463+00
98529dfb-d12d-40b6-bcce-e60c5ef34048	prod-agent	15.19	12.98	52.75	2025-11-28 17:32:03.365116+00
adc06c80-0d0b-40cb-925d-f70923e262ff	prod-agent	9.16	12.95	52.75	2025-11-28 17:32:13.370405+00
f6e1397b-ad42-4722-97f4-d50666a11832	prod-agent	8.38	12.99	52.75	2025-11-28 17:32:23.369142+00
8986208f-0afc-4ef4-97fb-f009fe07b789	prod-agent	50.10	16.08	52.77	2025-11-29 12:24:02.000932+00
e4221023-05e0-4f16-ae21-2b6dfdf52dd8	prod-agent	36.34	16.02	52.77	2025-11-29 12:24:11.993569+00
f7f90c36-0309-4a9c-8b5f-02c4efcffcb3	prod-agent	37.66	16.04	52.77	2025-11-29 12:24:21.985134+00
a348e4b3-129b-44c8-abb5-3f82595860eb	prod-agent	36.45	16.06	52.77	2025-11-29 12:24:31.990857+00
c59a1297-26e2-4a2d-8e32-d2ed8da31e31	prod-agent	35.76	16.09	52.77	2025-11-29 12:24:41.986124+00
61adc033-c3bc-47fa-af38-c990cbd76203	prod-agent	35.81	16.03	52.77	2025-11-29 12:24:51.99232+00
89c8071f-d430-4158-a4ef-d05b54370b28	prod-agent	39.82	15.95	52.77	2025-11-29 12:25:01.990864+00
127ddd31-de07-49b2-b5eb-40aa4492ec40	prod-agent	39.27	16.03	52.77	2025-11-29 12:25:11.987045+00
7576e317-6b51-437f-9013-524db889a5b1	prod-agent	35.97	15.89	52.77	2025-11-29 12:25:21.990981+00
d2fddda6-d63a-43fa-b39f-82fc9433d4c8	prod-agent	38.06	15.92	52.77	2025-11-29 12:25:31.992787+00
a12048d3-3055-4f9f-9b13-c73573714263	prod-agent	30.27	15.92	52.77	2025-11-29 12:25:41.98375+00
a185681a-13de-4fb3-9211-8d64b25db92e	prod-agent	29.53	15.95	52.77	2025-11-29 12:25:51.990758+00
ab3f7b29-6c0b-4c7a-9d1a-9657531b077e	prod-agent	8.04	15.99	52.77	2025-11-29 12:26:01.991925+00
c71627e7-ff28-4af5-943f-529b9f7630f2	prod-agent	12.20	15.96	52.77	2025-11-29 12:26:11.985295+00
09847eea-ca88-4ad9-b074-1da9c78179a4	prod-agent	7.33	15.93	52.77	2025-11-29 12:26:21.985354+00
d8724fa6-a831-4dfc-b2a4-d086b711b107	prod-agent	10.37	15.95	52.77	2025-11-29 12:26:31.993391+00
54c73845-1805-4cfb-bc30-03c8350614ed	prod-agent	7.93	15.92	52.77	2025-11-29 12:26:41.991757+00
80cf9653-b68f-4eeb-b1c8-39e87d969a32	prod-agent	8.56	15.95	52.77	2025-11-29 12:26:51.991896+00
45dfc1a1-dbc9-431c-b8c5-d87836102c44	prod-agent	9.51	16.00	52.77	2025-11-29 12:27:01.992176+00
6dfd1469-8639-4384-8156-3156740357c6	prod-agent	17.17	15.88	52.77	2025-11-29 12:27:11.991312+00
841f9b1a-2f17-483d-95fb-25378522e88d	prod-agent	12.03	15.33	52.77	2025-11-29 12:27:21.984679+00
5f446b0d-3a49-45f5-8382-b91502e6c5cd	prod-agent	14.13	15.09	52.77	2025-11-29 12:27:31.988672+00
57186cd7-7826-4816-8ee4-873be005b037	prod-agent	10.83	15.01	52.77	2025-11-29 12:27:41.997186+00
ea492cf2-844a-4cc3-9105-cfb27dd2a4af	prod-agent	11.32	15.05	52.77	2025-11-29 12:27:51.987522+00
6e6657a5-68ee-4202-a388-e03a63724110	prod-agent	12.99	15.11	52.77	2025-11-29 12:28:01.985344+00
c861a634-54c4-4243-a434-6cb11342a8f0	prod-agent	16.73	15.16	52.77	2025-11-29 12:28:11.990625+00
5950367b-fa15-4a5b-ba74-837dad36290c	prod-agent	11.30	15.08	52.77	2025-11-29 12:28:21.991852+00
b6e892d1-125e-45d3-ad49-e25797a3b3a9	prod-agent	14.91	15.15	52.77	2025-11-29 12:28:31.986402+00
697a9b2d-58eb-4aa7-8739-5c613e923af2	prod-agent	10.41	15.17	52.77	2025-11-29 12:28:41.997219+00
b17b0d6d-09fd-4e35-a2b9-cc54e0ec6538	prod-agent	12.03	15.13	52.77	2025-11-29 12:28:51.990857+00
8e4ec2cf-2c21-4e54-bb83-102c6819c690	prod-agent	10.67	15.09	52.77	2025-11-29 12:29:02.041006+00
51e02ebd-3fe7-4bda-8b57-0cd25b6d2997	prod-agent	13.79	15.19	52.77	2025-11-29 12:29:11.989946+00
4b00e52e-6caf-4aac-825a-a0eb2f3019ea	prod-agent	7.55	15.08	52.77	2025-11-29 12:29:21.98641+00
f5e04f39-a886-46a4-9dfa-3eb895558017	prod-agent	10.99	15.01	52.77	2025-11-29 12:29:31.990775+00
48e8b41b-47da-416a-bb10-e6cfaa684824	prod-agent	9.80	15.03	52.77	2025-11-29 12:29:41.984816+00
b077099f-f1b4-407f-9763-4d41b1f954c4	prod-agent	9.75	15.08	52.77	2025-11-29 12:29:51.990729+00
aadd86d4-f1f6-4a3e-8d0e-ea3d337a4525	prod-agent	15.23	15.05	52.77	2025-11-29 12:30:01.985354+00
8f229dc5-7b24-4dc7-8407-6441b869c33b	prod-agent	13.15	15.19	52.77	2025-11-29 12:30:11.991679+00
40167ae2-fd84-458e-96ed-9e133eb28455	prod-agent	7.77	15.06	52.77	2025-11-29 12:30:21.990995+00
22af5da0-593f-4d92-8840-9a9d64867aa4	prod-agent	9.45	15.10	52.77	2025-11-29 12:30:31.991432+00
bdc8ed17-4588-4ce2-9baa-c51695dc4765	prod-agent	7.32	15.06	52.77	2025-11-29 12:30:41.985319+00
2606c314-94dc-4786-be1a-b95434dc4742	prod-agent	8.89	14.65	52.77	2025-11-29 12:30:51.991916+00
2e97a019-d880-403c-a125-54429cc55f38	prod-agent	14.23	14.90	52.77	2025-11-29 12:31:01.990957+00
ae5adead-747b-4e4f-b568-aa324719e024	prod-agent	27.22	14.88	52.77	2025-11-29 12:31:11.985512+00
537dd982-18a5-4f7e-be56-bbc3b0e44210	prod-agent	11.72	14.65	52.77	2025-11-29 12:31:21.990973+00
fb8b03c9-b2d6-43c1-9459-5fdf044a6e60	prod-agent	8.48	14.75	52.77	2025-11-29 12:31:31.991607+00
dfb37b0d-7f1a-4ac3-9e62-b2f5333754bd	prod-agent	7.11	14.72	52.77	2025-11-29 12:31:41.985753+00
fdd4a9ea-b60e-4552-bf8c-ab5285a4a6f9	prod-agent	8.51	14.73	52.77	2025-11-29 12:31:51.995999+00
121a28cd-4e9b-46fc-a1fd-a0cc8b93d03c	prod-agent	16.73	15.07	52.77	2025-11-29 12:32:01.990651+00
81e0f371-6310-41a1-9b0b-50ce56c7a5af	prod-agent	45.97	15.57	52.77	2025-11-29 12:34:11.994744+00
6d839c5a-5927-4817-b4f2-f77cbaa1f649	prod-agent	38.58	15.68	52.77	2025-11-29 12:34:21.986411+00
eb6428bf-6a66-4e9b-9c67-985e180ccfba	prod-agent	33.36	15.61	52.77	2025-11-29 12:34:32.00055+00
4a8df8f5-7c4b-450e-85f3-e72e76f8194a	prod-agent	38.37	15.68	52.77	2025-11-29 12:34:41.990216+00
b9bace54-d552-428a-92f1-e3c52c0f6163	prod-agent	36.34	15.69	52.77	2025-11-29 12:34:51.991376+00
d93486a2-d332-4229-b1fe-b589a1d0db90	prod-agent	40.22	15.79	52.77	2025-11-29 12:35:01.992119+00
7b7af866-9f95-4b28-b6d2-7902a1fe87a2	prod-agent	39.29	15.74	52.77	2025-11-29 12:35:11.990277+00
483d6b18-12e6-4903-bf5c-ad05f78c5d02	prod-agent	40.85	15.83	52.77	2025-11-29 12:35:21.994143+00
3743057f-4ead-46f0-bd75-315f40baa345	prod-agent	40.39	15.76	52.77	2025-11-29 12:35:31.988106+00
aeed0cfa-2e8a-4ccd-9ea4-c339844d67c8	prod-agent	38.82	15.73	52.77	2025-11-29 12:35:41.995372+00
b854f4e4-a519-4397-993c-1444df3f6e7f	prod-agent	30.48	15.78	52.77	2025-11-29 12:35:51.990597+00
ba668552-a2fa-4b9d-8364-70badfdc4a15	prod-agent	36.65	15.77	52.77	2025-11-29 12:36:11.993002+00
ee5764ca-3d75-4c42-9f13-a9acfe0f8afd	prod-agent	35.85	15.91	52.77	2025-11-29 12:36:21.990191+00
2908c793-e443-414e-8ebd-4afccf760689	prod-agent	33.00	15.90	52.77	2025-11-29 12:36:31.983819+00
e82a7756-4261-4b75-94a1-e3233567e444	prod-agent	37.45	16.06	52.77	2025-11-29 12:36:41.98604+00
2ddbbd8a-2dc9-4dd2-9f19-b850e688d47b	prod-agent	79.51	16.36	52.74	2025-11-29 12:36:51.993338+00
a20663ec-3b71-48bd-b00e-f9809290451c	prod-agent	93.44	16.69	52.79	2025-11-29 12:37:01.990818+00
46326fe9-882b-41f9-ae53-19b0276fe2b5	prod-agent	67.03	17.00	52.87	2025-11-29 12:37:11.989768+00
e2b00012-c325-42b9-9399-beca061c45c9	prod-agent	18.95	17.14	52.87	2025-11-29 12:37:22.000766+00
9ce22cea-7440-479b-b80e-6c02320de7f8	prod-agent	77.93	17.23	52.96	2025-11-29 12:37:31.992207+00
1f34f960-4924-4164-af61-00728abccff3	prod-agent	74.34	17.17	53.12	2025-11-29 12:37:41.982967+00
b2a05e09-9fd3-46d1-b939-2bb9711515bf	prod-agent	36.09	17.20	53.14	2025-11-29 12:37:51.994686+00
19ab157d-900c-4294-b9c1-5e791c8d35f1	prod-agent	52.62	18.11	53.53	2025-11-29 12:38:01.985123+00
d06f7527-ab77-4a69-9855-097549fdec04	prod-agent	44.90	18.23	53.71	2025-11-29 12:38:11.993973+00
4e2a5005-a97d-4ff3-9dff-cc7dcc0a7726	prod-agent	29.55	17.85	53.74	2025-11-29 12:38:21.985304+00
c9127e1d-2c55-47c5-8649-f62a3be70c7b	prod-agent	34.53	17.61	53.74	2025-11-29 12:38:31.991653+00
b7a74d70-63a4-4bdf-9cd2-0146cbfd2712	prod-agent	46.75	18.03	53.75	2025-11-29 12:38:41.999359+00
fb0c2fe9-413e-48c3-90c5-43ac48dbda70	prod-agent	75.83	19.76	53.76	2025-11-29 12:38:51.994004+00
142abcee-958b-44b2-b961-198ec7932bb9	prod-agent	90.80	20.20	53.79	2025-11-29 12:39:02.000648+00
9a50e871-1f7f-42fc-a996-63f1793fb91c	prod-agent	99.82	19.92	53.82	2025-11-29 12:39:12.007453+00
258d215f-f860-4d42-bdfa-9fadb02f2996	prod-agent	99.67	20.28	53.85	2025-11-29 12:39:22.031184+00
d6ca5bf5-918f-42c9-9c18-4c1617571ef6	prod-agent	99.84	20.59	53.88	2025-11-29 12:39:31.99508+00
a9eca489-90e7-4fcc-a0cc-a34601031bca	prod-agent	88.02	21.61	53.90	2025-11-29 12:39:41.99259+00
348470c4-8931-47df-a6e8-ff94e6193ca3	prod-agent	99.92	20.72	53.93	2025-11-29 12:39:51.991931+00
f544ba84-0995-4186-b8ff-642fd81bf9e6	prod-agent	99.84	21.05	53.97	2025-11-29 12:40:01.993936+00
074a38a0-b8cb-4961-88fe-ffd777b26e16	prod-agent	99.82	21.03	54.00	2025-11-29 12:40:11.996613+00
f9b05692-e34c-42d1-9421-614acd610146	prod-agent	99.87	21.52	54.03	2025-11-29 12:40:21.99469+00
3b56a171-58a6-4652-9a71-33c3658356a3	prod-agent	99.82	22.74	54.05	2025-11-29 12:40:32.014504+00
2dffab39-94dd-4fad-b903-52546247429d	prod-agent	99.12	21.97	54.08	2025-11-29 12:40:41.998522+00
f0f785e3-8f29-4722-981d-45abdf925cf5	prod-agent	99.85	24.47	54.09	2025-11-29 12:40:52.026482+00
40d76267-da0d-49d3-b164-5086fd6f3fd7	prod-agent	99.89	25.63	54.11	2025-11-29 12:41:01.99378+00
365f6c91-d862-4131-b200-e04ccfdf7553	prod-agent	99.82	27.23	54.12	2025-11-29 12:41:12.079706+00
f20ec8c9-43e3-4307-bef5-0ca885e9bac0	prod-agent	99.84	28.39	54.14	2025-11-29 12:41:22.002803+00
47500791-9958-49a2-93ee-48f5cb7419a7	prod-agent	99.94	30.70	54.14	2025-11-29 12:41:32.000222+00
e90f3200-315a-4db5-9df3-8cd5a7e396ec	prod-agent	99.89	34.72	54.19	2025-11-29 12:41:41.99465+00
f30cea97-ced9-42cf-b484-4bf167ff543a	prod-agent	99.89	30.89	54.24	2025-11-29 12:41:51.998368+00
18df6605-a7be-41e0-b404-d1d5987603a7	prod-agent	99.85	30.90	54.20	2025-11-29 12:42:01.99662+00
6ff62473-140e-4b1b-8637-64418a3845f4	prod-agent	83.09	18.25	54.25	2025-11-29 12:42:11.991875+00
ab26270d-ec48-4bb2-8dee-5995f65889a8	prod-agent	68.61	20.12	54.30	2025-11-29 12:42:21.992541+00
efb687fe-354d-40b6-bdb5-52bda0528fee	prod-agent	39.84	17.82	54.10	2025-11-29 12:42:31.984971+00
9bb53dfa-e122-458d-8ac9-2214c698aa59	prod-agent	59.30	17.61	54.11	2025-11-29 12:42:41.991293+00
3b4a07e9-35a5-48bb-8dac-401fb2d8b88a	prod-agent	32.62	17.39	54.11	2025-11-29 12:42:51.986987+00
dbef3d44-a279-4bbc-b8c8-4df6efdab68c	prod-agent	34.45	17.40	54.11	2025-11-29 12:43:01.984576+00
04fc1e7f-0129-4aac-a0c5-e3afdf3e42a3	prod-agent	34.94	17.28	54.11	2025-11-29 12:43:11.983017+00
c59c6e7d-7d31-490f-bd40-ba330ffc8ff6	prod-agent	29.79	17.35	54.11	2025-11-29 12:43:21.990407+00
93364964-265c-4eac-9e3c-e526eff03dc5	prod-agent	31.81	17.31	54.11	2025-11-29 12:43:31.989684+00
b092ce74-da52-4457-9f3d-3a90c07c0e08	prod-agent	34.10	17.37	54.11	2025-11-29 12:43:41.991548+00
3251f98e-0e75-427f-8c70-b2a507fe0903	prod-agent	32.21	17.32	54.11	2025-11-29 12:43:51.983799+00
cc2d69ce-ebab-456e-8f50-244a2cb765ec	prod-agent	32.06	17.23	54.11	2025-11-29 12:44:01.993348+00
43076c2a-0036-4e0f-92e6-8851bf364c5b	prod-agent	32.70	17.19	54.11	2025-11-29 12:44:11.983949+00
d41ecf7b-b9c3-4832-a4ef-684d8deaa263	prod-agent	30.47	17.21	54.11	2025-11-29 12:44:21.99115+00
0ccbd5ce-a47c-459e-a967-1ec0f5e1a9f5	prod-agent	30.98	17.16	54.11	2025-11-29 12:44:31.984496+00
77a936af-fbff-4821-a0a2-3f78a5d5db47	prod-agent	31.84	17.19	54.11	2025-11-29 12:44:41.990778+00
77c1678d-7e81-4e2e-92b0-c4f2a11989d0	prod-agent	35.28	17.06	54.11	2025-11-29 12:45:01.992047+00
4c8b4f19-cebd-4896-acc6-687e5cd3bc1f	prod-agent	33.75	17.07	54.11	2025-11-29 12:45:11.992553+00
1c0e6bed-ff55-4ee6-bbd0-04e6fdee81cf	prod-agent	37.28	17.04	54.11	2025-11-29 12:45:21.982971+00
a861bada-8cec-40c6-9847-0704dfcff859	prod-agent	37.65	17.15	54.11	2025-11-29 12:45:31.988448+00
41db0cb7-6fff-40bb-a00e-c68478e68885	prod-agent	34.24	17.06	54.11	2025-11-29 12:45:41.990919+00
f8fbc9d1-6f86-4751-97f1-fccf2b0c463e	prod-agent	40.21	17.13	54.11	2025-11-29 12:45:51.985894+00
984868b0-86e3-4732-b6b5-6bfb92b75736	prod-agent	36.76	17.15	54.11	2025-11-29 12:46:01.990572+00
d5801c5c-fac6-4ad7-83f2-02f6daf8e2dc	prod-agent	36.57	17.21	54.11	2025-11-29 12:46:11.983954+00
6c4b2ec7-99a1-444a-8e01-0b8f16fe731f	prod-agent	35.72	17.21	54.11	2025-11-29 12:46:21.984454+00
f6a80c8e-74d0-4722-bfb7-125dae35f4b3	prod-agent	36.32	17.18	54.11	2025-11-29 12:46:31.991365+00
8e130b95-931e-42eb-9bf5-493e619a75e5	prod-agent	35.32	17.06	54.11	2025-11-29 12:46:41.988804+00
546d2d46-b95a-4441-93e0-c239068f499b	prod-agent	39.26	17.16	54.11	2025-11-29 12:46:51.989342+00
60321c9e-b2e7-42bd-8780-569893b685ea	prod-agent	36.91	17.12	54.11	2025-11-29 12:47:01.989586+00
ce77c373-a3c5-4036-a0d9-baef1a789279	prod-agent	35.90	17.20	54.11	2025-11-29 12:47:11.995313+00
7354ca14-4afb-4714-9146-8116d94ce052	prod-agent	34.50	17.11	54.11	2025-11-29 12:47:21.991123+00
a9e0275a-8334-42d0-8a4a-c59ab6b9fef4	prod-agent	65.51	17.43	54.11	2025-11-29 12:47:31.991244+00
cbae17a5-9b8e-45d6-afd5-6c545b2b3bb1	prod-agent	55.66	17.25	54.11	2025-11-29 12:47:41.99025+00
89f9a95e-1589-4640-9c57-bfc14243f3ae	prod-agent	32.71	17.31	54.11	2025-11-29 12:47:51.985485+00
c9bbe162-aef4-488a-a256-0e64ebf97886	prod-agent	34.87	17.35	54.11	2025-11-29 12:48:01.98354+00
0f027098-3e75-4436-aa47-ecf43fdb9d26	prod-agent	37.22	17.32	54.11	2025-11-29 12:48:11.997866+00
e96f15b1-0834-4546-b199-1394b0b9fcbb	prod-agent	34.30	17.35	54.11	2025-11-29 12:48:21.986486+00
c10f60f8-2afa-42bf-921c-5c96634395b6	prod-agent	47.41	17.49	54.11	2025-11-29 12:48:31.990938+00
6674ce45-3303-437f-8dbf-da5a9dde73d9	prod-agent	36.45	17.34	54.11	2025-11-29 12:48:41.989576+00
c741f5e4-9813-43a9-816f-4c0ab8936a2e	prod-agent	39.10	17.40	54.11	2025-11-29 12:48:51.988369+00
64886fc2-00c4-4085-ae49-650f25610876	prod-agent	36.46	17.46	54.11	2025-11-29 12:49:02.000829+00
d2d1053a-911c-41fc-bee7-75a4199adfad	prod-agent	36.18	17.37	54.11	2025-11-29 12:49:11.988924+00
d6123a54-1e7e-488a-876c-1dbc64e1ad53	prod-agent	32.39	17.32	54.11	2025-11-29 12:49:21.983539+00
67c7216b-adbf-4105-9749-7fff36fb0a9d	prod-agent	28.90	17.33	54.11	2025-11-29 12:49:31.989813+00
adf9dd00-1fd6-49bc-9e20-6269bec15c49	prod-agent	29.30	17.47	54.11	2025-11-29 12:49:41.982839+00
34508fd6-b5ca-4e39-8395-e62efb1d95d0	prod-agent	9.78	17.56	54.11	2025-11-29 12:49:52.00099+00
245d1cb5-2343-47cc-b962-d9fc4483397b	prod-agent	30.05	17.07	54.11	2025-11-29 12:50:02.024139+00
61e4ab10-dfe4-4e1f-8c94-aeefc3c6de77	prod-agent	7.15	17.08	54.11	2025-11-29 12:50:11.985969+00
058a27da-5fc8-48fc-8f72-5fd002485152	prod-agent	12.57	17.14	54.11	2025-11-29 12:50:21.984998+00
ebfc219d-fdaa-43e5-bac1-dc4a4a647b95	prod-agent	9.09	17.15	54.11	2025-11-29 12:50:31.990504+00
09319a48-6b02-4cdb-9a56-d237f89c3c95	prod-agent	9.28	17.07	54.11	2025-11-29 12:50:41.993772+00
8b2086d8-a2ef-4c49-bfbe-5e1e984ffbfe	prod-agent	9.64	17.09	54.11	2025-11-29 12:50:51.986986+00
e1478d24-125a-4b2d-9e4c-5284640ae3ee	prod-agent	14.38	17.18	54.11	2025-11-29 12:51:01.984822+00
1ffd70a7-c40a-4337-91ed-22c558410d34	prod-agent	15.68	16.57	54.11	2025-11-29 12:51:11.99115+00
fa178efd-6cbc-47d0-aaff-2d1b596255cc	prod-agent	8.80	16.50	54.11	2025-11-29 12:51:21.985369+00
e230458d-7480-459a-8006-7a809e31e93a	prod-agent	8.52	16.50	54.11	2025-11-29 12:51:31.991459+00
1d877f5c-d28c-4abb-8b0f-7115ad747262	prod-agent	9.73	16.48	54.11	2025-11-29 12:51:41.985885+00
8a3008c9-b002-4f9b-a869-53804970b9e7	prod-agent	10.29	16.62	54.11	2025-11-29 12:51:51.986932+00
b703458a-91e0-4410-8461-7475a1b9e98c	prod-agent	14.81	16.42	54.11	2025-11-29 12:52:01.994816+00
4b5e140a-49fe-4d6f-9acb-a6716db66328	prod-agent	9.19	16.38	54.11	2025-11-29 12:52:11.991133+00
eda9d857-41e5-41f3-9d5e-1af17b3269af	prod-agent	8.69	16.33	54.11	2025-11-29 12:52:21.985172+00
540e4e14-0bc0-4437-939c-7ccb73a1277b	prod-agent	9.39	16.34	54.11	2025-11-29 12:52:31.988972+00
879dd90d-36ea-431f-b033-0966b5740653	prod-agent	9.04	16.41	54.11	2025-11-29 12:52:41.985874+00
c9d39030-748e-4a11-99b4-0d56a6734628	prod-agent	8.55	16.33	54.11	2025-11-29 12:52:51.992408+00
f01a6a7c-76eb-4c66-8ef5-3fb42fffcd52	prod-agent	16.16	16.48	54.11	2025-11-29 12:53:01.986669+00
765c3297-994c-4281-8fe1-6a2c8d050c24	prod-agent	11.97	16.45	54.11	2025-11-29 12:53:11.994342+00
51420d72-ca17-4b67-a165-d8a6a9883a39	prod-agent	12.29	16.43	54.11	2025-11-29 12:53:21.987275+00
7927e4fc-ca14-4dfb-8e03-c78ed49f4230	prod-agent	13.14	16.34	54.11	2025-11-29 12:53:31.991296+00
041a83b1-bf75-4f62-a5d6-c769fb4af87a	prod-agent	11.00	16.44	54.11	2025-11-29 12:53:41.985595+00
6f1d2fd0-66d3-481c-85e3-9447e8f53720	prod-agent	11.96	16.59	54.11	2025-11-29 12:53:51.985175+00
87425be6-a8bc-4e31-978f-d6ca37beeb21	prod-agent	17.64	16.56	54.11	2025-11-29 12:54:01.984926+00
09e8f6c8-f158-4015-a7c7-fe8614485136	prod-agent	15.58	16.53	54.11	2025-11-29 12:54:11.989506+00
d0c2c356-f8c7-46dd-a58e-a92dc5ac1981	prod-agent	27.92	16.85	54.11	2025-11-29 12:54:21.985752+00
7f061430-53bc-4c32-9116-f7d3be303a31	prod-agent	32.08	16.97	54.11	2025-11-29 12:54:31.991354+00
aca14c7e-204c-4314-b0ec-8f9631d66eab	prod-agent	37.53	17.07	54.11	2025-11-29 12:54:41.986263+00
fd2e3fdd-419a-48e5-a293-5f0d0bfb8030	prod-agent	37.82	17.19	54.11	2025-11-29 12:54:52.001239+00
80425660-a5e4-4a0a-be03-74c24c54b310	prod-agent	42.47	17.35	54.11	2025-11-29 12:55:01.991242+00
2dc69dbb-e080-43e7-9799-ffe3270412e8	prod-agent	39.70	17.20	54.11	2025-11-29 12:55:11.995501+00
488cbdab-b83c-439b-bcc8-18882c19129e	prod-agent	36.87	17.26	54.11	2025-11-29 12:55:21.992941+00
c69977e3-9e18-4a3e-afa1-02732fd9135a	prod-agent	36.19	17.35	54.11	2025-11-29 12:55:31.99334+00
7cec1dbd-cf3f-47d1-b736-8b04aa1fd02f	prod-agent	37.80	17.43	54.11	2025-11-29 12:55:41.995262+00
972d4a2d-5161-43dd-9579-ebc9550a0f4c	prod-agent	33.43	17.37	54.11	2025-11-29 12:55:51.989628+00
2a9dd768-6afe-45b3-95c4-c79490ecd5a4	prod-agent	41.45	17.38	54.11	2025-11-29 12:56:01.983417+00
d76fd6c2-bd57-4f6a-806b-59c7d7eb4791	prod-agent	39.48	17.36	54.11	2025-11-29 12:56:11.988234+00
5fff9e6f-3ddc-427f-be39-7d817b8a749b	prod-agent	19.68	17.27	54.11	2025-11-29 12:56:21.992914+00
9116678c-2f2a-4d1d-8529-3d14cb3b3f56	prod-agent	24.15	17.19	54.11	2025-11-29 12:56:31.989148+00
cd24359f-1cf1-4489-8ab1-76487d14ebc1	prod-agent	23.32	17.10	54.11	2025-11-29 12:56:41.991249+00
57512088-836b-4b70-815d-422dd66729be	prod-agent	29.33	17.15	54.11	2025-11-29 12:56:51.993443+00
acf8bf3e-1271-4ea9-a685-4558f07e0c99	prod-agent	38.31	17.06	54.11	2025-11-29 12:57:01.983342+00
d7aa8b3c-9d36-4e66-bfa3-d40ccdf4ae09	prod-agent	40.25	17.46	54.11	2025-11-29 12:57:11.990104+00
68b27080-4eaf-4f71-b38b-2d020bdff64a	prod-agent	36.90	17.31	54.11	2025-11-29 12:57:21.984377+00
42f05cf7-84b9-4fa2-a79c-d7a899d67664	prod-agent	38.09	17.35	54.11	2025-11-29 12:57:31.984463+00
8c4f1c7f-3964-4693-8dd0-142af666c8cb	prod-agent	35.72	17.40	54.11	2025-11-29 12:57:41.987921+00
6009048f-c6a1-4344-ad3f-48f4e06f6270	prod-agent	33.36	17.25	54.11	2025-11-29 12:57:51.985588+00
9d576294-b9e8-4b5b-be20-fb76382e456d	prod-agent	42.65	17.45	54.11	2025-11-29 12:58:01.993255+00
c4ba1f3c-b673-469d-a3c2-256aacd07ff8	prod-agent	36.21	17.55	54.11	2025-11-29 12:58:11.990365+00
91a9795d-b15b-45f0-b26f-7996a723cb9e	prod-agent	33.94	17.52	54.11	2025-11-29 12:58:21.992816+00
adadb79a-d9b7-410c-9324-670eaddc2d5f	prod-agent	35.11	17.53	54.11	2025-11-29 12:58:31.993725+00
78836053-2af2-4a27-8035-31399b768856	prod-agent	32.10	17.49	54.11	2025-11-29 12:58:41.986432+00
62b5b1c6-3210-4029-9355-c4cf5a3c4ddb	prod-agent	33.75	17.51	54.11	2025-11-29 12:58:51.984405+00
00f9ec4b-99d9-40e8-a016-94175a504e59	prod-agent	33.55	17.54	54.11	2025-11-29 12:59:01.986395+00
000f10f5-856c-4dda-bb14-64999a986e8e	prod-agent	39.58	17.40	54.11	2025-11-29 12:59:11.989264+00
66a6fa4f-0cd2-4917-94fc-3575474f0a5b	prod-agent	33.37	17.52	54.11	2025-11-29 12:59:21.990638+00
5c4c7735-51bd-4e77-9e5a-73a70900d809	prod-agent	36.22	17.44	54.11	2025-11-29 12:59:31.991593+00
ee598847-55cd-4ed9-b116-534fa7da2f46	prod-agent	35.07	17.36	54.11	2025-11-29 12:59:41.990349+00
5d8ff31a-ed47-41a7-9750-335e37c72ef6	prod-agent	43.91	17.44	54.11	2025-11-29 12:59:51.990047+00
a5751332-12de-4878-82e1-06bfd6c8014f	prod-agent	37.23	17.50	54.11	2025-11-29 13:00:01.983739+00
a0f8045e-19e0-44f1-a246-7e8cc13cf5cf	prod-agent	37.51	17.63	54.11	2025-11-29 13:00:12.020615+00
cd43c3bb-6e9d-4fc7-8767-9725dc854887	prod-agent	36.57	17.47	54.11	2025-11-29 13:00:21.984359+00
5128b1a5-e9be-4d6c-83c2-11b188caadab	prod-agent	35.98	17.51	54.11	2025-11-29 13:00:31.989801+00
bc1db1d5-d233-4699-8c70-213a4e0c1b7c	prod-agent	42.05	17.49	54.11	2025-11-29 13:00:41.985136+00
fb5410b1-bc1d-4338-919e-8ad6e7e5cba2	prod-agent	38.32	17.59	54.11	2025-11-29 13:00:51.988355+00
65bf8813-6e6a-4352-9816-83d7dfbce866	prod-agent	35.88	17.54	54.11	2025-11-29 13:01:01.987374+00
9af5e76e-dd7c-4184-9183-2a79805c6093	prod-agent	41.46	17.89	56.22	2025-11-29 13:17:51.992288+00
6d72da29-bfdb-4911-b82b-f68126c26db8	prod-agent	32.81	17.81	56.22	2025-11-29 13:18:01.994093+00
a19b5137-6baf-417d-8a56-7ed98223deaa	prod-agent	33.94	17.80	56.22	2025-11-29 13:18:11.992783+00
f4f69981-75a3-424b-b33b-11f39ae1c8fc	prod-agent	35.47	17.82	56.22	2025-11-29 13:18:21.990656+00
8046a530-0786-4f2c-b106-2b57a5011d5d	prod-agent	43.49	17.99	56.22	2025-11-29 13:18:31.988585+00
73f258e9-d5e2-4ee0-957c-8060a2588349	prod-agent	37.20	17.79	56.22	2025-11-29 13:18:41.984237+00
656a2779-0c9f-47f3-86bb-f2c4c095b3b7	prod-agent	34.77	17.95	56.22	2025-11-29 13:18:51.989435+00
8edf9e87-1a57-4c13-9705-6e76cd04b0e0	prod-agent	43.39	18.03	56.22	2025-11-29 13:19:01.99248+00
d93bab96-d5de-499e-9196-c9419d04f94e	prod-agent	32.64	17.95	56.23	2025-11-29 13:19:11.986313+00
e6079f1c-811c-4d48-aacc-5fd5e010165a	prod-agent	37.00	17.94	56.23	2025-11-29 13:19:21.98561+00
3e5607fd-d93b-44f0-a523-bcf450cab8bd	prod-agent	40.74	17.87	56.23	2025-11-29 13:19:31.990627+00
c3d8c7a6-454b-4e02-b8f5-e9556755ccc1	prod-agent	36.22	17.93	56.23	2025-11-29 13:19:41.984697+00
f56733c1-a77d-4618-94c3-fae3f09f0800	prod-agent	34.90	18.00	56.23	2025-11-29 13:19:51.999427+00
d04d10d3-fc1e-469d-8154-4fdb004d9cb0	prod-agent	49.49	18.09	56.23	2025-11-29 13:20:01.989866+00
35fa5160-5fd6-43df-b021-9df86efbd6e0	prod-agent	38.87	18.05	56.23	2025-11-29 13:20:11.991886+00
5764b093-7118-4e66-9dc1-619ddc1b9370	prod-agent	35.88	18.00	56.23	2025-11-29 13:20:21.990294+00
91b95b91-c86f-4b1e-aaea-46378c56107d	prod-agent	36.69	17.97	56.23	2025-11-29 13:20:31.990108+00
aee7a51d-51d2-44ba-9923-c18841d37e53	prod-agent	37.26	18.00	56.23	2025-11-29 13:20:41.985509+00
29bea3fe-3488-4ca6-8104-f35f40d77902	prod-agent	35.65	18.01	56.23	2025-11-29 13:20:51.990408+00
e7ccad37-e6fc-4e5b-9b77-8de7f2c1d494	prod-agent	36.29	17.96	56.23	2025-11-29 13:21:01.996042+00
f0951f31-316b-4a5a-b16c-f0ffff98177c	prod-agent	37.36	17.99	56.23	2025-11-29 13:21:11.983786+00
6b3d754e-2b2d-4e23-9803-e6b64137823e	prod-agent	38.75	18.01	56.23	2025-11-29 13:21:21.991758+00
4ea81777-8e4f-4feb-b4c4-cd06d3a1103e	prod-agent	38.69	17.93	56.23	2025-11-29 13:21:31.995696+00
06128a14-644e-4eaa-9222-c0235e4bc17b	prod-agent	23.76	17.99	56.23	2025-11-29 13:21:41.995599+00
cb1cdb1b-c529-4017-9480-9754f3ae73d7	prod-agent	21.12	18.10	56.23	2025-11-29 13:21:51.990441+00
2a069bdc-a3a5-465e-8b92-8dae93416448	prod-agent	12.25	18.00	56.23	2025-11-29 13:22:01.991826+00
a169769d-bcd3-4ed9-95d7-c942f6cb2207	prod-agent	17.49	17.86	56.23	2025-11-29 13:22:11.990469+00
44468d61-af68-489f-b523-0262941ca80e	prod-agent	8.94	17.89	56.23	2025-11-29 13:22:21.984341+00
6905c662-c694-4b3b-90bc-d5eb9292f223	prod-agent	7.48	17.87	56.23	2025-11-29 13:22:31.991297+00
5fe98999-1656-47e3-ab9a-e975c0f54dbe	prod-agent	8.54	17.85	56.23	2025-11-29 13:22:41.986774+00
41202b7c-8e6f-47fd-95df-afb3390951f0	prod-agent	8.58	17.89	56.23	2025-11-29 13:22:52.040182+00
a830cb29-1e96-4c85-88d3-e0bb8c751ad2	prod-agent	9.69	17.82	56.23	2025-11-29 13:23:01.991355+00
16f744da-4520-467e-9a8d-d4ffeb8070a4	prod-agent	14.83	17.88	56.23	2025-11-29 13:23:11.99493+00
964385b9-3112-4aab-afdb-03b317f12e3b	prod-agent	8.74	17.82	56.23	2025-11-29 13:23:21.988311+00
eee2a1bb-0b8f-4a4e-8d73-870e8759327e	prod-agent	14.82	16.95	56.23	2025-11-29 13:23:31.993578+00
3c18e914-1d7a-4b6f-9269-3d26cd54c068	prod-agent	8.30	16.81	56.23	2025-11-29 13:23:41.987246+00
f33d7fd9-616b-4ea7-939f-a9fbc1907e4e	prod-agent	9.93	16.67	56.23	2025-11-29 13:23:52.038202+00
af2c8cce-e10f-4e5b-8638-4d187aa0c880	prod-agent	7.68	16.70	56.23	2025-11-29 13:24:01.992859+00
50b051d3-2155-459e-9495-17bbcc8c9528	prod-agent	14.13	16.81	56.23	2025-11-29 13:24:11.990489+00
59247c8e-b5ec-4ee7-92c3-f87b238bad00	prod-agent	8.85	16.70	56.23	2025-11-29 13:24:21.992329+00
1d2bd763-cfcb-49e1-81dc-02a95a36b184	prod-agent	8.82	16.74	56.23	2025-11-29 13:24:31.985123+00
503f0ba9-97ab-445c-a5a7-18302729d8d6	prod-agent	9.48	16.77	56.23	2025-11-29 13:24:41.992852+00
3bf04a64-0311-4457-8b10-adfc5028305d	prod-agent	9.92	16.73	56.23	2025-11-29 13:24:51.992021+00
b9eb1503-1907-4cda-85c4-4add823b35b8	prod-agent	7.81	16.75	56.23	2025-11-29 13:25:01.992382+00
a98ada77-3a24-4c0e-982d-29ed3f60ee2d	prod-agent	12.36	16.83	56.23	2025-11-29 13:25:11.991528+00
59baa607-9fae-4ed9-94c0-d1e9fb99181a	prod-agent	11.22	16.70	56.23	2025-11-29 13:25:21.992269+00
b740461d-02e1-43f8-9b4c-55393a64ca78	prod-agent	8.93	16.70	56.23	2025-11-29 13:25:31.991588+00
fb96c836-3350-49aa-8b7c-cc0de6c5c6d2	prod-agent	9.20	16.72	56.23	2025-11-29 13:25:41.988797+00
b0cfd521-eeea-4a90-a084-47a4f17b0787	prod-agent	12.93	16.60	56.23	2025-11-29 13:25:51.985648+00
9c6733e6-e42e-4e05-9407-a01bf636c48e	prod-agent	12.26	16.85	56.23	2025-11-29 13:26:01.986353+00
e888ea8f-a9ae-4edc-9f33-dd48b80525da	prod-agent	10.26	16.83	56.23	2025-11-29 13:26:11.991704+00
9c4e54ed-1c0d-45ec-b111-40a3f3ae9ccb	prod-agent	12.60	16.74	56.23	2025-11-29 13:26:21.988293+00
d9dd3a21-3e19-47cc-b038-0b7f9cbdb054	prod-agent	6.58	16.78	56.23	2025-11-29 13:26:31.990896+00
57865525-9dcb-4ece-ac5a-916c4a2c06e7	prod-agent	8.50	16.79	56.23	2025-11-29 13:26:41.991938+00
13d1cdc4-3a35-497e-b7a7-5c1ab8f81513	prod-agent	7.92	16.79	56.23	2025-11-29 13:26:51.984928+00
077b2843-d038-46e2-9767-05622cd92ea6	prod-agent	8.16	16.80	56.23	2025-11-29 13:27:01.991215+00
69dcc25d-467d-443f-b204-dfcc60fd0930	prod-agent	9.38	16.83	56.23	2025-11-29 13:27:11.9952+00
718d2e29-10c5-4241-b364-9b7aede835b0	prod-agent	21.11	16.87	56.23	2025-11-29 13:27:21.988889+00
65c9dfb1-c948-4426-b879-384d304fe874	prod-agent	8.02	16.88	56.23	2025-11-29 13:27:31.992356+00
22d1c6dd-14e1-4d0e-b1a7-a4d81744fe60	prod-agent	8.84	16.84	56.23	2025-11-29 13:27:41.991561+00
59e8e319-fbf7-4f04-a915-6031d9b6b6ea	prod-agent	12.08	16.85	56.23	2025-11-29 13:27:51.988826+00
ec4c5d1e-8474-4b24-80cd-24ba1a3a7596	prod-agent	13.22	16.81	56.23	2025-11-29 13:28:01.995502+00
7a04ba48-b389-4353-86e2-1b7de48bb59a	prod-agent	13.43	16.83	56.23	2025-11-29 13:28:11.995524+00
64e1cb47-5762-4c91-a7e2-06643677d3d9	prod-agent	13.73	16.85	56.23	2025-11-29 13:28:21.994263+00
8ff728a7-3c11-4c47-93b0-7adaf7e8a337	prod-agent	7.75	16.82	56.23	2025-11-29 13:28:31.989228+00
33097859-7d86-486d-856b-4cf3b97b6607	prod-agent	8.84	16.77	56.23	2025-11-29 13:28:41.993943+00
a5c96ae3-db41-41bc-b60b-329f2565a1cb	prod-agent	8.46	16.80	56.23	2025-11-29 13:28:51.988785+00
9f676e96-4bfe-486d-aba2-181bd1692f47	prod-agent	9.21	16.83	56.23	2025-11-29 13:29:01.992341+00
b80ab58f-4ce6-44f7-bb4e-e1af2f6547a4	prod-agent	9.42	16.87	56.23	2025-11-29 13:29:11.991397+00
36d3f37f-f62c-4629-9b0e-ed515a9aaf0d	prod-agent	13.56	16.85	56.23	2025-11-29 13:29:21.994903+00
106f9a4e-e92b-4871-98e6-5604716ca454	prod-agent	7.26	16.79	56.23	2025-11-29 13:29:31.991667+00
3c44e190-4017-4f94-82c5-b78611d7597b	prod-agent	8.51	16.85	56.23	2025-11-29 13:29:41.993649+00
47c5df38-6a8b-48c6-ae51-ba56cea98dbb	prod-agent	8.84	16.84	56.23	2025-11-29 13:29:51.992583+00
28e79942-47b0-48ed-8f94-48b383083c4e	prod-agent	7.90	16.84	56.23	2025-11-29 13:30:01.991118+00
4d0f007c-303f-4c63-95a6-10b8e0f95e91	prod-agent	8.95	16.85	56.23	2025-11-29 13:30:11.991963+00
a3270b6d-e87e-4ccb-84a0-ec144dcf1660	prod-agent	13.19	16.95	56.23	2025-11-29 13:30:21.988159+00
35d16050-93e0-422c-8b13-2cdb88a11274	prod-agent	9.63	16.80	56.23	2025-11-29 13:30:31.991485+00
8d8001db-25ce-407d-9665-00b3293c05d0	prod-agent	7.82	16.83	56.23	2025-11-29 13:30:41.994203+00
58dda08f-868c-4cc6-9d2f-03ef8daaf2eb	prod-agent	7.95	16.84	56.23	2025-11-29 13:30:51.985443+00
0a51347f-dcb6-4853-8d6b-bf1c491808f8	prod-agent	8.52	16.75	56.23	2025-11-29 13:31:01.992227+00
5e96ec3c-014d-43cc-afbb-ef02db7c77e0	prod-agent	8.23	16.86	56.23	2025-11-29 13:31:11.987032+00
b8b3973c-b5bf-4ceb-92fd-68f3be54c6eb	prod-agent	9.52	16.80	56.23	2025-11-29 13:31:21.992288+00
5c46c3a8-1982-4b16-83f1-493c014b8f4e	prod-agent	12.46	16.84	56.23	2025-11-29 13:31:31.98564+00
6f7a7f49-a214-46cd-940a-95ef5130ae47	prod-agent	8.30	16.81	56.23	2025-11-29 13:31:41.989725+00
d0ec55d3-016d-45dc-878a-f26b5f7a118e	prod-agent	9.27	16.85	56.23	2025-11-29 13:31:51.985748+00
a857c338-743f-4f57-abc1-31c64053c556	prod-agent	7.86	16.87	56.23	2025-11-29 13:32:01.991624+00
f88bca83-6366-424b-abd2-c89c5be3a0a0	prod-agent	9.81	16.79	56.23	2025-11-29 13:32:11.989758+00
ca3e995f-a8e0-40ea-99e5-383c37d60bdb	prod-agent	13.30	16.77	56.23	2025-11-29 13:32:21.991311+00
31cf4151-cb53-42e5-9c7e-a494e0ce5a6f	prod-agent	12.94	16.80	56.23	2025-11-29 13:32:31.994047+00
ec87fd3e-68cf-4dd6-927b-7ae1b4734efb	prod-agent	8.10	16.77	56.23	2025-11-29 13:32:41.992133+00
e844f1df-642a-456c-ba6b-7165165a47b1	prod-agent	13.28	16.73	56.23	2025-11-29 13:32:51.985408+00
e58e47d9-bcfc-46c7-a75d-bed0e04ed732	prod-agent	10.25	16.80	56.23	2025-11-29 13:33:02.0549+00
8ad3cb4a-45a1-4121-b5f6-65079f31ddc6	prod-agent	12.08	16.73	56.23	2025-11-29 13:33:11.990905+00
c7fd88d1-d80c-4be1-a914-b2e9e17193d2	prod-agent	28.00	17.42	56.23	2025-11-29 13:33:21.994496+00
816caea5-3b06-4814-be6e-615282e4407b	prod-agent	57.99	17.52	56.23	2025-11-29 13:33:31.996366+00
\.


--
-- Data for Name: system_settings; Type: TABLE DATA; Schema: public; Owner: nas_user
--

COPY public.system_settings (key, value, updated_at) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: nas_user
--

COPY public.users (id, username, email, password_hash, email_verified, verified_at, created_at, updated_at, role) FROM stdin;
54d22f82-ac3b-4316-9ebc-5052ca22bcb2	felix	freund_felix@icloud.com	$2a$12$SenAlfD52vKiWf/R4MxcmO0MRF1YAX5ZOuBoLjOFPZ3T8IHNgVu7.	f	\N	2025-11-27 20:16:30.339376+00	2025-11-27 20:16:30.339376+00	user
f4dd36e7-0329-4838-b062-4c34b7579ec3	testuser	test@example.com	$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyE3H.6K0s3i	f	\N	2025-11-27 20:11:24.903834+00	2025-11-29 13:00:48.312181+00	admin
\.


--
-- Name: monitoring_samples monitoring_samples_pkey; Type: CONSTRAINT; Schema: public; Owner: nas_user
--

ALTER TABLE ONLY public.monitoring_samples
    ADD CONSTRAINT monitoring_samples_pkey PRIMARY KEY (id);


--
-- Name: refresh_tokens refresh_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: nas_user
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id);


--
-- Name: system_alerts system_alerts_pkey; Type: CONSTRAINT; Schema: public; Owner: nas_user
--

ALTER TABLE ONLY public.system_alerts
    ADD CONSTRAINT system_alerts_pkey PRIMARY KEY (id);


--
-- Name: system_metrics system_metrics_pkey; Type: CONSTRAINT; Schema: public; Owner: nas_user
--

ALTER TABLE ONLY public.system_metrics
    ADD CONSTRAINT system_metrics_pkey PRIMARY KEY (id);


--
-- Name: system_settings system_settings_pkey; Type: CONSTRAINT; Schema: public; Owner: nas_user
--

ALTER TABLE ONLY public.system_settings
    ADD CONSTRAINT system_settings_pkey PRIMARY KEY (key);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: nas_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: nas_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: nas_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: idx_monitoring_samples_created_at; Type: INDEX; Schema: public; Owner: nas_user
--

CREATE INDEX idx_monitoring_samples_created_at ON public.monitoring_samples USING btree (created_at DESC);


--
-- Name: idx_refresh_tokens_token_hash; Type: INDEX; Schema: public; Owner: nas_user
--

CREATE INDEX idx_refresh_tokens_token_hash ON public.refresh_tokens USING btree (token_hash);


--
-- Name: idx_refresh_tokens_user_id; Type: INDEX; Schema: public; Owner: nas_user
--

CREATE INDEX idx_refresh_tokens_user_id ON public.refresh_tokens USING btree (user_id);


--
-- Name: idx_system_alerts_open; Type: INDEX; Schema: public; Owner: nas_user
--

CREATE INDEX idx_system_alerts_open ON public.system_alerts USING btree (is_resolved, created_at DESC);


--
-- Name: idx_system_metrics_created_at; Type: INDEX; Schema: public; Owner: nas_user
--

CREATE INDEX idx_system_metrics_created_at ON public.system_metrics USING btree (created_at DESC);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: nas_user
--

CREATE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: idx_users_email_verified; Type: INDEX; Schema: public; Owner: nas_user
--

CREATE INDEX idx_users_email_verified ON public.users USING btree (email_verified);


--
-- Name: idx_users_role; Type: INDEX; Schema: public; Owner: nas_user
--

CREATE INDEX idx_users_role ON public.users USING btree (role);


--
-- Name: users update_users_updated_at; Type: TRIGGER; Schema: public; Owner: nas_user
--

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: refresh_tokens fk_user; Type: FK CONSTRAINT; Schema: public; Owner: nas_user
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: refresh_tokens refresh_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: nas_user
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict X1vg6fAut3rTk8KhZt0ByhpGUo9w9V89EyWOOpEaHhcGcpyCUR9kRZvW5mok4lR

