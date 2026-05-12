---
paths: "**/*"
---

# Git Push Restrictions — Protected Branches

## Rule: No Direct Push to Master/Main

**Severity:** CRITICAL (enforced by git pre-push hook + future GitHub branch protection)

### Protected Branches

- `master`
- `main`

### What This Means

❌ **Não é permitido:**
```bash
git checkout master
git push origin master         # BLOQUEADO pelo hook
git push --force origin master # BLOQUEADO pelo hook
```

✅ **Fluxo correto:**
```bash
# 1. Desenvolver em feature branch
git checkout -b feature/minha-feature
git push origin feature/minha-feature

# 2. Criar PR (apenas @devops)
gh pr create --title "Minha PR" --body "Descrição"

# 3. @devops valida e faz merge/push
@devops *push
```

### Por Que?

- **Integridade:** Master/main são branches de produção
- **Rastreabilidade:** Todas as mudanças devem ter histórico de PR
- **Qualidade:** PRs garantem revisão + testes antes de merge
- **Safety:** Git hook local + futura branch protection GitHub = defesa em profundidade

### Enforcement Mechanism

**Local (Git Pre-Push Hook):**
- Arquivo: `.git/hooks/pre-push`
- Ativado: Automaticamente em todo `git push`
- Instalado: 2026-05-12
- Comportamento: Bloqueia push e exibe mensagem educativa

**Future (GitHub Branch Protection):**
- Será configurado quando versão Pro estiver disponível
- Ativa: Impossível contornar mesmo com force push

### Se Você Tentar Push para Master

```
❌ PUSH BLOQUEADO: Não é permitido push direto para 'master'

   ℹ️  Fluxo correto:
      1. git push origin feature/seu-branch
      2. gh pr create --title 'Sua PR' --body 'Descrição'
      3. @devops faz merge após validação

   📚 Referência: .claude/rules/git-push-restrictions.md
```

### Agent Authority

**ONLY @devops can push** (conforme `.claude/rules/agent-authority.md`):
- Responsável por `git push` (exclusive)
- Responsável por `gh pr create` / `gh pr merge` (exclusive)
- Valida quality gates antes de qualquer push
- Documenta todas as operações

### What Happens When You Violate

1. **Git hook blocks the push:** ❌ Não deixa acontecer
2. **Local:** Se usar `git push --force`: Aviso + bloqueio
3. **Future (GitHub Pro):** Branch protection rules bloqueia no servidor

### Exception Process

Se há razão **legítima** para bypassar (emergência, etc):
1. Contactar @devops explicitamente
2. @devops evalua motivo
3. Se aprovado: manual push com documentação
4. Sempre log a exceção em `EXCEPTIONS.md`

---

## Related Documents

- `.claude/rules/agent-authority.md` — Agent push authority
- `.claude/CLAUDE.md` — Push Authority section
- `.git/hooks/pre-push` — Implementation

## Timeline

| Data | Evento |
|------|--------|
| 2026-05-12 | Git pre-push hook instalado |
| TBD | GitHub branch protection rules (versão Pro) |
| TBD | Workflow CI/CD com status checks obrigatórios |

---

**Última Atualização:** 2026-05-12  
**Responsável:** @devops  
**Status:** ✅ Ativo
