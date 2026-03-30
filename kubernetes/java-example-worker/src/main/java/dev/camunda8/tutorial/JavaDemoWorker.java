package dev.camunda8.tutorial;

import io.camunda.zeebe.client.ZeebeClient;
import io.camunda.zeebe.client.api.worker.JobClient;
import io.camunda.zeebe.client.api.response.ActivatedJob;

import java.util.LinkedHashMap;
import java.util.Map;
import java.util.concurrent.CountDownLatch;

/**
 * Minimalny job worker: domyślny typ zadania {@code c8jw-java} (zmienna środowiskowa JOB_TYPE).
 * W Modeler pole Service Task Type = c8jw-java (lub wartość JOB_TYPE).
 */
public final class JavaDemoWorker {

  public static void main(String[] args) throws InterruptedException {
    String gateway = System.getenv("ZEEBE_ADDRESS");
    if (gateway == null || gateway.isBlank()) {
      System.err.println("ZEEBE_ADDRESS is required");
      System.exit(1);
    }
    String jobType = System.getenv().getOrDefault("JOB_TYPE", "c8jw-java");

    CountDownLatch shutdown = new CountDownLatch(1);
    Runtime.getRuntime().addShutdownHook(new Thread(shutdown::countDown));

    try (ZeebeClient client =
        ZeebeClient.newClientBuilder().gatewayAddress(gateway).usePlaintext().build()) {

      client
          .newWorker()
          .jobType(jobType)
          .handler(JavaDemoWorker::handle)
          .name("c8jw-java")
          .open();

      System.out.println("listening gateway=" + gateway + " job_type=" + jobType);
      shutdown.await();
    }
  }

  private static void handle(JobClient jobClient, ActivatedJob job) {
    System.out.println(
        "job activated key="
            + job.getKey()
            + " processInstanceKey="
            + job.getProcessInstanceKey()
            + " variables="
            + job.getVariables());

    Map<String, Object> out = new LinkedHashMap<>();
    out.put("fromJava", true);
    out.put("worker", "c8jw-java");

    jobClient.newCompleteCommand(job.getKey()).variables(out).send().join();
    System.out.println("job completed key=" + job.getKey());
  }
}
