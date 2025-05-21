package config

import (
	"log"
	"os"
	"path/filepath" // Para manipulação de caminhos de forma portável
	"runtime"       // Para obter informações sobre o ambiente de execução
	"strconv"       // Para converter strings para inteiros

	"github.com/joho/godotenv" // Para carregar variáveis de ambiente de um arquivo .env
)

// Config armazena todas as configurações da aplicação.
// Os valores são lidos de variáveis de ambiente, com a possibilidade de usar um arquivo .env para desenvolvimento.
type Config struct {
	DBDriver          string // Driver do banco de dados (ex: "postgres")
	DBSource          string // String de conexão com o banco de dados (DSN)
	ServerPort        string // Porta em que o servidor HTTP vai rodar
	JWTSecret         string // Segredo usado para assinar e verificar tokens JWT
	JWTExpirationHours int    // Tempo de expiração para tokens JWT em horas
	// Adicione outras configurações que sua aplicação possa precisar aqui
	// Ex: LogLevel string, ApiKeyExterna string, etc.
}

// LoadConfig carrega as configurações da aplicação.
// Ele tenta carregar um arquivo .env da raiz do projeto Go primeiro,
// e depois lê as variáveis de ambiente do sistema (que podem sobrescrever as do .env).
func LoadConfig(customEnvPath ...string) *Config {
	var envPath string

	if len(customEnvPath) > 0 && customEnvPath[0] != "" {
		// Se um caminho customizado para o .env for fornecido, usa ele
		envPath = customEnvPath[0]
	} else {
		// Tenta determinar o caminho para o arquivo .env na raiz do projeto Go.
		// Isso é útil porque 'go run' pode ser executado de subdiretórios (ex: cmd/myapp).
		// runtime.Caller(0) retorna o path do arquivo atual (config.go).
		_, b, _, ok := runtime.Caller(0)
		if !ok {
			log.Println("Aviso: Não foi possível determinar o caminho do arquivo config.go para localizar .env. Tentando carregar .env do diretório atual.")
			// Se não conseguir determinar, tenta carregar ".env" do diretório de trabalho atual
			// Isso pode ou não ser o que você quer, dependendo de onde 'go run' é chamado.
			cwd, err := os.Getwd()
			if err != nil {
				log.Printf("Aviso: Não foi possível obter o diretório de trabalho atual: %v. Tentando '.env' relativo.", err)
				envPath = ".env"
			} else {
				envPath = filepath.Join(cwd, ".env")
			}
		} else {
			// b é o path completo para este arquivo (internal/config/config.go)
			// Queremos subir dois níveis para chegar à raiz do projeto Go (onde o .env geralmente está)
			// Ex: backend_go/internal/config/config.go -> backend_go/
			projectRoot := filepath.Join(filepath.Dir(b), "..", "..")
			envPath = filepath.Join(projectRoot, ".env")
		}
	}

	log.Printf("Tentando carregar arquivo .env de: %s", envPath)

	// Tenta carregar o arquivo .env. Não é um erro fatal se não encontrar,
	// pois as variáveis podem estar definidas diretamente no ambiente do sistema.
	err := godotenv.Load(envPath)
	if err == nil {
		log.Printf("Arquivo .env carregado com sucesso de: %s", envPath)
	} else {
		// Se o erro NÃO for "arquivo não encontrado", pode ser um problema de permissão ou formato.
		// Se for "arquivo não encontrado", é apenas um aviso.
		if !os.IsNotExist(err) {
			log.Printf("Aviso: Erro ao carregar arquivo %s: %v. As variáveis de ambiente do sistema ainda podem ser usadas.", envPath, err)
		} else {
			log.Printf("Aviso: Arquivo .env não encontrado em %s. Usando apenas variáveis de ambiente do sistema (se definidas).", envPath)
		}
	}

	// Lê as variáveis de ambiente (ou os valores padrão se não encontradas)
	cfg := &Config{
		DBDriver:   getEnv("DB_DRIVER", "postgres"),
		DBSource:   getEnv("DB_SOURCE", "host=localhost user=postgres password=secret dbname=bizly_db port=5432 sslmode=disable TimeZone=America/Sao_Paulo"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		JWTSecret:  getEnv("JWT_SECRET", "seu-jwt-segredo-muito-secreto-e-longo-e-aleatorio"),
		JWTExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 72), // Padrão de 72 horas (3 dias)
		// Adicione aqui a leitura de outras variáveis de ambiente
	}

	// Você pode adicionar logs aqui para verificar os valores carregados (CUIDADO para não logar segredos em produção)
	// log.Printf("Configurações carregadas: ServerPort=%s, DBDriver=%s", cfg.ServerPort, cfg.DBDriver)

	return cfg
}

// getEnv é uma função helper para ler uma variável de ambiente.
// Se a variável não estiver definida, retorna o valor fallback fornecido.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// getEnvAsInt é uma função helper para ler uma variável de ambiente como um inteiro.
// Se a variável não estiver definida ou não for um inteiro válido, retorna o valor fallback.
func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	log.Printf("Aviso: Variável de ambiente '%s' não é um inteiro válido ('%s') ou não está definida. Usando fallback: %d", key, valueStr, fallback)
	return fallback
}

// getEnvAsBool é uma função helper para ler uma variável de ambiente como um booleano.
// Se a variável não estiver definida, retorna o valor fallback.
// Considera "true", "1", "yes" como true (case-insensitive).
func getEnvAsBool(key string, fallback bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return fallback
	}
	val, err := strconv.ParseBool(valueStr)
	if err != nil {
		log.Printf("Aviso: Variável de ambiente '%s' não é um booleano válido ('%s'). Usando fallback: %t", key, valueStr, fallback)
		return fallback
	}
	return val
}