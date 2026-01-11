-- Create sys_settings table
-- Stores application-wide settings and configurations
CREATE TABLE IF NOT EXISTS sys_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key using UUID
    
    -- Basic Info
    setting_no VARCHAR(30) UNIQUE, -- Unique setting identifier
    name VARCHAR(150),             -- Organization name
    phone VARCHAR(100),            -- Primary contact phone
    fax VARCHAR(100),              -- Organization fax number
    address TEXT,                   -- Physical office address
    image_path VARCHAR(255),        -- General image storage path
    logo_file_name TEXT,            -- Organization logo filename
    logo_file_path TEXT,            -- Storage path for the logo
    
    -- Member Card Images
    card_member_image_name TEXT,     -- Filename for member card design
    card_member_image_path TEXT,     -- Storage path for member card design
    card_member_min_age INTEGER DEFAULT 0, -- Minimum age for card eligibility
    card_member_max_age INTEGER DEFAULT 0, -- Maximum age for card eligibility
    
    -- Target Settings
    target_enabled BOOLEAN DEFAULT false, -- Whether business targets are active
    target_month DECIMAL(15,2) DEFAULT 0, -- Monthly target goal
    target_year DECIMAL(15,2) DEFAULT 0,  -- Yearly target goal
    target_quarter1 DECIMAL(15,2) DEFAULT 0, -- Q1 target goal
    target_quarter2 DECIMAL(15,2) DEFAULT 0, -- Q2 target goal
    target_quarter3 DECIMAL(15,2) DEFAULT 0, -- Q3 target goal
    target_quarter4 DECIMAL(15,2) DEFAULT 0, -- Q4 target goal
    
    -- NAB Switching Settings (Net Asset Value Switching)
    nab_switching_max_years SMALLINT DEFAULT 0,      -- Maximum years allowed for NAV switching
    nab_switching_max_year_amount INTEGER DEFAULT 0, -- Maximum annual amount for NAV switching
    nab_switching_enabled BOOLEAN DEFAULT false,     -- Whether NAV switching is enabled
    nab_switching_amount INTEGER DEFAULT 0,          -- Default amount for NAV switching
    nab_switching_add_enabled BOOLEAN DEFAULT false, -- Whether additional NAV switching is enabled
    nab_switching_add_amount INTEGER DEFAULT 0,      -- Additional amount for NAV switching
    
    -- Member Approval Settings (Approval Peserta)
    member_approval_effective_date DATE,          -- Effective date for member approval rules
    member_approval_skip_holiday BOOLEAN DEFAULT false, -- Whether to skip holidays in approval time
    member_approval_limit VARCHAR(10),            -- Daily/count limit for member approvals
    
    -- Company Verification (Verifikasi Perusahaan)
    company_verify_enabled BOOLEAN DEFAULT false, -- Whether company verification is enabled
    company_verify_checkers VARCHAR(100),         -- Role/User ID list for checkers
    company_verify_approvers VARCHAR(100),        -- Role/User ID list for approvers
    
    -- Company Member Verification (Verifikasi Peserta Perusahaan)
    company_member_verify_enabled BOOLEAN DEFAULT false, -- Whether company member verification is enabled
    company_member_verify_checkers VARCHAR(100),         -- Checker roles for company members
    company_member_verify_approvers VARCHAR(100),        -- Approver roles for company members
    
    -- Individual Verification (Verifikasi Peserta Individu)
    individual_verify_enabled BOOLEAN DEFAULT false, -- Whether individual verification is enabled
    individual_verify_checkers VARCHAR(100),         -- Checker roles for individuals
    individual_verify_approvers VARCHAR(100),        -- Approver roles for individuals
    
    -- Data Change Verification (Verifikasi Perubahan Data)
    data_change_verify_enabled BOOLEAN DEFAULT false, -- Whether data change verification is enabled
    data_change_verify_checkers VARCHAR(100),         -- Checker roles for data changes
    data_change_verify_approvers VARCHAR(100),        -- Approver roles for data changes
    
    -- Product Approval Settings
    product_approval_effective_date DATE,
    product_approval_skip_holiday BOOLEAN DEFAULT false,
    product_approval_limit VARCHAR(10),
    product_nav_effective_date DATE,
    product_nav_skip_holiday BOOLEAN DEFAULT false,
    product_nav_limit VARCHAR(10),
    product_verify_enabled BOOLEAN DEFAULT false,
    product_verify_checkers VARCHAR(100),
    product_verify_approvers VARCHAR(100),
    
    -- Benefit Approval Settings
    benefit_approval_effective_date DATE,
    benefit_approval_skip_holiday BOOLEAN DEFAULT false,
    benefit_approval_limit VARCHAR(10),
    benefit_company_verify_enabled BOOLEAN DEFAULT false,
    benefit_company_verify_checkers VARCHAR(100),
    benefit_company_verify_approvers VARCHAR(100),
    benefit_company_member_verify_enabled BOOLEAN DEFAULT false,
    benefit_company_member_verify_checkers VARCHAR(100),
    benefit_company_member_verify_approvers VARCHAR(100),
    benefit_individual_verify_enabled BOOLEAN DEFAULT false,
    benefit_individual_verify_checkers VARCHAR(100),
    benefit_individual_verify_approvers VARCHAR(100),
    
    -- Dues Approval Settings
    dues_approval_effective_date DATE,
    dues_approval_skip_holiday BOOLEAN DEFAULT false,
    dues_approval_limit VARCHAR(10),
    dues_company_verify_enabled BOOLEAN DEFAULT false,
    dues_company_verify_checkers VARCHAR(100),
    dues_company_verify_approvers VARCHAR(100),
    dues_company_member_verify_enabled BOOLEAN DEFAULT false,
    dues_company_member_verify_checkers VARCHAR(100),
    dues_company_member_verify_approvers VARCHAR(100),
    dues_individual_verify_enabled BOOLEAN DEFAULT false,
    dues_individual_verify_checkers VARCHAR(100),
    dues_individual_verify_approvers VARCHAR(100),
    
    -- Entertainment Budget Approval (0=disabled, 1=1 level, 2=2 levels)
    entertainment_budget_approval_level SMALLINT DEFAULT 0,
    
    -- Investment Products Settings (Pengaturan Produk Investasi)
    investment_products_enabled BOOLEAN DEFAULT false,       -- Overall investment feature flag
    investment_products_flag BOOLEAN DEFAULT false,          -- Internal processing flag
    investment_products_verify_enabled BOOLEAN DEFAULT false, -- Verification for investment changes
    investment_products_verify_checkers VARCHAR(100),        -- Checker role IDs
    investment_products_approval_enabled BOOLEAN DEFAULT false, -- Approval for investment changes
    investment_products_verify_approvers VARCHAR(100),       -- Approver role IDs
    investment_defined_contribution_pension_enabled BOOLEAN DEFAULT false, -- PPIP (Program Pensiun Iuran Pasti)
    investment_min_enabled BOOLEAN DEFAULT false,           -- Minimum investment threshold active
    investment_min_amount DECIMAL(15,2) DEFAULT 0,          -- Minimum amount for investment
    investment_topup_enabled BOOLEAN DEFAULT false,          -- Top-up capability enabled
    investment_topup_amount DECIMAL(15,2) DEFAULT 0,         -- Minimum top-up amount
    
    -- Investment Switching Settings (Pengaturan Switching Investasi)
    investment_switching_max_years SMALLINT DEFAULT 0,      -- Max years allowed for switching
    investment_switching_max_year_amount INTEGER DEFAULT 0, -- Annual count limit
    investment_switching_enabled BOOLEAN DEFAULT false,    -- Switching feature active
    investment_switching_checked BOOLEAN DEFAULT false,    -- Mandatory verification flag
    investment_switching_amount INTEGER DEFAULT 0,          -- Default switching fee/count
    investment_switching_add_enabled BOOLEAN DEFAULT false, -- Additional switching active
    investment_switching_add_checked BOOLEAN DEFAULT false, -- Mandatory additional check
    investment_switching_add_amount INTEGER DEFAULT 0,      -- Additional switching amount
    
    -- Individual Withdrawal Settings (Pengaturan Penarikan Individu)
    individual_withdrawal_enabled BOOLEAN DEFAULT false,     -- Individual withdrawal active
    individual_withdrawal_amount DECIMAL(15,2) DEFAULT 0,    -- Max withdrawal amount
    individual_withdrawal_time_days INTEGER DEFAULT 0,       -- Cooldown period in days
    individual_withdrawal_min_age INTEGER DEFAULT 0,         -- Min age for withdrawal
    individual_withdrawal_max_age INTEGER DEFAULT 0,         -- Max age for withdrawal
    individual_withdrawal_type_enabled BOOLEAN DEFAULT false, -- Tiered withdrawal active
    individual_withdrawal_type SMALLINT DEFAULT 0,           -- Type of tiered withdrawal
    individual_withdrawal_min INTEGER DEFAULT 0,             -- Tiered min amount
    individual_withdrawal_max INTEGER DEFAULT 0,             -- Tiered max amount
    individual_withdrawal_thawing_enabled BOOLEAN DEFAULT false, -- Fund thawing period active
    individual_withdrawal_thawing SMALLINT DEFAULT 0,          -- Thawing period length
    individual_withdrawal_join_account BOOLEAN DEFAULT false, -- Joint account withdrawal flag
    individual_withdrawal_members_min SMALLINT DEFAULT 0,    -- Min member duration for withdrawal
    corporate_withdrawal_reason SMALLINT DEFAULT 0,          -- Default reason for corporate withdrawal
    
    -- Managed Funds Cost Settings (Biaya Pengelolaan Dana)
    managed_funds_cost_apply BOOLEAN DEFAULT false,         -- Apply management fee
    managed_funds_cost_period SMALLINT DEFAULT 0,          -- Billing period (e.g., monthly)
    managed_funds_cost_type SMALLINT DEFAULT 0,            -- Fee calculation type
    managed_funds_cost_amount DECIMAL(15,2) DEFAULT 0,     -- Fee amount or rate
    managed_funds_withdrawal_date DATE,                      -- Date for cost withdrawal
    managed_funds_system_date DATE,                          -- System reference date for costs
    managed_funds_withdrawal_type SMALLINT DEFAULT 0,        -- Type of cost withdrawal
    managed_funds_withdrawal_percent_type SMALLINT DEFAULT 0, -- Rate type
    managed_funds_withdrawal_amount DECIMAL(15,2) DEFAULT 0,  -- Amount to withdraw
    managed_funds_buy_type SMALLINT DEFAULT 0,             -- Subscription fee type
    managed_funds_buy_percent SMALLINT DEFAULT 0,          -- Subscription rate
    managed_funds_buy_amount DECIMAL(15,2) DEFAULT 0,      -- Subscription amount
    managed_funds_diversion_type SMALLINT DEFAULT 0,       -- Diversion fee type
    managed_funds_diversion_percent_type SMALLINT DEFAULT 0, -- Diversion rate type
    managed_funds_diversion_amount DECIMAL(15,2) DEFAULT 0,  -- Diversion fee amount
    managed_funds_others_type SMALLINT DEFAULT 0,          -- Miscellaneous fee type
    managed_funds_others_percent_type SMALLINT DEFAULT 0,  -- Misc rate type
    managed_funds_others_amount DECIMAL(15,2) DEFAULT 0,   -- Misc fee amount
    
    -- Special Purpose Cost Settings (Biaya Khusus / SPC)
    special_purpose_fee_withdrawal_date DATE,         -- Date for SPC fee withdrawal
    special_purpose_apply_system_date DATE,           -- Effective date for SPC calculations
    special_purpose_period SMALLINT DEFAULT 0,        -- Recurring period for SPC fees
    special_purpose_payment_type SMALLINT DEFAULT 0,  -- Method of SPC payment
    special_purpose_payment_amount DECIMAL(15,2),     -- Fixed amount for SPC payment
    special_purpose_admin_type SMALLINT DEFAULT 0,    -- Calculation type for admin fees
    special_purpose_admin_enabled BOOLEAN DEFAULT false, -- Enable SPC admin fees
    special_purpose_admin_amount DECIMAL(15,2) DEFAULT 0, -- SPC admin fee amount
    special_purpose_admin_fee_date DATE,              -- Date for admin fee charge
    special_purpose_admin_apply_date DATE,            -- Effective date for admin fee rules
    special_purpose_withdrawal_funds_enabled BOOLEAN DEFAULT false, -- Fund withdrawal fee active
    special_purpose_withdrawal_funds_type SMALLINT DEFAULT 0,      -- Type of withdrawal fee
    special_purpose_withdrawal_funds_amount DECIMAL(15,2) DEFAULT 0, -- Amount of withdrawal fee
    special_purpose_withdrawal_funds_fee_date DATE,                 -- Date for withdrawal fee charge
    special_purpose_withdrawal_funds_apply_date DATE,               -- Effective date for withdrawal fee
    special_purpose_transfer_funds_enabled BOOLEAN DEFAULT false, -- Fund transfer fee active
    special_purpose_transfer_funds_type SMALLINT DEFAULT 0,       -- Type of transfer fee
    special_purpose_transfer_funds_amount DECIMAL(15,2) DEFAULT 0, -- Amount of transfer fee
    special_purpose_transfer_funds_fee_date DATE,                 -- Date for transfer fee charge
    special_purpose_transfer_funds_apply_date DATE,               -- Effective date for transfer fee
    special_purpose_dues_enabled BOOLEAN DEFAULT false, -- Dues/contribution fee active
    special_purpose_dues_type SMALLINT DEFAULT 0,       -- Type of dues fee
    special_purpose_dues_amount DECIMAL(15,2) DEFAULT 0, -- Amount of dues fee
    
    -- Alert Settings - Member Termination (Notifikasi Berhenti Peserta)
    alert_member_termination BOOLEAN DEFAULT false,         -- Alert for member termination
    alert_member_termination_choice SMALLINT DEFAULT 0,      -- Delivery channel choice
    alert_member_termination_approved BOOLEAN DEFAULT false, -- Alert when termination approved
    alert_member_termination_approved_choice SMALLINT DEFAULT 0, -- Channel for approval alert
    alert_member_checked BOOLEAN DEFAULT false,             -- Mandatory check alert
    alert_member_checked_choice SMALLINT DEFAULT 0,          -- Channel for check alert
    alert_member_checked_approved BOOLEAN DEFAULT false,     -- Alert when check approved
    alert_member_checked_approved_choice SMALLINT DEFAULT 0, -- Channel for check approval
    
    -- Alert Settings - Transactions (Notifikasi Transaksi)
    alert_member_transaction BOOLEAN DEFAULT false,         -- Alert for general transactions
    alert_member_transaction_min DECIMAL(15,2) DEFAULT 0,    -- Min amount to trigger alert
    alert_member_transaction_max DECIMAL(15,2) DEFAULT 0,    -- Max amount for alert trigger
    alert_member_withdrawal BOOLEAN DEFAULT false,          -- Alert for withdrawals
    alert_member_withdrawal_min DECIMAL(15,2) DEFAULT 0,     -- Min withdrawal amount for alert
    alert_member_withdrawal_max DECIMAL(15,2) DEFAULT 0,     -- Max withdrawal amount for alert
    alert_member_email_time_checked INTEGER DEFAULT 0,      -- Frequency of email alerts
    alert_member_email_time_choice INTEGER DEFAULT 0,       -- Time of day for alerts
    alert_member_email_process_checked INTEGER DEFAULT 0,   -- Background process check
    alert_member_email_process VARCHAR(255),               -- Process status name
    
    -- Alert Settings - Investment (Notifikasi Investasi)
    alert_investment_termination BOOLEAN DEFAULT false,     -- Alert for investment termination
    alert_investment_termination_choice SMALLINT DEFAULT 0,  -- Channel for investment termination
    alert_investment_cutoff BOOLEAN DEFAULT false,          -- Cut-off time alert active
    alert_investment_cutoff_time TIME,                       -- daily cut-off time for NAV
    
    -- Alert Settings - Benefits (Notifikasi Manfaat)
    alert_benefit_min_enabled BOOLEAN DEFAULT false,        -- Min benefit amount alert active
    alert_benefit_min_amount DECIMAL(15,2) DEFAULT 0,       -- Min threshold for alert
    alert_benefit_max_enabled BOOLEAN DEFAULT false,        -- Max benefit amount alert active
    alert_benefit_max_amount DECIMAL(15,2) DEFAULT 0,       -- Max threshold for alert
    
    -- Terms & Privacy
    terms_and_conditions TEXT, -- Legal terms and conditions content
    privacy_policy TEXT,       -- Data privacy policy content
    
    -- Audit fields
    created_by UUID, -- User who created this record
    updated_by UUID, -- User who last updated this record
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP  -- Record last update timestamp
);

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_sys_settings_setting_no ON sys_settings(setting_no);
