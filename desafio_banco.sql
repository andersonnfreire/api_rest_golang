

CREATE DATABASE banco_desafio;


CREATE TABLE IF NOT EXISTS `banco_desafio`.`dados`(
  `id` int(11) PRIMARY KEY NOT NULL AUTO_INCREMENT ,
  `nome` varchar(100) NOT NULL,
  `dt_nascimento` date DEFAULT NULL,
  `salario` decimal(15,2) DEFAULT 0.00,
  `qtd_filhos` int(11) DEFAULT 0,
  `sexo` enum('M','F') DEFAULT NULL
);


create table `pacientes` (
    `id` integer primary key,
    `nome` varchar (100) not null,
    `sexo` char (1) not null,
    `dt_nascimento` date
);
 
create table `medicos` (
    `id` integer primary key,
    `nome` varchar (100) not null,
    `sexo` char (1) not null,
    `dt_nascimento` date,
    `salario` decimal (15,2) default 0.00
);
 
create table `consultas` (
    `id` integer primary key,
    `paciente_id` integer not null,
    `medico_id` integer not null,
    `data` date not null,
    `valor` decimal (10,2) default 0.00,
    `pago` tinyint (1) not null,
    `observacao` varchar (500)
);

use banco_desafio;
 
alter table `consultas` add constraint `fk_medico_id` foreign key (`medico_id`)
references `medicos` (`id`);
 
 
alter table `consultas` add constraint `fk_paciente_id` foreign key (`paciente_id`)
references `pacientes` (`id`);
 
 
insert into `pacientes` (`id`, `nome`, `sexo`, `dt_nascimento`)
values (1, 'joao freitas', 'm', '1984-12-12');
 
insert into `pacientes` (`id`, `nome`, `sexo`, `dt_nascimento`)
values (2, 'maria do carmo rodrigues' ,'f', '1986-04-08');
 
insert into `pacientes` (`id`, `nome`, `sexo`, `dt_nascimento`)
values (3, 'suellen rodrigues' ,'f', '1998-06-12');
 
insert into `pacientes` (`id`, `nome`, `sexo`, `dt_nascimento`)
values (4, 'francisco rocha' ,'m', '1964-04-04');
 
insert into `pacientes` (`id`, `nome`, `sexo`, `dt_nascimento`)
values (5, 'francisco cunha' ,'m', '1961-02-03');
 
insert into `pacientes` (`id`, `nome`, `sexo`, `dt_nascimento`)
values (6, 'paulo ferreira' ,'m', '1986-04-24');
 
insert into `medicos` (`id`, `nome`, `sexo`, `dt_nascimento`, `salario`)
values (1, 'joaquim lopes', 'm', '1974-08-08', 12000.68);
 
insert into `medicos` (`id`, `nome`, `sexo`, `dt_nascimento`, `salario`)
values (2,'carla maria correa','f','1977-12-22',10851.90);
 
insert into `medicos` (`id`, `nome`, `sexo`, `dt_nascimento`, `salario`)
values (3,'fernanda abreu','f','1980-12-13',12000);
 
insert into `medicos` (`id`, `nome`, `sexo`, `dt_nascimento`, `salario`)
values (4,'lucas moraes','m','1974-10-10',15004.30);
 
insert into `medicos` (`id`, `nome`, `sexo`, `dt_nascimento`, `salario`)
values (5,'maria das graças foster','f','1975-12-01',10124.13);
 
insert into `medicos` (`id`, `nome`, `sexo`, `dt_nascimento`, `salario`)
values (6,'paulo roberto araujo','m','1961-07-22',9734.80);
 
insert into `consultas`
    (`id`, `paciente_id`, `medico_id`, `data`, `valor`, `pago`, `observacao`)
    values(1,6 ,2 , '2015-08-04',120 , 1,'retorno 3 meses');
 
insert into `consultas`
    (`id`, `paciente_id`, `medico_id`, `data`, `valor`, `pago`)
    values(2,1 ,3 , '2015-08-05',200 , 0);
 
insert into `consultas`
    (`id`, `paciente_id`, `medico_id`, `data`, `valor`, `pago`, `observacao`)
    values(3,2 ,5 , '2015-08-22',120 , 1,'retorno 6 meses');
 
insert into `consultas`
    (`id`, `paciente_id`, `medico_id`, `data`, `valor`, `pago`)
    values(4,5 ,4 , '2015-09-13',200 , 1);
 
insert into `consultas`
    (`id`, `paciente_id`, `medico_id`, `data`, `valor`, `pago`)
    values(5,3 ,3 , '2015-09-13',200 , 0);
 
insert into `consultas`
    (`id`, `paciente_id`, `medico_id`, `data`, `valor`, `pago`)
    values(6,6 ,2 , '2015-08-04',120 , 1);

INSERT INTO dados (nome,dt_nascimento,salario,qtd_filhos,sexo)
VALUES("João",'1977-10-10','7000',2,'M');
INSERT INTO dados (nome,dt_nascimento,salario,qtd_filhos,sexo)
VALUES("Maria",'1967-05-20','15000',0,'F');
INSERT INTO dados (nome,dt_nascimento,salario,qtd_filhos,sexo)
VALUES("José",'1981-06-01','4306.60',4,'M');
INSERT INTO dados (nome,dt_nascimento,salario,qtd_filhos,sexo)
VALUES("Antonio",'1988-08-08','9780',5,'M');


select * from dados;


/* Alterar o salário para R$ 9000,00 de todos do sexo feminino;*/
UPDATE dados
SET salario = '9000'
WHERE sexo = 'F';

/* Alterar o nome Antonio para Caio; */
UPDATE dados
SET nome = 'Caio'
WHERE nome = 'Antonio';

/* Alterar a quantidade de filhos de quem tem 2 ou 5 filhos para 3 filhos; */
UPDATE dados
SET qtd_filhos = 3
WHERE qtd_filhos IN(2,5);

/*Alterar o salário para R$ 13700,00 de quem tem 4 filhos e nasceu em 15/07/1980*/

UPDATE dados
SET salario = '13700'
WHERE qtd_filhos = 4 and dt_nascimento = '1980-07-15';

/* Apagar todos os dados de quem é do sexo feminino; */

DELETE FROM dados WHERE sexo = "F";

/* Apagar todos os dados de quem nasceu em 10/10/1977 */

DELETE FROM dados WHERE dt_nascimento = "1977-10-10";


/* Exibir quais médicos realizaram as consultas mostrando o código e o nome do
mesmo além do código e data data da consulta; */

SELECT m.id as CodigoMedico, m.nome as NomeMedico, c.id as CodigoConsulta, c.data DataConsulta 
FROM medicos as m INNER JOIN consultas as c on c.medico_id = m.id

/* Exibir quais pacientes consultaram, mostrando o código e nome do paciente, código,
data da consulta, se foi paga ou não e o nome do médico que consultou; */

SELECT p.id as CodigoPaciente, p.nome as NomePaciente, c.id as CodigoConsulta, c.data as DataConsulta, c.pago as StatusPagamento, m.nome as NomeMedico 
FROM medicos as m INNER JOIN consultas as c INNER JOIN pacientes as p on c.medico_id = m.id and c.paciente_id = p.id

/* Exibir quantas consultas cada médico tem e o valor total pago das mesmas
juntamente com o código e nome do médico; */

SELECT COUNT(c.medico_id) as QuantidadeDeConsultas, SUM(c.valor) as ValorTotal, c.medico_id as CodigoMedico, m.nome as NomeMedico 
FROM medicos as m INNER JOIN consultas as c on c.medico_id = m.id and c.pago = 1 group by c.medico_id















