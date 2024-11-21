package settings_loader

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"log"
)

// SecurityConfig armazena configurações relacionadas à segurança (chave de criptografia do cookie e segredo do JWT)
type SecurityConfig struct {
	CookieEncryptionKey string
	JWTSecret           string
}

// DatabaseConfig armazena a URL do banco de dados
type DatabaseConfig struct {
	DatabaseURL string
}

type PathConfig struct {
	LogoPath            string
	FavIconPath         string
	AdminWebPanelPath   string
	StudentWebPanelPath string
	EmailImagesRootPath string
	FileServerRootPath  string
	RedirectUrl         string
	HTMLRootPath        string
}

type TLSConfig struct {
	IsTLS bool
	Cert  string
	Key   string
}

type ServerDomainConfig struct {
	ServerDomain   string
	UseCloudSql    string
	EnableRedirect bool
}

// SettingsLoader carrega todas as configurações do ambiente
type SettingsLoader struct {
	SecurityConfig     SecurityConfig
	DatabaseConfig     DatabaseConfig
	PathConfig         PathConfig
	TLSConfig          TLSConfig
	ServerDomainConfig ServerDomainConfig
}

// NewSettingsLoader cria uma nova instância do SettingsLoader e carrega as configurações do ambiente
func NewSettingsLoader() *SettingsLoader {
	config, err := toml.LoadFile("settings.toml")
	if err != nil {
		log.Fatalf("Erro ao carregar settings.toml: %v", err)
	}

	return &SettingsLoader{
		SecurityConfig: SecurityConfig{
			CookieEncryptionKey: config.Get("SecurityConfig.COOKIE_ENCRYPTION_KEY").(string),
			JWTSecret:           config.Get("SecurityConfig.JWT_SECRET").(string),
		},
		DatabaseConfig: DatabaseConfig{
			DatabaseURL: config.Get("DatabaseConfig.DATABASE_URL").(string),
		},
	}
}

// GetSecurityConfig retorna as configurações de segurança
func (s *SettingsLoader) GetSecurityConfig() SecurityConfig {
	return s.SecurityConfig
}

// GetDatabaseConfig retorna as configurações do banco de dados
func (s *SettingsLoader) GetDatabaseConfig() DatabaseConfig {
	return s.DatabaseConfig
}

func (s *SettingsLoader) GetPathConfig() (PathConfig, error) {
	return s.PathConfig, nil
}

func (s *SettingsLoader) GetFullDomain() (string, error) {
	if s.TLSConfig.IsTLS {
		return fmt.Sprintf("https://%s", s.ServerDomainConfig.ServerDomain), nil
	} else {
		return fmt.Sprintf("http://%s", s.ServerDomainConfig.ServerDomain), nil
	}
}
