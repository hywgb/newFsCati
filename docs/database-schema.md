# 数据库模式与DDL（核心表）

```sql
create table tenants (
  id bigserial primary key,
  name text not null,
  created_at timestamptz default now()
);

create table users (
  id bigserial primary key,
  tenant_id bigint not null references tenants(id),
  username text unique not null,
  role text not null,
  created_at timestamptz default now()
);

create table surveys (
  id bigserial primary key,
  tenant_id bigint not null references tenants(id),
  name text not null,
  version int not null default 1,
  schema_json jsonb not null,
  status text not null default 'draft',
  created_at timestamptz default now()
);

create table campaigns (
  id bigserial primary key,
  tenant_id bigint not null references tenants(id),
  survey_id bigint not null references surveys(id),
  mode varchar(16) not null check (mode in ('preview','predictive','manual')),
  cps int not null default 5,
  concurrency int not null default 50,
  time_windows jsonb not null default '[]'::jsonb,
  status varchar(16) not null default 'stopped',
  created_at timestamptz default now()
);

create table samples (
  id bigserial primary key,
  tenant_id bigint not null references tenants(id),
  phone varchar(32) not null,
  attrs_json jsonb not null default '{}'::jsonb,
  dnc boolean not null default false,
  timezone text,
  priority int not null default 0,
  unique(tenant_id, phone)
);

create table campaign_samples (
  id bigserial primary key,
  campaign_id bigint not null references campaigns(id),
  sample_id bigint not null references samples(id),
  state varchar(24) not null default 'new',
  last_result varchar(32),
  attempt_count int not null default 0,
  next_time timestamptz,
  created_at timestamptz default now()
);
create index idx_cs_camp_state on campaign_samples(campaign_id, state, next_time);

create table call_attempts (
  id bigserial primary key,
  campaign_id bigint not null references campaigns(id),
  sample_id bigint not null references samples(id),
  agent_id bigint,
  uuid uuid,
  started_at timestamptz,
  ended_at timestamptz,
  result_code varchar(32),
  recording_uri text,
  paradata_json jsonb,
  asr_result varchar(32),
  asr_confidence numeric(5,4),
  asr_latency_ms int,
  asr_transcript text,
  asr_mode varchar(16),
  asr_fallback boolean default false
);

create table quotas (
  id bigserial primary key,
  campaign_id bigint not null references campaigns(id),
  dimension_json jsonb not null,
  target int not null,
  filled int not null default 0
);

create table agents (
  id bigserial primary key,
  tenant_id bigint not null references tenants(id),
  sip_user text not null,
  skills jsonb,
  state varchar(16) not null default 'offline'
);

create table dispositions (
  id bigserial primary key,
  code text unique not null,
  description text,
  final boolean not null default false,
  callbackable boolean not null default false
);

create table qa_scores (
  id bigserial primary key,
  call_attempt_id bigint not null references call_attempts(id),
  score int not null,
  tags text[],
  comments text
);

create table audit_logs (
  id bigserial primary key,
  tenant_id bigint,
  user_id bigint,
  action text not null,
  entity text,
  entity_id text,
  detail jsonb,
  created_at timestamptz default now()
);
```